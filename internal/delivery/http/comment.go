package http

import (
	"errors"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		h.errorHandler(w, 403, user)
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, 405, user)
		return
	}
	referer := r.Header.Get("Referer")
	err := r.ParseForm()
	if err != nil {
		h.errorHandler(w, 500, nil)
		return
	}

	comment := r.FormValue("comments")
	postIDstr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		h.errorHandler(w, 400, user)
		return
	}

	err = h.services.Comments.Create(models.Comment{
		UserID:  user.Id,
		Content: comment,
		PostID:  postID,
	})
	if err != nil {
		log.Println(err)
		if errors.Is(err, models.ErrPostNotExists) || errors.Is(err, models.ErrCommentLength) {
			h.errorHandler(w, 400, user)
		} else {
			h.errorHandler(w, 500, user)
		}
		return
	}
	http.Redirect(w, r, referer, 303)
}
