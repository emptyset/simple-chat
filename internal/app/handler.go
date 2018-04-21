package app

import (
	"io/ioutil"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"github.com/emptyset/simple-chat/internal/models"
)

type Handler struct {
	model models.ChatModel
	endpoints map[string]func(http.ResponseWriter, *http.Request)
}

func NewHandler(model models.ChatModel) (*Handler, error) {
	handler := &Handler{
		model: model,
		endpoints: make(map[string]func(http.ResponseWriter, *http.Request)),
	}

	handler.endpoints["/user"] = handler.userEndpoint
	handler.endpoints["/message"] = handler.messageEndpoint

	return handler, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"request": r.URL,
	}).Debug("incoming request")
	endpoint, ok := h.endpoints[r.URL.String()]
	if ok {
		endpoint(w, r)
		return
	} else {
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
	}
}

func (h *Handler) userEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Debug("invoking createUser handler")
		h.createUser(w, r)
	default:
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) messageEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Debug("invoking createMessage handler")
		h.createMessage(w, r)
	case "GET":
		log.Debug("invoking getMessages handler")
		h.getMessages(w, r)
	default:
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
	}
}

func (h *Handler) createMessage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	senderId, err := strconv.Atoi(query.Get("s"))
	if err != nil {
		log.Errorf("error when converting sender id from request: %s", err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	recipientId, err := strconv.Atoi(query.Get("r"))
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
	err = h.model.SendMessage(senderId, recipientId, string(content), mediaType)
	if err != nil {
		log.Errorf("error when sending message: %s", err)
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	senderId, err := strconv.Atoi(query.Get("s"))
	if err != nil {
		log.Errorf("error when converting sender id from request: %s", err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	recipientId, err := strconv.Atoi(query.Get("r"))
	if err != nil {
		log.Errorf("error when converting recipient id from request: %s", err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	count, err := strconv.Atoi(query.Get("c"))
	if err != nil {
		log.Errorf("error when converting count from request: %s", err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	offset, err := strconv.Atoi(query.Get("o"))
	if err != nil {
		log.Errorf("error when converting offset from request: %s", err)
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	log.Debug("invoking model GetMessages")
	messages, err := h.model.GetMessages(senderId, recipientId, count, offset)
	if err != nil {
		log.Errorf("error when getting messages: %s", err)
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
