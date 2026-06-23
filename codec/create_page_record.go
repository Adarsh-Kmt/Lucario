package codec

import "encoding/binary"

type CreatePageRecord struct {
	PageId   uint64
	PageType byte
}

func EncodeCreatePageRecord(record CreatePageRecord) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, record.PageId)

	data = append(data, record.PageType)
	return data
}

func DecodeCreatePageRecord(data []byte) CreatePageRecord {

	record := CreatePageRecord{}

	record.PageId = binary.BigEndian.Uint64(data[:8])
	record.PageType = data[8]

	return record
}
