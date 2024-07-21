package client_test

import (
	"ignite"
	"ignite/client"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/rand"
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
func Test_SQL(t *testing.T) {

	type Person struct {
		id          int64
		orgId       int16
		groupId     int32
		name        string
		createdAt   time.Time
		salary      decimal.Decimal
		floatValue  float32
		doubleValue float64
		isDead      bool
	}

	c := client.NewClient("ignitea", "ignite")
	err := c.Connect(&ignite.ConnectionOpts{
		Host: "127.0.0.1",
		Port: 10855,
	})
	if err != nil {
		panic(err)
	}
	p := Person{
		id:          rand.Int63(),
		orgId:       4,
		groupId:     57,
		name:        "mamad sibil",
		createdAt:   time.Now(),
		salary:      decimal.NewFromFloat(470.014),
		floatValue:  123.12,
		doubleValue: 123.512,
		isDead:      false,
	}
	_, err = c.QuerySqlFields("default", &client.QueryFieldsRequest{
		Schema: "default",
		Statement: "INSERT INTO Person(_key, id, orgId, groupId, name, createdAt, " +
			"salary, " +
			" floatValue, doubleValue, isDead) VALUES" +
			"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		Args: []any{p.id, p.id, p.orgId, p.groupId, p.name, p.createdAt,
			p.salary,
			p.floatValue, p.doubleValue, p.isDead},
		MaxRows:  100,
		PageSize: 100,
		Timeout:  100,
		Type:     client.ANY,
	})
	if err != nil {
		t.Fail()
	}
	result, err := c.QuerySqlFields("default", &client.QueryFieldsRequest{
		Schema: "default",
		Statement: "Select id, orgId, groupId, name, createdAt," +
			" salary," +
			"floatValue, doubleValue, isDead from Person order by createdAt desc",
		Args:     []any{},
		MaxRows:  100,
		PageSize: 100,
		Timeout:  100,
		Type:     client.SELECT,
	})
	if err != nil {
		t.Fail()
	}
	if p.id != result.Rows[0][0].(int64) {
		t.Fail()
	}
	if p.orgId != result.Rows[0][1].(int16) {
		t.Fail()
	}
	if p.groupId != result.Rows[0][2].(int32) {
		t.Fail()
	}
	if p.name != result.Rows[0][3] {
		t.Fail()
	}
	if !p.createdAt.Equal(result.Rows[0][4].(time.Time)) {
		t.Fail()
	}
	if !p.salary.Equal(result.Rows[0][5].(decimal.Decimal)) {
		t.Fail()
	}
	if p.floatValue != result.Rows[0][6] {
		t.Fail()
	}
	if p.doubleValue != result.Rows[0][7] {
		t.Fail()
	}
	if p.isDead != result.Rows[0][8] {
		t.Fail()
	}
}
