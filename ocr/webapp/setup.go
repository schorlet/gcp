package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

type configuration struct {
	UploadBucket string
	UploadCreds  string
}

type credentials struct {
	Email      string
	PrivateKey []byte
}

var global = struct {
	Config        *configuration
	Creds         *credentials
	StorageClient *storage.Client
}{}

func setup() error {
	if global.Config == nil {
		uploadBucket := os.Getenv("UPLOAD_BUCKET")
		uploadCreds := os.Getenv("UPLOAD_CREDS")

		if uploadBucket == "" {
			return fmt.Errorf("UPLOAD_BUCKET environment variable is missing")
		}
		if uploadCreds == "" {
			return fmt.Errorf("UPLOAD_CREDS environment variable is missing")
		}

		global.Config = &configuration{
			UploadBucket: uploadBucket,
			UploadCreds:  uploadCreds,
		}
	}

	if global.StorageClient == nil {
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("creating storage client: %v", err)
		}

		global.StorageClient = client
	}

	if global.Creds == nil {
		creds, err := readCredentials()
		if err != nil {
			return err
		}

		global.Creds = creds
	}

	return nil
}

func readCredentials() (*credentials, error) {
	ctx := context.Background()

	credsObject := global.Config.UploadCreds
	split := strings.SplitN(credsObject, "/", 2)
	bucketName, objectName := split[0], split[1]

	object := global.StorageClient.Bucket(bucketName).Object(objectName)
	if _, err := object.Attrs(ctx); err != nil {
		return nil, fmt.Errorf("get object %q: %v", credsObject, err)
	}

	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("get reader on credentials: %v", err)
	}
	defer reader.Close()

	jsonKey, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read credentials: %v", err)
	}

	jwtConf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return nil, fmt.Errorf("decode credentials: %v", err)
	}

	return &credentials{
		Email:      jwtConf.Email,
		PrivateKey: jwtConf.PrivateKey,
	}, nil
}
