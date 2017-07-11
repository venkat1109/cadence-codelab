package restaurant

import (
	"time"

	"go.uber.org/cadence"
	"go.uber.org/zap"

	"github.com/venkat1109/cadence-codelab/eatsapp/worker/activity/restaurant"
)

func init() {
	cadence.RegisterWorkflow(OrderWorkflow)
}

// OrderWorkflow implements the restaurant order workflow.
func OrderWorkflow(ctx cadence.Context, wfRunID string, orderID string, items []string) (time.Duration, error) {

	ao := cadence.ActivityOptions{
		ScheduleToStartTimeout: time.Minute * 5,
		StartToCloseTimeout:    time.Minute * 15,
	}

	ctx = cadence.WithActivityOptions(ctx, ao)
	err := cadence.ExecuteActivity(ctx, restaurant.PlaceOrderActivity, wfRunID, orderID, items).Get(ctx, nil)
	if err != nil {
		cadence.GetLogger(ctx).Error("Failed to send order to restaurant", zap.Error(err))
		return time.Minute * 0, err
	}

	var eta time.Duration
	err = cadence.ExecuteActivity(ctx, restaurant.EstimateETAActivity, orderID).Get(ctx, &eta)
	if err != nil {
		cadence.GetLogger(ctx).Error("Failed to estimate ETA for order ready", zap.Error(err))
		return time.Minute * 0, err
	}

	cadence.GetLogger(ctx).Info("Completed PlaceOrder!")
	return eta, err
}
