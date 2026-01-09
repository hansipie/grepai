---
title: Embedders
description: Configure embedding providers for grepai
---

Embedders convert text (code chunks) into vector representations that enable semantic search.

## Available Embedders

| Provider | Type | Pros | Cons |
|----------|------|------|------|
| Ollama | Local | Privacy, free, no internet | Requires local resources |
| OpenAI | Cloud | High quality, fast | Costs money, sends code to cloud |

## Ollama (Local)

### Setup

1. Install Ollama:
```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.com/install.sh | sh
```

2. Start the server:
```bash
ollama serve
```

3. Pull an embedding model:
```bash
ollama pull nomic-embed-text
```

### Configuration

```yaml
embedder:
  provider: ollama
  ollama:
    url: http://localhost:11434
    model: nomic-embed-text
```

### Available Models

| Model | Dimensions | Speed | Quality |
|-------|------------|-------|---------|
| `nomic-embed-text` | 768 | Fast | Good |
| `mxbai-embed-large` | 1024 | Medium | Better |
| `all-minilm` | 384 | Very Fast | Basic |

### Troubleshooting

```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Test embedding
curl http://localhost:11434/api/embeddings -d '{
  "model": "nomic-embed-text",
  "prompt": "Hello world"
}'
```

## OpenAI (Cloud)

### Setup

1. Get an API key from [OpenAI Platform](https://platform.openai.com/api-keys)

2. Set the environment variable:
```bash
export OPENAI_API_KEY=sk-...
```

### Configuration

```yaml
embedder:
  provider: openai
  openai:
    api_key: ${OPENAI_API_KEY}
    model: text-embedding-3-small
```

### Available Models

| Model | Dimensions | Price (per 1M tokens) |
|-------|------------|----------------------|
| `text-embedding-3-small` | 1536 | $0.02 |
| `text-embedding-3-large` | 3072 | $0.13 |

### Cost Estimation

For a typical codebase:
- 10,000 lines of code ≈ 50,000 tokens
- Initial index: ~$0.001 with `text-embedding-3-small`
- Ongoing updates: negligible

## Adding a New Embedder

To add a new embedding provider:

1. Implement the `Embedder` interface in `embedder/`:

```go
type Embedder interface {
    Embed(ctx context.Context, text string) ([]float32, error)
    EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
    Dimensions() int
}
```

2. Add configuration in `config/config.go`

3. Wire it up in the CLI commands

See [Contributing](/grepai/contributing/) for more details.
