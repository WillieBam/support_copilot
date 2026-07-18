package interfaces

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/types/requests"
)

type IMCPClient interface {
	DetectAnomalies(ctx context.Context, anomalyReq requests.AnomalyDetectionRequest) (*requests.AnomalyDetectionResponse, error)
}
