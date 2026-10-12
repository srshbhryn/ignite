package ignite

import (
	"bytes"
	"encoding/binary"
	"github.com/srshbhryn/ignite/datatypes"
	"time"
)

type QueryFieldStatementType byte

const (
	ANY    = QueryFieldStatementType(0)
	SELECT = QueryFieldStatementType(1)
	UPDATE = QueryFieldStatementType(2)
)

type QueryFieldsRequest struct {
	Schema    string
	Statement string
	Args      []any
	Type      QueryFieldStatementType
	PageSize  int32
	MaxRows   int32
	Timeout   time.Duration
}

type QueryFieldsResult struct {
	CursorId     int64
	ColumnsCount int32
	RowsCount    int32
	Rows         [][]any
	HasMore      bool
}

func (c *Client) QuerySqlFields(cacheName string, request *QueryFieldsRequest) (*QueryFieldsResult, error) {
	requestPayload := &bytes.Buffer{}
	code := hashCode([]byte(cacheName))
	if err := binary.Write(requestPayload, binary.LittleEndian, &code); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, byte(0)); err != nil {
		return nil, err
	}
	if err := datatypes.Write(requestPayload, &request.Schema); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, &request.PageSize); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, &request.MaxRows); err != nil {
		return nil, err
	}
	if err := datatypes.Write(requestPayload, &request.Statement); err != nil {
		return nil, err
	}
	argsCount := int32(len(request.Args))
	if err := binary.Write(requestPayload, binary.LittleEndian, &argsCount); err != nil {
		return nil, err
	}
	for _, arg := range request.Args {
		if err := datatypes.Write(requestPayload, &arg); err != nil {
			return nil, err
		}
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, byte(request.Type)); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, false); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, false); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, false); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, false); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, false); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, false); err != nil {
		return nil, err
	}
	timeout := int64(request.Timeout / time.Millisecond)
	if err := binary.Write(requestPayload, binary.LittleEndian, &timeout); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, false); err != nil {
		return nil, err
	}
	responsePayload, err := c.sendRequest(_OP_QUERY_SQL_FIELDS, requestPayload)
	if err != nil {
		return nil, err
	}
	result := &QueryFieldsResult{}
	if err := binary.Read(responsePayload, binary.LittleEndian, &result.CursorId); err != nil {
		return nil, err
	}
	if err := binary.Read(responsePayload, binary.LittleEndian, &result.ColumnsCount); err != nil {
		return nil, err
	}
	if err := binary.Read(responsePayload, binary.LittleEndian, &result.RowsCount); err != nil {
		return nil, err
	}

	result.Rows = make([][]any, result.RowsCount)
	for i := 0; i < int(result.RowsCount); i++ {
		result.Rows[i] = make([]any, result.ColumnsCount)
		for j := 0; j < int(result.ColumnsCount); j++ {
			value, err := datatypes.Read(responsePayload)
			if err != nil {
				return nil, err
			}
			result.Rows[i][j] = value
		}
	}
	if err := binary.Read(responsePayload, binary.LittleEndian, &result.HasMore); err != nil {
		return nil, err
	}
	return result, nil
}
