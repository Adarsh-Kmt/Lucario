package lucario

import (
	"encoding/binary"
	"errors"
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
	file             *os.File
	currLSN          uint64
	checkpointOffset uint64
	walMetadataCodec WALMetadataCodec
	mutex            *sync.Mutex
}

func NewWAL() (*WAL, error) {
	filePath := "/lucario.wal"

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
		file:             file,
		walMetadataCodec: NewWALMetadataCodec(),
		mutex:            &sync.Mutex{},
	}

	if !fileExists || stat.Size() == 0 {

		metadataBytes := wal.walMetadataCodec.EncodeWALMetadata(0, uint64(wal.walMetadataCodec.MetadataLength))
		if _, err := file.Write(metadataBytes); err != nil {
			return nil, err
		}
		wal.currLSN = 0
		wal.checkpointOffset = uint64(wal.walMetadataCodec.MetadataLength)
	} else {

		metadataBytes := make([]byte, wal.walMetadataCodec.MetadataLength)
		if _, err := wal.file.ReadAt(metadataBytes, 0); err != nil {
			return nil, err
		}
		wal.currLSN, wal.checkpointOffset = wal.walMetadataCodec.DecodeWALMetadata(metadataBytes)

		if _, err := file.Seek(0, 2); err != nil {
			return nil, err
		}
	}

	return wal, nil
}

func (wal *WAL) Close() error {

	metadataBytes := wal.walMetadataCodec.EncodeWALMetadata(wal.currLSN, wal.checkpointOffset)

	if _, err := wal.file.Seek(0, 0); err != nil {
		return err
	}
	if _, err := wal.file.Write(metadataBytes); err != nil {
		return err
	}

	if err := wal.file.Sync(); err != nil {
		return err
	}
	return wal.file.Close()
}

func (wal *WAL) NewWALIterator() (*WALIterator, error) {

	wal.mutex.Lock()

	_, err := wal.file.Seek(int64(wal.checkpointOffset), 0)

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
		currOffset:  wal.checkpointOffset,
		walFileSize: uint64(info.Size()),
	}, nil
}

func (wal *WAL) Log(operation Operation, payload []byte) (LSN uint64, err error) {
	wal.mutex.Lock()
	defer wal.mutex.Unlock()

	LSN = atomic.AddUint64(&wal.currLSN, 1)

	record := WALRecord{
		LSN:       LSN,
		Operation: operation,
		Payload:   payload,
	}

	walRecordBytes := EncodeWALRecord(record)

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, uint64(len(walRecordBytes)))
	data = append(data, walRecordBytes...)

	if _, err := wal.file.Write(data); err != nil {
		slog.Error(err.Error(), "at", "Log")
		return 0, ErrWrite
	}
	if err := wal.file.Sync(); err != nil {
		return 0, err
	}
	return LSN, nil

}
