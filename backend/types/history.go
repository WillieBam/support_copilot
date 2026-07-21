package types

// HistoryMessage represents a single turn in a conversation.
// Role is either "user" or "assistant".
type HistoryMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
