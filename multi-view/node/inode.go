package node

import (
	"github.com/implement-pbft-pos/multi-view/utils"
)

type NodeInfo struct {
	Host string
	Id int32
	ListenPort uint16
	KeyPair *utils.KeyPair
}
