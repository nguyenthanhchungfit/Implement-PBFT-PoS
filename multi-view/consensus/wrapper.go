package consensus

import (
	"fmt"
	core_rpc "github.com/implement-pbft-pos/multi-view/corerpc"
)

type NeighborNode struct {
	NodeId int32
	PublicKey []byte
	Client *core_rpc.GClient
}

func (neighborNode NeighborNode) String() string {
	return fmt.Sprintf("NeighborNode(Id: %d, Host: %s, Port: %d, pubKey: %x), ", neighborNode.NodeId,
		neighborNode.Client.ClientConfig.Host,  neighborNode.Client.ClientConfig.Port, neighborNode.PublicKey)
}
