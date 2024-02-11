package core

import "math/rand"

type CommitteeSelectionMode int

const CommitteeRandomSelection CommitteeSelectionMode = 0 // 0: Random

type Committee struct {
	members []*Node
}

func (c *Committee) GetMembers() []*Node {
	return c.members
}

func (c *Committee) SelectCommittee(nodes []*Node, mode CommitteeSelectionMode, size int) {
	switch mode {
	case CommitteeRandomSelection:
		c.selectRandomCommittee(nodes, size)
	default:
		c.selectRandomCommittee(nodes, size)
	}
}

func (c *Committee) selectRandomCommittee(nodes []*Node, size int) {
	c.members = make([]*Node, 0, size) // Initialize members with size 3

	rand.Shuffle(len(nodes), func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })

	for i := 0; i < min(len(nodes), size); i++ {
		c.members = append(c.members, nodes[i])
	}
}
