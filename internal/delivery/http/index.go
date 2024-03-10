package http

import (
	"forum/internal/models"
	"forum/internal/service"
	"log"
	"net/http"
)

type indexData struct {
	User       *models.User       `json:"user"`
	ErrorText  string             `json:"error_text" json:"error_text,omitempty"`
	Categories service.Categories `json:"categories"`
}

const (
	authTokenCookie = "auth_token"
)

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	var err error
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		user = nil
	}
	if r.URL.Path != "/" {
		h.errorHandler(w, 404, user)
		return
	}
	if r.Method != http.MethodGet {
		h.errorHandler(w, 405, user)
		return
	}

	categories, err := h.services.Posts.GetCategoriesWithInfo()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := htmlResponse(w, "index.html", indexData{user, "", categories}, 200); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
