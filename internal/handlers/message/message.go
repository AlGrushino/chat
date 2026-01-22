package message

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AlGrushino/chat/internal/handlers/models"
	"github.com/AlGrushino/chat/internal/service"
	"github.com/sirupsen/logrus"
)

type Message struct {
	service *service.Service
	mux     *http.ServeMux
	log     *logrus.Logger
}

func NewMessage(service *service.Service, mux *http.ServeMux, log *logrus.Logger) *Message {
	return &Message{
		service: service,
		mux:     mux,
		log:     log,
	}
}

func (h *Message) AddMessage(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log := h.log.WithFields(
		logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		},
	)

	log.Info("Incoming request")

	if r.Method != http.MethodPost {
		log.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var req models.CreateMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Warn("Invalid JSON")
		http.Error(w, "InvalidJSON: ", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log = log.WithField("title", req.Text)

	text, err := h.service.AddMessage(r.Context(), id, req.Text)
	if err != nil {
		if "chat does not exist" == err.Error() {
			log.WithError(err).Error("Chat does not exist")
			http.Error(w, "Chat does not exist", http.StatusNotFound)
			return
		}
		log.WithError(err).Error("Service error")
		http.Error(w, "Failed to add message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := models.CreateMessageResponse{
		Status: "success",
		Text:   text,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(resp); err != nil {
		log.WithError(err).Error("Failed to encode response")
	}

	log.WithField("duration_ms", time.Since(start).Milliseconds()).Info("Request completed")
}

func (h *Message) GetMessages(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log := h.log.WithFields(
		logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		},
	)

	log.Info("Incoming request")

	if r.Method != http.MethodGet {
		log.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.WithError(err).Warn("Invalid limit parameter")
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	messages, err := h.service.Message.GetMessages(r.Context(), id, limit)
	if err != nil {
		if "chat does not exist" == err.Error() {
			log.WithError(err).Error("Chat does not exist")
			http.Error(w, "Chat does not exist", http.StatusNotFound)
			return
		}
		log.WithError(err).Error("Service error")
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := models.GetMessagesResponse{
		Status:   "success",
		ID:       id,
		Messages: messages,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(resp); err != nil {
		log.WithError(err).Error("Failed to encode response")
	}

	log.WithField("duration_ms", time.Since(start).Milliseconds()).Info("Request completed")
}
