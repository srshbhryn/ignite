package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"ignite/client/datatypes"
)

func (c *Client) doHandShake() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	buffer := &bytes.Buffer{}
	if err := binary.Write(buffer, binary.LittleEndian, byte(1)); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.LittleEndian, _BINARY_PROTO_VERSION_MAJOR); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.LittleEndian, _BINARY_PROTO_VERSION_MINOR); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.LittleEndian, _BINARY_PROTO_VERSION_PATCH); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.LittleEndian, byte(2)); err != nil {
		return err
	}
	if err := datatypes.Write(buffer, c.username); err != nil {
		return err
	}
	if err := datatypes.Write(buffer, c.password); err != nil {
		return err
	}
	payloadLength := int32(buffer.Len())
	if err := binary.Write(c.conn, binary.LittleEndian, payloadLength); err != nil {
		return err
	}
	_, err := c.conn.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	payloadLength = 0
	if err := binary.Read(c.conn, binary.LittleEndian, &payloadLength); err != nil {
		return err
	}
	if payloadLength < 1 {
		return errors.New("invalid handshake response size")
	}
	response := make([]byte, int(payloadLength))
	if err := binary.Read(c.conn, binary.LittleEndian, &response); err != nil {
		return err
	}
	if response[0] != 1 {
		return errors.New("handshake failed")
	}
	return nil
}
