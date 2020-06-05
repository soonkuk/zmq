package node

type Node interface {
	Init() error
	Run()
}
