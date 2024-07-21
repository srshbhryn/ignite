package ignite

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"golang.org/x/exp/rand"
)

func (c *Client) sendRequest(opCode int16, payload *bytes.Buffer) (*bytes.Buffer, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if err := binary.Write(c.conn, binary.LittleEndian, int32(10+payload.Len())); err != nil {
		return nil, err
	}
	if err := binary.Write(c.conn, binary.LittleEndian, &opCode); err != nil {
		return nil, err
	}
	requestId := rand.Int63()
	if err := binary.Write(c.conn, binary.LittleEndian, &requestId); err != nil {
		return nil, err
	}
	if _, err := c.conn.Write(payload.Bytes()); err != nil {
		return nil, err
	}
	var payloadLength int32
	if err := binary.Read(c.conn, binary.LittleEndian, &payloadLength); err != nil {
		return nil, err
	}
	if payloadLength < 1 {
		return nil, errors.New("invalid response size")
	}
	var receivedRequestId int64
	if err := binary.Read(c.conn, binary.LittleEndian, &receivedRequestId); err != nil {
		return nil, errors.New("failed to read request id")
	}
	var statusCode int32
	if err := binary.Read(c.conn, binary.LittleEndian, &statusCode); err != nil {
		return nil, errors.New("failed to read status code")
	}
	if statusCode != 0 {
		var errorsMessage = make([]byte, payloadLength-8-4)
		if err := binary.Read(c.conn, binary.LittleEndian, errorsMessage); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("request failed: %s", string(errorsMessage))
	}
	responseBuffer := &bytes.Buffer{}
	responsePayload := make([]byte, payloadLength-12)
	_, err := c.conn.Read(responsePayload)
	if err != nil {
		return nil, err
	}
	responseBuffer.Write(responsePayload)
	return responseBuffer, nil
}
