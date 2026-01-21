package chat

import (
	"chat/internal/handlers/models"
	"chat/internal/service"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Chat struct {
	service *service.Service
	mux     *http.ServeMux
	log     *logrus.Logger
}

func NewChat(service *service.Service, mux *http.ServeMux, log *logrus.Logger) *Chat {
	return &Chat{
		service: service,
		mux:     mux,
		log:     log,
	}
}

func (h *Chat) CreateChat(w http.ResponseWriter, r *http.Request) {
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

	var req models.CreateChat
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Warn("Invalid JSON")
		http.Error(w, "InvalidJSON: ", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log = log.WithField("title", req.Title)

	title, err := h.service.Chat.CreateChat(r.Context(), req.Title)
	if err != nil {
		log.WithError(err).Error("Service error")
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := models.CreateChatReposnse{
		Status: "success",
		Title:  title,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(resp); err != nil {
		log.WithError(err).Error("Failed to encode response")
	}

	log.WithField("duration_ms", time.Since(start).Milliseconds()).Info("Request completed")
}
