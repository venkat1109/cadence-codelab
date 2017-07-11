package eats

import (
	"go.uber.org/cadence"
	"go.uber.org/zap"
)

func init() {
	cadence.RegisterWorkflow(OrderWorkflow)
}

// OrderWorkflow implements the eats order workflow.
func OrderWorkflow(ctx cadence.Context, orderID string, items []string) error {

	cadence.GetLogger(ctx).Info("Received order", zap.Strings("items", items))

	restaurantEta, err := placeRestaurantOrder(ctx, orderID, items)
	if err != nil {
		return err
	}

	err = waitForRestaurant(ctx, orderID, restaurantEta)
	if err != nil {
		return err
	}

	err = deliverOrder(ctx, orderID)
	if err != nil {
		return err
	}

	err = chargeOrder(ctx, orderID)
	if err != nil {
		return err
	}

	cadence.GetLogger(ctx).Info("Completed order", zap.String("order", orderID))
	return nil
}
