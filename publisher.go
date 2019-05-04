package main

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/mchmarny/stocker/pkg/object"
	"github.com/mchmarny/stocker/pkg/util"
)

var (
	projectID = util.MustEnvVar("GCP_PROJECT", "")
)

type queuePublisher struct {
	topic *pubsub.Topic
}

// Publish provides generic publish capability
func (p *queuePublisher) publish(ctx context.Context, content *object.TextContent) error {

	c, err := json.Marshal(content)
	if err != nil {
		return err
	}

	result := p.topic.Publish(ctx, &pubsub.Message{
		Data: c,
	})

	// block and wait for the result
	_, err = result.Get(ctx)
	return err

}

func newQueuePublisher(ctx context.Context, topicName string) (q *queuePublisher, err error) {

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		logger.Printf("Error creating pubsub client: %v", err)
		return nil, err
	}

	// Creates the new topic.
	t, err := client.CreateTopic(ctx, topicName)
	if err != nil {
		logger.Printf("Error creating pubsub topic: %v", err)
		return nil, err
	}

	queue := &QueuePublisher{
		topic: t,
	}

	return queue, nil

}
