package types

type Client interface {
	Connect(ConnectionOpts) error
}
