package main

import (
	"fmt"
	"messages-ms/src/config"
	config_db "messages-ms/src/config/db"
	setupJaeger "messages-ms/src/config/jaeger"
	"messages-ms/src/controller"
	"messages-ms/src/rabbitmq"
	"messages-ms/src/repository"
	"messages-ms/src/route"
	"messages-ms/src/service"
	"messages-ms/src/utils"
	"net/http"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/cors"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	logger := utils.Logger()

	logger.Info("Connecting with DB")

	dataBase, _ := config_db.SetupDB()

	tracer, trCloser, err := setupJaeger.InitJaeger()

	if err != nil {
		logger.Debug(err.Error())

		fmt.Printf("error init jaeger %v", err)
	} else {
		defer trCloser.Close()

		opentracing.SetGlobalTracer(tracer)
	}

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	rabbit := rabbitmq.RMQProducer{
		ConnectionString: amqpServerURL,
	}

	channel, _ := rabbit.StartRabbitMQ()

	defer channel.Close()

	repositoryContainer := initializeRepositories(dataBase)
	serviceContainer := initializeServices(repositoryContainer, channel)
	controllerContainer := initializeControllers(serviceContainer)

	router := route.SetupRoutes(controllerContainer)

	port := os.Getenv("SERVER_PORT")

	logger.Info("Starting server")

	http.ListenAndServe(fmt.Sprintf(":%s", port), cors.AllowAll().Handler(router))
}

func initializeControllers(serviceContainer config.ServiceContainer) config.ControllerContainer {
	messageController := controller.NewMessageController(serviceContainer.MessageService)

	container := config.NewControllerContainer(
		messageController,
	)

	return container
}

func initializeServices(repositoryContainer config.RepositoryContainer, channel *amqp.Channel) config.ServiceContainer {
	messageService := service.MessageService{
		MessageRepository:      repositoryContainer.MessageRepository,
		ConversationRepository: repositoryContainer.ConversationRepository,
		Logger:                 utils.Logger(),
		RabbitMQChannel:        channel,
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
