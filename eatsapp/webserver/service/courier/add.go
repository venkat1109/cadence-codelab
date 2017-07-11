package courier

import (
	common "github.com/venkat1109/cadence-codelab/eatsapp/webserver/service"
	"net/http"
)

func (h *CourierService) addJob(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// create order object
	job := DeliveryJob{
		OrderID:         r.Form.Get("id"),
		AcceptTaskToken: []byte(r.Form.Get("task_token")),
		Status:          djPending,
	}

	// store order
	h.DeliveryQueue.Jobs[job.OrderID] = &job
	common.ViewHandler(w, r, &h.DeliveryQueue)
}
