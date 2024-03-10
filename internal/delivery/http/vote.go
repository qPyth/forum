package http

import (
	"encoding/json"
	"errors"
	"forum/internal/models"
	"forum/internal/service"
	"log"
	"net/http"
)

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		http.Error(w, http.StatusText(403), 403)
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, 405, user)
		return
	}
	var v service.VoteInput

	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		log.Println(err.Error())
		h.errorHandler(w, 400, user)
		return
	}
	v.UserID = user.Id
	output, err := h.services.Votes.MakeVote(v)
	if err != nil {
		log.Println(err)
		if errors.Is(err, models.ErrWrongAction) ||
			errors.Is(err, models.ErrWrongVoteItem) {
			http.Error(w, http.StatusText(400), 400)
		} else {
			http.Error(w, http.StatusText(500), 500)
		}
	}
	err = jsonResponse(w, output, 200)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
