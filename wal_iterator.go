package lucario

import (
	"encoding/binary"

	"github.com/Adarsh-Kmt/Lucario/codec"
)

type WALIterator struct {
	wal         *WAL
	currOffset  uint64
	walFileSize uint64
}

func (iterator *WALIterator) HasNext() bool {

	if iterator.currOffset == iterator.walFileSize {
		return false
	}
	return true
}

func (iterator *WALIterator) GetRecord() (codec.WALRecord, error) {

	walRecordLengthBytes := make([]byte, 8)

	_, err := iterator.wal.file.Read(walRecordLengthBytes)

	if err != nil {
		return codec.WALRecord{}, err
	}

	walRecordLength := binary.BigEndian.Uint64(walRecordLengthBytes)

	walRecordBytes := make([]byte, walRecordLength)

	_, err = iterator.wal.file.Read(walRecordBytes)

	if err != nil {
		return codec.WALRecord{}, err
	}

	iterator.currOffset += uint64(8 + len(walRecordBytes))

	return codec.DecodeWALRecord(walRecordBytes), nil
}

func (iterator *WALIterator) Close() {
	iterator.wal.mutex.Unlock()

}
