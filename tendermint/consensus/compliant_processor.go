package consensus

import (
	"bytes"
	"fmt"
	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"github.com/implement-pbft-pos/tendermint/utils"
	"sync"
	"time"
)

const MaxHeight = 1000
const DefaultValue = -1
const ValidValue = 1
const NilValue = 0


type CompliantProcessor struct {
	Height             int32
	Round              int32
	Step               int8
	Decisions          []int32
	NeighborNodes      []*NeighborNode
	NeighborNodeIds    []int32
	Selector           *ProposerSelector
	ConsensusConfig    *ConsensusConfig
	KeyPair            *utils.KeyPair
	NodeId             int32
	TotalNodes         int32
	lockedValue        string
	lockedRound        int32
	tmpValue           string
	validValue         string
	validRound         int32
	mapPreVoteValues   map[int32][]byte
	mapPreCommitValues map[int32][]byte
	defaultVoteFlag	  []byte
	mu                 sync.Mutex
}

func (processor *CompliantProcessor) StartConsensus() {
	processor.initConsensus()
	round := int32(0)
	for {
		processor.StartRound(round)
		proposerId := processor.Selector.SelectProposer(processor.Height, processor.Round)
		//fmt.Printf("Select proposerId: %d\n", proposerId)
		if proposerId != processor.NodeId {
			// schedule
			cScheduleWaitingPropose := make(chan int, 1)
			go processor.onWaitingProposeMsg(cScheduleWaitingPropose)
			select {
			case validValue := <-cScheduleWaitingPropose:
				if validValue == NilValue { // invalid value
					processor.broadcastPreVoteMsg(processor.Height, processor.Height, nil)
				} else if validValue == ValidValue { // valid value
					hashData := utils.Hash(processor.validValue)
					processor.broadcastPreVoteMsg(processor.Height, processor.Height, hashData)
				}
				processor.Step = STEP_PRE_VOTE
			case <-time.After(processor.ConsensusConfig.TimeoutPropose): // timeout propose
				fmt.Printf("******************** 1. Node %d TimeoutPropose\n", processor.NodeId)
				processor.onTimeoutPropose(processor.Height, processor.Round)
			}
		}

		// set timeout PreVote step
		cScheduleWaitingPreVote := make(chan int, 1)
		go processor.onWaitingPreVoteMsg(cScheduleWaitingPreVote)
		select {
		case validValue := <-cScheduleWaitingPreVote:
			if validValue == NilValue {
				processor.broadcastPreCommitMsg(processor.Height, processor.Round, nil)
				processor.Step = STEP_PRE_COMMIT
			} else if validValue == ValidValue {
				if processor.Step == STEP_PRE_VOTE {
					processor.lockedValue = processor.tmpValue
					processor.lockedRound = processor.Round
					processor.broadcastPreCommitMsg(processor.Height, processor.Round, utils.Hash(processor.tmpValue))
					processor.Step = STEP_PRE_COMMIT
				}
				processor.validValue = processor.tmpValue
				processor.validRound = processor.Round

			}
		case <-time.After(processor.ConsensusConfig.TimeoutPreVote):
			fmt.Printf("******************** 2. Node %d Timeout Prevote\n", processor.NodeId)
			processor.onTimeoutPreVote(processor.Height, processor.Round)
		}

		// Set timeout PreCommit
		cScheduleWaitingPreCommit := make(chan int, 1)
		go processor.onWaitingPreCommitMsg(cScheduleWaitingPreCommit)
		select {
		case validValue := <-cScheduleWaitingPreCommit:
			if validValue == NilValue {
				round++
			} else if validValue == ValidValue {
				round = 0
				processor.commitValue(ValidValue)
			}
		case <-time.After(processor.ConsensusConfig.TimeoutPreCommit):
			fmt.Printf("******************** 3. Node %d Timeout PreCommit\n", processor.NodeId)
			round++
		}

	}
}

