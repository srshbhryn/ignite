package datatypes

import (
	"encoding/binary"
	"errors"
	"io"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

var ErrInvalidTypeToRead = errors.New("invalid type to read")

func Read(reader io.Reader) (any, error) {
	var typeCode byte
	if err := binary.Read(reader, binary.LittleEndian, &typeCode); err != nil {
		return nil, err
	}
	switch typeCode {
	case _INT16:
		return readInt16(reader)
	case _INT32:
		return readInt32(reader)
	case _INT64:
		return readInt64(reader)
	case _FLOAT32:
		return readFloat32(reader)
	case _FLOAT64:
		return readFloat64(reader)
	case _BOOL:
		return readBool(reader)
	case _STRING:
		return readString(reader)
	case _TIMESTAMP:
		return readTimestamp(reader)
	case _DECIMAL:
		return readDecimal(reader)
	case _NULL:
		return nil, nil
	default:
		return nil, ErrInvalidTypeToRead
	}
}

func readInt16(reader io.Reader) (int16, error) {
	var value int16
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func readInt32(reader io.Reader) (int32, error) {
	var value int32
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func readInt64(reader io.Reader) (int64, error) {
	var value int64
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil

}

func readFloat32(reader io.Reader) (float32, error) {
	var value float32
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func readFloat64(reader io.Reader) (float64, error) {
	var value float64
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func readBool(reader io.Reader) (bool, error) {
	var value bool
	if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return false, err
	}
	return value, nil
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

func readTimestamp(reader io.Reader) (time.Time, error) {
	var msecsSinceEpoch int64
	var msecFractionInNsecs int32
	if err := binary.Read(reader, binary.LittleEndian, &msecsSinceEpoch); err != nil {
		return time.Time{}, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &msecFractionInNsecs); err != nil {
		return time.Time{}, err
	}
	secondsSinceEpoch := int64(msecsSinceEpoch / 1e3)
	nanoSeconds := (msecsSinceEpoch-secondsSinceEpoch*1e3)*1e6 + int64(msecFractionInNsecs)
	timestamp := time.Unix(secondsSinceEpoch, nanoSeconds)
	return timestamp, nil
}

func readDecimal(reader io.Reader) (decimal.Decimal, error) {
	var scale int32
	if err := binary.Read(reader, binary.LittleEndian, &scale); err != nil {
		return decimal.Decimal{}, err
	}
	var intPartLen int32
	if err := binary.Read(reader, binary.LittleEndian, &intPartLen); err != nil {
		return decimal.Decimal{}, err
	}
	var signByte byte
	if err := binary.Read(reader, binary.LittleEndian, &signByte); err != nil {
		return decimal.Decimal{}, err
	}
	intPartBytes := make([]byte, intPartLen-1)
	if err := binary.Read(reader, binary.BigEndian, intPartBytes); err != nil {
		return decimal.Decimal{}, err
	}
	intPart := big.NewInt(0).SetBytes(intPartBytes)
	if signByte == byte(1) {
		intPart = intPart.Neg(intPart)
	}
	decimalValue := decimal.NewFromBigInt(intPart, scale)
	return decimalValue, nil
}
