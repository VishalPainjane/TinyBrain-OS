package llama

// generateStats holds optional timing data from a generate call (microseconds).
// Both the CGO native path and the Windows dynamic DLL path populate this.
type generateStats struct {
	TTFTMicros   int64 // time-to-first-token
	DecodeMicros int64 // time for all tokens after the first
}
