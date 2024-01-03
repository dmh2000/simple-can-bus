package can

import (
	"encoding/binary"
)

func BytesToUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func BytesToInt32(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}

func Uint16ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	return b
}

func Int32ToBytes(v int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	return b
}
