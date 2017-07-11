package courier

import (
	"time"

	"github.com/venkat1109/cadence-codelab/eatsapp/worker/activity/courier"

	"go.uber.org/cadence"
	"go.uber.org/zap"
)

func init() {
	cadence.RegisterWorkflow(OrderWorkflow)
}

// OrderWorkflow implements the deliver order workflow.
func OrderWorkflow(ctx cadence.Context, orderID string) error {

	ao := cadence.ActivityOptions{
		ScheduleToStartTimeout: time.Minute * 5,
		StartToCloseTimeout:    time.Minute * 15,
	}
	ctx = cadence.WithActivityOptions(ctx, ao)

	for {
		err := cadence.ExecuteActivity(ctx, courier.DispatchCourierActivity, orderID).Get(ctx, nil)
		if err != nil {
			// retry forever until a driver accepts the trip
			cadence.GetLogger(ctx).Error("Failed to dispatch courier", zap.Error(err))
			continue
		}
		break
	}

	execution := cadence.GetWorkflowInfo(ctx).WorkflowExecution
	err := cadence.ExecuteActivity(ctx, courier.PickUpOrderActivity, execution, orderID).Get(ctx, nil)
	if err != nil {
		cadence.GetLogger(ctx).Error("Failed to pick up order from restaurant", zap.Error(err))
		return err
	}

	err = waitForRestaurantPickupConfirmation(ctx, orderID)
	if err != nil {
		cadence.GetLogger(ctx).Error("Failed to confirm pickup with restaurant", zap.Error(err))
		return err
	}

	err = cadence.ExecuteActivity(ctx, courier.DeliverOrderActivity, orderID).Get(ctx, nil)
	if err != nil {
		cadence.GetLogger(ctx).Error("Failed to complete delivery", zap.Error(err))
		return err
	}

	return nil
}
