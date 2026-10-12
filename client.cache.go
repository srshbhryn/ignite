package ignite

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/srshbhryn/ignite/datatypes"
)

var ErrDoesNotExists = errors.New("entry does not exists")

func (c *Client) Get(cacheName string, key string) (string, error) {
	requestPayload := &bytes.Buffer{}
	code := hashCode([]byte(cacheName))
	if err := binary.Write(requestPayload, binary.LittleEndian, &code); err != nil {
		return "", err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, byte(0)); err != nil {
		return "", err
	}
	if err := datatypes.Write(requestPayload, key); err != nil {
		return "", err
	}
	responsePayload, err := c.sendRequest(_OP_CACHE_GET, requestPayload)
	if err != nil {
		return "", err
	}
	value, err := datatypes.Read(responsePayload)
	if err != nil {
		return "", nil
	}
	if value == nil {
		return "", ErrDoesNotExists
	}
	return value.(string), err
}

func (c *Client) Put(cacheName string, key string, value string) error {
	requestPayload := &bytes.Buffer{}
	code := hashCode([]byte(cacheName))
	if err := binary.Write(requestPayload, binary.LittleEndian, &code); err != nil {
		return err
	}
	if err := binary.Write(requestPayload, binary.LittleEndian, byte(0)); err != nil {
		return err
	}
	if err := datatypes.Write(requestPayload, key); err != nil {
		return err
	}
	if err := datatypes.Write(requestPayload, value); err != nil {
		return err
	}
	_, err := c.sendRequest(_OP_CACHE_PUT, requestPayload)
	if err != nil {
		return err
	}
	return nil
}
