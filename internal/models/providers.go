// Package models provides interfaces and types used when connecting with the selected provider
package models

import "context"

type Provider interface {
	Stream(ctx context.Context, req Request) (<-chan Response, error)
	// ListModels() []Model
}

type Request struct {
	Model        string
	Messages     []Message
	SystemPrompt string
	//Tools []Tool - TODO
}

type Response struct {
	Type    ResponseType
	Content string
	Error   error
}

type Message struct {
	Role    string
	Content string
}

type Model struct {
	ModelFamily string
	ModelName   string
}

type ResponseType int

const (
	ResponseTypeText ResponseType = iota
	ResponseTypeToolUse
	ResponseTypeError
	ResponseTypeDone
)
