package node

import (
	"fmt"
	"github.com/implement-pbft-pos/multi-view/consensus"
	core_rpc "github.com/implement-pbft-pos/multi-view/corerpc"
	"github.com/implement-pbft-pos/multi-view/middleware"
	"github.com/implement-pbft-pos/multi-view/utils"
	"sync"
)

type CompliantNode struct {
	Id int32
	Port uint16
	Server core_rpc.GServer
	Processor *consensus.CompliantProcessor
	NeighborNodes []* consensus.NeighborNode
}


func (node* CompliantNode) InitNode(id int32, listenPort uint16, keyPair *utils.KeyPair, neighborNodes []*consensus.NeighborNode,
	selector *consensus.ProposerSelector, consensusCfg *consensus.ConsensusConfig){
	node.Id = id
	node.Port = listenPort
	node.NeighborNodes = neighborNodes
	neighborNodesId := make([]int32, len(neighborNodes))
	for idx, neighborNode := range neighborNodes {
		neighborNodesId[idx] = neighborNode.NodeId
	}
	gServerCfg := core_rpc.GServerConfig{Port : node.Port, Id: id}
	node.Processor = &consensus.CompliantProcessor{Selector: selector, NodeId: id, ConsensusConfig: consensusCfg, KeyPair: keyPair,
		NeighborNodeIds: neighborNodesId, NeighborNodes: neighborNodes}
	middleware := middleware.ConsensusMiddleware{Processor: node.Processor}
	node.Server = core_rpc.GServer{ServerCfg: gServerCfg, Processor: middleware }

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
	go func() {
		node.Processor.StartConsensus()
	}()
}

func (node* CompliantNode) StartPing(){
	go func() {
		fmt.Println("alo alo")
		node.Processor.TestPing()
	}()
}