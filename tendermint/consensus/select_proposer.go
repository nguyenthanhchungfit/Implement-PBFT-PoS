package consensus

type ProposerSelector struct {
	NodeIds []int32
	count int
}

func (selector *ProposerSelector) SelectProposer(height, round int32) int32{
	idx := selector.count % len(selector.NodeIds)
	selector.count++
	return selector.NodeIds[idx]
}