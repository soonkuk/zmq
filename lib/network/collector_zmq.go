package network

import (
	"log"

	"github.com/soonkuk/zmq/lib/common"

	zmq "github.com/pebbe/zmq4"
)

// CollectorZmq is a zmq router wrapper
type CollectorZmq struct {
	Router *zmq.Socket
}

// NewCollectorZmq is CollectorZmq Constructor
func NewCollectorZmq() (*CollectorZmq, error) {
	var err error
	var router *zmq.Socket
	if router, err = zmq.NewSocket(zmq.ROUTER); err != nil {
		log.Print("#collector_zmq: ", err)
		return nil, err
	}
	return &CollectorZmq{Router: router}, nil
}

// Bind is zmq router binder
func (c *CollectorZmq) Bind(endpoint string) error {
	var err error
	err = c.Router.Bind(endpoint)
	return common.HandleError(err)
}

// Send
func (c *CollectorZmq) Send(m string) error {
	_, err := c.Router.SendMessage(m)
	return common.HandleError(err)
}

func (c *CollectorZmq) Receive() (head []string, tail []string, err error) {
	var msg []string
	if msg, err = c.Router.RecvMessage(0); err != nil {
		log.Print("#collector_zmq: ", err)
		return nil, nil, err
	}
	if msg[1] == "" {
		head = msg[:2]
		tail = msg[2:]
	} else {
		head = msg[:1]
		tail = msg[1:]
	}
	return
}

func (c *CollectorZmq) Close() {
	c.Router.Close()
}
