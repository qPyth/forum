package http

import (
	"context"
	"errors"
	"forum/internal/models"
	"net/http"
)

func (h *Handler) AuthCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		authToken, err := r.Cookie(authTokenCookie)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if authToken != nil {
			user, _ = h.services.Users.GetByToken(authToken.Value)
		}
		if user != nil {
			_, ok, err := h.services.Users.CheckSession(user.Id)
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}
			if !ok {
				user = nil
			} else {
				ctx := context.WithValue(r.Context(), "user", user)
				next(w, r.WithContext(ctx))
			}
		} else {
			next(w, r)
		}
	}
}
