# GitHub Copilot API

This repository provides an API for GitHub Copilot.

## Completion

### Example

```go
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
	// you can get the content directly
	// fmt.Println(response.Choices[0].Message.Content)
}
```

### Running the Example

To run the example, use the following command:

```bash
go run main.go
```

### Expected Output

The output will be similar to this:

```
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}

This program first imports the `fmt` package, which contains functions for formatting text, including printing to the console. Then it defines a `main` function. When the program runs, it starts by executing this function. The `main` function uses `fmt.Println` to print "Hello, World!" to the console.

```

This example demonstrates how to use the GitHub Copilot API to generate a simple "Hello, World!" program in Go.

## Embedding

### Example

```go
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
```

### Run

```bash
go run main.go
```

### Expected Output

```bash
0: {[-0.0026773715 -0.0018009724 0.010035048 ...]}
1: {[-0.0312465 -0.0329363 0.020233147 ...]}
```
