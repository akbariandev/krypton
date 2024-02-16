package core

import (
	"errors"
	"math/rand"
)

type CommitteeSelectionMode int

const CommitteeRandomSelection CommitteeSelectionMode = 0 // 0: Random

type Committee struct {
	members []*Node
}

func (c *Committee) GetMembers() []*Node {
	return c.members
}

func (c *Committee) SelectCommittee(nodes []*Node, mode CommitteeSelectionMode, size int) error {
	switch mode {
	case CommitteeRandomSelection:
		c.selectRandomCommittee(nodes, size)
	default:
		return errors.New("invalid committee selection mode")
	}

	return nil
}

func (c *Committee) selectRandomCommittee(nodes []*Node, size int) {
	c.members = make([]*Node, 0, size)

	rand.Shuffle(len(nodes), func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })

	for i := 0; i < min(len(nodes), size); i++ {
		c.members = append(c.members, nodes[i])
	}
}
