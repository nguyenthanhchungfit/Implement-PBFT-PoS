package node

import core_rpc "github.com/implement-pbft-pos/multi-view/corerpc"

type MaliciousNode struct {
	Id int32
	Server core_rpc.GServer
	NeighborPorts []int32
	NeighborNodes []core_rpc.GClient
}

func (node* MaliciousNode) addNeighborNode(neighborNode core_rpc.GClient){

}

func (node* MaliciousNode) InitNode(id int32, listenPort int32, ){
	node.Id = id

}

func (node* MaliciousNode) StartServer(){
	go func() {
		node.Server.StartServer()
	}()
}

func (node* MaliciousNode) connectNeighborNodes(){

}