package main

import (
	"fmt"
	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"github.com/implement-pbft-pos/tendermint/node"
	"log"
	"time"
)

func main() {
	fmt.Println("Starting...")

	numNodes := 3
	beginPort := 9000
	ports := make([]uint16, numNodes)
	nodes := make([]node.CompliantNode, numNodes)

	for idx := 0; idx < numNodes; idx++ {
		ports[idx] = uint16(beginPort + idx)
	}

	for idx := 0; idx < numNodes; idx++ {
		var compliantNode = node.CompliantNode{}
		log.Printf("compliantNode %p\n", &compliantNode)
		listenPort := ports[idx]
		neighborCfgs := make([]core_rpc.GClientConfig, 0)
		for idy := 0; idy < numNodes; idy++ {
			if idx != idy {
				neighborCfg := core_rpc.GClientConfig{Host: "0.0.0.0", Port: ports[idy]}
				neighborCfgs = append(neighborCfgs, neighborCfg)
			}
		}
		compliantNode.InitNode(idx, listenPort, neighborCfgs)
		nodes[idx] = compliantNode
	}

	for idx, _ := range nodes {
		compliantNode := nodes[idx]
		//log.Printf("Node: %d, %d, %v, %p\n", compliantNode.Id, compliantNode.Port, compliantNode.Server, &compliantNode)
		compliantNode.StartServer()
	}

	time.Sleep(3 * time.Second)

	for idx, _ := range nodes {
		compliantNode := nodes[idx]
		compliantNode.ConnectNeighborNodes()
	}

	time.Sleep(10 * time.Second)
	//result := new corerpc.GResult{Error: 0, Data: "hehe"}
	//fmt.Println("result: ", result)
	//hello()
}
