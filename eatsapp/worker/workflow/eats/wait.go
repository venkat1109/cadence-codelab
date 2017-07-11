package eats

import (
	"errors"
	"go.uber.org/cadence"
	"time"
)

func waitForRestaurant(ctx cadence.Context, signalName string, eta time.Duration) error {
	return errors.New("not implemented")
}
