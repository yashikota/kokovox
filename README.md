# KokoVox

[日本語版README](README.ja.md)  

An OpenAI-compatible Text-to-Speech API that routes to VOICEVOX (Japanese) or Kokoro (English) based on the `model` field.  

## Quick Start

```sh
docker compose up --build -d
```

The server will be available at `http://localhost:5108`

## API

### Health Check

```sh
curl http://localhost:5108/health
```

### Text-to-Speech

**Endpoint:** `POST /v1/audio/speech`

This endpoint is compatible with the OpenAI Text-to-Speech API request shape.  

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `model` | string | Yes | `"voicevox"` for Japanese or `"kokoro"` for English |
| `input` | string | Yes | Text to synthesize |
| `voice` | string | No | Voice ID. For Japanese: speaker ID (default: `"3"` = Zundamon). For English: voice name (default: `"af_heart"`) |
| `response_format` | string | No | Audio format. Currently only `"wav"` is supported (default: `"wav"`) |
| `speed` | number | No | Speech speed (default: `1.0`) |

**Response:** `audio/wav`

### Examples

**Japanese (VOICEVOX):**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"model": "voicevox", "input": "こんにちは、世界！"}' \
  --output hello_ja.wav
```

**English (Kokoro):**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"model": "kokoro", "input": "Hello, world!"}' \
  --output hello_en.wav
```

**With custom voice and speed:**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"model": "voicevox", "input": "速く話します", "voice": "1", "speed": 1.5}' \
  --output fast.wav
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `VOICEVOX_URL` | `http://localhost:50021` | VOICEVOX Engine URL |
| `KOKORO_URL` | `http://localhost:8880` | Kokoro FastAPI URL |
| `PORT` | `:5108` | Server port |
