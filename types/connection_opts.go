package types

import (
	"crypto/tls"
)

type ConnectionOpts struct {
	Host      string
	Port      int
	TLSConfig *tls.Config
}
