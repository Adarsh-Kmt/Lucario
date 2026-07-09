package lucario

import "encoding/binary"

type BeginOperationPayload struct {
	BPlusTreeId uint64
}

func EncodeBeginOperationPayload(payload BeginOperationPayload) []byte {

	data := make([]byte, 0)
	data = binary.BigEndian.AppendUint64(data, payload.BPlusTreeId)
	return data
}

func DecodeBeginOperationPayload(data []byte) BeginOperationPayload {

	payload := BeginOperationPayload{}
	payload.BPlusTreeId = binary.BigEndian.Uint64(data[:])
	return payload
}
