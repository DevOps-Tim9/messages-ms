package controller

import (
	"bytes"
	"encoding/json"
	"messages-ms/src/dto"
	"messages-ms/src/entity"
	"messages-ms/src/service"
	"messages-ms/src/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v8"
)

type MessageController struct {
	MessageService service.IMessageService
	validate       *validator.Validate
	logger         *logrus.Entry
}

func NewMessageController(messageService service.IMessageService) MessageController {
	config := &validator.Config{TagName: "validate"}
	logger := utils.Logger()

	return MessageController{MessageService: messageService, validate: validator.New(config), logger: logger}
}

func (c MessageController) CreateNewMessage(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /api/messages")

	defer span.Finish()

	c.logger.Info("Creating new message request received")

	var messageDto dto.MessageDto

	json.NewDecoder(r.Body).Decode(&messageDto)

	message, error := c.MessageService.CreateNewMessage(messageDto, ctx)

	if error != nil {
		AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), "New message created unsuccessfully")

		handleMessageError(error, w)

		return
	}

	c.logger.Info("Message created successfully")

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), "New message created successfully")

	payload, _ := json.Marshal(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))
}

func (c MessageController) GetMesssagesByConversation(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /messages/conversations/{conversation}")

	defer span.Finish()

	c.logger.Info("getting messages for specified conversation request received")

	params := mux.Vars(r)

	conversation := params["conversation"]

	if conversation == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	messages := c.MessageService.GetMesssagesByConversation(conversation, ctx)

	if messages == nil {
		messages = []entity.Message{}
	}

	payload, _ := json.Marshal(messages)

	c.logger.Info("Returning found messages")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (c MessageController) GetConversationsByUser(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /messages/users/{user}")

	defer span.Finish()

	c.logger.Info("Getting conversations for specified user request received")

	params := mux.Vars(r)

	user, err := strconv.Atoi(params["user"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	conversations := c.MessageService.GetConversationsByUser(uint(user), ctx)

	if conversations == nil {
		conversations = []entity.Conversation{}
	}

	payload, _ := json.Marshal(conversations)

	c.logger.Info("Returning found conversations")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func handleMessageError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusInternalServerError)

	return w
}

func AddSystemEvent(time string, message string) error {
	logger := utils.Logger()
	event := dto.EventRequestDTO{
		Timestamp: time,
		Message:   message,
	}

	b, _ := json.Marshal(&event)
	endpoint := os.Getenv("EVENTS_MS")
	logger.Info("Sending system event to events-ms")
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Debug("Error happened during sending system event")
		return err
	}

	return nil
}
