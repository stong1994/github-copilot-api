package main

import (
	"context"
	"fmt"

	copilot "github.com/stong1994/github-copilot-api"
)

func main() {
	token, err := copilot.DeviceLogin(context.TODO(), "Iv23lisucfy8lSS2j0sn")
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
}
