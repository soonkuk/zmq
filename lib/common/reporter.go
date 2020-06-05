package common

type Reporter interface {
	Send() error
	Receive() (string, error)
	Close()
}
