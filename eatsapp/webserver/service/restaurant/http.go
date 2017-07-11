package restaurant

import (
	common "github.com/venkat1109/cadence-codelab/eatsapp/webserver/service"
	"go.uber.org/cadence"
	"net/http"
)

type (

	// RestaurantService implements handlers for requests sent
	// to the restaurant http service
	RestaurantService struct {
		client cadence.Client
		state  RestaurantState
	}

	// RestaurantState models a restaurant order wheel.
	RestaurantState struct {
		menu   *common.Menu
		Orders map[string]*Order
	}

	// Order models a restaurant order.
	Order struct {
		ID           string
		RunID        string
		ShortID      string
		Items        []*common.Item
		TaskToken    []byte
		Status       OrderStatus
		ReadySignal  *SignalParam
		PickUpSignal *SignalParam
	}

	// SignalParam stores the value needed to send a signal to a workflow.
	SignalParam struct {
		WorkflowID string
		RunID      string
	}

	// OrderStatus is the type that represents the status of an eats order
	OrderStatus string
)

// Values representing order status.
const (
	OSPending   OrderStatus = "PENDING"
	OSRejected              = "REJECTED"
	OSPreparing             = "PREPARING"
	OSReady                 = "READY"
	OSSent                  = "SENT"
)

// NewService returns a new instance of the RestaurantService object.
func NewService(c cadence.Client, menuFile string) *RestaurantService {
	menu, err := common.NewMenu(menuFile)
	if err != nil {
		panic("error loading menu file")
	}
	return &RestaurantService{
		client: c,
		state: RestaurantState{
			menu:   menu,
			Orders: make(map[string]*Order),
		},
	}
}

func (h *RestaurantService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.showOrders(w, r)
	case "POST":
		h.addOrder(w, r)
	case "PATCH":
		h.updateOrder(w, r)
	default:
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (h *RestaurantService) GetMenu() *common.Menu {
	return h.state.menu
}
