# KokoVox

[日本語版README](README.ja.md)  

A unified Text-to-Speech API that routes to VOICEVOX (Japanese) or Kokoro (English) based on language selection.

## Quick Start

```sh
docker compose up --build
```

The server will be available at `http://localhost:5108`

## API

### Health Check

```sh
curl http://localhost:5108/health
```

### Text-to-Speech

**Endpoint:** `POST /v1/audio/speech`

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `language` | string | Yes | `"ja"` for Japanese (VOICEVOX) or `"en"` for English (Kokoro) |
| `text` | string | Yes | Text to synthesize |
| `voice` | string | No | Voice ID. For Japanese: speaker ID (default: `"3"` = Zundamon). For English: voice name (default: `"af_heart"`) |
| `speed` | number | No | Speech speed (default: `1.0`) |

**Response:** `audio/wav`

### Examples

**Japanese (VOICEVOX):**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"language": "ja", "text": "こんにちは、世界！"}' \
  --output hello_ja.wav
```

**English (Kokoro):**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"language": "en", "text": "Hello, world!"}' \
  --output hello_en.wav
```

**With custom voice and speed:**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"language": "ja", "text": "速く話します", "voice": "1", "speed": 1.5}' \
  --output fast.wav
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `VOICEVOX_URL` | `http://localhost:50021` | VOICEVOX Engine URL |
| `KOKORO_URL` | `http://localhost:8880` | Kokoro FastAPI URL |
| `PORT` | `:5108` | Server port |
