package codec

import "encoding/binary"

type InsertLeafNodeEntryRecord struct {
	PageId uint64
	Key    []byte
	Value  []byte
}

func EncodeInsertLeafNodeEntryRecord(record InsertLeafNodeEntryRecord) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, record.PageId)

	data = binary.BigEndian.AppendUint16(data, uint16(len(record.Key)))
	data = append(data, record.Key...)

	data = binary.BigEndian.AppendUint16(data, uint16(len(record.Value)))
	data = append(data, record.Value...)

	return data
}

func DecodeInsertLeafNodeEntryRecord(data []byte) InsertLeafNodeEntryRecord {

	record := InsertLeafNodeEntryRecord{}

	pointer := 0

	record.PageId = binary.BigEndian.Uint64(data[pointer : pointer+8])

	pointer += 8

	keyLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	record.Key = data[pointer : pointer+int(keyLength)]

	pointer += pointer + int(keyLength)

	valueLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	record.Value = data[pointer : pointer+int(valueLength)]

	return record
}
