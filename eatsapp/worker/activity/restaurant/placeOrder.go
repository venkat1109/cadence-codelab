package restaurant

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"go.uber.org/cadence"
)

func init() {
	cadence.RegisterActivity(PlaceOrderActivity)
}

// PlaceOrderActivity implements of send order activity.
func PlaceOrderActivity(ctx context.Context, wfRunID string, orderID string, items []string) (string, error) {
	return "", errors.New("not implemented")
}

func sendOrder(wfRunID string, orderID string, items []string, taskToken string) error {
	formData := url.Values{}
	formData.Add("id", orderID)
	formData.Add("workflow_id", orderID)
	formData.Add("run_id", wfRunID)
	formData.Add("task_token", taskToken)
	for _, item := range items {
		formData.Add("item", item)
	}
	url := "http://localhost:8090/restaurant"
	_, err := http.PostForm(url, formData)
	if err != nil {
		return err
	}
	return nil
}
