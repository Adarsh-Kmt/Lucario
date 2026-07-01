package lucario

import (
	"encoding/binary"
	"io"
)

type WALIterator struct {
	wal         *WAL
	CurrOffset  uint64
	WalFileSize uint64
}

func (iterator *WALIterator) HasNext() bool {

	if iterator.CurrOffset < iterator.WalFileSize {
		return true
	}
	return false
}

func (iterator *WALIterator) GetCurrOffset() uint64 {
	return iterator.CurrOffset
}

func (iterator *WALIterator) GetWalFileSize() uint64 {
	return iterator.WalFileSize
}

func (iterator *WALIterator) GetRecord() (WALRecord, error) {

	walRecordLengthBytes := make([]byte, 8)

	_, err := io.ReadFull(iterator.wal.file, walRecordLengthBytes)

	if err != nil {
		return WALRecord{}, err
	}

	walRecordLength := binary.BigEndian.Uint64(walRecordLengthBytes)

	walRecordBytes := make([]byte, walRecordLength)

	_, err = io.ReadFull(iterator.wal.file, walRecordBytes)

	if err != nil {
		return WALRecord{}, err
	}

	iterator.CurrOffset += uint64(8 + len(walRecordBytes))

	return DecodeWALRecord(walRecordBytes), nil
}

func (iterator *WALIterator) Close() {
	iterator.wal.mutex.Unlock()

}
