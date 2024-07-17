package client

import (
	"crypto/tls"
	"fmt"
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
