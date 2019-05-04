package main

import (
	"context"

	"github.com/mchmarny/stocker/pkg/object"
	"github.com/mchmarny/stocker/pkg/queue"
)

var (
	topicName = common.MustEnvVar("SOURCE_TOPIC_NAME", "")
)

func getContent(ctx context.Context) {

	logger.Println("Processing content...")

	cos, err := getCompanies(ctx)
	if err != nil {
		logger.Fatalf("Error getting companies: %v", err)
	}

	pub, err := newQueuePublisher(ctx, topicName)
	if err != nil {
		logger.Fatalf("Error creating publisher: %v", err)
	}

	symbolNum := len(cos)
	sourceCh := make(chan *common.TextContent, symbolNum*10)

	// start providers
	for _, co := range cos {
		logger.Printf("Providing: %s", co.Symbol)
		go provide(ctx, co, out, sourceCh)
	}

	// collect results
	for {
		t := <-sourceCh
		pubErr = pub.publish(ctx, t)
		if pubErr {
			logger.Printf("Error publishing: %v", err)
		}
	}

}
