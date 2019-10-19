package image

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func publish(ctx context.Context, dt DetectText) error {
	if err := setup(ctx); err != nil {
		return fmt.Errorf("setup: %v", err)
	}

	data, err := json.Marshal(dt)
	if err != nil {
		return fmt.Errorf("json marshal: %v", err)
	}

	textTopic := os.Getenv("TEXT_TOPIC")
	topic := global.PubsubClient.Topic(textTopic)

	ok, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("topic %q, error: %v", textTopic, err)
	}
	if !ok {
		return fmt.Errorf("topic %q does not exists", textTopic)
	}

	// Publish a batch when it has this many messages
	topic.PublishSettings.CountThreshold = 1

	message := pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"object": dt.object,
		},
	}

	// publish messages synchronously
	id, err := topic.Publish(ctx, &message).Get(ctx)
	if err != nil {
		return fmt.Errorf("topic %q, publish error: %v", textTopic, err)
	}

	log.Printf("published message ID: %q", id)
	return nil
}
