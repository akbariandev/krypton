package core

import (
	"context"
	"github.com/akbariandev/krypton/internal/p2p"
)

type Role struct {
	Name string
}

type Node struct {
	ID         string
	Roles      []Role
	PeerStream *p2p.PeerStream
	IsPrimary  bool
}

func NewNode(ctx context.Context, ID string, networkGroupName string, listenPort int, isPrimary bool) (*Node, error) {
	ps, err := p2p.NewPeerStream(listenPort)
	if err != nil {
		return nil, err
	}

	ps.Run(ctx, networkGroupName)
	return &Node{
		ID:         ID,
		PeerStream: ps,
		IsPrimary:  isPrimary,
	}, nil
}

func (n *Node) AddRole(role Role) {
	n.Roles = append(n.Roles, role)
}

func (n *Node) AddRoles(roles []Role) {
	n.Roles = append(n.Roles, roles...)
}

func (n *Node) GetRoles() []Role {
	return n.Roles
}

func (n *Node) GetID() string {
	return n.ID
}
