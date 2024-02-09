package core

import "github.com/akbariandev/krypton/internal/p2p"

type Role struct {
	Name string
}

type Node struct {
	ID         string
	Roles      []Role
	PeerStream *p2p.PeerStream
}

func NewNode(ID string, listenPort int) (*Node, error) {
	ps, err := p2p.NewPeerStream(listenPort)
	if err != nil {
		return nil, err
	}

	return &Node{
		ID:         ID,
		PeerStream: ps,
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
