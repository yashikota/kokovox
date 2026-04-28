package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSpeechEndpointWithVoiceVoxModel(t *testing.T) {
	oldVoiceVoxURL := voicevoxURL
	t.Cleanup(func() { voicevoxURL = oldVoiceVoxURL })

	const audioBody = "RIFFvoicevox-mock-WAVE"

	voiceVoxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/audio_query":
			if r.URL.Query().Get("text") != "こんにちは、世界！" {
				t.Errorf("unexpected text query: %q", r.URL.Query().Get("text"))
			}
			if r.URL.Query().Get("speaker") != "1" {
				t.Errorf("unexpected audio_query speaker: %q", r.URL.Query().Get("speaker"))
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"speedScale":1}`))
		case "/synthesis":
			if r.URL.Query().Get("speaker") != "1" {
				t.Errorf("unexpected synthesis speaker: %q", r.URL.Query().Get("speaker"))
			}

			var body map[string]float64
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("failed to decode synthesis body: %v", err)
			}
			if body["speedScale"] != 1.5 {
				t.Errorf("unexpected speedScale: %v", body["speedScale"])
			}

			w.Header().Set("Content-Type", "audio/wav")
			_, _ = w.Write([]byte(audioBody))
		default:
			http.NotFound(w, r)
		}
	}))
	defer voiceVoxServer.Close()

	voicevoxURL = voiceVoxServer.URL
	appServer := newTestAppServer()
	defer appServer.Close()

	resp := postSpeech(t, appServer.URL, map[string]any{
		"model":           "voicevox",
		"input":           "こんにちは、世界！",
		"voice":           "1",
		"response_format": "wav",
		"speed":           1.5,
	})
	defer resp.Body.Close()

	assertStatus(t, resp, http.StatusOK)
	assertContentType(t, resp, "audio/wav")
	assertBody(t, resp, audioBody)
}

func TestSpeechEndpointWithKokoroModel(t *testing.T) {
	oldKokoroURL := kokoroURL
	t.Cleanup(func() { kokoroURL = oldKokoroURL })

	const audioBody = "RIFFkokoro-mock-WAVE"

	kokoroServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/audio/speech" {
			http.NotFound(w, r)
			return
		}

		var req KokoroRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode kokoro request: %v", err)
		}

		if req != (KokoroRequest{
			Model:          "kokoro",
			Input:          "Hello, world!",
			Voice:          "af_heart",
			ResponseFormat: "wav",
			Speed:          0.9,
		}) {
			t.Errorf("unexpected kokoro request: %+v", req)
		}

		w.Header().Set("Content-Type", "audio/wav")
		_, _ = w.Write([]byte(audioBody))
	}))
	defer kokoroServer.Close()

	kokoroURL = kokoroServer.URL
	appServer := newTestAppServer()
	defer appServer.Close()

	resp := postSpeech(t, appServer.URL, map[string]any{
		"model":           "kokoro",
		"input":           "Hello, world!",
		"voice":           "af_heart",
		"response_format": "wav",
		"speed":           0.9,
	})
	defer resp.Body.Close()

	assertStatus(t, resp, http.StatusOK)
	assertContentType(t, resp, "audio/wav")
	assertBody(t, resp, audioBody)
}

func TestSpeechEndpointRejectsUnsupportedResponseFormat(t *testing.T) {
	appServer := newTestAppServer()
	defer appServer.Close()

	resp := postSpeech(t, appServer.URL, map[string]any{
		"model":           "voicevox",
		"input":           "こんにちは",
		"response_format": "mp3",
	})
	defer resp.Body.Close()

	assertStatus(t, resp, http.StatusBadRequest)
}

func newTestAppServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/audio/speech", handleSpeech)
	return httptest.NewServer(mux)
}

func postSpeech(t *testing.T, baseURL string, payload map[string]any) *http.Response {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	endpoint, err := url.JoinPath(baseURL, "/v1/audio/speech")
	if err != nil {
		t.Fatalf("failed to build endpoint: %v", err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to post speech request: %v", err)
	}

	return resp
}

func assertStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()

	if resp.StatusCode != want {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("unexpected status: got %d, want %d, body: %s", resp.StatusCode, want, body)
	}
}

func assertContentType(t *testing.T, resp *http.Response, want string) {
	t.Helper()

	if got := resp.Header.Get("Content-Type"); got != want {
		t.Fatalf("unexpected content type: got %q, want %q", got, want)
	}
}

func assertBody(t *testing.T, resp *http.Response, want string) {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if string(body) != want {
		t.Fatalf("unexpected body: got %q, want %q", string(body), want)
	}
}
