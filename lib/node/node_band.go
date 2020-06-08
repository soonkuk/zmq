package node

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/soonkuk/zmq/lib/common"
)

var mu sync.Mutex

type NodeBand struct {
	config      ConfigNode
	location    common.Location
	temperature float32
	hrv         int16
	ecg         int16
	timestamp   time.Time
	reporter    *common.ReporterZmq
}

func NewNodeBand(config ConfigNode) (*NodeBand, error) {
	var reporter *common.ReporterZmq
	var err error
	node := &NodeBand{
		config: config,
	}
	reporter, err = common.NewReporterZmq(config.deviceID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	node.reporter = reporter

	return node, nil
}

func (n *NodeBand) Init() {
	n.reporter.Bind(common.DefaultReporterEndPoint)
}

func (n *NodeBand) Run() {
	go n.Report("hello")
	// go n.handleMessage()
	// n.Stop()
}

func (n *NodeBand) InitAndRun() {
	// go n.reporter.BindAndSendAndReceive(common.DefaultReporterEndPoint)
	// go n.reporter.ClientTask(common.DefaultReporterEndPoint)
	go n.Report("hello")
	go n.HandleMessage()
	// n.Stop()
}

func (n *NodeBand) HandleMessage() {
	for {
		var m, id string
		m, id = n.reporter.Receive()
		log.Print("#node_band: ", id, m)
	}
}

func (n *NodeBand) Report(m string) {

	var err error
	for {
		mu.Lock()
		if err = n.reporter.Send(time.Now().String()); err != nil {
			log.Print("#node_band: ", err)
		}
		mu.Unlock()
		id, _ := n.reporter.Dealer.GetIdentity()
		log.Print(id)
		sleep_time := rand.Intn(10000)
		time.Sleep(time.Duration(sleep_time) * time.Millisecond)
	}

	/*
		for {
			var m, id string
			var err error
			mu.Lock()
			if m, id, err = n.reporter.Receive(); err != nil {
				log.Print("#node_band: ", err)
			}
			log.Print("#node_band: ", id, m)
			mu.Unlock()
		}
	*/
}

func (n *NodeBand) Stop() {
	n.reporter.Close()
}
