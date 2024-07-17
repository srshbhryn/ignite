package types

type Client interface {
	Connect(ConnectionOpts) error
	Get(string, string) (string, error)
	Put(string, string, string) error
}
