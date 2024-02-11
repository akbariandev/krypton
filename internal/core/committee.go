package core

import "math/rand"

type CommitteeSelectionMode int

const CommitteeRandomSelection CommitteeSelectionMode = 0 // 0: Random

type Committee struct {
	members []*Node
}

func InitCommittee(nodes []*Node) *Committee {
	return &Committee{
		members: nodes,
	}
}

func (c *Committee) SelectCommittee(mode CommitteeSelectionMode, size int) []*Node {
	switch mode {
	case CommitteeRandomSelection:
		return c.selectRandomCommittee(size)
	default:
		return c.selectRandomCommittee(size)
	}
}

func (c *Committee) selectRandomCommittee(size int) []*Node {
	members := make([]*Node, 0, size) // Initialize members with size 3

	rand.Shuffle(len(c.members), func(i, j int) { c.members[i], c.members[j] = c.members[j], c.members[i] })

	for i := 0; i < min(len(c.members), size); i++ {
		members = append(members, c.members[i])
	}

	return members
}
