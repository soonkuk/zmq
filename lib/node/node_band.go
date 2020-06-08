package node

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/soonkuk/zmq/lib/common"
	"github.com/soonkuk/zmq/lib/network"
)

var mu sync.Mutex

type NodeBand struct {
	config         ConfigNode
	timestamp      time.Time
	queryResponser *common.QueryResponser
	reporter       *network.ReporterZmq
}

func NewNodeBand(config ConfigNode) (*NodeBand, error) {
	var err error
	node := &NodeBand{
		config: config,
	}
	node.reporter, err = network.NewReporterZmq(config.deviceID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	node.queryResponser = common.NewQueryResponser()
	return node, nil
}

func (n *NodeBand) Init() {
	n.reporter.Bind(common.DefaultReporterEndPoint)
}

func (n *NodeBand) Run() {
	go n.Report()
	n.HandleMessage()
	// n.Stop()
}

func (n *NodeBand) InitAndRun() {
	go n.reporter.BindAndSendAndReceive(common.DefaultReporterEndPoint)
	// go n.reporter.ClientTask(common.DefaultReporterEndPoint)
	// go n.Report("hello")
	// go n.HandleMessage()
	// n.Stop()
}

func (n *NodeBand) HandleMessage() {
	go n.reporter.Receive()
}

func (n *NodeBand) Report() {
	var err error
	var b []byte
	c := make(chan common.Query, 100)
	go n.queryResponser.Run(c)

	for {
		data := <-c
		mu.Lock()
		b, err = data.ToJson()
		if err != nil {
			continue
		}
		id, _ := n.reporter.Dealer.GetIdentity()
		log.Println(id, ":", data)
		if err = n.reporter.Send(b); err != nil {
			log.Print("#node_band: ", err)
		}
		mu.Unlock()
		sleep_time := rand.Intn(1000)
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
