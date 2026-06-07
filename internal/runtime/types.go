package runtime

// GenerateRequest is structured input for inference.
// See docs/contracts/runtime.md.
type GenerateRequest struct {
	ModelID string
	Prompt  string
}

// GenerateResponse is structured output from inference.
// See docs/contracts/runtime.md.
type GenerateResponse struct {
	ModelID        string
	Output         string
	TokensProduced int
}
