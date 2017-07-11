package courier

import (
	"net/http"

	"go.uber.org/cadence"
)

type (
	// JobStatus is the custom type to record status of a job
	JobStatus string

	// DeliveryJob is the struct storing metadata about a delivery job
	DeliveryJob struct {
		OrderID          string
		Status           JobStatus
		AcceptTaskToken  []byte
		PickupTaskToken  []byte
		CompletTaskToken []byte
	}

	// DeliveryQueue is the struct modeling the list of jobs to be delivered.
	DeliveryQueue struct {
		Jobs map[string]*DeliveryJob
	}

	// CourierService implements the handlers for requests
	// sent to the courier http service
	CourierService struct {
		client        cadence.Client
		DeliveryQueue DeliveryQueue
	}
)

const (
	djPending   JobStatus = "PENDING"
	djRejected            = "REJECTED"
	djAccepted            = "ACCEPTED"
	djPickedUp            = "PICKED_UP"
	djCompleted           = "COMPLETED"
)

// NewService returns a new instance of the CourierService object.
func NewService(c cadence.Client) *CourierService {
	return &CourierService{
		client: c,
		DeliveryQueue: DeliveryQueue{
			Jobs: make(map[string]*DeliveryJob),
		},
	}
}

func (h *CourierService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.showJobs(w, r)
	case "POST":
		h.addJob(w, r)
	case "PATCH":
		h.updateJob(w, r)
	default:
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
