package client

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/soonkuk/zmq/lib/common"
	"github.com/soonkuk/zmq/lib/network"
)

var mu sync.Mutex

type ClientImpl struct {
	config         *ConfigClient
	timestamp      time.Time
	queryResponser *common.QueryResponser
	reporter       *network.ReporterZmq
}

func NewClientImpl(config *ConfigClient) (*ClientImpl, error) {
	var err error
	clnt := &ClientImpl{
		config: config,
	}
	clnt.reporter, err = network.NewReporterZmq(config.deviceID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if config.status == CorrectClnt {
		clnt.queryResponser = common.NewQueryResponser(common.CorrectResponser)
	} else if config.status == FailClnt {
		clnt.queryResponser = common.NewQueryResponser(common.FailResponser)
	} else {
		clnt.queryResponser = common.NewQueryResponser(common.TestResponser)
	}

	return clnt, nil
}

func (c *ClientImpl) Init() {
	c.reporter.Bind(c.config.endpoint)
}

func (c *ClientImpl) Run() {
	if c.queryResponser.Status == common.TestResponser {
		go c.SendHrtbt()
	} else {
		go c.Report()
	}
	go c.HandleMessage()
	// c.Stop()
}

func (c *ClientImpl) InitAndRun() {
	go c.reporter.BindAndSendAndReceive(c.config.endpoint)
	// go n.reporter.ClientTask(common.DefaultReporterEndPoint)
	// go n.Report("hello")
	// go n.HandleMessage()
	// n.Stop()
}

func (c *ClientImpl) HandleMessage() {
	for {
		mu.Lock()
		id, msg, err := c.reporter.Receive()
		if err != nil {
			//time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
		if len(msg) != 0 {
			log.Println(id, " : ", msg)
			c.config.counter.Recv++
			log.Printf("Count messages - Send message : %d #### Received message : %d", c.config.counter.Send, c.config.counter.Recv)
		}
		mu.Unlock()
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}

func (c *ClientImpl) SendHrtbt() {
	var err error
	ch := make(chan interface{}, 100)
	go c.queryResponser.Run(ch)

	for {
		data := <-ch
		val, ok := data.(string)
		if !ok {
			log.Print("#client: ", err)
		}
		mu.Lock()
		id, _ := c.reporter.Dealer.GetIdentity()
		log.Println(id, ":", val)
		if err = c.reporter.Send([]byte(val)); err != nil {
			log.Print("#client: ", err)
		}
		c.config.counter.Send++
		mu.Unlock()
		time.Sleep(time.Duration(rand.Intn(c.config.interval)) * time.Second)
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

func (c *ClientImpl) Report() {
	var err error
	var b []byte
	ch := make(chan interface{}, 100)
	go c.queryResponser.Run(ch)

	for {
		time.Sleep(time.Duration(rand.Intn(c.config.interval)) * time.Second)
		data := <-ch
		val, ok := data.(common.Query)
		if !ok {
			log.Print("#client: ", err)
		}
		mu.Lock()
		b, err = val.ToJson()
		if err != nil {
			continue
		}
		// id, _ := c.reporter.Dealer.GetIdentity()
		// log.Println(id, ":", data)
		if err = c.reporter.Send(b); err != nil {
			log.Print("#client: ", err)
		}
		c.config.counter.Send++
		// log.Println("Count send message : ", c.config.counter.Send)
		mu.Unlock()
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

func (c *ClientImpl) Stop() {
	c.reporter.Close()
}
