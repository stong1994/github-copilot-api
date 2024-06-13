# GitHub Copilot API

This repository provides an API for GitHub Copilot.

## Example

Here's an example of how to use the API:

```go
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

	prompt := "you are a great developer!"
	request := "give me a code to print hello world with golang"
	opts := copilot.CompletionOpts{
		Model:       "gpt-4",
		N:           1,
		TopP:        1,
		Stream:      false,
		Temperature: 0.1,
	}

	response, err := client.Complete(prompt, request, opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Choices[0].Message.Content)
}
```

To run the example, use the following command:

```bash
go run main.go
```

The output will be similar to this:

```
Sure, here is a simple "Hello, World!" program in Go:

package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

This program will print "Hello, World!" to the console.

```

This example demonstrates how to use the GitHub Copilot API to generate a simple "Hello, World!" program in Go.
```
