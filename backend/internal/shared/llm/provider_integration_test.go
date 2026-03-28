// Package llm integration tests.
// These tests verify the LLM provider works end-to-end with real API calls.
// They are skipped by default; set MINIMAX_API_KEY to run.
package llm

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/openai/openai-go"
)

func TestIntegrationMiniMaxChatCompletion(t *testing.T) {
	apiKey := os.Getenv("MINIMAX_API_KEY")
	if apiKey == "" {
		t.Skip("MINIMAX_API_KEY not set; skipping integration test")
	}

	client, err := NewClient(ProviderMiniMax, apiKey, "")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.Background()
	model := ModelOrDefault("", ProviderMiniMax)
	temp := ClampTemperature(0.1, ProviderMiniMax)

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       model,
		Temperature: openai.Float(temp),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Say 'hello' and nothing else."),
		},
	})
	if err != nil {
		t.Fatalf("Chat completion error: %v", err)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("No choices returned")
	}

	content := strings.ToLower(resp.Choices[0].Message.Content)
	if !strings.Contains(content, "hello") {
		t.Errorf("Expected 'hello' in response, got: %s", resp.Choices[0].Message.Content)
	}
}

func TestIntegrationMiniMaxStreaming(t *testing.T) {
	apiKey := os.Getenv("MINIMAX_API_KEY")
	if apiKey == "" {
		t.Skip("MINIMAX_API_KEY not set; skipping integration test")
	}

	client, err := NewClient(ProviderMiniMax, apiKey, "")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.Background()
	model := ModelOrDefault("", ProviderMiniMax)

	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model: model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Count from 1 to 3."),
		},
	})

	var chunks int
	var fullContent strings.Builder

	for stream.Next() {
		chunk := stream.Current()
		chunks++
		for _, choice := range chunk.Choices {
			fullContent.WriteString(choice.Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		t.Fatalf("Stream error: %v", err)
	}

	if chunks == 0 {
		t.Fatal("No streaming chunks received")
	}

	if fullContent.Len() == 0 {
		t.Fatal("Streaming produced empty content")
	}
}

func TestIntegrationMiniMaxHighspeedModel(t *testing.T) {
	apiKey := os.Getenv("MINIMAX_API_KEY")
	if apiKey == "" {
		t.Skip("MINIMAX_API_KEY not set; skipping integration test")
	}

	client, err := NewClient(ProviderMiniMax, apiKey, "")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.Background()

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       "MiniMax-M2.7-highspeed",
		Temperature: openai.Float(0.5),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("What is 2+2? Reply with just the number."),
		},
	})
	if err != nil {
		t.Fatalf("Chat completion error: %v", err)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("No choices returned")
	}

	content := resp.Choices[0].Message.Content
	if !strings.Contains(content, "4") {
		t.Errorf("Expected '4' in response, got: %s", content)
	}
}
