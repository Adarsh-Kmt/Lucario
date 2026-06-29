package lucario

import "encoding/binary"

type SplitLeafNodePayload struct {
	LeftLeafNodePageId  uint64
	RightLeafNodePageId uint64
	ParentNodePageId    uint64
	SeparatorKeyIndex   uint16

	InsertKey   []byte
	InsertValue []byte

	ElementsLength uint16
	Elements       []byte
}

func EncodeSplitLeafNodePayload(payload SplitLeafNodePayload) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, payload.LeftLeafNodePageId)
	data = binary.BigEndian.AppendUint64(data, payload.RightLeafNodePageId)
	data = binary.BigEndian.AppendUint64(data, payload.ParentNodePageId)
	data = binary.BigEndian.AppendUint16(data, payload.SeparatorKeyIndex)
	data = binary.BigEndian.AppendUint16(data, uint16(len(payload.InsertKey)))
	data = append(data, payload.InsertKey...)
	data = binary.BigEndian.AppendUint16(data, uint16(len(payload.InsertValue)))
	data = append(data, payload.InsertValue...)

	data = binary.BigEndian.AppendUint16(data, uint16(len(payload.Elements)))
	data = append(data, payload.Elements...)
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

	payload.SeparatorKeyIndex = binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	insertKeyLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	payload.InsertKey = data[pointer : pointer+int(insertKeyLength)]
	pointer += int(insertKeyLength)

	insertValueLength := binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	payload.InsertValue = data[pointer : pointer+int(insertValueLength)]
	pointer += int(insertValueLength)

	payload.ElementsLength = binary.BigEndian.Uint16(data[pointer:])
	pointer += 2

	payload.Elements = data[pointer:]

	return payload
}
