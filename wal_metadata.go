package lucario

import "encoding/binary"

type WALMetadataCodec struct {
	MetadataLength   int
	currLSNOffset    int
	checkpointOffset int
}

func NewWALMetadataCodec() WALMetadataCodec {

	return WALMetadataCodec{
		MetadataLength:   16,
		currLSNOffset:    0,
		checkpointOffset: 8,
	}
}

func (codec WALMetadataCodec) EncodeWALMetadata(currLSN uint64, checkpointOffset uint64) []byte {

	b := make([]byte, 0)

	b = binary.LittleEndian.AppendUint64(b, currLSN)
	b = binary.LittleEndian.AppendUint64(b, checkpointOffset)

	return b
}

func (codec WALMetadataCodec) DecodeWALMetadata(metadataBytes []byte) (currLSN uint64, checkpointOffset uint64) {

	return binary.LittleEndian.Uint64(metadataBytes[codec.currLSNOffset:]), binary.LittleEndian.Uint64(metadataBytes[codec.checkpointOffset:])
}
