package client

type ConfigClient struct {
	deviceID   string
	deviceType string
	interval   int
	status     StatusClient
	endpoint   string
}

func NewConfigClient(id string, dtype string, interval int, status StatusClient, endpoint string) ConfigClient {
	return ConfigClient{
		deviceID:   id,
		deviceType: dtype,
		interval:   interval,
		status:     status,
		endpoint:   endpoint,
	}
}

type StatusClient string

const (
	CorrectClnt StatusClient = "Correct Client"
	FailClnt    StatusClient = "Fail Client"
	TestClnt    StatusClient = "Test Client"
)
