package consensus

import "time"

type ConsensusConfig struct {
	TimeoutPropose time.Duration
	TimeoutPreVote time.Duration
	TimeoutPreCommit time.Duration
}

const (
	STEP_PROPOSE int8 = 1
	STEP_PRE_VOTE		= 2
	STEP_PRE_COMMIT		= 3
)