package main

import (
	"fmt"
	"github.com/implement-pbft-pos/tendermint/consensus"
	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"github.com/implement-pbft-pos/tendermint/node"
	"github.com/implement-pbft-pos/tendermint/utils"
	config "gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	cfg, err := config.Load("config.ini")
	if err != nil {
		utils.ErrorStdOutLogger.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}
	// load node info
	cfgSectionNode := "Node"
	nodeId := int32(cfg.Section(cfgSectionNode).Key("id").MustInt(0))
	nodePort := uint16(cfg.Section(cfgSectionNode).Key("port").MustUint(9000))
	nodePubKeyPath := cfg.Section(cfgSectionNode).Key("pubKeyFile").String()
	nodePriKeyPath := cfg.Section(cfgSectionNode).Key("priKeyFile").String()

	pubKey, err := ioutil.ReadFile(nodePubKeyPath)
	if err != nil {
		utils.ErrorStdOutLogger.Printf("Read file failed: %v\n", err)
		os.Exit(1)
	}
	priKey, err := ioutil.ReadFile(nodePriKeyPath)
	if err != nil {
		utils.ErrorStdOutLogger.Printf("Read file failed: %v\n", err)
		os.Exit(1)
	}
	keyPair := utils.KeyPair{
		PublicKey:  pubKey,
		PrivateKey: priKey,
	}
	fmt.Printf("pubKey: %x, privKey: %x\n", pubKey, priKey)

	// load neighbor info
	cfgSectionNeighbors := "Neighbor_Nodes"
	neighborPubFolderPath := cfg.Section(cfgSectionNeighbors).Key("neighbor_pub_folder").String()
	neighborIds := cfg.Section(cfgSectionNeighbors).Key("ids").ValidInts(",")
	numNodes := len(neighborIds) + 1
	nodeIds := make([]int, numNodes)
	nodeIds[0] = int(nodeId)
	neighborNodes := make([]*consensus.NeighborNode, numNodes-1)
	for idx, neighborId := range neighborIds {
		neighborName := fmt.Sprintf("%s.%d", cfgSectionNeighbors, neighborId)
		neighborHost := cfg.Section(neighborName).Key("host").String()
		neighborPort := cfg.Section(neighborName).Key("port").MustInt(9000)
		neighborTimeout := cfg.Section(neighborName).Key("timeout").MustUint64(500)
		neighborPubPath := fmt.Sprintf("%s/node%d_pub", neighborPubFolderPath, neighborId)
		neighborPubKey, err := ioutil.ReadFile(neighborPubPath)
		if err != nil {
			utils.ErrorStdOutLogger.Printf("Read file failed: %v\n", err)
			os.Exit(1)
		}
		gClient := core_rpc.GClient{ClientConfig: core_rpc.GClientConfig{Host: neighborHost, Port: uint16(neighborPort), Timeout: neighborTimeout}}
		neighborNode := consensus.NeighborNode{NodeId: int32(neighborId), PublicKey: neighborPubKey, Client: &gClient}
		neighborNodes[idx] = &neighborNode
		nodeIds[idx+1] = neighborId
	}

	// load consensus properties
	cfgSectionConsensus := "Consensus"
	timeoutPropose := cfg.Section(cfgSectionConsensus).Key("timeout_propose").MustUint(1000)
	timeoutPreVote := cfg.Section(cfgSectionConsensus).Key("timeout_pre_vote").MustUint(1000)
	timeoutPreCommit := cfg.Section(cfgSectionConsensus).Key("timeout_pre_commit").MustUint(1000)
	consensusCfg := consensus.ConsensusConfig{
		TimeoutPropose:   time.Duration(timeoutPropose) * time.Millisecond,
		TimeoutPreVote:   time.Duration(timeoutPreVote) * time.Millisecond,
		TimeoutPreCommit: time.Duration(timeoutPreCommit) * time.Millisecond,
	}

	sort.Ints(nodeIds)
	allNodeIds := make([]int32, len(nodeIds))
	for idx,nodeId := range nodeIds {
		allNodeIds[idx] = int32(nodeId)
	}
	proposerSelector := consensus.ProposerSelector{NodeIds: allNodeIds}
	node := node.CompliantNode{}
	node.InitNode(nodeId, nodePort, &keyPair, neighborNodes, &proposerSelector, &consensusCfg)

	// start node
	var wg sync.WaitGroup
	node.StartServer(&wg)

	// connect to neighbor nodes
	time.Sleep(2 * time.Minute)
	node.ConnectNeighborNodes()

	// start consensus
	node.StartConsensus()

	wg.Wait()

	//folderPath := "./deploy/node_5/"
	//keyPair, err := utils.GenerateNewKeyPairAndDump(folderPath)
	////keyPair, err := utils.GenerateNewKeyPair()
	//if err != nil {
	//	utils.ErrorStdOutLogger.Printf("Failed to generate keyPair: %v\n", err)
	//	os.Exit(1)
	//}
	//fmt.Printf("pubKey: %x, privKey: %x\n", keyPair.PublicKey, keyPair.PrivateKey)

}
