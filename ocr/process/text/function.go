package text

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/functions/metadata"
)

type PubSubMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

func ProcessText(ctx context.Context, message PubSubMessage) error {
	// metadata
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		log.Printf("reading context: %v", err)
		return nil
	}

	// data
	if message.Data == nil {
		log.Printf("message %q, data is nil", meta.EventID)
		return nil
	}
	if len(message.Data) == 0 {
		log.Printf("message %q, data is empty", meta.EventID)
		return nil
	}

	// attributes
	if message.Attributes == nil {
		log.Printf("message %q, attributes is nil", meta.EventID)
		return nil
	}
	object := message.Attributes["object"]
	if object == "" {
		log.Printf("message %q, object is empty", meta.EventID)
		return nil
	}

	var dt DetectText
	if err := json.Unmarshal(message.Data, &dt); err != nil {
		log.Printf("message %q, unmarshal: %v", meta.EventID, err)
		return nil
	}

	if err := write(ctx, object, message.Data); err != nil {
		log.Printf("message %q, error: %v", meta.EventID, err)
		return nil
	}

	return nil
}

func write(ctx context.Context, path string, data []byte) error {
	if err := setup(ctx); err != nil {
		return fmt.Errorf("setup: %v", err)
	}

	// create context to cancel write operation in case of error
	writeCtx, cancelWrite := context.WithCancel(ctx)
	defer cancelWrite()

	split := strings.SplitN(path, "/", 2)
	bucket, name := split[0], split[1]

	// get an object writer
	object := global.StorageClient.Bucket(bucket).Object(name + ".json")
	writer := object.NewWriter(writeCtx)

	// set attributes
	writer.ContentType = "application/json"

	// do the write operation
	if _, err := writer.Write(data); err != nil {
		cancelWrite()
		return fmt.Errorf("could not write: %v", err)
	}
	if err := writer.Close(); err != nil {
		cancelWrite()
		return fmt.Errorf("could not close: %v", err)
	}

	log.Printf("created: %s/%s", object.BucketName(), object.ObjectName())

	return nil
}
