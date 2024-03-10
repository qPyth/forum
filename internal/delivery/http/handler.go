package http

import (
	"errors"
	"forum/internal/models"
	"forum/internal/service"
	"net/http"
)

type Handler struct {
	services *service.Services
}

type ErrorResponse struct {
	Code int
	Text string
	User *models.User
}

var (
	usernameFormatError = errors.New("username must contain more than 3 Latin letters, numbers")
	emailFormatError    = errors.New("email should be in the format like example@example.com")
	passMatchError      = errors.New("password mismatch")
	passFormatError     = errors.New("the password must contain one capital letter and at least 1 number")
	passLenError        = errors.New("password must be 8 characters or more")
)

func NewHandler(services *service.Services) *Handler {
	return &Handler{services}
}

func (h *Handler) InitRoutes() *http.ServeMux {

	r := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static/", http.StripPrefix("/static", fileServer))

	//home page
	r.HandleFunc("/", h.AuthCheck(h.index))

	// user group routes
	r.HandleFunc("/sign-in", h.SignIn)
	r.HandleFunc("/sign-up", h.SignUp)
	r.HandleFunc("/sign-out", h.SignOut)
	//TODO: add other routes

	// post routes
	r.HandleFunc("/sub-forum/create-post", h.AuthCheck(h.CreatePost))
	r.HandleFunc("/sub-forum/", h.AuthCheck(h.GetPostsByCategory))
	r.HandleFunc("/post", h.AuthCheck(h.PostView))
	r.HandleFunc("/profile/created-posts", h.AuthCheck(h.CreatedPost))

	//comments routes
	r.HandleFunc("/sub-forum/create-comment", h.AuthCheck(h.CreateComment))

	//vote
	r.HandleFunc("/vote", h.AuthCheck(h.Vote))
	r.HandleFunc("/profile/voted-posts", h.AuthCheck(h.VotedPosts))
	return r
}

func (h *Handler) errorHandler(w http.ResponseWriter, code int, user *models.User) {
	err := htmlResponse(w, "error.html", ErrorResponse{code, http.StatusText(code), user}, code)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}
