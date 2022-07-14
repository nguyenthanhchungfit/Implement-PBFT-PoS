package consensus

import (
	"fmt"
	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"github.com/implement-pbft-pos/tendermint/utils"
	"time"
)

type CompliantProcessor struct {
	Height          int32
	Round           int32
	Step            int8
	Decisions       []int32
	NeighborClients []core_rpc.GClient
	Selector        *ProposerSelector
	ConsensusConfig *ConsensusConfig
	KeyPair *utils.KeyPair
	NodeId          int32
	TotalNodes      int32
	lockedValue     string
	lockedRound     int32
	validValue      string
	validRound      int32
}

func (processor *CompliantProcessor) getProposeValue() string {
	return fmt.Sprintf("Value[Node%d-Height%d]", processor.NodeId, processor.Height)
}

func (processor *CompliantProcessor) isValidValue(value string) bool {
	return value != ""
}

func (processor *CompliantProcessor) StartRound(round int32) {
	utils.InfoStdOutLogger.Printf("Node %d start consensus round %d", processor.NodeId, round)
	processor.Round = round
	processor.Step = STEP_PROPOSE
	if processor.Selector.SelectProposer(processor.Height, round) == processor.NodeId {
		var proposalValue = ""
		if processor.validValue != "" {
			proposalValue = processor.validValue
		} else {
			proposalValue = processor.getProposeValue()
		}
		processor.broadcastProposeMsg(processor.Height, processor.Round, processor.validRound, proposalValue)
	} else {
		// schedule
		cScheduleWaiting := make(chan string, 1)
		go processor.onWaitingProposeMsg()
		select {
		case res := <-cScheduleWaiting:
			fmt.Println(res)
		case <-time.After(processor.ConsensusConfig.TimeoutPropose):
			processor.onTimeoutPropose(processor.Height, processor.Round)
		}
	}
}

func (processor *CompliantProcessor) onWaitingProposeMsg() {
	for {
		if processor.Step == STEP_PROPOSE {

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
	}
}

func (processor *CompliantProcessor) onTimeoutPreCommit(height, round int32) {
	if processor.Height == height && processor.Round == round {
		processor.StartRound(round + 1)
	}
}

func (processor *CompliantProcessor) broadcastProposeMsg(height, round, validRound int32, data string) {
	if len(data) > 0 {
		gdata := core_rpc.GData{Data: data}
		msg := core_rpc.GProposeMessage{Height: height, Round: round, ValidRound: validRound, Data: &gdata, NodeId: processor.NodeId}
		for _, neighbor := range processor.NeighborClients {
			neighbor.SendProposeMessage(&msg)
		}
	} else {
		msg := core_rpc.GProposeMessage{Height: height, Round: round, ValidRound: validRound, Data: nil}
		for _, neighbor := range processor.NeighborClients {
			neighbor.SendProposeMessage(&msg)
		}
	}
}

func (processor *CompliantProcessor) broadcastPreVoteMsg(height, round int32, hashValue []byte) {
	var signature []byte
	if hashValue != nil {
		signature = utils.SignData(processor.KeyPair.PrivateKey, hashValue)
	}
	msg := core_rpc.GPreVoteMessage{Height: height, Round: round, HashValue: hashValue, Signature: signature, NodeId: processor.NodeId}
	for _, neighbor := range processor.NeighborClients {
		neighbor.SendPreVoteMessage(&msg)
	}
}

func (processor *CompliantProcessor) broadcastPreCommitMsg(height, round int32, hashValue []byte) {
	var signature []byte
	if hashValue != nil {
		signature = utils.SignData(processor.KeyPair.PrivateKey, hashValue)
	}
	msg := core_rpc.GPreCommitMessage{Height: height, Round: round, HashValue: hashValue, Signature: signature, NodeId: processor.NodeId}
	for _, neighbor := range processor.NeighborClients {
		neighbor.SendPreCommitMessage(&msg)
	}
}

func (processor CompliantProcessor) ReceiveProposeMessage(height, round, validRound int32, data string, nodeId int32) {
	utils.InfoStdOutLogger.Printf("[ReceiveProposeMessage from %d] height: %d, round: %d, validRound: %d, data: %s\n", nodeId, height, round, validRound, data)
	if processor.Step == STEP_PROPOSE {
		if validRound == -1 { // initial round of new height

		} else {

		}
	}
}
func (processor CompliantProcessor) ReceivePreVoteMessage(height, round int32, hashValue, signature []byte, nodeId int32) {
	utils.InfoStdOutLogger.Printf("[ReceivePreVoteMessage from %d] height: %d, round: %d, hashValue: %x\n", nodeId, height, round, hashValue)
}
func (processor CompliantProcessor) ReceivePreCommitMessage(height, round int32, hashValue, signature []byte, nodeId int32) {
	utils.InfoStdOutLogger.Printf("[ReceivePreCommitMessage from %d] height: %d, round: %d, hashValue: %x\n", nodeId, height, round, hashValue)
}
