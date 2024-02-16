package p2p

import (
	"context"
	"fmt"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type TopicName string

const TopicCommitteeValidatorSeed TopicName = "/task/%d/global/committee/validator/seed"

func TopicNameGenerator(topicName TopicName, params ...interface{}) TopicName {
	if len(params) > 0 {
		s := string(topicName)
		s = fmt.Sprintf(s, params...)
		topicName = TopicName(s)
	}

	return topicName
}

func PublishOnTopic(ctx context.Context, topic *pubsub.Topic, message string) error {
	err := topic.Publish(ctx, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
