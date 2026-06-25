package codec

import "encoding/binary"

type InsertInternalNodePayload struct {
	PageId               uint64
	Key                  []byte
	LeftChildNodePageId  uint64
	RightChildNodePageId uint64
}

func EncodeInsertInternalNodePayload(payload InsertInternalNodePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, payload.PageId)

	data = binary.BigEndian.AppendUint16(data, uint16(len(payload.Key)))
	data = append(data, payload.Key...)

	data = binary.BigEndian.AppendUint64(data, payload.LeftChildNodePageId)
	data = binary.BigEndian.AppendUint64(data, payload.RightChildNodePageId)

	return data
}

func DecodeInsertInternalNodePayload(data []byte) InsertInternalNodePayload {

	payload := InsertInternalNodePayload{}

	pointer := 0

	payload.PageId = binary.BigEndian.Uint64(data[pointer : pointer+8])

	pointer += 8

	keyLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	payload.Key = data[pointer : pointer+int(keyLength)]

	pointer += int(keyLength)

	payload.LeftChildNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8

	payload.RightChildNodePageId = binary.BigEndian.Uint64(data[pointer:])

	return payload
}
