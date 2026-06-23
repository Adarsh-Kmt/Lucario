package codec

import "encoding/binary"

type InsertInternalNodePayload struct {
	PageId               uint64
	Key                  []byte
	LeftChildNodePageId  uint64
	RightChildNodePageId uint64
}

func EncodeInsertInternalNodePayload(Payload InsertInternalNodePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, Payload.PageId)

	data = binary.BigEndian.AppendUint16(data, uint16(len(Payload.Key)))
	data = append(data, Payload.Key...)

	data = binary.BigEndian.AppendUint64(data, Payload.LeftChildNodePageId)
	data = binary.BigEndian.AppendUint64(data, Payload.RightChildNodePageId)

	return data
}

func DecodeInsertInternalNodePayload(data []byte) InsertInternalNodePayload {

	Payload := InsertInternalNodePayload{}

	pointer := 0

	Payload.PageId = binary.BigEndian.Uint64(data[pointer : pointer+8])

	pointer += 8

	keyLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	Payload.Key = data[pointer : pointer+int(keyLength)]

	pointer += int(keyLength)

	Payload.LeftChildNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8

	Payload.RightChildNodePageId = binary.BigEndian.Uint64(data[pointer:])

	return Payload
}
