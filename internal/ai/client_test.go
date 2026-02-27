package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDetectProvider_Anthropic(t *testing.T) {
	cases := []string{
		"sk-ant-api03-abc123",
		"sk-ant-anything",
		"",
		"random-key",
	}
	for _, key := range cases {
		if got := detectProvider(key); got != providerAnthropic {
			t.Errorf("detectProvider(%q) = %q, want %q", key, got, providerAnthropic)
		}
	}
}

func TestDetectProvider_Gemini(t *testing.T) {
	cases := []string{
		"AIzaSyCftMxUVWgNu0fz8JgGsXnfEOd4MNtBDS0",
		"AIzaSyAnythingElse",
	}
	for _, key := range cases {
		if got := detectProvider(key); got != providerGemini {
			t.Errorf("detectProvider(%q) = %q, want %q", key, got, providerGemini)
		}
	}
}

func TestResolveGeminiModel_AlreadyGemini(t *testing.T) {
	cases := []string{"gemini-2.0-flash", "gemini-1.5-pro", "gemini-pro"}
	for _, m := range cases {
		if got := resolveGeminiModel(m); got != m {
			t.Errorf("resolveGeminiModel(%q) = %q, want same", m, got)
		}
	}
}

func TestResolveGeminiModel_ClaudeFallback(t *testing.T) {
	cases := []string{"claude-sonnet-4-6", "claude-opus-4-6", "claude-haiku-4-5-20251001"}
	for _, m := range cases {
		if got := resolveGeminiModel(m); got != geminiDefaultModel {
			t.Errorf("resolveGeminiModel(%q) = %q, want %q", m, got, geminiDefaultModel)
		}
	}
}

func TestParseSuggestions_ValidJSON(t *testing.T) {
	resp := AIResponse{
		Suggestions: []Suggestion{
			{Rank: 1, Confidence: "high", Message: "feat: add thing", Reasoning: "clear intent"},
			{Rank: 2, Confidence: "medium", Message: "chore: update thing", Reasoning: "could be either"},
			{Rank: 3, Confidence: "low", Message: "fix: adjust thing", Reasoning: "conservative"},
		},
		DetectedStyle: "conventional",
		Language:      "en",
	}
	raw, _ := json.Marshal(resp)

	suggestions, err := parseSuggestions(string(raw))
	if err != nil {
		t.Fatalf("parseSuggestions() unexpected error: %v", err)
	}
	if len(suggestions) != 3 {
		t.Errorf("len(suggestions) = %d, want 3", len(suggestions))
	}
	if suggestions[0].Message != "feat: add thing" {
		t.Errorf("suggestions[0].Message = %q", suggestions[0].Message)
	}
}

func TestParseSuggestions_StripsCodeFence(t *testing.T) {
	resp := AIResponse{
		Suggestions: []Suggestion{
			{Rank: 1, Confidence: "high", Message: "feat: fenced"},
		},
	}
	raw, _ := json.Marshal(resp)
	fenced := "```json\n" + string(raw) + "\n```"

	suggestions, err := parseSuggestions(fenced)
	if err != nil {
		t.Fatalf("parseSuggestions() with code fence: %v", err)
	}
	if len(suggestions) == 0 || suggestions[0].Message != "feat: fenced" {
		t.Errorf("unexpected suggestions: %v", suggestions)
	}
}

func TestParseSuggestions_InvalidJSON(t *testing.T) {
	_, err := parseSuggestions("not json at all")
	if err == nil {
		t.Error("parseSuggestions() should return error for invalid JSON")
	}
}

func TestParseSuggestions_EmptySuggestions(t *testing.T) {
	raw := `{"suggestions":[],"detected_style":"conventional","language":"en"}`
	_, err := parseSuggestions(raw)
	if err == nil {
		t.Error("parseSuggestions() should return error for empty suggestions array")
	}
}

func TestCallGemini_MockServer(t *testing.T) {
	resp := AIResponse{
		Suggestions: []Suggestion{
			{Rank: 1, Confidence: "high", Message: "feat: mock response", Reasoning: "from mock"},
			{Rank: 2, Confidence: "medium", Message: "chore: mock alt", Reasoning: "alternative"},
			{Rank: 3, Confidence: "low", Message: "fix: mock conservative", Reasoning: "conservative"},
		},
		DetectedStyle: "conventional",
		Language:      "en",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Error("missing Bearer token")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("missing Content-Type: application/json")
		}

		var req geminiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if len(req.Messages) != 2 {
			t.Errorf("expected 2 messages (system+user), got %d", len(req.Messages))
		}
		if req.Messages[0].Role != "system" {
			t.Errorf("first message role = %q, want system", req.Messages[0].Role)
		}

		jsonResp := geminiResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: mustMarshal(resp)}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonResp)
	}))
	defer server.Close()

	origEndpoint := geminiEndpoint
	_ = origEndpoint

	suggestions, err := callGeminiWithEndpoint("test prompt", "AIzaSy-test", "gemini-2.0-flash", server.URL)
	if err != nil {
		t.Fatalf("callGemini() error: %v", err)
	}
	if len(suggestions) != 3 {
		t.Errorf("len(suggestions) = %d, want 3", len(suggestions))
	}
	if suggestions[0].Message != "feat: mock response" {
		t.Errorf("unexpected message: %q", suggestions[0].Message)
	}
}

func TestCallGemini_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"message":"invalid API key"}}`))
	}))
	defer server.Close()

	_, err := callGeminiWithEndpoint("prompt", "bad-key", "gemini-2.0-flash", server.URL)
	if err == nil {
		t.Error("expected error for 401 response")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error should mention status code, got: %v", err)
	}
}

func mustMarshal(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
