package chat

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AlGrushino/chat/internal/handlers/models"
	"github.com/AlGrushino/chat/internal/service"
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

func (h *Chat) DeleteChat(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log := h.log.WithFields(
		logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		},
	)

	log.Info("Incoming request")

	if r.Method != http.MethodDelete {
		log.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid chat ID")
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteChat(r.Context(), id)
	if err != nil {
		log.WithError(err).Warn("Failed to delete chat")
		http.Error(w, "Failed to delete chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	log.WithField("duration_ms", time.Since(start).Milliseconds()).Info("Request completed")
}
