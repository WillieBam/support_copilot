package types

type CommandResult struct {
	Handled bool   `json:"handled"`
	Message string `json:"message"`
}
