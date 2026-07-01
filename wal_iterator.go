package lucario

import (
	"encoding/binary"
	"io"
)

type WALIterator struct {
	wal         *WAL
	currOffset  uint64
	walFileSize uint64
}

func (iterator *WALIterator) HasNext() bool {

	if iterator.currOffset <= iterator.walFileSize {
		return false
	}
	return true
}

func (iterator *WALIterator) GetRecord() (WALRecord, error) {

	walRecordLengthBytes := make([]byte, 8)

	_, err := io.ReadFull(iterator.wal.file, walRecordLengthBytes)

	if err != nil {
		return WALRecord{}, err
	}

	walRecordLength := binary.BigEndian.Uint64(walRecordLengthBytes) + 8

	walRecordBytes := make([]byte, walRecordLength)

	_, err = io.ReadFull(iterator.wal.file, walRecordBytes)

	if err != nil {
		return WALRecord{}, err
	}

	iterator.currOffset += uint64(8 + len(walRecordBytes))

	return DecodeWALRecord(walRecordBytes), nil
}

func (iterator *WALIterator) Close() {
	iterator.wal.mutex.Unlock()

}
