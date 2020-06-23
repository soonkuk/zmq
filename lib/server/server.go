package server

import (
	"log"
	"os"

	zmq "github.com/pebbe/zmq4"

	"github.com/soonkuk/zmq/lib/network"
)

type ServerImpl struct {
	config       ConfigServer
	collector    *network.CollectorZmq
	reporterList map[ /* reporter deviceID */ string] /* device type */ string
}

func NewServerImpl(config ConfigServer) (*ServerImpl, error) {
	var collector *network.CollectorZmq
	var err error
	collector, err = network.NewCollectorZmq()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	server := &ServerImpl{
		collector: collector,
		config:    config,
	}
	return server, nil
}

func (s *ServerImpl) Init() error {
	if err := s.collector.Bind(s.config.endpoint); err != nil {
		log.Print("#server: ", err)
		return err
	}
	return nil
}

func (s *ServerImpl) Run() {
	backend, _ := zmq.NewSocket(zmq.DEALER)
	defer backend.Close()
	backend.Bind("inproc://backend")
	for i := 0; i < s.config.workers; i++ {
		go server_worker(i)
	}
	err := zmq.Proxy(s.collector.Router, backend, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	// s.Listen()
}

func (s *ServerImpl) Listen() {
	for {
		if id, m, err := s.collector.Receive(); err != nil {
			log.Print("#server: ", err)
		} else {
			log.Print("#server: ", id, m)
		}
	}
}

func (s *ServerImpl) Stop() {
	log.Print("#################### server closed")
	s.collector.Close()
}
