package restaurant

import (
	common "github.com/venkat1109/cadence-codelab/eatsapp/webserver/service"
	"net/http"
)

func (h *RestaurantService) addOrder(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if len(r.Form["item"]) == 0 {
		http.Error(w, "Order constains no items!", http.StatusUnprocessableEntity)
		return
	}

	// create order object
	order := Order{
		ID:        r.Form.Get("id"),
		ShortID:   r.Form.Get("id"),
		TaskToken: []byte(r.Form.Get("task_token")),
		Status:    OSPending,
		ReadySignal: &SignalParam{
			WorkflowID: r.Form.Get("id"),
			RunID:      r.Form.Get("run_id"),
		},
	}
	for _, v := range r.Form["item"] {
		item, err := h.state.menu.GetItemByID(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		order.Items = append(order.Items, item)
	}

	// store order
	h.state.Orders[order.ID] = &order
	common.ViewHandler(w, r, &h.state)
}
