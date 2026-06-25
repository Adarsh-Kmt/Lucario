package codec

import "encoding/binary"

type DeletePagePayload struct {
	PageId uint64
}

func EncodeDeletePagePayload(payload DeletePagePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, payload.PageId)

	return data
}

func DecodeDeletePagePayload(data []byte) DeletePagePayload {

	payload := DeletePagePayload{}

	payload.PageId = binary.BigEndian.Uint64(data[:])

	return payload
}
