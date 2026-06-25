package codec

import "encoding/binary"

type InsertLeafNodeEntryPayload struct {
	PageId uint64
	Key    []byte
	Value  []byte
}

func EncodeInsertLeafNodeEntryPayload(payload InsertLeafNodeEntryPayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, payload.PageId)

	data = binary.BigEndian.AppendUint16(data, uint16(len(payload.Key)))
	data = append(data, payload.Key...)

	data = binary.BigEndian.AppendUint16(data, uint16(len(payload.Value)))
	data = append(data, payload.Value...)

	return data
}

func DecodeInsertLeafNodeEntryPayload(data []byte) InsertLeafNodeEntryPayload {

	payload := InsertLeafNodeEntryPayload{}

	pointer := 0

	payload.PageId = binary.BigEndian.Uint64(data[pointer : pointer+8])

	pointer += 8

	keyLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	payload.Key = data[pointer : pointer+int(keyLength)]

	pointer += int(keyLength)

	valueLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	payload.Value = data[pointer : pointer+int(valueLength)]

	return payload
}
