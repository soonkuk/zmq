package common

import (
	"log"
	"sync"

	zmq "github.com/pebbe/zmq4"
)

type ReporterZmq struct {
	mu     sync.Mutex
	Dealer *zmq.Socket
}

func NewReporterZmq(id string) (*ReporterZmq, error) {
	var err error
	var dealer *zmq.Socket
	if dealer, err = zmq.NewSocket(zmq.DEALER); err != nil {
		log.Print("#reporter_zmq: ", err)
		return nil, err
	}
	dealer.SetIdentity(id)
	return &ReporterZmq{Dealer: dealer}, nil
}

func (rpt *ReporterZmq) Bind(endpoint string) {
	rpt.Dealer.Connect(endpoint)
}

func (rpt *ReporterZmq) BindAndSendAndReceive(endpoint string) {
	rpt.Dealer.Connect(endpoint)
	for {
		rpt.mu.Lock()
		rpt.Dealer.SendMessage("hello")
		rpt.mu.Unlock()
		var reply []string
		var err error
		rpt.mu.Lock()
		if reply, err = rpt.Dealer.RecvMessage(zmq.DONTWAIT); err != nil {
			log.Print("#reporter_zmq: ", err)
		}
		rpt.mu.Unlock()
		id, _ := rpt.Dealer.GetIdentity()
		log.Print(reply, id)
	}
}

func (rpt *ReporterZmq) Send(m string) error {
	_, err := rpt.Dealer.SendMessage(m)
	return HandleError(err)
}

func (rpt *ReporterZmq) Receive() (string, string, error) {
	var reply []string
	var err error
	if reply, err = rpt.Dealer.RecvMessage(zmq.DONTWAIT); err != nil {
		log.Print("#reporter_zmq: ", err)
		return "", "", err
	}
	id, _ := rpt.Dealer.GetIdentity()
	return string(reply[0]), id, nil
}

func (rpt *ReporterZmq) Close() {
	rpt.Dealer.Close()
}
