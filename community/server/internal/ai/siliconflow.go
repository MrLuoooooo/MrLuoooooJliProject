package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type SiliconFlowEngine struct {
	ApiKey string
	Url    string
	Model  string
}

func NewSiliconFlowEngine(apiKey, url, model string) *SiliconFlowEngine {
	return &SiliconFlowEngine{
		ApiKey: apiKey,
		Url:    url,
		Model:  model,
	}
}

func (e *SiliconFlowEngine) Chat(ctx context.Context, messages []Message) (string, error) {
	req := &siliconFlowChatReq{
		Model:    e.Model,
		Messages: e.convertMessages(messages),
		Stream:   false,
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", e.Url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.ApiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("api return error, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp siliconFlowChatResp
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decode response failed: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (e *SiliconFlowEngine) StreamChat(ctx context.Context, messages []Message, callback func(chunk string, isFinish bool)) error {
	req := &siliconFlowChatReq{
		Model:    e.Model,
		Messages: e.convertMessages(messages),
		Stream:   true,
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", e.Url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.ApiKey))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api return error, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				callback("", true)
				return nil
			}
			return fmt.Errorf("read stream failed: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		}

		if line == "[DONE]" {
			callback("", true)
			return nil
		}

		var chunkResp siliconFlowChatResp
		if err := json.Unmarshal([]byte(line), &chunkResp); err != nil {
			continue
		}

		if len(chunkResp.Choices) > 0 {
			choice := chunkResp.Choices[0]
			if choice.FinishReason == "stop" {
				callback("", true)
				return nil
			}
			if choice.Delta.Content != "" {
				callback(choice.Delta.Content, false)
			}
		}
	}
}

func (e *SiliconFlowEngine) ExtractSearchIntent(ctx context.Context, question string) (*SearchIntent, error) {
	systemPrompt := `你是一个搜索意图分析助手。请分析用户的问题，提取出用于搜索帖子的关键词和意图。
请以JSON格式返回，包含keywords(字符串数组)和intent(字符串)。
例如：{"keywords": ["Go", "并发"], "intent": "想了解Go语言的并发编程"}`

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: question},
	}

	resp, err := e.Chat(ctx, messages)
	if err != nil {
		return nil, err
	}

	var intent SearchIntent
	if err := json.Unmarshal([]byte(resp), &intent); err != nil {
		intent.Keywords = []string{question}
		intent.Intent = question
	}

	if len(intent.Keywords) == 0 {
		intent.Keywords = []string{question}
	}

	return &intent, nil
}

func (e *SiliconFlowEngine) convertMessages(messages []Message) []siliconFlowChatMsg {
	result := make([]siliconFlowChatMsg, 0, len(messages))
	for _, msg := range messages {
		result = append(result, siliconFlowChatMsg{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	return result
}

type siliconFlowChatReq struct {
	Model    string               `json:"model"`
	Messages []siliconFlowChatMsg `json:"messages"`
	Stream   bool                 `json:"stream"`
}

type siliconFlowChatMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type siliconFlowChatResp struct {
	Choices []siliconFlowChoice `json:"choices"`
}

type siliconFlowChoice struct {
	Index        int                     `json:"index"`
	Delta        *siliconFlowChoiceDelta `json:"delta"`
	Message      *siliconFlowMessage     `json:"message"`
	FinishReason string                  `json:"finish_reason"`
}

type siliconFlowChoiceDelta struct {
	Content string `json:"content"`
}

type siliconFlowMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
