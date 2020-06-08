package network

type Reporter interface {
	Send() error
	Receive() (string, error)
	Close()
}
