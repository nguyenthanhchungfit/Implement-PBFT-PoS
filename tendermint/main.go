package main

import (
	"github.com/implement-pbft-pos/tendermint/consensus"
	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"github.com/implement-pbft-pos/tendermint/node"
	"github.com/implement-pbft-pos/tendermint/utils"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

func main() {
	fileLog := utils.SetupZeroLogFile("/data/log/implement-pbft/tendermint/logging")
	utils.InitLogConsole()

	utils.InfoStdOutLogger.Println("Starting ...")

	var wg sync.WaitGroup
	numNodes := 5
	beginPort := 9000
	ports := make([]uint16, numNodes)
	nodes := make([]node.CompliantNode, numNodes)
	allNodeIds := make([]int, numNodes)

	for idx := 0; idx < numNodes; idx++ {
		ports[idx] = uint16(beginPort + idx)
		allNodeIds[idx] = idx
	}

	proposerSelector := consensus.ProposerSelector{NodeIds: allNodeIds}

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
		compliantNode.InitNode(idx, listenPort, neighborCfgs, &proposerSelector)
		nodes[idx] = compliantNode
	}

	for idx, _ := range nodes {
		compliantNode := nodes[idx]
		compliantNode.StartServer(&wg)
	}

	time.Sleep(3 * time.Second)

	for idx, _ := range nodes {
		compliantNode := nodes[idx]
		compliantNode.ConnectNeighborNodes()
	}

	for idx, _ := range nodes {
		compliantNode := nodes[idx]
		compliantNode.StartConsensus()
	}

	wg.Wait()

	if fileLog != nil {
		fileLog.Close()
	}
	//result := new corerpc.GResult{Error: 0, Data: "hehe"}
	//fmt.Println("result: ", result)
	//hello()
}
