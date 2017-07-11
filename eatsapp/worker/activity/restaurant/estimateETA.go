package restaurant

import (
	"context"
	"time"

	"go.uber.org/cadence"
)

func init() {
	cadence.RegisterActivity(EstimateETAActivity)
}

// EstimateETAActivity implements the estimate eta activity.
func EstimateETAActivity(ctx context.Context, orderID string) (time.Duration, error) {
	return time.Minute, nil
}
