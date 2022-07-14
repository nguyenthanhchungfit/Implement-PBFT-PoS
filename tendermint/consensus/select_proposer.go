package consensus

type ProposerSelector struct {
	NodeIds []int
	count int
}

func (selector *ProposerSelector) SelectProposer(height, round int32) int{
	idx := int32(selector.count % len(selector.NodeIds))
	selector.count++
	return selector.NodeIds[idx]
}