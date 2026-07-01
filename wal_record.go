package lucario

import "encoding/binary"

type Operation uint16

const (
	CreatePage Operation = iota
	DeletePage
	InsertInternalNodeEntry
	InsertLeafNodeEntry
	UpdateLeafNodeEntry
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

	finalPayload := make([]byte, 0)
	finalPayload = binary.BigEndian.AppendUint64(finalPayload, uint64(len(data)))
	finalPayload = append(finalPayload, data...)
	finalPayload = binary.BigEndian.AppendUint64(finalPayload, uint64(len(data)))

	return finalPayload
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
