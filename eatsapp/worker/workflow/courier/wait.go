package courier

import (
	"errors"
	"go.uber.org/cadence"
)

func waitForRestaurantPickupConfirmation(ctx cadence.Context, signalName string) error {
	return errors.New("not implemented")
}
