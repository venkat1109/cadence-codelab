package courier

import (
	//"errors"
	"fmt"
	"net/http"
)

func (h *CourierService) updateJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("id")
	job, ok := h.DeliveryQueue.Jobs[jobID]
	if !ok {
		http.Error(w, "Order not found: "+jobID, http.StatusNotFound)
		return
	}

	action := r.URL.Query().Get("action")
	if len(action) == 0 {
		http.Error(w, "No update action specified! "+action, http.StatusUnprocessableEntity)
		return
	}

	h.handleAction(r, job, action)
	fmt.Fprintf(w, "%s", job)
}

// handleAction takes the action corresponding to the specified action type
func (h *CourierService) handleAction(r *http.Request, job *DeliveryJob, action string) {
	return
}
