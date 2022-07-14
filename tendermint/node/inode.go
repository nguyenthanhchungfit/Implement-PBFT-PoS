package node

import (
	"fmt"
	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"github.com/implement-pbft-pos/tendermint/utils"
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

type NodeInfo struct {
	Host string
	Id int32
	ListenPort uint16
	KeyPair *utils.KeyPair
}
