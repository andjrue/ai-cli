package models

import (
	"context"

	"github.com/ai-cli/internal/logger"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicProvider struct {
	apiKey string
	client *anthropic.Client
}

func NewAnthropicProvider(apiKey string) *AnthropicProvider {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	
	return &AnthropicProvider{
		apiKey: apiKey,
		client: &client,
	}
}

func (a *AnthropicProvider) Stream(ctx context.Context, req Request) (<- chan Response, error) {
	if a.apiKey == "" {
		logger.Log.Fatal("anthropic api key not available. did you set it?")
	}
	
	messages := make([]anthropic.MessageParam, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content))
	}
	
	stream := a.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model: anthropic.Model(req.Model),
		Messages: messages,
		MaxTokens: 4096,
	})
	
	respChan := make(chan Response)
	go func() {
		defer close(respChan)
		
		for stream.Next() {
			event := stream.Current()
			
			switch eventVariant := event.AsAny().(type) {
				case anthropic.ContentBlockDeltaEvent:
				switch deltaVariant := eventVariant.Delta.AsAny().(type) {
					case anthropic.TextDelta:
					respChan <- Response{
						Type: ResponseTypeText,
						Content: deltaVariant.Text,
					}
				}
			}
		}
		
		if err := stream.Err(); err != nil {
			respChan <- Response{
				Type: ResponseTypeError,
				Error: err,
			}
			
			return
		}
		
		respChan <- Response{Type: ResponseTypeDone}
	}()
	return respChan, nil
}