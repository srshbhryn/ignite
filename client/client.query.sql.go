package client

import (
	"bytes"
	"encoding/binary"
	"ignite/client/datatypes"
	"time"
)

type QueryRequest struct {
	Table     string
	Statement string
	Args      []any
	PageSize  int32
	Timeout   time.Duration
}

type QueryResult struct {
	CursorId   int64
	RowsCount  int32
	Rows       map[any]any
	IsLastPage bool
}

func (c *Client) QuerySql(cacheName string, request *QueryRequest) (*QueryResult, error) {
	requestPayload := &bytes.Buffer{}
	code := hashCode([]byte(cacheName))
	if err := binary.Write(requestPayload, binary.LittleEndian, &code); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, byte(0)); err != nil {
		return nil, err
	}
	if err := datatypes.Write(requestPayload, &request.Table); err != nil {
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
	if err := binary.Write(requestPayload, binary.LittleEndian, false); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, false); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, false); err != nil {
		return nil, err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, &request.PageSize); err != nil {
		return nil, err
	}
	timeout := int64(request.Timeout / time.Millisecond)
	if err := binary.Write(requestPayload, binary.LittleEndian, &timeout); err != nil {
		return nil, err
	}
	responsePayload, err := c.sendRequest(_OP_QUERY_SQL, requestPayload)
	if err != nil {
		return nil, err
	}
	return readQuerySqlResponse(responsePayload)
}

func (c *Client) GetQuerySqlPage(cacheName string, cursorID int64) (*QueryResult, error) {
	requestPayload := &bytes.Buffer{}
	if err := binary.Write(requestPayload, binary.LittleEndian, &cursorID); err != nil {
		return nil, err
	}
	responsePayload, err := c.sendRequest(_OP_QUERY_SQL_CURSOR_GET_PAGE, requestPayload)
	if err != nil {
		return nil, err
	}
	return readQuerySqlResponse(responsePayload)
}

func readQuerySqlResponse(responsePayload *bytes.Buffer) (*QueryResult, error) {
	result := &QueryResult{}
	if err := binary.Read(responsePayload, binary.LittleEndian, &result.CursorId); err != nil {
		return nil, err
	}
	if err := binary.Read(responsePayload, binary.LittleEndian, &result.RowsCount); err != nil {
		return nil, err
	}
	result.Rows = make(map[any]any)
	for i := 0; i < int(result.RowsCount); i++ {
		key, err := datatypes.Read(responsePayload)
		if err != nil {
			return nil, err
		}
		value, err := datatypes.Read(responsePayload)
		if err != nil {
			return nil, err
		}
		result.Rows[key] = value
	}
	if err := binary.Read(responsePayload, binary.LittleEndian, &result.IsLastPage); err != nil {
		return nil, err
	}
	return result, nil
}
