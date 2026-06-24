package codec

import "encoding/binary"

type SplitInternalNodePayload struct {
	LeftInternalNodePageId  uint64
	RightInternalNodePageId uint64
}

func EncodeSplitInternalNodePayload(Payload SplitInternalNodePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, Payload.LeftInternalNodePageId)
	data = binary.BigEndian.AppendUint64(data, Payload.RightInternalNodePageId)

	return data
}

func DecodeSplitInternalNodePayload(data []byte) SplitInternalNodePayload {

	payload := SplitInternalNodePayload{}

	pointer := 0
	payload.LeftInternalNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8
	payload.RightInternalNodePageId = binary.BigEndian.Uint64(data[pointer:])

	return payload
}
