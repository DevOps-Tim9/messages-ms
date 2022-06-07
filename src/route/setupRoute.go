package route

import (
	"messages-ms/src/config"

	"github.com/gorilla/mux"
)

func SetupRoutes(container config.ControllerContainer) *mux.Router {
	route := mux.NewRouter()

	routerWithApiAsPrefix := route.PathPrefix("/api").Subrouter()

	routerWithApiAsPrefix.HandleFunc("/messages", container.MessageController.CreateNewMessage).Methods("POST")
	routerWithApiAsPrefix.HandleFunc("/messages/users/{user}", container.MessageController.GetConversationsByUser).Methods("GET")
	routerWithApiAsPrefix.HandleFunc("/messages/conversations/{conversation}", container.MessageController.GetMesssagesByConversation).Methods("GET")

	return routerWithApiAsPrefix
}
