package text

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
)

var global = struct {
	StorageClient *storage.Client
}{}

func init() {
	log.SetFlags(0)
}

func setup(ctx context.Context) error {
	var err error

	if global.StorageClient == nil {
		global.StorageClient, err = storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("create storage client: %v", err)
		}
	}

	return nil
}
