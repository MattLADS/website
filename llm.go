package main

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// LLMClient wraps the OpenAI client for interaction with an LLM.
type LLMClient struct {
	client *openai.Client
}

// NewLLMClient initializes a new LLMClient with the provided OpenAI API key.
func NewLLMClient(apiKey string) *LLMClient {
	client := openai.NewClient(apiKey)
	return &LLMClient{client: client}
}

// ChatCompletion sends a user's query to the LLM and retrieves a response.
func (llm *LLMClient) ChatCompletion(question, systemMessage string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4, // Use GPT-4 model (or any other model supported by the library).
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemMessage,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
	}

	resp, err := llm.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion failed: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
