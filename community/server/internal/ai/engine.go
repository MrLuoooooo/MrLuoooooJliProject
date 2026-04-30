package ai

import (
	"context"

	"community-server/internal/config"

	"go.uber.org/fx"
)

type Engine interface {
	Chat(ctx context.Context, messages []Message) (string, error)
	StreamChat(ctx context.Context, messages []Message, callback func(chunk string, isFinish bool)) error
	ExtractSearchIntent(ctx context.Context, question string) (*SearchIntent, error)
}

type SearchIntent struct {
	Keywords []string `json:"keywords"`
	Intent   string   `json:"intent"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages []Message `json:"messages" binding:"required"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Content string `json:"content"`
}

type EngineParams struct {
	fx.In

	Config *config.Config
}

func NewEngine(params EngineParams) Engine {
	return NewSiliconFlowEngine(
		params.Config.AI.ApiKey,
		params.Config.AI.Url,
		params.Config.AI.Model,
	)
}
