package core

import (
	"context"
	"fmt"
	"github.com/akbariandev/krypton/internal/p2p"
	"github.com/libp2p/go-libp2p-pubsub"
	"math/rand"
	"strconv"
	"time"
)

type Node struct {
	ID         string
	PeerStream *p2p.PeerStream
	PubSub     *pubsub.PubSub
	Topics     map[p2p.TopicName]*pubsub.Topic
}

func NewNode(ctx context.Context, ID string, networkGroupName string, listenPort int) (*Node, error) {
	ps, err := p2p.NewPeerStream(listenPort)
	if err != nil {
		return nil, err
	}

	pSub, err := pubsub.NewGossipSub(ctx, ps.Host)
	if err != nil {
		return nil, err
	}

	ps.Run(ctx, networkGroupName)
	return &Node{
		ID:         ID,
		PeerStream: ps,
		PubSub:     pSub,
		Topics:     make(map[p2p.TopicName]*pubsub.Topic, 0),
	}, nil
}

func (n *Node) GetID() string {
	return n.ID
}

func (n *Node) GetPeerStream() *p2p.PeerStream {
	return n.PeerStream
}

func (n *Node) Run(ctx context.Context) {

	taskCh := make(chan struct{})
	taskID := 1
	go func(c chan struct{}) {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				c <- struct{}{}
			}
		}
	}(taskCh)

	for {
		select {
		case <-taskCh:
			go runTasks(ctx, n, taskID)
			taskID++
		}
	}
}

func runTasks(ctx context.Context, node *Node, taskID int) {

	committeeValidatorSeed := int64(0)
	randomSeedCh := make(chan struct{})
	seedTopicName := p2p.TopicNameGenerator(p2p.TopicCommitteeValidatorSeed, taskID)
	if err := node.JoinListenTopic(ctx, seedTopicName, handleCommitteeValidatorSeed, &committeeValidatorSeed); err != nil {
		fmt.Println(err)
		return
	}
	if err := p2p.PublishOnTopic(ctx, node.Topics[seedTopicName], fmt.Sprintf("%d", rand.Intn(1000))); err != nil {
		fmt.Println(err)
	}
	go func(c chan struct{}) {
		ticker := time.NewTicker(3 * time.Second)
		select {
		case <-ticker.C:
			c <- struct{}{}
		}
	}(randomSeedCh)
	go func(c chan struct{}, nodeID string, taskID int, seed *int64) {
		for {
			select {
			case <-c:
				validatedNode := []*Node{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}}
				rnd := rand.New(rand.NewSource(*seed))
				randomValidator := rnd.Intn(len(validatedNode))
				fmt.Println(fmt.Sprintf("Task: %d, Node: %s => %s", taskID, nodeID, validatedNode[randomValidator].ID))
			}
		}
	}(randomSeedCh, node.ID, taskID, &committeeValidatorSeed)

}

func (n *Node) JoinListenTopic(ctx context.Context, topicName p2p.TopicName, handleMessage func(*pubsub.Message, ...interface{}), params ...interface{}) error {
	topic, err := n.PubSub.Join(string(topicName))
	if err != nil {
		return err
	}

	n.Topics[topicName] = topic

	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Next(ctx)
			if err != nil {
				break
			}

			handleMessage(msg, params...)
		}
	}()

	return nil
}

func handleCommitteeValidatorSeed(msg *pubsub.Message, params ...interface{}) {
	if len(params) == 0 || len(params) > 1 {
		fmt.Println("params is not implemented")
		return
	}

	seedAddr := params[0].(*int64)
	randomSeed, _ := strconv.Atoi(string(msg.GetData()))
	*(seedAddr) += int64(randomSeed)
}
