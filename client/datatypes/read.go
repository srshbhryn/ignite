package datatypes

import (
	"encoding/binary"
	"errors"
	"io"
)

func Read(reader io.Reader) (any, error) {
	var typeCode byte
	if err := binary.Read(reader, binary.LittleEndian, &typeCode); err != nil {
		return nil, err
	}
	switch typeCode {
	case _STRING:
		return readString(reader)
	case _NULL:
		return nil, nil
	default:
		return nil, errors.New("invalid type to read")
	}
}

func readString(reader io.Reader) (string, error) {
	var len int32
	if err := binary.Read(reader, binary.LittleEndian, &len); err != nil {
		return "", err
	}
	payload := make([]byte, len)
	if err := binary.Read(reader, binary.LittleEndian, payload); err != nil {
		return "", err
	}
	return string(payload), nil
}
