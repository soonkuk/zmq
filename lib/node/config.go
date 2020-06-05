package node

type ConfigNode struct {
	deviceID   string
	deviceType string
}

func NewConfigNode(id string, dtype string) ConfigNode {
	return ConfigNode{
		deviceID:   id,
		deviceType: dtype,
	}
}
