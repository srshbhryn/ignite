package datatypes

import (
	"encoding/binary"
	"errors"
	"io"
)

func Write(writer io.Writer, data any) error {
	switch data.(type) {
	case string:
		return writeString(writer, data.(string))
	default:
		return errors.New("invalid type to write")
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