func (processor *CompliantProcessor) initConsensus() {
	processor.Height = 0
	processor.Round = 0
	processor.Step = 0
	processor.Decisions = make([]int32, MaxHeight)
	for idx := range processor.Decisions {
		processor.Decisions[idx] = DefaultValue
	}
	processor.lockedValue = ""
	processor.lockedRound = -1
	processor.tmpValue = ""
	processor.validValue = ""
	processor.validRound = -1
	processor.mapPreVoteValues = make(map[int32][]byte)
	processor.mapPreCommitValues = make(map[int32][]byte)
	processor.defaultVoteFlag = []byte{1}
	for _, neighborNodeId := range processor.NeighborNodeIds {
		processor.mapPreVoteValues[neighborNodeId] = processor.defaultVoteFlag
		processor.mapPreCommitValues[neighborNodeId] = processor.defaultVoteFlag
	}
}

func (processor *CompliantProcessor) dumpState() {
	utils.InfoStdOutLogger.Printf("Node %d=(ValidRound=%d, validValue=%s, lockedRound=%d, lockedValue=%s)\n", processor.NodeId,
		processor.validRound, processor.validValue, processor.lockedRound, processor.lockedValue)
}

func (processor *CompliantProcessor) getProposeValue() string {
	return fmt.Sprintf("Value[Node%d-Height%d]", processor.NodeId, processor.Height)
}

func (processor *CompliantProcessor) isValidValue(value string) bool {
	return value != ""
}

func (processor *CompliantProcessor) StartRound(round int32) {
	utils.InfoStdOutLogger.Printf("[New Round]Node %d start round %d at height %d", processor.NodeId, round, processor.Height)
	processor.Round = round
	processor.Step = STEP_PROPOSE
	if processor.Selector.SelectProposer(processor.Height, round) == processor.NodeId {
		var proposalValue = ""
		if processor.validValue != "" {
			proposalValue = processor.validValue
		} else {
			proposalValue = processor.getProposeValue()
		}
		processor.validValue = proposalValue
		processor.broadcastProposeMsg(processor.Height, processor.Round, processor.validRound, proposalValue)
	}
}

func (processor *CompliantProcessor) commitValue(value int32) {
	utils.InfoStdOutLogger.Printf("[Commit] Node %d commit value: %d at height: %d\n", processor.NodeId, value, processor.Height)
	if processor.Decisions[processor.Height] == DefaultValue {
		processor.Decisions[processor.Height] = value
		processor.Height++
		processor.resetState()
	}
}

func (processor *CompliantProcessor) resetState() {
	processor.lockedRound = -1
	processor.lockedValue = ""
	processor.validRound = -1
	processor.validValue = ""
	processor.tmpValue = ""
	processor.mu.Lock()
	for _, neighborNodeId := range processor.NeighborNodeIds {
		processor.mapPreVoteValues[neighborNodeId] = processor.defaultVoteFlag
		processor.mapPreCommitValues[neighborNodeId] = processor.defaultVoteFlag
	}
	processor.mu.Unlock()
}

func (processor *CompliantProcessor) getNumPreVotes(hashValue []byte) (int32, int32) {
	totValidVotes := 0
	totNilVotes := 0
	processor.mu.Lock()
	for _, flag := range processor.mapPreVoteValues {
		if hashValue != nil{
			if bytes.Compare(hashValue, flag) == 0 {
				totValidVotes++
			}
		} else {
			totNilVotes++
		}
	}
	processor.mu.Unlock()
	return int32(totValidVotes), int32(totNilVotes)
}

func (processor *CompliantProcessor) getNumPreCommits(hashValue []byte) (int32, int32) {
	totValidVotes := 0
	totNilVotes := 0
	processor.mu.Lock()
	for _, flag := range processor.mapPreCommitValues {
		if hashValue != nil{
			if  bytes.Compare(hashValue, flag) == 0 {
				totValidVotes++
			}
		} else {
			totNilVotes++
		}
	}
	processor.mu.Unlock()
	return int32(totValidVotes), int32(totNilVotes)
}

