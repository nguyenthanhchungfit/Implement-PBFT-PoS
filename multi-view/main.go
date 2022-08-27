//package main
//
//import (
//	"fmt"
//	"github.com/implement-pbft-pos/tendermint/consensus"
//	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
//	"github.com/implement-pbft-pos/tendermint/node"
//	"github.com/implement-pbft-pos/tendermint/utils"
//	"sync"
//	"time"
//)
//
//func main() {
//	fileLog := utils.SetupZeroLogFile("/data/log/implement-pbft/tendermint/logging")
//	utils.InitLogConsole()
//
//	utils.InfoStdOutLogger.Println("Starting ...")
//
//	var wg sync.WaitGroup
//	numNodes := 5
//	beginPort := 9000
//	nodeIds := make([]int32, numNodes)
//	nodeInfos := make([]*node.NodeInfo, numNodes)
//	compliantNodes := make([]node.CompliantNode, numNodes)
//
//	consensusCfg := consensus.ConsensusConfig{
//		TimeoutPropose:   1000 * time.Millisecond,
//		TimeoutPreVote:   1000 * time.Millisecond,
//		TimeoutPreCommit: 1000 * time.Millisecond,
//	}
//
//	// Init list node
//	for idx := 1; idx <= numNodes; idx++ {
//		host := "0.0.0.0"
//		nodeId := int32(idx)
//		listenPort := uint16(beginPort + idx - 1)
//		keyPair, _ := utils.GenerateNewKeyPair()
//		nodeInfo := node.NodeInfo{Host: host, Id: nodeId, ListenPort: listenPort, KeyPair: keyPair}
//		nodeInfos[idx - 1] = &nodeInfo
//		nodeIds[idx - 1] = nodeId
//	}
//	proposerSelector := consensus.ProposerSelector{NodeIds: nodeIds}
//
//	// Start Server and Init List Neighbors for nodes
//	idxCompliantNode := 0
//	for _, nodeInfo := range nodeInfos {
//		var compliantNode = node.CompliantNode{}
//
//		neighborNodes := make([]*consensus.NeighborNode, numNodes-1)
//		idxNeighbor := 0
//		for _, clientNodeInfo := range nodeInfos {
//			if clientNodeInfo.Id != nodeInfo.Id {
//				gClient := core_rpc.GClient{ClientConfig: core_rpc.GClientConfig{Host: clientNodeInfo.Host, Port: clientNodeInfo.ListenPort}}
//				neighborNode := consensus.NeighborNode{NodeId: clientNodeInfo.Id, PublicKey: clientNodeInfo.KeyPair.PublicKey,
//					Client: &gClient}
//				neighborNodes[idxNeighbor] = &neighborNode
//				idxNeighbor++
//			}
//		}
//		compliantNode.InitNode(nodeInfo.Id, nodeInfo.ListenPort, nodeInfo.KeyPair, neighborNodes, &proposerSelector, &consensusCfg)
//		compliantNode.StartServer(&wg)
//		compliantNodes[idxCompliantNode] = compliantNode
//		idxCompliantNode++
//	}
//
//	time.Sleep(1 * time.Second)
//
//	fmt.Println("********** Start connect to neighbor nodes **********")
//	for _, compliantNode := range compliantNodes {
//		node := compliantNode
//		node.ConnectNeighborNodes()
//	}
//
//	time.Sleep(1 * time.Second)
//	fmt.Println("********** Start Consensus **********")
//	for idx, _ := range compliantNodes {
//		compliantNode := compliantNodes[idx]
//		compliantNode.StartConsensus()
//		//compliantNode.StartPing()
//	}
//
//	wg.Wait()
//
//	if fileLog != nil {
//		fileLog.Close()
//	}
//	//result := new corerpc.GResult{Error: 0, Data: "hehe"}
//	//fmt.Println("result: ", result)
//	//hello()
//}
