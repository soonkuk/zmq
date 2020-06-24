package client

import (
	"github.com/soonkuk/zmq/lib/common"
)

type ConfigClient struct {
	deviceID   string
	deviceType string
	interval   int
	status     StatusClient
	endpoint   string
	counter    *common.Counter
}

func NewConfigClient(id string, dtype string, interval int, status StatusClient, endpoint string, counter *common.Counter) *ConfigClient {
	return &ConfigClient{
		deviceID:   id,
		deviceType: dtype,
		interval:   interval,
		status:     status,
		endpoint:   endpoint,
		counter:    counter,
	}
}

type StatusClient string

const (
	CorrectClnt StatusClient = "Correct Client"
	FailClnt    StatusClient = "Fail Client"
	TestClnt    StatusClient = "Test Client"
)