func (processor *CompliantProcessor) onWaitingProposeMsg(done chan<- int) {
	for {
		if processor.Step == STEP_PRE_VOTE {
			//fmt.Printf("Node %d onWaitingProposeMsg\n", processor.NodeId)
			if processor.isValidValue(processor.tmpValue) {
				done <- ValidValue
			} else {
				done <- NilValue
			}

		}
	}
}

func (processor *CompliantProcessor) onWaitingPreVoteMsg(done chan<- int) {
	for {
		if processor.Step == STEP_PRE_VOTE && processor.tmpValue != "" {
			hashValue := utils.Hash(processor.tmpValue)
			totValidVotes, totNilVotes := processor.getNumPreVotes(hashValue)

			//if totValidVotes > 0 || totNilVotes > 0 {
			//	utils.InfoStdOutLogger.Printf("[Check Node %d] totalValidVotes: %d, totNilVotes: %d\n", processor.NodeId, totValidVotes, totNilVotes)
			//}

			if totValidVotes*3 > 2*processor.TotalNodes {
				done <- ValidValue
			} else if totNilVotes*3 > 2*processor.TotalNodes {
				done <- NilValue
			}

		}
	}
}

func (processor *CompliantProcessor) onWaitingPreCommitMsg(done chan<- int) {
	for {
		if processor.Step == STEP_PRE_COMMIT {
			hashValue := utils.Hash(processor.tmpValue)
			totValidVotes, totNilVotes := processor.getNumPreCommits(hashValue)
			if totValidVotes*3 > 2*processor.TotalNodes {
				done <- ValidValue
			} else if totNilVotes*3 > 2*processor.TotalNodes {
				done <- NilValue
			}
		}
	}
}

func (processor *CompliantProcessor) onTimeoutPropose(height, round int32) {
	if processor.Height == height && processor.Round == round && processor.Step == STEP_PROPOSE {
		processor.broadcastPreVoteMsg(height, round, nil)
		processor.Step = STEP_PRE_VOTE
	}
}

func (processor *CompliantProcessor) onTimeoutPreVote(height, round int32) {
	if processor.Height == height && processor.Round == round && processor.Step == STEP_PRE_VOTE {
		processor.broadcastPreCommitMsg(height, round, nil)
		processor.Step = STEP_PRE_COMMIT
	}
}

func (processor *CompliantProcessor) onTimeoutPreCommit(height, round int32) {
	if processor.Height == height && processor.Round == round {
		processor.StartRound(round + 1)
	}
}

func (processor *CompliantProcessor) broadcastProposeMsg(height, round, validRound int32, data string) {
	if len(data) > 0 {
		utils.InfoStdOutLogger.Printf("[Propose] Node %d broadcastPropose value %s\n", processor.NodeId, data)
		gdata := core_rpc.GData{Data: data}
		msg := core_rpc.GProposeMessage{Height: height, Round: round, ValidRound: validRound, Data: &gdata, NodeId: processor.NodeId}
		for _, neighbor := range processor.NeighborNodes {
			neighbor.Client.SendProposeMessage(&msg)
		}
	} else {
		byteData := []byte(data)
		signature := utils.SignData(processor.KeyPair.PrivateKey, byteData)
		gData := core_rpc.GData{Data: data}
		msg := core_rpc.GProposeMessage{Height: height, Round: round, ValidRound: validRound, Data: &gData, Signature: signature}
		for _, neighbor := range processor.NeighborNodes {
			neighbor.Client.SendProposeMessage(&msg)
		}
	}
}

func (processor *CompliantProcessor) broadcastPreVoteMsg(height, round int32, hashValue []byte) {
	var signature []byte
	if hashValue != nil {
		signature = utils.SignData(processor.KeyPair.PrivateKey, hashValue)
	}
	//utils.InfoStdOutLogger.Printf("Node %d broadcastPreVoteMsg value %x\n", processor.NodeId, hashValue)
	msg := core_rpc.GPreVoteMessage{Height: height, Round: round, HashValue: hashValue, Signature: signature, NodeId: processor.NodeId}
	for _, neighbor := range processor.NeighborNodes {
		neighbor.Client.SendPreVoteMessage(&msg)
	}
}

