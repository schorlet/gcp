package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
)

var (
	functionURL         = flag.String("url", "", "the function url to invoke")
	impersonatedAccount = flag.String("account", "", "the service account for which the credentials are requested")
)

func main() {
	flag.Parse()
	if flag.NFlag() != 2 {
		flag.Usage()
		os.Exit(2)
	}

	ctx := context.Background()
	c, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		log.Fatalf("new credentials client: %v", err)
	}
	defer c.Close()

	req := &credentialspb.GenerateIdTokenRequest{
		Name:         "projects/-/serviceAccounts/" + *impersonatedAccount,
		Audience:     *functionURL,
		IncludeEmail: true,
	}

	resp, err := c.GenerateIdToken(ctx, req)
	if err != nil {
		log.Fatalf("generate oidc token: %v", err)
	}

	body, err := invoke(resp.Token)
	if err != nil {
		log.Fatalf("invoke: %v", err)
	}

	fmt.Println(body)
}

func invoke(token string) (string, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()

	client := oauth2.NewClient(ctx, ts)
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
