package common

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

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
	return &ReporterZmq{}, nil
}

func (rpt *ReporterZmq) Bind(endpoint string) {
	rpt.Dealer.Connect(endpoint)
}

func (rpt *ReporterZmq) BindAndSendAndReceive(endpoint string) {
	var mu sync.Mutex

	rpt.Dealer, _ = zmq.NewSocket(zmq.DEALER)
	defer rpt.Dealer.Close()

	set_id(rpt.Dealer)
	rpt.Dealer.Connect(endpoint)

	go func() {
		for {
			time.Sleep(time.Second)
			mu.Lock()
			rpt.Dealer.SendMessage("hello")
			mu.Unlock()
		}
	}()

	for {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		msg, err := rpt.Dealer.RecvMessage(zmq.DONTWAIT)
		if err != nil {
			log.Print("#reporter_zmq: ", err)
		}
		if err == nil {
			id, _ := rpt.Dealer.GetIdentity()
			fmt.Println(msg[0], id)
		}
		mu.Unlock()
		sleep_time := rand.Intn(10000)
		time.Sleep(time.Duration(sleep_time) * time.Millisecond)
	}
}

func (rpt *ReporterZmq) ClientTask(endpoint string) {
	var mu sync.Mutex

	rpt.Dealer, _ = zmq.NewSocket(zmq.DEALER)
	defer rpt.Dealer.Close()

	//  Set random identity to make tracing easier
	set_id(rpt.Dealer)
	rpt.Dealer.Connect(endpoint)

	go func() {
		for request_nbr := 1; true; request_nbr++ {
			time.Sleep(time.Second)
			mu.Lock()
			rpt.Dealer.SendMessage(fmt.Sprintf("request #%d", request_nbr))
			mu.Unlock()
		}
	}()

	for {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		msg, err := rpt.Dealer.RecvMessage(zmq.DONTWAIT)
		if err == nil {
			id, _ := rpt.Dealer.GetIdentity()
			fmt.Println(msg[0], id)
		}
		mu.Unlock()
	}
}

func (rpt *ReporterZmq) Send(m string) error {
	rpt.mu.Lock()
	_, err := rpt.Dealer.SendMessage(m)
	rpt.mu.Unlock()
	return HandleError(err)
}

func (rpt *ReporterZmq) Receive() (string, string) {
	rpt.mu.Lock()
	reply, _ := rpt.Dealer.RecvMessage(zmq.DONTWAIT)
	id, _ := rpt.Dealer.GetIdentity()
	rpt.mu.Unlock()
	return string(reply[0]), id
}

func (rpt *ReporterZmq) Close() {
	rpt.Dealer.Close()
}

func set_id(soc *zmq.Socket) {
	identity := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	soc.SetIdentity(identity)
}
