package datatypes

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"reflect"
	"time"

	"github.com/shopspring/decimal"
)

var ErrInvalidTypeToWrite = errors.New("invalid type to write")

func Write(writer io.Writer, data any) error {
	if data == nil {
		return writeNull(writer)
	}
	if v := reflect.ValueOf(data); v.Kind() == reflect.Ptr {
		data = v.Elem().Interface()
	}

	switch data := data.(type) {
	case int16:
		return writeInt16(writer, data)
	case int32:
		return writeInt32(writer, data)
	case int64:
		return writeInt64(writer, data)
	case float32:
		return writeFloat32(writer, data)
	case float64:
		return writeFloat64(writer, data)
	case bool:
		return writeBool(writer, data)
	case string:
		return writeString(writer, data)
	case time.Time:
		return writeTimestamp(writer, data)
	case decimal.Decimal:
		return writeDecimal(writer, data)
	default:
		return ErrInvalidTypeToWrite
	}
}

func writeInt16(writer io.Writer, value int16) error {
	if err := binary.Write(writer, binary.LittleEndian, _INT16); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, value); err != nil {
		return err
	}
	return nil
}

func writeInt32(writer io.Writer, value int32) error {
	if err := binary.Write(writer, binary.LittleEndian, _INT32); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, value); err != nil {
		return err
	}
	return nil

}
func writeInt64(writer io.Writer, value int64) error {
	if err := binary.Write(writer, binary.LittleEndian, _INT64); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, value); err != nil {
		return err
	}
	return nil
}

func writeFloat32(writer io.Writer, value float32) error {
	if err := binary.Write(writer, binary.LittleEndian, _FLOAT32); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, value); err != nil {
		return err
	}
	return nil
}

func writeFloat64(writer io.Writer, value float64) error {
	if err := binary.Write(writer, binary.LittleEndian, _FLOAT64); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, value); err != nil {
		return err
	}
	return nil
}

func writeBool(writer io.Writer, value bool) error {
	if err := binary.Write(writer, binary.LittleEndian, _BOOL); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, value); err != nil {
		return err
	}
	return nil
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

func writeTimestamp(writer io.Writer, value time.Time) error {
	msecsSinceEpoch := value.UnixMilli()
	msecFractionInNsecs := int32(value.UnixNano() - value.UnixMilli()*1e6)
	if err := binary.Write(writer, binary.LittleEndian, _TIMESTAMP); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, msecsSinceEpoch); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, msecFractionInNsecs); err != nil {
		return err
	}
	return nil

}

func writeDecimal(writer io.Writer, value decimal.Decimal) error {
	scale := value.Exponent()
	signByte := byte(0)
	if value.IsNegative() {
		signByte = byte(1)
	}
	if err := binary.Write(writer, binary.LittleEndian, _DECIMAL); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, scale); err != nil {
		return err
	}
	intPartBytes := value.Mul(decimal.NewFromFloat(math.Pow10(int(scale)))).BigInt().Bytes()
	intPartLen := int32(len(intPartBytes) + 1)
	if err := binary.Write(writer, binary.LittleEndian, &intPartLen); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, &signByte); err != nil {
		return err
	}
	if _, err := writer.Write(intPartBytes); err != nil {
		return err
	}
	return nil
}

func writeNull(writer io.Writer) error {
	if err := binary.Write(writer, binary.LittleEndian, _NULL); err != nil {
		return err
	}
	return nil
}
