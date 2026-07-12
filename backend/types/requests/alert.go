package requests

import (
	"encoding/json"

	"github.com/google/uuid"
)

type AlertIngestRequest struct {
	IncidentID  uuid.UUID       `json:"incident_id"`
	ServiceName string          `json:"service_name"`
	Severity    string          `json:"severity"`
	Metrics     json.RawMessage `json:"metrics"`
}
