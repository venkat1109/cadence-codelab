package eats

import (
	"go.uber.org/cadence"
	"time"
)

func placeRestaurantOrder(ctx cadence.Context, orderID string, items []string) (time.Duration, error) {
	return time.Minute, nil
}
