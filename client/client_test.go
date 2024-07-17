package client_test

import (
	"ignite"
	"ignite/client"
	"testing"
)

func Test_Connect(t *testing.T) {
	c := client.NewClient("ignite", "ignite")
	err := c.Connect(&ignite.ConnectionOpts{
		Host: "127.0.0.1",
		Port: 10855,
	})
	if err != nil {
		t.Error(err)
	}
}

func Test_Get_Put(t *testing.T) {

	c := client.NewClient("ignite", "ignite")
	err := c.Connect(&ignite.ConnectionOpts{
		Host: "127.0.0.1",
		Port: 10855,
	})
	if err != nil {
		t.Error(err)
	}
	_, err = c.Get("default", "test_key")
	if err != client.ErrDoesNotExists {
		t.Error(err)
	}
	err = c.Put("default", "test_key", "test_value")
	if err != nil {
		t.Error(err)
	}
	value, err := c.Get("default", "test_key")
	if err != nil {
		t.Error(err)
	}
	if value != "test_value" {
		t.Fail()
	}
}
