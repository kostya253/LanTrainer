package translator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// requestBody is the JSON payload sent to the local LLM endpoint.
type requestBody struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// responseBody is the JSON structure returned by the LLM.
type responseBody struct {
	Response string `json:"response"`
}

// TranslateUsingLLM sends the prompt to the locally‑running LLM
// (http://localhost:11434/api/generate) and returns the Czech translation
// of the supplied English text.
//
// The function respects the provided context for cancellation and timeout
// and uses a 5‑minute deadline similar to the original Python code.
func TranslateUsingLLM(ctx context.Context, text string) (string, error) {
	// Enforce a sane timeout – the Python code uses 300 s.
	const defaultTimeout = 5 * time.Minute

	// Create a child context if the caller didn't supply one
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// Build the request JSON.
	promptUsed := os.Getenv("TRANSLATE_PROMPT")

	// The two new lines are a requirements for translategemma model, if you are using another
	// consider changing the prompt
	prompt := fmt.Sprintf(`%s\n\n%s`, promptUsed, text)

	modelUsed := os.Getenv("TRANSLATE_MODEL")

	payload := requestBody{
		Model:  modelUsed,
		Prompt: prompt,
		Stream: false,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshaling request body: %w", err)
	}

	// Build the request.
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:11434/api/generate", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request.
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request to LLM: %w", err)
	}
	defer resp.Body.Close()

	// Handle non‑200 status codes.
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("non‑OK HTTP status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode the JSON response.
	var res responseBody
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("decoding LLM response: %w", err)
	}

	return res.Response, nil
}
