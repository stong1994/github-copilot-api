package main

import (
	"context"
	"fmt"

	copilot "github.com/stong1994/github-copilot-api"
)

func main() {
	client, err := copilot.NewCopilot()
	if err != nil {
		panic(err)
	}
	response, err := client.CreateEmbedding(context.Background(), &copilot.EmbeddingRequest{
		Model: "copilot-text-embedding-ada-002",
		Input: []string{
			"you are a great developer!",
			"thanks for your help",
		},
	})
	if err != nil {
		panic(err)
	}
	for i, embedding := range response.Data {
		fmt.Printf("%d: %v\n", i, embedding)
	}
}
