package models

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ai-cli/internal/logger"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
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

// Stream handles method switching depending on OAI model selected
// 5 series model use response streaming
// <= 4 use chat completion (I think, I doubt we'll ever need to go below 4)
func (o *OpenAIProvider) Stream(ctx context.Context, req Request) (<-chan Response, error) {
	if o.apiKey == "" {
		log.Fatal("open ai api key not available. did you set it?")
	}
	
	logger.Log.Printf("[OAI] Model: %s", req.Model)
	
	resp, err := o.streamResponeCompletions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot stream response: %w", err)
	}
	
	return resp, nil
}

func buildSinglePromptFromMessages(messages []Message) string {
	var sb strings.Builder
	for _, msg := range messages {
		sb.WriteString(fmt.Sprintf("%s: %s\n", strings.Title(msg.Role), msg.Content))
	}
	
	return sb.String()
}

func (o *OpenAIProvider) streamResponeCompletions(ctx context.Context, req Request) (<-chan Response, error) {
	messages := buildSinglePromptFromMessages(req.Messages)
	stream := o.client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(messages)},
		Model: req.Model,
	})
	
	respChan := make(chan Response)
	
	// TODO: Figure out streaming for oai. Can't get it to work w/ the response api
	// Rn, it dumps everything, which is fine for mvp.
	go func() {
		defer close(respChan)
		
		for stream.Next() {
			data := stream.Current()
			
			if data.JSON.Text.Valid() {
				respChan <- Response{
					Type: ResponseTypeText,
					Content: data.Text,
				}
			}
		}
		
		if stream.Err() != nil {
			respChan <- Response{
				Type: ResponseTypeError,
				Error: stream.Err(),
			}
			
			return
		}
		
		respChan <- Response{
			Type: ResponseTypeDone,
		}
	}()

	return respChan, nil
}
