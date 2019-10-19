package image

import (
	"context"
	"log"
	"strings"
	"time"
)

type GCSEvent struct {
	ID            string    `json:"id"`
	Bucket        string    `json:"bucket"`
	Name          string    `json:"name"`
	ResourceState string    `json:"resourceState"`
	TimeCreated   time.Time `json:"timeCreated"`
	ContentType   string    `json:"contentType"`
	Size          string    `json:"size"`
}

func ProcessImage(ctx context.Context, event GCSEvent) error {
	if event.ResourceState == "not_exists" {
		log.Printf("event %q, does not exists", event.ID)
		return nil
	}

	if !strings.HasPrefix(event.ContentType, "image/") {
		log.Printf("event %q, skip content type: %q", event.ID, event.ContentType)
		return nil
	}

	// ignore objects that are too old
	expiration := event.TimeCreated.Add(5 * time.Minute)
	if time.Now().After(expiration) {
		log.Printf("event %q, too old: %s", event.ID, event.TimeCreated.Format(time.Stamp))
		return nil
	}

	annotations, err := detect(ctx, event.Bucket, event.Name)
	if err != nil {
		log.Printf("event %q, detect error: %v", event.ID, err)
		return nil
	}

	if err := publish(ctx, annotations); err != nil {
		log.Printf("event %q, publish error: %v", event.ID, err)
		return nil
	}

	return nil
}
