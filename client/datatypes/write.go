package datatypes

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

var ErrInvalidTypeToWrite = errors.New("invalid type to write")

func Write(writer io.Writer, data any) error {
	if data == nil {
		return ErrInvalidTypeToWrite
	}
	if v := reflect.ValueOf(data); v.Kind() == reflect.Ptr {
		data = v.Elem().Interface()
	}

	switch data := data.(type) {
	case string:
		return writeString(writer, data)
	default:
		return ErrInvalidTypeToWrite
	}
}

func writeString(writer io.Writer, value string) error {
	if err := binary.Write(writer, binary.LittleEndian, _STRING); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, int32(len(value))); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, []byte(value)); err != nil {
		return err
	}
	return nil
}
