package restaurant

import (
	//"errors"
	"fmt"
	"net/http"
)

func (h *RestaurantService) updateOrder(w http.ResponseWriter, r *http.Request) {

	orderID := r.URL.Query().Get("id")
	order, ok := h.state.Orders[orderID]
	if !ok {
		http.Error(w, "Order not found: "+orderID, http.StatusNotFound)
		return
	}

	action := r.URL.Query().Get("action")
	if len(action) == 0 {
		http.Error(w, "No update action specified! "+action, http.StatusUnprocessableEntity)
		return
	}

	h.handleAction(r, order, action)
	fmt.Fprintf(w, "%+v", order)
}

func (h *RestaurantService) handleAction(r *http.Request, order *Order, action string) {
	return
}

func getSignalParams(r *http.Request) *SignalParam {
	return &SignalParam{
		WorkflowID: r.URL.Query().Get("workflow_id"),
		RunID:      r.URL.Query().Get("run_id"),
	}
}
