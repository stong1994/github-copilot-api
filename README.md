# GitHub Copilot API

This repository provides an API for GitHub Copilot. The API includes features for completion, embedding, and device authentication.

## Table of Contents

- [Completion](#completion)
- [Embedding](#embedding)
- [Device Authentication](#device-authentication)

## Completion

The completion feature allows you to generate code completions.

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

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

This example demonstrates how to use the GitHub Copilot API to generate a simple "Hello, World!" program in Go.

## Embedding

The embedding feature allows you to create embeddings for text.

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

### Running the Example

```bash
go run main.go
```

### Expected Output

```bash
0: {[-0.0026773715 -0.0018009724 0.010035048 ...]}
1: {[-0.0312465 -0.0329363 0.020233147 ...]}
```

## Device Authentication

You can get the GitHub Copilot token with device authentication.

### Example

```go
package main

import (
	"context"
	"fmt"

	copilot "github.com/stong1994/github-copilot-api"
)

func main() {
	token, err := copilot.DeviceLogin(context.TODO(), "your GITHUB_CLIENT_ID")
	if err != nil {
		panic(err)
	}
	fmt.Printf("GITHUB_OAUTH_TOKEN is %s, you can set it into environment: export GITHUB_OATUH_TOKEN=%s", token, token)
}
```

To get the `GITHUB_CLIENT_ID`, you can follow [this blog](https://support.heateor.com/get-github-client-id-client-secret/).

### Running the Example

```go
go run main.go
```

### Expected output

```bash
Please take this code "B8F8-AF41" to authenticate at https://github.com/login/device.
Press 'y' to continue, or any other key to abort.
y # after enter 'y', you should goto the web page and paste the code above.
Authenticating, please wait...
GITHUB_OAUTH_TOKEN is xxxxx, you can set it into environment: export GITHUB_OAUTH_TOKEN=xxx
```
