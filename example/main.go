package main

import (
	"context"
	"fmt"
	"os"

	copilot "github.com/stong1994/github-copilot-api"
)

func main() {
	token, err := copilot.DeviceLogin(context.TODO(), os.Getenv("GITHUB_CLIENT_ID"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("GITHUB_COPILOT_TOKEN is %s, you can set it into environment: export GITHUB_COPILOT_TOKEN=%s", token, token)
}
