# Go Ernie

This library provides unofficial Go clients for [Ernie API](https://ctssmartprogram.baidu.com/doc/pages/hotRecommend/detail?id=12927). We support:

* Ernie-Bot
* Ernie-Bot-turbo
* BLOOMZ-7B

## Installation

```
go get github.com/LazyLinchen/go-ernie
```
go-ernie requires Go version 1.20 or greater.

## Usage

### Ernie-Bot example usage:

```go
package main

import (
	"context"
	"fmt"
	goernie "github.com/LazyLinchen/go-ernie"
)

func main() {
	client, err := goernie.NewClient("ak", "sk")
	if err != nil {
		fmt.Printf("NewClient error: %v\n", err)
		return
    }
	resp, err := client.CreateChatCompletion(context.Background(), goernie.ChatCompletionRequest{
		Model: goernie.ErnieBot,
		Messages: []goernie.ChatCompletionMessage{
			{
				Role:    goernie.ChatMessageRoleUser,
				Content: "你好!",
			},
		},
	})

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Result)
}

```

### Other examples:

<details>
<summary>Ernie-Bot streaming completion</summary>

```go
package main

import (
	"context"
	"errors"
	"fmt"
	goernie "github.com/LazyLinchen/go-ernie"
	"io"
)

func main() {
	client, err := goernie.NewClient("ak", "sk")
	if err != nil {
		fmt.Printf("NewClient error: %v\n", err)
		return
	}
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		goernie.ChatCompletionRequest{
			Model: goernie.ErnieBot,
			Messages: []goernie.ChatCompletionMessage{
				{
					Role:    goernie.ChatMessageRoleUser,
					Content: "你好!",
				},
			},
			Stream: true,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(response.Result)
	}
}
```
</details>

<details>
<summary>Ernie-Bot Turbo streaming completion</summary>

```go
package main

import (
	"context"
	"errors"
	"fmt"
	goernie "github.com/LazyLinchen/go-ernie"
	"io"
)

func main() {
	client, err := goernie.NewClient("ak", "sk")
	if err != nil {
		fmt.Printf("NewClient error: %v\n", err)
		return
	}
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		goernie.ChatCompletionRequest{
			Model: goernie.ErnieBotTurBo,
			Messages: []goernie.ChatCompletionMessage{
				{
					Role:    goernie.ChatMessageRoleUser,
					Content: "你好!",
				},
			},
			Stream: true,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(response.Result)
	}
}
```
</details>

<details>
<summary>Ernie-Bot TurBo completion</summary>

```go
package main

import (
	"context"
	"fmt"
	goernie "github.com/LazyLinchen/go-ernie"
)

func main() {
	client, err := goernie.NewClient("ak", "sk")
	if err != nil {
		fmt.Printf("NewClient error: %v\n", err)
		return
	}
	resp, err := client.CreateChatCompletion(
		context.Background(),
		goernie.ChatCompletionRequest{
			Model: goernie.ErnieBotTurBo,
			Messages: []goernie.ChatCompletionMessage{
				{
					Role:    goernie.ChatMessageRoleUser,
					Content: "你好!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Result)
}
```
</details>

<details>
<summary>BLOOMZ-7B streaming completion</summary>

```go
package main

import (
	"context"
	"errors"
	"fmt"
	goernie "github.com/LazyLinchen/go-ernie"
	"io"
)

func main() {
	client, err := goernie.NewClient("ak", "sk")
	if err != nil {
		fmt.Printf("NewClient error: %v\n", err)
		return
	}
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		goernie.ChatCompletionRequest{
			Model: goernie.BLOOMZ7B,
			Messages: []goernie.ChatCompletionMessage{
				{
					Role:    goernie.ChatMessageRoleUser,
					Content: "你好!",
				},
			},
			Stream: true,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(response.Result)
	}
}
```
</details>

<details>
<summary>BLOOMZ-7B completion</summary>

```go
package main

import (
	"context"
	"fmt"
	goernie "github.com/LazyLinchen/go-ernie"
)

func main() {
	client, err := goernie.NewClient("ak", "sk")
	if err != nil {
		fmt.Printf("NewClient error: %v\n", err)
		return
	}
	resp, err := client.CreateChatCompletion(
		context.Background(),
		goernie.ChatCompletionRequest{
			Model: goernie.BLOOMZ7B,
			Messages: []goernie.ChatCompletionMessage{
				{
					Role:    goernie.ChatMessageRoleUser,
					Content: "你好!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Result)
}
```
</details>