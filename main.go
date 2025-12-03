package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	voicevoxURL = getEnv("VOICEVOX_URL", "http://localhost:50021")
	kokoroURL   = getEnv("KOKORO_URL", "http://localhost:8880")
	serverPort  = getEnv("PORT", ":5108")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type SpeechRequest struct {
	Language string  `json:"language"` // "ja" or "en"
	Text     string  `json:"text"`
	Voice    string  `json:"voice,omitempty"`
	Speed    float64 `json:"speed,omitempty"`
}

type KokoroRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format"`
	Speed          float64 `json:"speed"`
}

func main() {
	http.HandleFunc("GET /health", handleHealth)
	http.HandleFunc("POST /v1/audio/speech", handleSpeech)

	fmt.Printf("KokoVox server starting on %s\n", serverPort)
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleSpeech(w http.ResponseWriter, r *http.Request) {
	var req SpeechRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// デフォルト値の設定
	if req.Speed == 0 {
		req.Speed = 1.0
	}

	var audioData []byte
	var err error

	switch req.Language {
	case "ja":
		audioData, err = synthesizeWithVoiceVox(req)
	case "en":
		audioData, err = synthesizeWithKokoro(req)
	default:
		http.Error(w, "Invalid language. Use 'ja' or 'en'", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Synthesis error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(audioData)))
	w.Write(audioData)
}

func synthesizeWithVoiceVox(req SpeechRequest) ([]byte, error) {
	// デフォルトのspeaker ID（ずんだもん）
	speakerID := "3"
	if req.Voice != "" {
		speakerID = req.Voice
	}

	queryURL := fmt.Sprintf("%s/audio_query?text=%s&speaker=%s",
		voicevoxURL,
		url.QueryEscape(req.Text),
		speakerID,
	)

	queryResp, err := http.Post(queryURL, "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("audio_query request failed: %w", err)
	}
	defer queryResp.Body.Close()

	if queryResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(queryResp.Body)
		return nil, fmt.Errorf("audio_query failed with status %d: %s", queryResp.StatusCode, string(body))
	}

	queryData, err := io.ReadAll(queryResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio_query response: %w", err)
	}

	if req.Speed != 1.0 {
		var queryJSON map[string]interface{}
		if err := json.Unmarshal(queryData, &queryJSON); err == nil {
			queryJSON["speedScale"] = req.Speed
			queryData, _ = json.Marshal(queryJSON)
		}
	}

	synthesisURL := fmt.Sprintf("%s/synthesis?speaker=%s", voicevoxURL, speakerID)

	synthesisResp, err := http.Post(synthesisURL, "application/json", bytes.NewReader(queryData))
	if err != nil {
		return nil, fmt.Errorf("synthesis request failed: %w", err)
	}
	defer synthesisResp.Body.Close()

	if synthesisResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(synthesisResp.Body)
		return nil, fmt.Errorf("synthesis failed with status %d: %s", synthesisResp.StatusCode, string(body))
	}

	audioData, err := io.ReadAll(synthesisResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read synthesis response: %w", err)
	}

	return audioData, nil
}

func synthesizeWithKokoro(req SpeechRequest) ([]byte, error) {
	voice := "af_heart"
	if req.Voice != "" {
		voice = req.Voice
	}

	kokoroReq := KokoroRequest{
		Model:          "kokoro",
		Input:          req.Text,
		Voice:          voice,
		ResponseFormat: "wav",
		Speed:          req.Speed,
	}

	reqBody, err := json.Marshal(kokoroReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	speechURL := fmt.Sprintf("%s/v1/audio/speech", kokoroURL)
	resp, err := http.Post(speechURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("kokoro request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("kokoro failed with status %d: %s", resp.StatusCode, string(body))
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read kokoro response: %w", err)
	}

	return audioData, nil
}
