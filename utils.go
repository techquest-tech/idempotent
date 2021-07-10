package idempotent

import (
	"encoding/binary"
	"errors"
)

var errTypeNotSupport = errors.New("type is not supported")

// func NumToBytes(num interface{}) []byte {

// }

// Convert Object to bytes, support int
func ToBytes(obj interface{}) ([]byte, error) {
	if str, ok := obj.(string); ok {
		return []byte(str), nil
	}

	if buf, ok := obj.([]byte); ok {
		return buf, nil
	}
	if b, ok := obj.(uint8); ok {
		return []byte{b}, nil
	}
	if i16, ok := obj.(uint16); ok {
		bs := make([]byte, 2)
		binary.LittleEndian.PutUint16(bs, i16)
		return bs, nil
	}
	if i32, ok := obj.(uint32); ok {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, i32)
		return bs, nil
	}
	if i64, ok := obj.(uint64); ok {
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, i64)
		return bs, nil
	}
	if i, ok := obj.(int); ok {
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, uint64(i))
		return bs, nil
	}

	return []byte{}, errTypeNotSupport
}
