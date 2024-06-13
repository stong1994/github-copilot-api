package main

import (
	"fmt"

	copilot "github.com/stong1994/github-copilot-api"
)

func main() {
	client, err := copilot.NewCopilot()
	if err != nil {
		panic(err)
	}
	response, err := client.Complete("you are a great developer!", "give me a code to print hello world with golang", copilot.CompletionOpts{
		Model:       "gpt-4",
		N:           1,
		TopP:        1,
		Stream:      false,
		Temperature: 0.1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Choices[0].Message.Content)
}
