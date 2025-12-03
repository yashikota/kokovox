# KokoVox

言語に応じてVOICEVOX（日本語）またはKokoro（英語）に振り分ける統合音声合成APIです。

## クイックスタート

```sh
docker compose up --build
```

サーバーは `http://localhost:5108` で起動します。

## API

### ヘルスチェック

```sh
curl http://localhost:5108/health
```

### 音声合成

**エンドポイント:** `POST /v1/audio/speech`

**リクエストボディ:**

| フィールド | 型 | 必須 | 説明 |
|-----------|------|------|------|
| `language` | string | はい | `"ja"` で日本語（VOICEVOX）、`"en"` で英語（Kokoro） |
| `text` | string | はい | 合成するテキスト |
| `voice` | string | いいえ | 音声ID。日本語: スピーカーID（デフォルト: `"3"` = ずんだもん）。英語: 音声名（デフォルト: `"af_heart"`） |
| `speed` | number | いいえ | 話速（デフォルト: `1.0`） |

**レスポンス:** `audio/wav`

### 使用例

**日本語（VOICEVOX）:**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"language": "ja", "text": "こんにちは、世界！"}' \
  --output hello_ja.wav
```

**英語（Kokoro）:**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"language": "en", "text": "Hello, world!"}' \
  --output hello_en.wav
```

**音声と速度を指定:**

```sh
curl -X POST http://localhost:5108/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"language": "ja", "text": "速く話します", "voice": "1", "speed": 1.5}' \
  --output fast.wav
```

## 環境変数

| 変数 | デフォルト | 説明 |
|------|-----------|------|
| `VOICEVOX_URL` | `http://localhost:50021` | VOICEVOX EngineのURL |
| `KOKORO_URL` | `http://localhost:8880` | Kokoro FastAPIのURL |
| `PORT` | `:5108` | サーバーポート |
