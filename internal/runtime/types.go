package runtime

// GenerateRequest is structured input for inference.
// See docs/contracts/runtime.md.
type GenerateRequest struct {
	ModelID   string
	Prompt    string
	Grammar   string // Optional GBNF grammar for constrained sampling
	MaxTokens int32  // Optional token limit for this request
}

// MessagePart represents a single piece of a multimodal or tool-augmented message.
type MessagePart struct {
	Type     string // e.g., "text", "image", "tool_call", "tool_result"
	MIMEType string // e.g., "image/jpeg"
	Data     []byte // Raw payload
	Text     string // Text content if type is "text"
}

// ChatMessage represents a single message in a conversation.
type ChatMessage struct {
	Role       string
	Name       string        // Optional: for distinguishing multiple tools or agents
	Content    string        // Raw text content (can be empty if Parts is used)
	ToolCallID string        // Optional: for tool execution tracking
	ToolName   string        // Optional: the name of the tool called/resulting
	Parts      []MessagePart // For multimodal or structured messages
}

// FormatChatOpts specifies constraints and options for formatting prompts.
type FormatChatOpts struct {
	TemplateName        string // Optional: Override default template
	AddGenerationPrompt bool   // True if generating the next turn
	ToolMode            string // Optional: e.g., "json_schema"
	MultimodalMode      string // Optional: e.g., "base64"
}

// ModelCapabilities defines the metadata and supported features of a loaded model.
type ModelCapabilities struct {
	ModelID            string
	TokenizerHash      string
	ContextLength      uint32
	ChatTemplate       string
	AvailableTemplates []string
	SupportsTools      bool
	SupportsMultimodal bool
	SupportsGrammar    bool
}

// GenerateResponse is structured output from inference.
// See docs/contracts/runtime.md.
type GenerateResponse struct {
	ModelID        string
	Output         string
	TokensProduced int
}
