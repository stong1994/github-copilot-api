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
	response, err := client.CreateCompletion(context.Background(), &copilot.CompletionRequest{
		Messages: []copilot.Message{
			{
				Role:    "system",
				Content: "you are a great developer!",
			},
			{
				Role:    "user",
				Content: "give me a code to print hello world with golang",
			},
		},
		StreamingFunc: func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\n======================\nResponse finished, above is the stream response, below is the full content\n======================")
	fmt.Println(response.Choices[0].Message.Content)
}
