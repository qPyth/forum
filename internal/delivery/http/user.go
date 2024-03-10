package http

import (
	"errors"
	"forum/internal/models"
	"forum/internal/service"
	"log"
	"net/http"
	"strings"
)

type AuthResponse struct {
	Success   bool
	ErrorText string `json:"errorText"`
	Referer   string `json:"referer"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		username := strings.TrimSpace(r.Form.Get("username"))
		email := strings.ToLower(strings.TrimSpace(r.Form.Get("email")))
		password := strings.TrimSpace(r.Form.Get("password"))
		passConfirm := strings.TrimSpace(r.Form.Get("password_confirm"))
		response := AuthResponse{Success: true}
		err = signUpValidation(username, email, password, passConfirm)
		if err != nil {
			response.Success = true
			response.ErrorText = err.Error()
			err = jsonResponse(w, response, 400)
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
			}
			return
		}
		inp := service.UserSignUpInput{
			Username: username,
			Email:    email,
			Password: password,
		}
		err = h.services.Users.SignUp(inp)
		if err != nil {
			if errors.Is(err, models.ErrUsernameExist) || errors.Is(err, models.ErrEmailExist) {
				response.Success = false
				response.ErrorText = err.Error()
				err = jsonResponse(w, response, 400)
				if err != nil {
					http.Error(w, http.StatusText(500), 500)
				}
				return
			}
			http.Error(w, http.StatusText(500), 500)
			return
		}
		err = jsonResponse(w, response, 200)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
	default:
		err := htmlResponse(w, "error.html", ErrorResponse{405, http.StatusText(405), nil}, 405)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			h.errorHandler(w, 403, nil)
			return
		}
		referer := r.Header.Get("Referer")
		signInInput := service.UserSignInInput{
			Email:    strings.ToLower(r.Form.Get("email")),
			Password: r.Form.Get("password"),
		}
		token, expiredTime, err := h.services.Users.SignIn(signInInput)
		if err != nil {
			log.Println(err)
			if !errors.Is(err, models.ErrUserNotFound) && !errors.Is(err, models.ErrMissMatchEmailOrPass) {
				h.errorHandler(w, 500, nil)
				return
			}
			if err = jsonResponse(w, AuthResponse{false, err.Error(), ""}, 401); err != nil {
				h.errorHandler(w, 500, nil)
				return
			}
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "auth_token",
			Value:   token,
			Expires: expiredTime,
		})
		err = jsonResponse(w, AuthResponse{true, "", referer}, 200)
		if err != nil {
			h.errorHandler(w, 500, nil)
			return
		}
	default:
		h.errorHandler(w, 405, nil)
		return
	}
}

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(authTokenCookie)
	if err != nil {
		log.Println(err)
		if errors.Is(err, http.ErrNoCookie) {
			h.errorHandler(w, 403, nil)
		} else {
			h.errorHandler(w, 500, nil)
		}
		return
	}
	err = h.services.Users.Logout(cookie.Value)
	if err != nil {
		log.Println(err)
		h.errorHandler(w, 500, nil)
		return
	}
	http.Redirect(w, r, "/", 303)
}
