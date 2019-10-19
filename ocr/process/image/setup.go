package image

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	vision "cloud.google.com/go/vision/apiv1"
)

var global = struct {
	VisionClient *vision.ImageAnnotatorClient
	PubsubClient *pubsub.Client
}{}

func init() {
	log.SetFlags(0)
}

func setup(ctx context.Context) error {
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		return fmt.Errorf("GCP_PROJECT environment variable is missing")
	}

	topicName := os.Getenv("TEXT_TOPIC")
	if topicName == "" {
		return fmt.Errorf("TEXT_TOPIC environment variable is missing")
	}

	var err error

	if global.VisionClient == nil {
		global.VisionClient, err = vision.NewImageAnnotatorClient(ctx)
		if err != nil {
			return fmt.Errorf("create image annotator client: %v", err)
		}
	}

	if global.PubsubClient == nil {
		global.PubsubClient, err = pubsub.NewClient(ctx, projectID)
		if err != nil {
			return fmt.Errorf("create pubsub client: %v", err)
		}
	}

	return nil
}
