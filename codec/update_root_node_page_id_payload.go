package codec

import "encoding/binary"

type UpdateRootNodePageIdPayload struct {
	BPlusTreeId    uint64
	RootNodePageId uint64
}

func EncodeUpdateRootNodePageIdPayload(payload UpdateRootNodePageIdPayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, payload.BPlusTreeId)
	data = binary.BigEndian.AppendUint64(data, payload.RootNodePageId)

	return data
}

func DecodeUpdateRootNodePageIdPayload(data []byte) UpdateRootNodePageIdPayload {

	payload := UpdateRootNodePageIdPayload{}

	pointer := 0

	payload.BPlusTreeId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8
	payload.RootNodePageId = binary.BigEndian.Uint64(data[pointer:])

	return payload
}
