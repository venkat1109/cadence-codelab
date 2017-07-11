package main

import (
	"fmt"
	"net/http"

	"github.com/venkat1109/cadence-codelab/common"
	"github.com/venkat1109/cadence-codelab/eatsapp/webserver/service"
	"github.com/venkat1109/cadence-codelab/eatsapp/webserver/service/courier"
	"github.com/venkat1109/cadence-codelab/eatsapp/webserver/service/eats"
	"github.com/venkat1109/cadence-codelab/eatsapp/webserver/service/restaurant"
)

func main() {

	runtime := common.NewRuntime()
	workflowClient, err := runtime.Builder.BuildCadenceClient()
	if err != nil {
		panic(err)
	}

	service.LoadTemplates()

	restaurant := restaurant.NewService(workflowClient, "eatsapp/webserver/assets/data/menu.yaml")

	http.Handle("/restaurant", restaurant)
	http.Handle("/courier", courier.NewService(workflowClient))
	http.Handle("/eats-orders", eats.NewService(workflowClient, restaurant.GetMenu()))
	http.Handle("/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/eats-menu", func(w http.ResponseWriter, r *http.Request) {
		service.ViewHandler(w, r, restaurant.GetMenu())
	})
	// setup & start server
	http.HandleFunc("/bistro", func(w http.ResponseWriter, r *http.Request) {
		service.ViewHandler(w, r, nil)
	})

	fmt.Println("Starting Webserver")
	http.ListenAndServe(":8090", nil)
}
