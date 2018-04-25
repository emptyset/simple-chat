package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/emptyset/simple-chat/internal/models"
)

// Handler routes the endpoints to specific handlers
type Handler struct {
	model models.ChatModel
}

// NewHandler returns a Handler type based on the ChatModel
func NewHandler(model models.ChatModel) (*Handler, error) {
	handler := &Handler{
		model: model,
	}

	return handler, nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("Username")
	password := r.Header.Get("Password")

	log.Debug("invoking model SignupUser")
	user, err := h.model.SignupUser(username, password)
	if err != nil {
		log.Errorf("error when signing up user: %s", err)
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	senderID, err := strconv.Atoi(query.Get("s"))
	if err != nil {
		log.Errorf("error when converting sender id from request: %s", err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	recipientID, err := strconv.Atoi(query.Get("r"))
	if err != nil {
		log.Errorf("error when converting recipient id from request: %s", err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	mediaType := query.Get("t")

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("error when reading body of request: %s", err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	log.Debug("invoking model SendMessage")
	err = h.model.SendMessage(senderID, recipientID, string(content), mediaType)
	if err != nil {
		log.Errorf("error when sending message: %s", err)
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ReadMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	senderID, err := strconv.Atoi(query.Get("s"))
	if err != nil {
		log.Errorf("error when converting sender id from request: %s", err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	recipientID, err := strconv.Atoi(query.Get("r"))
	if err != nil {
		log.Errorf("error when converting recipient id from request: %s", err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	var count int
	rawCount := query.Get("c")
	if rawCount == "" {
		count = 100
	} else {
		count, err = strconv.Atoi(rawCount)
		if err != nil {
			log.Errorf("error when converting count from request: %s", err)
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}
	}

	var offset int
	rawOffset := query.Get("o")
	if rawOffset == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(rawOffset)
		if err != nil {
			log.Errorf("error when converting offset from request: %s", err)
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
			return
		}
	}

	log.Debug("invoking model GetMessages")
	messages, err := h.model.GetMessages(senderID, recipientID, count, offset)
	if err != nil {
		log.Errorf("error when getting messages: %s", err)
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
