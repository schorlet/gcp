package main

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func detectText(w http.ResponseWriter, r *http.Request) {
	if err := setup(); err != nil {
		http.Error(w, "Detect Text", http.StatusInternalServerError)
		return
	}

	bucket := global.Config.UploadBucket
	name := filepath.Base(r.URL.Path) + ".json"

	object := global.StorageClient.Bucket(bucket).Object(name)
	if _, err := object.Attrs(r.Context()); err != nil {
		http.NotFound(w, r)
		return
	}

	reader, err := object.NewReader(r.Context())
	if err != nil {
		log.Printf("reading object: %v", err)
		http.Error(w, "Reading object", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", reader.Attrs.ContentType)

	if _, err = io.Copy(w, reader); err != nil {
		log.Printf("copying object: %v", err)
		http.Error(w, "Copying object", http.StatusInternalServerError)
		return
	}
}
