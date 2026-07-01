package lucario

import (
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
)

var (
	ErrWrite    = errors.New("error while writing to log file")
	ErrEndOfWAL = errors.New("end of WAL")
)

type WAL struct {
	file    *os.File
	currLSN uint64
	mutex   *sync.Mutex
}

func NewWAL(filePath string) (*WAL, error) {

	var file *os.File
	var fileExists bool

	_, err := os.Stat(filePath)

	if err != nil {
		if os.IsNotExist(err) {

			file, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				return nil, err
			}
			fileExists = false
		} else {

			return nil, err
		}
	} else {

		file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		fileExists = true
	}

	stat, err := os.Stat(filePath)

	if err != nil {
		return nil, err
	}

	wal := &WAL{
		file:  file,
		mutex: &sync.Mutex{},
	}

	if !fileExists || stat.Size() == 0 {
		slog.Info("WAL file does not exist")
		wal.currLSN = 0
	} else {
		slog.Info("WAL file exists")
		lastWALRecordLengthBytes := make([]byte, 8)

		slog.Info("SEEK", "TO", stat.Size()-8)
		_, err = wal.file.Seek(stat.Size()-8, io.SeekStart)
		if err != nil {
			return nil, err
		}

		_, err = io.ReadFull(wal.file, lastWALRecordLengthBytes)
		if err != nil {
			return nil, err
		}
		slog.Info("BYTE ARRAY CONTENTS ", "last WAL record length", lastWALRecordLengthBytes)
		lastWALRecordLength := binary.BigEndian.Uint64(lastWALRecordLengthBytes)

		lastWALRecordBytes := make([]byte, int(lastWALRecordLength)+8)

		_, err = wal.file.Seek(stat.Size()-int64(lastWALRecordLength)-8, io.SeekStart)
		if err != nil {
			return nil, err
		}
		_, err = io.ReadFull(wal.file, lastWALRecordBytes)
		if err != nil {
			return nil, err
		}

		lastWALRecord := DecodeWALRecord(lastWALRecordBytes)

		wal.currLSN = lastWALRecord.LSN
	}

	return wal, nil
}

func (wal *WAL) Close() error {

	if err := wal.file.Sync(); err != nil {
		return err
	}
	return wal.file.Close()
}

func (wal *WAL) NewWALIterator() (*WALIterator, error) {

	wal.mutex.Lock()

	_, err := wal.file.Seek(0, 0)

	if err != nil {
		wal.mutex.Unlock()
		return nil, err
	}

	info, err := wal.file.Stat()

	if err != nil {
		wal.mutex.Unlock()
		return nil, err
	}

	return &WALIterator{
		wal:         wal,
		CurrOffset:  0,
		WalFileSize: uint64(info.Size()),
	}, nil
}

func (wal *WAL) log(operation Operation, payload []byte) (LSN uint64, err error) {
	wal.mutex.Lock()
	defer wal.mutex.Unlock()

	LSN = atomic.AddUint64(&wal.currLSN, 1)

	record := WALRecord{
		LSN:       LSN,
		Operation: operation,
		Payload:   payload,
	}

	walRecordBytes := EncodeWALRecord(record)

	if _, err := wal.file.Write(walRecordBytes); err != nil {
		slog.Error(err.Error(), "at", "Log")
		return 0, ErrWrite
	}
	if err := wal.file.Sync(); err != nil {
		return 0, err
	}
	return LSN, nil

}

func (wal *WAL) LogCreatePageOperation(pageId uint64, pageType byte, allocationSource byte) (LSN uint64, err error) {

	payload := CreatePagePayload{
		PageId:           pageId,
		PageType:         pageType,
		AllocationSource: byte(allocationSource),
	}
	return wal.log(CreatePage, EncodeCreatePagePayload(payload))
}

func (wal *WAL) LogDeletePageOperation(pageId uint64) (LSN uint64, err error) {

	payload := DeletePagePayload{
		PageId: pageId,
	}
	return wal.log(DeletePage, EncodeDeletePagePayload(payload))
}

