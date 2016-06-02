package main

import (
	"fmt"
	"net/http"
)

type Handler struct {
	hook      string
	auth      *Auth
	commander *Commander
}

func NewHandler(hook string, auth *Auth, commander *Commander) *Handler {
	return &Handler{hook: hook, auth: auth, commander: commander}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 page not found", 404)
		return
	}

	if h.auth != nil {
		if user, pass, ok := r.BasicAuth(); !ok || !h.auth.Check(user, pass) {
			w.Header().Set("WWW-Authenticate",
				fmt.Sprintf("Basic realm=\"%s\"", h.hook))
			http.Error(w, "401 Unauthorized", 401)
			return
		}
	}

	h.commander.Exec()
	w.WriteHeader(204)
}
