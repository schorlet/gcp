package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSignedURL(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Content-Type", "image/jpg")
	req.Header.Set("X-Content-Length", "23096")

	rr := httptest.NewRecorder()

	signedURL(rr, req)
	resp := rr.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status code = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	contentType := resp.Header.Get("Content-Type")
	if want := "text/plain"; !strings.HasPrefix(contentType, want) {
		t.Fatalf("content type = %v, want %v", contentType, want)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	sbody := string(body)

	u, err := url.Parse(sbody)
	if err != nil {
		t.Fatalf("parsing signed url: %v", err)
	}
	fmt.Println(u.String())
}

func TestSignedURL_BadRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	signedURL(rr, req)
	resp := rr.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status code = %v, want %v", resp.StatusCode, http.StatusBadRequest)
	}
}
