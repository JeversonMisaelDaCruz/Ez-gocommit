package ai

type Suggestion struct {
	Rank       int    `json:"rank"`
	Confidence string `json:"confidence"`
	Message    string `json:"message"`
	Body       string `json:"body"`
	Reasoning  string `json:"reasoning"`
}

type AIResponse struct {
	Suggestions   []Suggestion `json:"suggestions"`
	DetectedStyle string       `json:"detected_style"`
	Language      string       `json:"language"`
}
