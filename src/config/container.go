package config

import (
	"messages-ms/src/controller"
	"messages-ms/src/repository"
	"messages-ms/src/service"
)

type ControllerContainer struct {
	MessageController controller.MessageController
}

type ServiceContainer struct {
	MessageService service.IMessageService
}

type RepositoryContainer struct {
	MessageRepository      repository.IMessageRepository
	ConversationRepository repository.IConversationRepository
}

func NewControllerContainer(
	messageController controller.MessageController,

) ControllerContainer {
	return ControllerContainer{
		MessageController: messageController,
	}
}

func NewServiceContainer(messageService service.IMessageService) ServiceContainer {
	return ServiceContainer{
		MessageService: messageService,
	}
}

func NewRepositoryContainer(
	messageRepository repository.IMessageRepository,
	conversationRepository repository.IConversationRepository,
) RepositoryContainer {
	return RepositoryContainer{
		MessageRepository:      messageRepository,
		ConversationRepository: conversationRepository,
	}
}
