package consensus

type ConsensusConfig struct {
	TimeoutPropose uint32
	TimeoutPreVote uint32
	TimeoutPreCommit uint32
}

const (
	STEP_PROPOSE uint8 = 1
	STEP_PRE_VOTE		= 2
	STEP_PRE_COMMIT		= 3
)