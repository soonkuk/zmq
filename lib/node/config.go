package node

type ConfigNode struct {
	deviceID   string
	deviceType string
	interval   int
	status     NodeStatus
}

func NewConfigNode(id string, dtype string, interval int, status NodeStatus) ConfigNode {
	return ConfigNode{
		deviceID:   id,
		deviceType: dtype,
		interval:   interval,
		status:     status,
	}
}

type NodeStatus string

const (
	CorrectNode NodeStatus = "CorrectNode"
	FailNode    NodeStatus = "FailNode"
)
