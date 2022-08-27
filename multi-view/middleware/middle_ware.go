package middleware

import (
	"github.com/implement-pbft-pos/multi-view/consensus"
)

type ConsensusMiddleware struct {
	Processor* consensus.CompliantProcessor
}

func (middleware ConsensusMiddleware) ReceiveProposeMessage(height, round, validRound int32, data string, nodeId int32) {
	middleware.Processor.ReceiveProposeMessage(height, round, validRound, data, nodeId)
}

func (middleware ConsensusMiddleware) ReceiveVoteMessage(height, round int32, hashValue, signature []byte, nodeId int32) {
	middleware.Processor.ReceiveVoteMessage(height, round, hashValue, signature, nodeId)
}
