package node

import (
	"fmt"
	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"log"
)

type CompliantNode struct {
	Id int
	Port uint16
	Server core_rpc.GServer
	NeighborConfigs []core_rpc.GClientConfig
	NeighborNodes []core_rpc.GClient
}

func (node* CompliantNode) addNeighborNode(idx int, neighborNode core_rpc.GClient){
	node.NeighborNodes[idx] = neighborNode
}

func (node* CompliantNode) InitNode(id int, listenPort uint16, gClientConfigs []core_rpc.GClientConfig){
	node.Id = id
	node.Port = listenPort
	gServerCfg := core_rpc.GServerConfig{Port : node.Port, Id: id}
	node.Server = core_rpc.GServer{ServerCfg: gServerCfg}
	node.NeighborConfigs = gClientConfigs
	node.NeighborNodes = make([]core_rpc.GClient, len(node.NeighborConfigs))
}

func (node* CompliantNode) StartServer(){
	log.Printf("Node %d start Server at %d\n", node.Id, node.Port)
	go func() {
		log.Printf("go Node %d start Server at %d\n", node.Id, node.Port)
		node.Server.StartServer()
	}()
}

func (node* CompliantNode) ConnectNeighborNodes(){
	log.Printf("Node %d connectNeighborNodes\n", node.Id)
	for idx, neighborCfg := range node.NeighborConfigs {
		neighbor := core_rpc.NewGClient(neighborCfg)
		node.addNeighborNode(idx, *neighbor)
	}

	count := 0
	for _, neighbor := range node.NeighborNodes {
		count++
		msgData := fmt.Sprintf("msg %d from node %d\n", count, node.Id)
		msg := core_rpc.GProposeMessage{Height: 0, Round: 0, ValidRound: -1, Data: &core_rpc.GData{Data: msgData}}
		neighbor.SendProposeMessage(&msg)
	}
}