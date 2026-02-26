package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func GenerateSuggestions(userPrompt, apiKey, model string) ([]Suggestion, error) {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: SystemPrompt()},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Claude API error: %w", err)
	}

	if len(msg.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude API")
	}

	rawText := msg.Content[0].Text
	rawText = strings.TrimSpace(rawText)

	if strings.HasPrefix(rawText, "```") {
		lines := strings.Split(rawText, "\n")
		if len(lines) >= 3 {
			rawText = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	var response AIResponse
	if err := json.Unmarshal([]byte(rawText), &response); err != nil {
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w\n\nRaw response:\n%s", err, rawText)
	}

	if len(response.Suggestions) == 0 {
		return nil, fmt.Errorf("AI returned no suggestions")
	}

	return response.Suggestions, nil
}