func (wal *WAL) LogInsertInternalNodeEntryOperation(pageId uint64,
	key []byte,
	leftChildNodePageId uint64,
	rightChildNodePageId uint64) (LSN uint64, err error) {

	payload := InsertInternalNodePayload{
		PageId:               pageId,
		Key:                  key,
		LeftChildNodePageId:  leftChildNodePageId,
		RightChildNodePageId: rightChildNodePageId,
	}
	return wal.log(InsertInternalNodeEntry, EncodeInsertInternalNodePayload(payload))
}

func (wal *WAL) LogInsertLeafNodeEntryOperation(pageId uint64, key []byte, value []byte) (LSN uint64, err error) {

	payload := InsertLeafNodeEntryPayload{
		PageId: pageId,
		Key:    key,
		Value:  value,
	}
	return wal.log(InsertLeafNodeEntry, EncodeInsertLeafNodeEntryPayload(payload))
}

func (wal *WAL) LogUpdateLeafNodeEntryOperation(pageId uint64, key []byte, value []byte) (LSN uint64, err error) {

	payload := UpdateLeafNodeEntryPayload{
		PageId: pageId,
		Key:    key,
		Value:  value,
	}
	return wal.log(UpdateLeafNodeEntry, EncodeUpdateLeafNodeEntryPayload(payload))
}

func (wal *WAL) LogUpdateRootNodePageIdOperation(bPlusTreeId uint64, rootNodePageId uint64) (LSN uint64, err error) {

	payload := UpdateRootNodePageIdPayload{
		BPlusTreeId:    bPlusTreeId,
		RootNodePageId: rootNodePageId,
	}
	return wal.log(UpdateRootNodePageId, EncodeUpdateRootNodePageIdPayload(payload))
}

func (wal *WAL) LogUpdateFirstLeafNodePageIdOperation(bPlusTreeId uint64, firstLeafNodePageId uint64) (LSN uint64, err error) {

	payload := UpdateFirstLeafNodePageIdPayload{
		BPlusTreeId:         bPlusTreeId,
		FirstLeafNodePageId: firstLeafNodePageId,
	}
	return wal.log(UpdateFirstLeafNodePageId, EncodeUpdateFirstLeafNodePageIdPayload(payload))
}

func (wal *WAL) LogSplitInternalNodeOperation(leftInternalNodePageId uint64,
	rightInternalNodePageId uint64,
	parentNodePageId uint64,
	separatorKeyIndex uint16,
	insertKey []byte,
	insertLeftNodePageId uint64,
	insertRightNodePageId uint64,
	elementsLength uint16,
	elements []byte) (LSN uint64, err error) {

	payload := SplitInternalNodePayload{
		LeftInternalNodePageId:  leftInternalNodePageId,
		RightInternalNodePageId: rightInternalNodePageId,
		ParentNodePageId:        parentNodePageId,
		SeparatorKeyIndex:       separatorKeyIndex,
		InsertKey:               insertKey,
		InsertLeftNodePageId:    insertLeftNodePageId,
		InsertRightNodePageId:   insertRightNodePageId,
		ElementsLength:          elementsLength,
		Elements:                elements,
	}
	return wal.log(SplitInternalNode, EncodeSplitInternalNodePayload(payload))
}

func (wal *WAL) LogSplitLeafNodeOperation(leftLeafNodePageId uint64,
	rightLeafNodePageId uint64,
	parentNodePageId uint64,
	separatorKeyIndex uint16,
	nextLeafNodePageId uint64,
	insertKey []byte,
	insertValue []byte,
	elementsLength uint16,
	elements []byte) (LSN uint64, err error) {

	payload := SplitLeafNodePayload{
		LeftLeafNodePageId:  leftLeafNodePageId,
		RightLeafNodePageId: rightLeafNodePageId,
		ParentNodePageId:    parentNodePageId,
		SeparatorKeyIndex:   separatorKeyIndex,
		NextLeafNodePageId:  nextLeafNodePageId,
		InsertKey:           insertKey,
		InsertValue:         insertValue,
		ElementsLength:      elementsLength,
		Elements:            elements,
	}
	return wal.log(SplitLeafNode, EncodeSplitLeafNodePayload(payload))
}
