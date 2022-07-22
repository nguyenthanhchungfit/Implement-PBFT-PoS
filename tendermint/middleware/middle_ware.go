package middleware

import (
	"github.com/implement-pbft-pos/tendermint/consensus"
)

type ConsensusMiddleware struct {
	Processor* consensus.CompliantProcessor
}

func (middleware ConsensusMiddleware) ReceiveProposeMessage(height, round, validRound int32, data string, nodeId int32) {
	middleware.Processor.ReceiveProposeMessage(height, round, validRound, data, nodeId)
}

func (middleware ConsensusMiddleware) ReceivePreVoteMessage(height, round int32, hashValue, signature []byte, nodeId int32) {
	middleware.Processor.ReceivePreVoteMessage(height, round, hashValue, signature, nodeId)
}

func (middleware ConsensusMiddleware) ReceivePreCommitMessage(height, round int32, hashValue, signature []byte, nodeId int32) {
	middleware.Processor.ReceivePreCommitMessage(height, round, hashValue, signature, nodeId)
}