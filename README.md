
# AI Gateway (Go)

A simple Go-based AI proxy gateway. It centralizes AI model calls and routes requests to multiple providers such as Gemini and Ollama under a single `/v1` API.

## Supported Endpoints

- `GET /healthz`
- `GET /v1/models` (enabled models)
- `GET /v1/models?all` (available models from providers, with `enabled` flag)
- `GET /v1/models/{provider}` (enabled models for one provider)
- `GET /v1/models/{provider}?all` (available models for one provider, with `enabled` flag)
- `POST /v1/chat/completions`

## Configuration Example

```yaml
server:
  port: 8080
  logLevel: info
  validateModelsOnStartup: true

fallback:
  maxFallbacks: -1
  retryOnCodes: [429, 500, 502, 503, 504]
  perAttemptTimeout: 30s

models:
# format: {provider}/{model}
- gemini/gemini-2.5-flash
- ollama/gemma4:31b

providers:
  gemini:
    apiKey: YOUR_GEMINI_API_KEY
  ollama:
    apiKey: YOUR_OLLAMA_API_KEY
```

Provider endpoints are fixed in code for `gemini` and `ollama`.

`server.port` defaults to `8080` when omitted.
`server.logLevel` supports `debug`, `info`, `error` (default: `info`).
Use `debug` to print upstream provider request/response logs.
The gateway loads each provider's available model list once during startup and keeps it in memory for `/v1/models?all`.
`server.validateModelsOnStartup` verifies each configured model against that startup model list (default: `true`).
`models` uses `{provider}/{model}` and splits on the first `/`, so model ids containing `/` are supported.
`model` omission is always allowed and routes by configured model order.
When `model` is specified, fallback is always allowed within the resolved model chain.
`fallback.maxFallbacks: -1` means unlimited, `0` disables fallback, and `>0` caps fallback count.
`fallback.perAttemptTimeout` is per-attempt HTTP timeout and timeout errors are fallback-eligible.

## Usage Examples

### List available models

```bash
curl http://localhost:8080/v1/models?all
```

### Create a chat completion

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"gemini/gemini-2.5-flash","messages":[{"role":"user","content":"Hello"}]}'
```
