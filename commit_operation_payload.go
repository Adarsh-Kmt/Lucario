package lucario

import "encoding/binary"

type CommitOperationPayload struct {
	BPlusTreeId uint64
}

func EncodeCommitOperationPayload(payload CommitOperationPayload) []byte {

	data := make([]byte, 0)
	data = binary.BigEndian.AppendUint64(data, payload.BPlusTreeId)
	return data
}

func DecodeCommitOperationPayload(data []byte) CommitOperationPayload {

	payload := CommitOperationPayload{}
	payload.BPlusTreeId = binary.BigEndian.Uint64(data[:])
	return payload
}
