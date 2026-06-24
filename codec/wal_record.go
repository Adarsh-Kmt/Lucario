package codec

import "encoding/binary"

type Operation uint16

const (
	CreatePage Operation = iota
	InsertInternalNodeEntry
	InsertLeafNodeEntry
	SplitInternalNode
	SplitLeafNode
	UpdateRootNodePageId
	UpdateFirstLeafNodePageId
)

type WALRecord struct {
	LSN       uint64
	Operation Operation
	Payload   []byte
}

func EncodeWALRecord(record WALRecord) []byte {

	data := make([]byte, 0)

	data = binary.BigEndian.AppendUint64(data, record.LSN)
	data = binary.BigEndian.AppendUint16(data, uint16(record.Operation))

	data = binary.BigEndian.AppendUint64(data, uint64(len(record.Payload)))
	data = append(data, record.Payload...)

	return data
}

func DecodeWALRecord(data []byte) WALRecord {

	record := WALRecord{}

	pointer := 0

	record.LSN = binary.BigEndian.Uint64(data[pointer:])
	pointer += 8

	record.Operation = Operation(binary.BigEndian.Uint16(data[pointer:]))
	pointer += 2

	payloadLength := binary.BigEndian.Uint64(data[pointer:])
	pointer += 8

	record.Payload = data[pointer : pointer+int(payloadLength)]

	return record
}
