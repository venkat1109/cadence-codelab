package courier

import (
	"context"
	"errors"

	"go.uber.org/cadence"
)

func init() {
	cadence.RegisterActivity(DeliverOrderActivity)
}

// DeliverOrderActivity implements the devliver order activity.
func DeliverOrderActivity(ctx context.Context, orderID string) (string, error) {
	return "", errors.New("not implemented")
}

func deliver(orderID string, taskToken string) error {
	url := "http://localhost:8090/courier?action=c_token&id=" + orderID + "&task_token=" + taskToken
	return sendPatch(url)
}