func (processor *CompliantProcessor) broadcastPreCommitMsg(height, round int32, hashValue []byte) {
	var signature []byte
	if hashValue != nil {
		signature = utils.SignData(processor.KeyPair.PrivateKey, hashValue)
	}
	msg := core_rpc.GPreCommitMessage{Height: height, Round: round, HashValue: hashValue, Signature: signature, NodeId: processor.NodeId}
	for _, neighbor := range processor.NeighborNodes {
		neighbor.Client.SendPreCommitMessage(&msg)
	}
}

func (processor *CompliantProcessor) ReceiveProposeMessage(height, round, validRound int32, data string, nodeId int32) {
	//utils.InfoStdOutLogger.Printf("[Node %d ReceiveProposeMessage from %d] height: %d, round: %d, validRound: %d, data: %s\n", processor.NodeId, nodeId, height, round, validRound, data)
	if processor.Step == STEP_PROPOSE && nodeId == processor.Selector.SelectProposer(processor.Height, processor.Round) {
		if validRound == -1 { // initial round of new height
			if processor.isValidValue(data) && (processor.lockedRound == -1 || processor.lockedValue == data) {
				processor.tmpValue = data
				hashValue := utils.Hash(data)
				processor.broadcastPreVoteMsg(height, round, hashValue)
			} else {
				processor.broadcastPreVoteMsg(height, round, nil)
			}
			processor.Step = STEP_PRE_VOTE
		} else {
			if validRound >= 0 || validRound < round {
				if processor.isValidValue(data) && (processor.lockedRound <= validRound || processor.lockedValue == data) {
					processor.tmpValue = data
					hashValue := utils.Hash(data)
					processor.broadcastPreVoteMsg(height, round, hashValue)
				} else {
					processor.broadcastPreVoteMsg(height, round, nil)
				}
				processor.Step = STEP_PRE_VOTE
			}
		}
	}
}

func (processor *CompliantProcessor) ReceivePreVoteMessage(height, round int32, hashValue, signature []byte, nodeId int32) {
	//utils.InfoStdOutLogger.Printf("[Node %d ReceivePreVoteMessage from %d] height: %d, round: %d, hashValue: %x\n",
	//	processor.NodeId, nodeId, height, round, hashValue)

	//if processor.Step >= STEP_PRE_VOTE {
	if processor.Height == height && processor.Round == round {
		curVal := processor.mapPreVoteValues[nodeId]
		//utils.InfoStdOutLogger.Printf("Test Node%d receive node %d curVal: %x\n", processor.NodeId, nodeId, curVal)
		if bytes.Compare(curVal, processor.defaultVoteFlag) == 0  {
			processor.mu.Lock()
			processor.mapPreVoteValues[nodeId] = hashValue
			processor.mu.Unlock()
		}
	}
	//}
}

func (processor *CompliantProcessor) ReceivePreCommitMessage(height, round int32, hashValue, signature []byte, nodeId int32) {
	//utils.InfoStdOutLogger.Printf("[ReceivePreCommitMessage from %d] height: %d, round: %d, hashValue: %x\n", nodeId, height, round, hashValue)

	//if processor.Step >= STEP_PRE_VOTE {
	if processor.Height == height && processor.Round == round {
		curVal := processor.mapPreCommitValues[nodeId]
		if bytes.Compare(curVal, processor.defaultVoteFlag) == 0 {
			processor.mu.Lock()
			processor.mapPreCommitValues[nodeId] = hashValue
			processor.mu.Unlock()
		}

	}
	//}
}

func (processor *CompliantProcessor) TestPing() {
	gdata := core_rpc.GData{Data: "test"}
	msg := core_rpc.GProposeMessage{Height: 0, Round: 0, ValidRound: -1, Data: &gdata, NodeId: processor.NodeId}
	fmt.Println("neighbor: ", processor.NeighborNodes)
	for _, neighbor := range processor.NeighborNodes {
		ret := neighbor.Client.SendProposeMessage(&msg)
		fmt.Println("Ret: ", ret)
	}
}
