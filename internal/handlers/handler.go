package handlers

import (
	"net/http"

	"github.com/AlGrushino/chat/internal/handlers/chat"
	"github.com/AlGrushino/chat/internal/handlers/message"
	"github.com/AlGrushino/chat/internal/service"
	"github.com/sirupsen/logrus"
)

type Chat interface {
	CreateChat(w http.ResponseWriter, r *http.Request)
}

type Message interface {
	AddMessage(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	service *service.Service
	chat    Chat
	message Message
	log     *logrus.Logger
	mux     *http.ServeMux
}

func NewHandler(service *service.Service, log *logrus.Logger) *Handler {
	log.WithFields(logrus.Fields{
		"layer":  "handler",
		"method": "NewHandler",
	}).Info("Create new handler")

	mux := http.NewServeMux()

	chatHandler := chat.NewChat(service, mux, log)
	messageHandler := message.NewMessage(service, mux, log)

	return &Handler{
		service: service,
		chat:    chatHandler,
		message: messageHandler,
		log:     log,
		mux:     mux,
	}
}

func (h *Handler) InitRoutes() {
	h.log.WithFields(logrus.Fields{
		"layer":  "handler",
		"method": "InitRoutes",
	}).Info("Initing routes")

	h.mux.HandleFunc("POST /chats", h.chat.CreateChat)
	h.mux.HandleFunc("POST /chats/{id}/messages", h.message.AddMessage)
	h.mux.HandleFunc("GET /chats/{id}", h.message.GetMessages)

	h.log.Info("Routes initialized successfully")
}

func (h *Handler) RunServer(addr string) error {
	h.log.WithField("address", addr).Info("Starting HTTP server")

	server := &http.Server{
		Addr:    addr,
		Handler: h.mux,
	}

	return server.ListenAndServe()
}

func (h *Handler) GetMux() *http.ServeMux {
	return h.mux
}
