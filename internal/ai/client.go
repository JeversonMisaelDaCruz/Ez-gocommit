package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const (
	providerAnthropic = "anthropic"
	providerGemini    = "gemini"

	geminiEndpoint     = "https://generativelanguage.googleapis.com/v1beta/openai/chat/completions"
	geminiDefaultModel = "gemini-2.0-flash"

	maxTokens   = 1024
	maxRetries  = 3
)

func GenerateSuggestions(userPrompt, apiKey, model string) ([]Suggestion, error) {
	switch detectProvider(apiKey) {
	case providerGemini:
		return withRetry(maxRetries, func() ([]Suggestion, error) {
			return callGemini(userPrompt, apiKey, resolveGeminiModel(model))
		})
	default:
		return withRetry(maxRetries, func() ([]Suggestion, error) {
			return callAnthropic(userPrompt, apiKey, model)
		})
	}
}

func detectProvider(apiKey string) string {
	if strings.HasPrefix(apiKey, "AIzaSy") {
		return providerGemini
	}
	return providerAnthropic
}

func resolveGeminiModel(model string) string {
	if strings.HasPrefix(model, "gemini-") {
		return model
	}
	return geminiDefaultModel
}

func withRetry(maxAttempts int, fn func() ([]Suggestion, error)) ([]Suggestion, error) {
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
		result, err := fn()
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !isTransientError(err) {
			return nil, err
		}
	}
	return nil, fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
}

func isTransientError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "429") ||
		strings.Contains(msg, "500") ||
		strings.Contains(msg, "502") ||
		strings.Contains(msg, "503") ||
		strings.Contains(msg, "504")
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

func callAnthropic(userPrompt, apiKey, model string) ([]Suggestion, error) {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: maxTokens,
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

	return parseSuggestions(msg.Content[0].Text)
}

type geminiRequest struct {
	Model     string          `json:"model"`
	Messages  []geminiMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens"`
}

type geminiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type geminiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func callGemini(userPrompt, apiKey, model string) ([]Suggestion, error) {
	return callGeminiWithEndpoint(userPrompt, apiKey, model, geminiEndpoint)
}

func callGeminiWithEndpoint(userPrompt, apiKey, model, endpoint string) ([]Suggestion, error) {
	payload := geminiRequest{
		Model: model,
		Messages: []geminiMessage{
			{Role: "system", Content: SystemPrompt()},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens: maxTokens,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Gemini API error: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Gemini response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API error %d: %s", resp.StatusCode, string(rawBody))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(rawBody, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	if geminiResp.Error != nil {
		return nil, fmt.Errorf("Gemini API error: %s", geminiResp.Error.Message)
	}

	if len(geminiResp.Choices) == 0 {
		return nil, fmt.Errorf("empty response from Gemini API")
	}

	return parseSuggestions(geminiResp.Choices[0].Message.Content)
}

func parseSuggestions(raw string) ([]Suggestion, error) {
	raw = strings.TrimSpace(raw)

	if strings.HasPrefix(raw, "```") {
		lines := strings.Split(raw, "\n")
		if len(lines) >= 3 {
			raw = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	var response AIResponse
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		preview := raw
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w\n\nResponse preview:\n%s", err, preview)
	}

	if len(response.Suggestions) == 0 {
		return nil, fmt.Errorf("AI returned no suggestions")
	}

	return response.Suggestions, nil
}
