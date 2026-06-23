package codec

import "encoding/binary"

type CreatePagePayload struct {
	PageId   uint64
	PageType byte
}

func EncodeCreatePagePayload(Payload CreatePagePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, Payload.PageId)

	data = append(data, Payload.PageType)
	return data
}

func DecodeCreatePagePayload(data []byte) CreatePagePayload {

	Payload := CreatePagePayload{}

	Payload.PageId = binary.BigEndian.Uint64(data[:8])
	Payload.PageType = data[8]

	return Payload
}
