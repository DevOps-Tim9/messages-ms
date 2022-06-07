package controller

import (
	"encoding/json"
	"messages-ms/src/dto"
	"messages-ms/src/entity"
	"messages-ms/src/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v8"
)

type MessageController struct {
	MessageService service.IMessageService
	validate       *validator.Validate
}

func NewMessageController(messageService service.IMessageService) MessageController {
	config := &validator.Config{TagName: "validate"}

	return MessageController{MessageService: messageService, validate: validator.New(config)}
}

func (c MessageController) CreateNewMessage(w http.ResponseWriter, r *http.Request) {
	var messageDto dto.MessageDto

	json.NewDecoder(r.Body).Decode(&messageDto)

	message, error := c.MessageService.CreateNewMessage(messageDto)

	if error != nil {
		handleMessageError(error, w)

		return
	}

	payload, _ := json.Marshal(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))
}

func (c MessageController) GetMesssagesByConversation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	conversation := params["conversation"]

	if conversation == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	messages := c.MessageService.GetMesssagesByConversation(conversation)

	if messages == nil {
		messages = []entity.Message{}
	}

	payload, _ := json.Marshal(messages)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (c MessageController) GetConversationsByUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	user, err := strconv.Atoi(params["user"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	conversations := c.MessageService.GetConversationsByUser(uint(user))

	if conversations == nil {
		conversations = []entity.Conversation{}
	}

	payload, _ := json.Marshal(conversations)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func handleMessageError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusInternalServerError)

	return w
}
