package server

type ConfigServer struct {
	workers  int
	endpoint string
}

func NewConfigServer(workers int, endpoint string) ConfigServer {
	return ConfigServer{
		workers:  workers,
		endpoint: endpoint,
	}
}
