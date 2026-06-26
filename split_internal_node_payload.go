package lucario

import "encoding/binary"

type SplitInternalNodePayload struct {
	LeftInternalNodePageId  uint64
	RightInternalNodePageId uint64
	ParentNodePageId        uint64
	SeparatorKeyIndex       uint16
}

func EncodeSplitInternalNodePayload(payload SplitInternalNodePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, payload.LeftInternalNodePageId)
	data = binary.BigEndian.AppendUint64(data, payload.RightInternalNodePageId)
	data = binary.BigEndian.AppendUint64(data, payload.ParentNodePageId)
	data = binary.BigEndian.AppendUint16(data, payload.SeparatorKeyIndex)

	return data
}

func DecodeSplitInternalNodePayload(data []byte) SplitInternalNodePayload {

	payload := SplitInternalNodePayload{}

	pointer := 0
	payload.LeftInternalNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8
	payload.RightInternalNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8
	payload.ParentNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8

	payload.SeparatorKeyIndex = binary.BigEndian.Uint16(data[pointer:])

	return payload
}
