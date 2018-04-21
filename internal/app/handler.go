package app

import (
	"database/sql"
	"io"
	"net/http"
)

type Handler struct {
	database *sql.DB
	endpoints map[string]func(http.ResponseWriter, *http.Request)
}

func NewHandler(database *sql.DB) (*Handler, error) {
	handler := &Handler{
		database: database,
		endpoints: make(map[string]func(http.ResponseWriter, *http.Request)),
	}

	handler.endpoints["/user"] = handler.userEndpoint
	handler.endpoints["/message"] = handler.messageEndpoint

	return handler, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint, ok := h.endpoints[r.URL.String()]
	if ok {
		endpoint(w, r)
		return
	}
}

func (h *Handler) userEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.createUser(w, r)
	case "PUT":
		h.updateUser(w, r)
	default:
		io.WriteString(w, "unsupported method")
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	// TODO: detect if username already exists
	// TODO: create new record to user table

	// TODO: actually set response correctly
	io.WriteString(w, "POST /user")
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "PUT /user")
}

func (h *Handler) messageEndpoint(w http.ResponseWriter, r *http.Request) {
	// TODO: switch on r.Method
}
