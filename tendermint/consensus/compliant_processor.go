package consensus

import (
	core_rpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"github.com/implement-pbft-pos/tendermint/utils"
)

type CompliantProcessor struct {
	Height uint32
	Round int32
	Step uint8
	Decisions[] int32
	NeighborClients [] core_rpc.GClient
	Selector* ProposerSelector
	NodeId int
	TotalNodes int32
	lockedValue string
	lockedRound int32
	validValue string
	validRound int32
}

func (processor *CompliantProcessor) StartRound(round int32) {
	utils.InfoStdOutLogger.Printf("Node %d start consensus round %d", processor.NodeId, round)
	processor.Round = round
	processor.Step = STEP_PROPOSE

}

func (processor CompliantProcessor) ReceiveProposeMessage(height, round, validRound int32, data string) {
	utils.InfoStdOutLogger.Printf("[ReceiveProposeMessage] height: %d, round: %d, validRound: %d, data: %s\n", height, round, validRound, data)
}
func (processor CompliantProcessor) ReceivePreVoteMessage(height, round int32, hashValue []byte){
	utils.InfoStdOutLogger.Printf("[ReceivePreVoteMessage] height: %d, round: %d, hashValue: %x\n", height, round, hashValue)
}
func (processor CompliantProcessor) ReceivePreCommitMessage(height, round int32, hashValue []byte){
	utils.InfoStdOutLogger.Printf("[ReceivePreCommitMessage] height: %d, round: %d, hashValue: %x\n", height, round, hashValue)
}