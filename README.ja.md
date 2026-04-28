# KokoVox

`model` に応じてVOICEVOX（日本語）またはKokoro（英語）に振り分けるOpenAI互換の統合音声合成APIです。

## クイックスタート

```sh
docker compose up --build -d
```

サーバーは `http://localhost:5108` で起動します。

## API

### ヘルスチェック

```sh
curl http://localhost:5108/health
```

### 音声合成

**エンドポイント:** `POST /v1/audio/speech`

このエンドポイントはOpenAI Text-to-Speech APIのリクエスト形式と互換性があります。  

**リクエストボディ:**

| フィールド | 型 | 必須 | 説明 |
|-----------|------|------|------|
| `model` | string | はい | `"voicevox"` で日本語、`"kokoro"` で英語 |
| `input` | string | はい | 合成するテキスト |
| `voice` | string | いいえ | 音声ID。日本語: スピーカーID（デフォルト: `"3"` = ずんだもん）。英語: 音声名（デフォルト: `"af_heart"`） |
| `response_format` | string | いいえ | 音声形式。現在は `"wav"` のみ対応（デフォルト: `"wav"`） |
| `speed` | number | いいえ | 話速（デフォルト: `1.0`） |

**レスポンス:** `audio/wav`

### 使用例

**日本語（VOICEVOX）:**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"model": "voicevox", "input": "こんにちは、世界！"}' \
  --output hello_ja.wav
```

**英語（Kokoro）:**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"model": "kokoro", "input": "Hello, world!"}' \
  --output hello_en.wav
```

**音声と速度を指定:**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"model": "voicevox", "input": "速く話します", "voice": "1", "speed": 1.5}' \
  --output fast.wav
```

## 環境変数

| 変数 | デフォルト | 説明 |
|------|-----------|------|
| `VOICEVOX_URL` | `http://localhost:50021` | VOICEVOX EngineのURL |
| `KOKORO_URL` | `http://localhost:8880` | Kokoro FastAPIのURL |
| `PORT` | `:5108` | サーバーポート |
