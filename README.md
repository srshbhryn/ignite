# ignite

A lightweight [Apache Ignite](https://ignite.apache.org/) thin client for Go.

`ignite` talks to an Ignite node directly over its **binary thin client protocol**
on a raw TCP socket (with optional TLS). The wire format — handshake, request
framing, and type serialization — is implemented from scratch, with no cgo and
no JDBC/ODBC bridge.

> **Status:** early / work in progress. Connecting, authenticated handshake,
> key/value `Get`/`Put`, and SQL field queries work and are covered by
> integration tests. Most operations listed in [operation.codes.go](operation.codes.go)
> are not implemented yet.

## Features

- TCP and TLS connections
- Username / password handshake (binary protocol v1.1.0)
- Key/value cache operations: `Get`, `Put`
- SQL queries and DML via `QuerySqlFields` (with page size / max rows / timeout)
- Typed value (de)serialization for: `int16`, `int32`, `int64`, `float32`,
  `float64`, `bool`, `string`, `time.Time` (timestamp), `decimal.Decimal`, and `null`

## Requirements

- Go 1.22+
- A reachable Apache Ignite node (tested against Ignite 2.16)

## Installation

```bash
go get ignite
```

> The module is named `ignite` (see [go.mod](go.mod)). Adjust the import path if
> you publish it under a hosted module path.

## Usage

### Connect

```go
c := ignite.NewClient("ignite", "ignite") // username, password

err := c.Connect(&ignite.ConnectionOpts{
    Host: "127.0.0.1",
    Port: 10800,
    // TLSConfig: &tls.Config{...}, // optional
})
if err != nil {
    log.Fatal(err)
}
```

### Key/value cache

```go
if err := c.Put("default", "test_key", "test_value"); err != nil {
    log.Fatal(err)
}

value, err := c.Get("default", "test_key")
if err == ignite.ErrDoesNotExists {
    // key not present
} else if err != nil {
    log.Fatal(err)
}
fmt.Println(value) // "test_value"
```

### SQL fields query

```go
// INSERT / DML
_, err = c.QuerySqlFields("default", &ignite.QueryFieldsRequest{
    Schema:    "default",
    Statement: "INSERT INTO Person(_key, id, name) VALUES (?, ?, ?)",
    Args:      []any{int64(1), int64(1), "alice"},
    PageSize:  100,
    MaxRows:   100,
    Timeout:   time.Second,
    Type:      ignite.ANY,
})

// SELECT
result, err := c.QuerySqlFields("default", &ignite.QueryFieldsRequest{
    Schema:    "default",
    Statement: "SELECT id, name FROM Person ORDER BY id",
    Args:      []any{},
    PageSize:  100,
    MaxRows:   100,
    Timeout:   time.Second,
    Type:      ignite.SELECT,
})
if err != nil {
    log.Fatal(err)
}

for _, row := range result.Rows {
    id := row[0].(int64)
    name := row[1].(string)
    fmt.Println(id, name)
}
// result.HasMore indicates more pages are available on the server cursor.
```

`Type` can be `ignite.ANY`, `ignite.SELECT`, or `ignite.UPDATE`.

## Project layout

| Path | Purpose |
| --- | --- |
| [client.go](client.go) | `Client` type and TCP/TLS connection setup |
| [client.handshake.go](client.handshake.go) | Binary protocol handshake + auth |
| [client.send_request.go](client.send_request.go) | Request/response framing and error handling |
| [client.cache.go](client.cache.go) | `Get` / `Put` cache operations |
| [client.query.sql.fields.go](client.query.sql.fields.go) | `QuerySqlFields` SQL/DML queries |
| [operation.codes.go](operation.codes.go) | Ignite operation code constants |
| [utils.go](utils.go) | Java-compatible `hashCode` for cache name → ID |
| [datatypes/](datatypes/) | Typed value read/write (the wire type system) |
| [test-compose.yml](test-compose.yml), [test-config/](test-config/) | Dockerized Ignite node for tests |

## Testing

The tests are integration tests that run against a real Ignite node. The
provided `Makefile` spins one up with Docker Compose, runs the suite, and tears
it down:

```bash
make test
```

This starts Ignite 2.16 (mapping container port `10800` → host `10855`, see
[test-compose.yml](test-compose.yml)) with the schema in
[test-config/config.xml](test-config/config.xml), which defines a `default`
cache and a `Person` SQL table.

To run the tests manually against an already-running node:

```bash
go test ./... -v
```

> Note: [client_test.go](client_test.go) connects to port `10855` (the host port
> from the compose file).

## Notes & limitations

- Only a subset of cache operations is implemented; the full op-code table exists
  but most ops are not yet wired up.
- `Get` / `Put` currently handle `string` keys and values only.
- A single connection is used per client, guarded by a mutex (no connection
  pooling or cluster discovery / partition awareness).
