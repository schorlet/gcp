package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

var functionURL = flag.String("url", "", "function url")

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(-2)
	}

	ctx := context.Background()
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatalf("find default credentials: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(creds.JSON)
	if err != nil {
		log.Fatalf("create jwt config: %v", err)
	}

	body, err := invoke(*conf)
	if err != nil {
		log.Fatalf("invoke: %v", err)
	}

	fmt.Println(body)
}

func invoke(conf jwt.Config) (string, error) {
	conf.PrivateClaims = map[string]interface{}{"target_audience": *functionURL}
	conf.UseIDToken = true

	client := conf.Client(context.Background())
	resp, err := client.Get(*functionURL)
	if err != nil {
		return "", fmt.Errorf("client get: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %v", err)
	}

	return string(body), nil
}
