package codec

import "encoding/binary"

type InsertLeafNodeEntryPayload struct {
	PageId uint64
	Key    []byte
	Value  []byte
}

func EncodeInsertLeafNodeEntryPayload(Payload InsertLeafNodeEntryPayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, Payload.PageId)

	data = binary.BigEndian.AppendUint16(data, uint16(len(Payload.Key)))
	data = append(data, Payload.Key...)

	data = binary.BigEndian.AppendUint16(data, uint16(len(Payload.Value)))
	data = append(data, Payload.Value...)

	return data
}

func DecodeInsertLeafNodeEntryPayload(data []byte) InsertLeafNodeEntryPayload {

	Payload := InsertLeafNodeEntryPayload{}

	pointer := 0

	Payload.PageId = binary.BigEndian.Uint64(data[pointer : pointer+8])

	pointer += 8

	keyLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	Payload.Key = data[pointer : pointer+int(keyLength)]

	pointer += int(keyLength)

	valueLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	Payload.Value = data[pointer : pointer+int(valueLength)]

	return Payload
}
