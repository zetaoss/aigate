APP := aigate
CONFIG ?= config.yaml
BASE_URL ?= http://localhost:8080
MODEL ?= gemini/gemini-3-flash-preview
PROVIDER ?= ollama
PROMPT ?= Hello
escape_squote = $(subst ','\'',$(1))

.PHONY: help run curl-health curl-models curl-models-all curl-chat curl-chat-with-model curl-chat-with-provider

help:
	@echo "Targets:"
	@echo "  make run                - run gateway server"
	@echo "  make curl-health        - GET /healthz"
	@echo "  make curl-models        - GET /v1/models"
	@echo "  make curl-models-all    - GET /v1/models?all"
	@echo "  make curl-chat            - POST /v1/chat/completions (without model)"
	@echo "  make curl-chat-with-model - POST /v1/chat/completions (with model)"
	@echo "  make curl-chat-with-provider - POST /v1/chat/completions (with provider-only model)"
	@echo ""
	@echo "Variables:"
	@echo "  CONFIG=<path>           (default: config.yaml)"
	@echo "  BASE_URL=<url>          (default: http://localhost:8080)"
	@echo "  MODEL=<provider/model>  (default: gemini/gemini-2.5-flash)"
	@echo "  PROVIDER=<provider>     (default: ollama)"
	@echo "  PROMPT=<text>           (default: Hello)"

run:
	go run ./cmd/$(APP) -config $(CONFIG)

curl-health:
	curl -sS $(BASE_URL)/healthz | jq

curl-models:
	curl -sS $(BASE_URL)/v1/models | jq

curl-models-all:
	curl -sS "$(BASE_URL)/v1/models?all" | jq

curl-chat:
	@PROMPT_JSON="$$(printf '%s' '$(call escape_squote,$(PROMPT))' | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))')"; \
	curl -sS -X POST $(BASE_URL)/v1/chat/completions \
		-H "Content-Type: application/json" \
		-d "{\"messages\":[{\"role\":\"user\",\"content\":$${PROMPT_JSON}}]}" | jq

curl-chat-with-model:
	@PROMPT_JSON="$$(printf '%s' '$(call escape_squote,$(PROMPT))' | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))')"; \
	curl -sS -X POST $(BASE_URL)/v1/chat/completions \
		-H "Content-Type: application/json" \
		-d "{\"model\":\"$(MODEL)\",\"messages\":[{\"role\":\"user\",\"content\":$${PROMPT_JSON}}]}" | jq

curl-chat-with-provider:
	@PROMPT_JSON="$$(printf '%s' '$(call escape_squote,$(PROMPT))' | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))')"; \
	curl -sS -X POST $(BASE_URL)/v1/chat/completions \
		-H "Content-Type: application/json" \
		-d "{\"model\":\"$(PROVIDER)\",\"messages\":[{\"role\":\"user\",\"content\":$${PROMPT_JSON}}]}" | jq
