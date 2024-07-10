package client

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"ignite/client/datatypes"
	"ignite/types"
	"net"
	"sync"
)

const _NETWORK = "tcp"
const _BINARY_PROTO_VERSION_MAJOR = int16(1)
const _BINARY_PROTO_VERSION_MINOR = int16(1)
const _BINARY_PROTO_VERSION_PATCH = int16(0)

type Client struct {
	username string
	password string
	conn     net.Conn
	mutex    sync.Mutex
}

func NewClient(username string, password string) *Client {
	return &Client{
		username: username,
		password: password,
	}
}

func (c *Client) Connect(opts *types.ConnectionOpts) error {
	err := c.createConnection(opts)
	if err != nil {
		return err
	}
	return c.doHandShake()
}

func (c *Client) createConnection(opts *types.ConnectionOpts) error {
	address := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	var conn net.Conn
	var err error
	dialer := &net.Dialer{}
	if opts.TLSConfig != nil {
		conn, err = tls.DialWithDialer(dialer, _NETWORK, address, opts.TLSConfig)
	} else {
		conn, err = dialer.Dial(_NETWORK, address)
	}
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

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
