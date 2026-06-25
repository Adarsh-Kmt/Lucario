package codec

import "encoding/binary"

type UpdateFirstLeafNodePageIdPayload struct {
	BPlusTreeId         uint64
	FirstLeafNodePageId uint64
}

func EncodeUpdateFirstLeafNodePageIdPayload(payload UpdateFirstLeafNodePageIdPayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, payload.BPlusTreeId)
	data = binary.BigEndian.AppendUint64(data, payload.FirstLeafNodePageId)

	return data
}

func DecodeUpdateFirstLeafNodePageIdPayload(data []byte) UpdateFirstLeafNodePageIdPayload {

	payload := UpdateFirstLeafNodePageIdPayload{}

	pointer := 0

	payload.BPlusTreeId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8
	payload.FirstLeafNodePageId = binary.BigEndian.Uint64(data[pointer:])

	return payload
}
