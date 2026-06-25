package codec

import "encoding/binary"

type SplitLeafNodePayload struct {
	LeftLeafNodePageId  uint64
	RightLeafNodePageId uint64
}

func EncodeSplitLeafNodePayload(Payload SplitLeafNodePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, Payload.LeftLeafNodePageId)
	data = binary.BigEndian.AppendUint64(data, Payload.RightLeafNodePageId)

	return data
}

func DecodeSplitLeafNodePayload(data []byte) SplitLeafNodePayload {

	payload := SplitLeafNodePayload{}

	pointer := 0
	payload.LeftLeafNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8
	payload.RightLeafNodePageId = binary.BigEndian.Uint64(data[pointer:])

	return payload
}
