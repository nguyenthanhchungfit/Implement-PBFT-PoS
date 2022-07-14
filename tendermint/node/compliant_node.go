package node

import (
	"github.com/implement-pbft-pos/tendermint/consensus"
	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"github.com/implement-pbft-pos/tendermint/utils"
	"sync"
)

type CompliantNode struct {
	Id int32
	Port uint16
	Server core_rpc.GServer
	Processor consensus.CompliantProcessor
	NeighborNodes []* NeighborNode
}


func (node* CompliantNode) InitNode(id int32, listenPort uint16, keyPair *utils.KeyPair, neighborNodes []*NeighborNode,
	selector *consensus.ProposerSelector, consensusCfg *consensus.ConsensusConfig){
	node.Id = id
	node.Port = listenPort
	gServerCfg := core_rpc.GServerConfig{Port : node.Port, Id: id}
	node.Processor = consensus.CompliantProcessor{Selector: selector, NodeId: id, ConsensusConfig: consensusCfg, KeyPair: keyPair}
	node.Server = core_rpc.GServer{ServerCfg: gServerCfg, Processor: node.Processor }
	node.NeighborNodes = neighborNodes
}

func (node* CompliantNode) StartServer(wg * sync.WaitGroup){
	wg.Add(1)
	go func() {
		defer wg.Done()
		node.Server.StartServer()
	}()
}

func (node* CompliantNode) ConnectNeighborNodes(){
	for _, neighborNode := range node.NeighborNodes {
		neighborNode.Client.ConnectToRemote()
	}
	utils.InfoStdOutLogger.Printf("Node %d connect to neighborNodes: %v\n", node.Id, node.NeighborNodes)
}

func (node* CompliantNode) StartConsensus(){
	node.Processor.StartRound(0)
}