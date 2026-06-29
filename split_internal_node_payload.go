package lucario

import "encoding/binary"

type SplitInternalNodePayload struct {
	LeftInternalNodePageId  uint64
	RightInternalNodePageId uint64
	ParentNodePageId        uint64
	SeparatorKeyIndex       uint16

	InsertKey             []byte
	InsertLeftNodePageId  uint64
	InsertRightNodePageId uint64

	ElementsLength uint16
	Elements       []byte
}

func EncodeSplitInternalNodePayload(payload SplitInternalNodePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, payload.LeftInternalNodePageId)
	data = binary.BigEndian.AppendUint64(data, payload.RightInternalNodePageId)
	data = binary.BigEndian.AppendUint64(data, payload.ParentNodePageId)
	data = binary.BigEndian.AppendUint16(data, payload.SeparatorKeyIndex)
	data = binary.BigEndian.AppendUint16(data, uint16(len(payload.InsertKey)))
	data = append(data, payload.InsertKey...)
	data = binary.BigEndian.AppendUint64(data, payload.InsertLeftNodePageId)
	data = binary.BigEndian.AppendUint64(data, payload.InsertRightNodePageId)

	data = binary.BigEndian.AppendUint16(data, uint16(len(payload.Elements)))
	data = append(data, payload.Elements...)
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
	pointer += 2

	insertKeyLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	payload.InsertKey = data[pointer : pointer+int(insertKeyLength)]
	pointer += int(insertKeyLength)

	payload.InsertLeftNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8
	payload.InsertRightNodePageId = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8

	payload.ElementsLength = binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	payload.Elements = data[pointer:]

	return payload
}
