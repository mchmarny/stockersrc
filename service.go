package main

import (
	"context"
)

var (
	topicName = mustEnvVar("SOURCE_TOPIC_NAME", "stocker-source")
)

func getContent(ctx context.Context) {

	logger.Println("Processing content...")

	cos, err := getCompanies(ctx)
	if err != nil {
		logger.Fatalf("Error getting companies: %v", err)
	}

	if len(cos) == 0 {
		logger.Println("No companies to process")
		return
	}

	pub, err := newQueuePublisher(ctx, topicName)
	if err != nil {
		logger.Fatalf("Error creating publisher: %v", err)
	}

	symbolNum := len(cos)
	sourceCh := make(chan *TextContent, symbolNum*10)

	// start providers
	for _, co := range cos {
		logger.Printf("Providing: %s", co.Symbol)
		go provide(ctx, co, sourceCh)
	}

	// collect results
	for {
		t := <-sourceCh
		pubErr := pub.publish(ctx, t)
		if pubErr != nil {
			logger.Printf("Error publishing: %v", pubErr)
		}
	}

}
