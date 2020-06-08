package server

import (
	"log"
	"os"

	zmq "github.com/pebbe/zmq4"

	"github.com/soonkuk/zmq/lib/common"
	"github.com/soonkuk/zmq/lib/network"
)

type ServerZmq struct {
	collector    *network.CollectorZmq
	reporterList map[ /* reporter deviceID */ string] /* device type */ string
}

func NewServerZmq() (*ServerZmq, error) {
	var collector *network.CollectorZmq
	var err error
	collector, err = network.NewCollectorZmq()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	server := &ServerZmq{
		collector: collector,
	}
	return server, nil
}

func (s *ServerZmq) Init() error {
	if err := s.collector.Bind(common.DefaultCollectorEndPoint); err != nil {
		log.Print("#server: ", err)
		return err
	}
	return nil
}

func (s *ServerZmq) Run() {
	backend, _ := zmq.NewSocket(zmq.DEALER)
	defer backend.Close()
	backend.Bind("inproc://backend")
	for i := 0; i < 5; i++ {
		go server_worker(i)
	}
	err := zmq.Proxy(s.collector.Router, backend, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	// s.Listen()
}

func (s *ServerZmq) Listen() {
	for {
		if id, m, err := s.collector.Receive(); err != nil {
			log.Print("#server: ", err)
		} else {
			log.Print("#server: ", id, m)
		}
	}
}

func (s *ServerZmq) Stop() {
	log.Print("#################### server closed")
	s.collector.Close()
}
