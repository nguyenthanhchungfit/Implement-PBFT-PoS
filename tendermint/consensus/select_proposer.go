package consensus

import "sync"

type ProposerSelector struct {
	NodeIds   []int32
	count     int
	mu        sync.Mutex
	curHeight int32
	curRound  int32
}

func (selector *ProposerSelector) SelectProposer(height, round int32) int32 {
	idx := 0
	if !(selector.curHeight == height && selector.curRound == round) {
		selector.mu.Lock()
		if height > selector.curHeight || round > selector.curRound {
			selector.count++
			selector.curHeight = height
			selector.curRound = round
		}
		selector.mu.Unlock()
	}
	idx = selector.count % len(selector.NodeIds)
	return selector.NodeIds[idx]
}
