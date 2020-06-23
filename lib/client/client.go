package client

type Client interface {
	Init() error
	Run()
}
