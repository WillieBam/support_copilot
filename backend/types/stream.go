package types

type StreamEvent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
