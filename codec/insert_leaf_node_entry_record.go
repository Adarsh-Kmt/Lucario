package codec

import "encoding/binary"

type InsertInternalNodeRecord struct {
	PageId               uint64
	Key                  []byte
	LeftChildNodePageId  uint64
	RightChildNodePageId uint64
}

func EncodeInsertInternalNodeRecord(record InsertInternalNodeRecord) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, record.PageId)

	data = binary.BigEndian.AppendUint16(data, uint16(len(record.Key)))
	data = append(data, record.Key...)

	data = binary.BigEndian.AppendUint64(data, record.LeftChildNodePageId)
	data = binary.BigEndian.AppendUint64(data, record.RightChildNodePageId)

	return data
}

func DecodeInsertInternalNodeRecord(data []byte) InsertInternalNodeRecord {

	record := InsertInternalNodeRecord{}

	pointer := 0

	record.PageId = binary.BigEndian.Uint64(data[pointer : pointer+8])

	pointer += 8

	keyLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	record.Key = data[pointer : pointer+int(keyLength)]

	pointer += pointer + int(keyLength)

	record.LeftChildNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8

	record.RightChildNodePageId = binary.BigEndian.Uint64(data[pointer:])

	return record
}
