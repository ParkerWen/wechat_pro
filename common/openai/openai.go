package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	config := openai.DefaultConfig("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")
	proxyUrl, err := url.Parse("http://192.168.56.1:7890")
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}
	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "买来的steamdack已经被激活过",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
