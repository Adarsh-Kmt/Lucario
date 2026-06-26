package lucario

import "encoding/binary"

type CreatePagePayload struct {
	PageId           uint64
	PageType         byte
	AllocationSource byte
}

func EncodeCreatePagePayload(payload CreatePagePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, payload.PageId)

	data = append(data, payload.PageType)
	data = append(data, payload.AllocationSource)
	return data
}

func DecodeCreatePagePayload(data []byte) CreatePagePayload {

	payload := CreatePagePayload{}

	payload.PageId = binary.BigEndian.Uint64(data[:8])
	payload.PageType = data[8]
	payload.AllocationSource = data[9]

	return payload
}
