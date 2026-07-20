package responses

import (
	"time"

	"github.com/WillieBam/support_copilot/backend/types/requests"
)

type CombinedValidationResult struct {
	AlertID      string                            `json:"alert_id"`
	ServiceName  string                            `json:"service_name"`
	Severity     string                            `json:"severity"`
	ReceivedAt   time.Time                         `json:"received_at"`
	Metrics      requests.AnomalyDetectionRequest  `json:"metrics"`
	MLPrediction requests.AnomalyDetectionResponse `json:"ml_prediction"`
}
