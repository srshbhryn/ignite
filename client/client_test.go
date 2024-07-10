package client_test

import (
	"ignite"
	"ignite/client"
	"testing"
)

func Test_Handshake(t *testing.T) {
	c := client.NewClient("ignite", "ignite")
	err := c.Connect(&ignite.ConnectionOpts{
		Host: "127.0.0.1",
		Port: 10855,
	})
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(err)
}
