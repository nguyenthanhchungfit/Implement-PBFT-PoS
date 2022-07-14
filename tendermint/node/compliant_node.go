package node

import (
	"github.com/implement-pbft-pos/tendermint/consensus"
	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"github.com/implement-pbft-pos/tendermint/utils"
	"sync"
)

type CompliantNode struct {
	Id int
	Port uint16
	Server core_rpc.GServer
	Processor consensus.CompliantProcessor
	NeighborConfigs []core_rpc.GClientConfig
	NeighborNodes []core_rpc.GClient
}

func (node* CompliantNode) addNeighborNode(idx int, neighborNode core_rpc.GClient){
	node.NeighborNodes[idx] = neighborNode
}

func (node* CompliantNode) InitNode(id int, listenPort uint16, gClientConfigs []core_rpc.GClientConfig, selector *consensus.ProposerSelector){
	node.Id = id
	node.Port = listenPort
	gServerCfg := core_rpc.GServerConfig{Port : node.Port, Id: id}
	node.Processor = consensus.CompliantProcessor{Selector: selector, NodeId: id}
	node.Server = core_rpc.GServer{ServerCfg: gServerCfg, Processor: node.Processor }
	node.NeighborConfigs = gClientConfigs
	node.NeighborNodes = make([]core_rpc.GClient, len(node.NeighborConfigs))
}

func (node* CompliantNode) StartServer(wg * sync.WaitGroup){
	wg.Add(1)
	go func() {
		defer wg.Done()
		node.Server.StartServer()
	}()
}

func (node* CompliantNode) ConnectNeighborNodes(){
	for idx, neighborCfg := range node.NeighborConfigs {
		neighbor := core_rpc.NewGClient(neighborCfg)
		node.addNeighborNode(idx, *neighbor)
	}
	utils.InfoStdOutLogger.Printf("Node %d connect to neighborNodes: %v\n", node.Id, node.NeighborConfigs)
	node.Processor.NeighborClients = node.NeighborNodes

	//count := 0
	//for _, neighbor := range node.NeighborNodes {
	//	count++
	//	msgData := fmt.Sprintf("msg %d from node %d\n", count, node.Id)
	//	msg := core_rpc.GProposeMessage{Height: 0, Round: 0, ValidRound: -1, Data: &core_rpc.GData{Data: msgData}}
	//	neighbor.SendProposeMessage(&msg)
	//
	//	h := sha256.New()
	//	h.Write([]byte("hello world\n"))
	//	out := h.Sum(nil)
	//	fmt.Printf("%x", out)
	//
	//	msgPreVote := core_rpc.GPreVoteMessage{Height: 1, Round: 1, HashValue: out}
	//	neighbor.SendPreVoteMessage(&msgPreVote)
	//
	//	msgPreCommit := core_rpc.GPreCommitMessage{Height: 2, Round: 3, HashValue: out}
	//	neighbor.SendPreCommitMessage(&msgPreCommit)
	//}
}

func (node* CompliantNode) StartConsensus(){
	node.Processor.StartRound(0)
}