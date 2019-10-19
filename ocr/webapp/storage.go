package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

type signedHeader struct {
	ContentType   string
	ContentLength string
}

func signedURL(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("X-Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "Missing headers", http.StatusBadRequest)
		return
	}

	contentLength := r.Header.Get("X-Content-Length")
	if contentLength == "" {
		http.Error(w, "Missing headers", http.StatusBadRequest)
		return
	}

	filename := randomString(30)
	header := signedHeader{
		ContentType:   contentType,
		ContentLength: contentLength,
	}

	url, err := createSignedURL(filename, header)
	if err != nil {
		log.Printf("Creating signedURL: %v", err)
		http.Error(w, "Creating signedURL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, url)
}

func createSignedURL(filename string, header signedHeader) (string, error) {
	if err := setup(); err != nil {
		return "", fmt.Errorf("setup: %v", err)
	}

	if header.ContentLength == "" || header.ContentType == "" {
		return "", fmt.Errorf("missing header")
	}
	if !strings.HasPrefix(header.ContentType, "image/") {
		return "", fmt.Errorf("bad header")
	}

	options := storage.SignedURLOptions{
		Headers: []string{
			"Content-Length: " + header.ContentLength,
			// 5MiB
			"x-goog-content-length-range: 0,5242880",
			// the object does not currently exist
			"x-goog-if-generation-match: 0",
			// storage-class
			"x-goog-storage-class: STANDARD",
		},
		Method:         "PUT",
		ContentType:    header.ContentType,
		Expires:        time.Now().Add(5 * time.Minute),
		GoogleAccessID: global.Creds.Email,
		PrivateKey:     global.Creds.PrivateKey,
		Scheme:         storage.SigningSchemeV4,
	}

	return storage.SignedURL(global.Config.UploadBucket, filename, &options)
}
