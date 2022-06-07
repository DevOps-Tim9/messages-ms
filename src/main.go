package main

import (
	"fmt"
	"messages-ms/src/config"
	config_db "messages-ms/src/config/db"
	"messages-ms/src/controller"
	"messages-ms/src/repository"
	"messages-ms/src/route"
	"messages-ms/src/service"
	"net/http"
	"os"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	dataBase, _ := config_db.SetupDB()

	repositoryContainer := initializeRepositories(dataBase)
	serviceContainer := initializeServices(repositoryContainer)
	controllerContainer := initializeControllers(serviceContainer)

	router := route.SetupRoutes(controllerContainer)

	port := os.Getenv("SERVER_PORT")

	http.ListenAndServe(fmt.Sprintf(":%s", port), cors.AllowAll().Handler(router))
}

func initializeControllers(serviceContainer config.ServiceContainer) config.ControllerContainer {
	messageController := controller.NewMessageController(serviceContainer.MessageService)

	container := config.NewControllerContainer(
		messageController,
	)

	return container
}

func initializeServices(repositoryContainer config.RepositoryContainer) config.ServiceContainer {
	messageService := service.MessageService{
		MessageRepository:      repositoryContainer.MessageRepository,
		ConversationRepository: repositoryContainer.ConversationRepository,
	}

	container := config.NewServiceContainer(
		messageService,
	)

	return container
}

func initializeRepositories(dataBase *mongo.Database) config.RepositoryContainer {
	messageRepository := repository.MessageRepository{Database: dataBase}
	conversationRepository := repository.ConversationRepository{Database: dataBase}

	container := config.NewRepositoryContainer(
		messageRepository,
		conversationRepository,
	)

	return container
}
