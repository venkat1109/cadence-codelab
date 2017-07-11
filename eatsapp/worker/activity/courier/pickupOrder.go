package courier

import (
	"context"
	"errors"

	"go.uber.org/cadence"
)

func init() {
	cadence.RegisterActivity(PickUpOrderActivity)
}

// PickUpOrderActivity implements the pick-up order activity.
func PickUpOrderActivity(ctx context.Context, execution cadence.WorkflowExecution, orderID string) (string, error) {
	return "", errors.New("not implemented")
}

func notifyRestaurant(execution cadence.WorkflowExecution, orderID string) error {
	url := "http://localhost:8090/restaurant?action=p_sig&id=" + orderID +
		"&workflow_id=" + execution.ID + "&run_id=" + execution.RunID
	return sendPatch(url)
}

func pickup(orderID string, taskToken string) error {
	url := "http://localhost:8090/courier?action=p_token&id=" + orderID + "&task_token=" + taskToken
	return sendPatch(url)
}
