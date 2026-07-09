package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	v2 "github.com/VishalPainjane/TinyBrain-OS/internal/scheduler/v2"
	"github.com/sugarme/tokenizer"
)

// HTTPServer handles HTTP ingress for the engine via Server-Sent Events.
type HTTPServer struct {
	engine    *v2.Engine
	tokenizer *tokenizer.Tokenizer
}

// NewHTTPServer constructs a new HTTPServer bound to the given engine.
func NewHTTPServer(engine *v2.Engine, tk *tokenizer.Tokenizer) *HTTPServer {
	return &HTTPServer{
		engine:    engine,
		tokenizer: tk,
	}
}

// GenerateRequest represents the expected JSON payload for generating tokens.
type GenerateRequest struct {
	Prompt     string `json:"prompt"`
	MaxTokens  int32  `json:"max_tokens"`
	EosTokenID int32  `json:"eos_token_id"`
}

// HandleGenerate handles token generation requests and streams responses via SSE.
func (s *HTTPServer) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Setup Server-Sent Events headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	if s.tokenizer == nil {
		http.Error(w, "Tokenizer not configured on server", http.StatusInternalServerError)
		return
	}

	enc, err := s.tokenizer.EncodeSingle(req.Prompt)
	if err != nil {
		http.Error(w, "Failed to encode prompt", http.StatusInternalServerError)
		return
	}

	var promptIDs []int32
	// TinyLlama requires BOS token (1) at the start
	if len(enc.Ids) == 0 || enc.Ids[0] != 1 {
		promptIDs = append(promptIDs, 1)
	}
	
	for _, id := range enc.Ids {
		promptIDs = append(promptIDs, int32(id))
	}

	tokenChan, seqID := s.engine.Submit(promptIDs, req.MaxTokens, req.EosTokenID)

	// Stream tokens as they are generated
	for {
		select {
		case <-r.Context().Done():
			// Client disconnected early
			s.engine.Cancel(seqID)
			return
		case token, ok := <-tokenChan:
			if !ok {
				// Channel closed, generation complete
				return
			}

			if token == 2 || (req.EosTokenID != 0 && token == req.EosTokenID) {
				// Generation complete (EOS token)
				return
			}
			if token == 0 {
				fmt.Fprintf(w, "data: {\"text\": \"<unk>\"}\n\n")
				flusher.Flush()
				continue
			}
			textChunk := s.tokenizer.Decode([]int{int(token)}, true)
			fmt.Fprintf(w, "data: {\"text\": %q}\n\n", textChunk)
			flusher.Flush()
		}
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int32         `json:"max_tokens"`
	Temperature float32       `json:"temperature"`
}

// HandleChatCompletions implements an OpenAI-compatible ChatML endpoint.
func (s *HTTPServer) HandleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	if s.tokenizer == nil {
		http.Error(w, "Tokenizer not configured on server", http.StatusInternalServerError)
		return
	}

	estimatedSize := 0
	for _, msg := range req.Messages {
		estimatedSize += len(msg.Role) + len(msg.Content) + 20
	}
	estimatedSize += 20

	var sb strings.Builder
	sb.Grow(estimatedSize)

	for _, msg := range req.Messages {
		switch msg.Role {
		case "system":
			sb.WriteString("<|system|>\n")
		case "user":
			sb.WriteString("<|user|>\n")
		case "assistant":
			sb.WriteString("<|assistant|>\n")
		}
		sb.WriteString(msg.Content)
		sb.WriteString("\n")
	}
	sb.WriteString("<|assistant|>\n")

	prompt := sb.String()

	select {
	case <-r.Context().Done():
		return
	default:
	}

	enc, err := s.tokenizer.EncodeSingle(prompt)
	if err != nil {
		http.Error(w, "Failed to encode prompt", http.StatusInternalServerError)
		return
	}

	var promptIDs []int32
	if len(enc.Ids) == 0 || enc.Ids[0] != 1 {
		promptIDs = append(promptIDs, 1)
	}
	for _, id := range enc.Ids {
		promptIDs = append(promptIDs, int32(id))
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 256
	}

	tokenChan, seqID := s.engine.Submit(promptIDs, maxTokens, 2) // Defaults to EOS 2

	for {
		select {
		case <-r.Context().Done():
			s.engine.Cancel(seqID)
			return
		case token, ok := <-tokenChan:
			if !ok {
				return
			}
			// 2 is standard EOS, 32000 is common ChatML EOS
			if token == 2 || token == 32000 {
				fmt.Fprintf(w, "data: [DONE]\n\n")
				flusher.Flush()
				return
			}
			if token == 0 {
				continue
			}

			textChunk := s.tokenizer.Decode([]int{int(token)}, true)

			chunkData := map[string]interface{}{
				"id":      seqID,
				"object":  "chat.completion.chunk",
				"choices": []map[string]interface{}{
					{
						"index": 0,
						"delta": map[string]interface{}{
							"content": textChunk,
						},
					},
				},
			}

			chunkJSON, _ := json.Marshal(chunkData)
			fmt.Fprintf(w, "data: %s\n\n", chunkJSON)
			flusher.Flush()
		}
	}
}
