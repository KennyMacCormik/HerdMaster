package network

type Endpoint interface {
	Run() error
	Close() error
}
