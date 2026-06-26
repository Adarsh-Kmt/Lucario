package lucario

import "encoding/binary"

type SplitLeafNodePayload struct {
	LeftLeafNodePageId  uint64
	RightLeafNodePageId uint64
	ParentNodePageId    uint64
	SeparatorKey        []byte
}

func EncodeSplitLeafNodePayload(payload SplitLeafNodePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, payload.LeftLeafNodePageId)
	data = binary.BigEndian.AppendUint64(data, payload.RightLeafNodePageId)
	data = binary.BigEndian.AppendUint64(data, payload.ParentNodePageId)
	data = binary.BigEndian.AppendUint16(data, uint16(len(payload.SeparatorKey)))
	data = append(data, payload.SeparatorKey...)

	return data
}

func DecodeSplitLeafNodePayload(data []byte) SplitLeafNodePayload {

	payload := SplitLeafNodePayload{}

	pointer := 0
	payload.LeftLeafNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8
	payload.RightLeafNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8
	payload.ParentNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8

	separatorKeyLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	payload.SeparatorKey = data[pointer : pointer+int(separatorKeyLength)]

	return payload
}
