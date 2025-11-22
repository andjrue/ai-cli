package models

import (
	"context"
	"log"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAIProvider struct {
	apiKey string
	client *openai.Client
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIProvider{
		apiKey: apiKey,
		client: &client,
	}
}

func (o *OpenAIProvider) Stream(ctx context.Context, req Request) (<-chan Response, error) {
	if o.apiKey == "" {
		log.Fatal("open ai api key not available. did you set it?")
	}
	
	messages := make([]openai.ChatCompletionMessageParamUnion, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = openai.UserMessage(msg.Content)
	}
	
	stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model: req.Model,
	})
	
	respChan := make(chan Response)
	
	go func() {
		defer close(respChan)
		
		acc := openai.ChatCompletionAccumulator{}
		
		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)
			
			chunkChoice := chunk.Choices[0].Delta.Content
			
			if len(chunk.Choices) > 0 && chunkChoice != "" {
				respChan <- Response{
					Type: ResponseTypeText,
					Content: chunkChoice,
				}
			}
			
			if err := stream.Err(); err != nil {
				respChan <- Response{
					Type: ResponseTypeError,
					Error: err,
				}
				
				return
			}
		}
		
		respChan <- Response{Type: ResponseTypeDone}
	} ()
	
	return respChan, nil
}

