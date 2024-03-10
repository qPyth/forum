package http

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PostsPageData struct {
	User  *models.User             `json:"user,omitempty"`
	Posts *service.PostsByCategory `json:"posts,omitempty"`
	Tags  []models.Tag
}

type CreatePostData struct {
	User        *models.User `json:"user,omitempty"`
	CategoryURL string       `json:"categoryURL"`
	Tags        []models.Tag
	ErrorText   string `json:"errorText,omitempty"`
}

type PostData struct {
	User *models.User
	Post *service.PostViewData
}

type CreatedPostsData struct {
	User       *models.User
	Posts      []models.Post
	Categories map[int]string
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		h.errorHandler(w, 403, nil)
		return
	}
	switch r.Method {
	case http.MethodGet:
		params := r.URL.Query()
		category := params.Get("category")
		tags, err := h.services.Posts.GetTagsByCategory(category)
		if err != nil {
			log.Println(err)
			if errors.Is(err, models.ErrWrongCategory) {
				h.errorHandler(w, 400, user)
			} else {
				h.errorHandler(w, 500, user)
			}
			return
		}
		err = htmlResponse(w, "create-post.html", CreatePostData{User: user, CategoryURL: category, Tags: tags}, 200)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}
		params := r.URL.Query()
		category := params.Get("category")
		tags, err := h.services.Posts.GetTagsByCategory(category)
		if err != nil {
			log.Println(err)
			if errors.Is(err, models.ErrWrongCategory) {
				h.errorHandler(w, 400, user)
			} else {
				h.errorHandler(w, 500, user)
			}
			return
		}

		var input service.PostCreateInput

		file, FileHeader, err := r.FormFile("image")
		if err != nil {
			if !errors.Is(err, http.ErrMissingFile) {
				h.errorHandler(w, 400, nil)
				return
			}
		} else {
			input.HasImage = true
		}
		url := r.URL.Query()
		categoryUrl := url.Get("category")
		if input.HasImage {
			defer file.Close()
			input.Image = service.Image{
				File:       file,
				FileHeader: FileHeader,
			}
		}
		input.Title = r.Form.Get("title")
		input.Content = r.Form.Get("content")
		input.Tags = r.PostForm["tag"]
		if len(input.Content) == 0 || len(input.Title) == 0 {
			err = htmlResponse(w, "create-post.html", CreatePostData{CategoryURL: categoryUrl, ErrorText: "please write something", Tags: tags}, 400)
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
			}
			log.Println(errors.New("handler CreatePost: empty content or title"))
			return
		}
		input.Category = categoryUrl
		input.UserId = user.Id
		postID, err := h.services.Posts.Create(input)
		if err != nil {
			log.Println(err)
			if errors.Is(err, models.ErrImageSize) ||
				errors.Is(err, models.ErrImageExtension) ||
				errors.Is(err, models.ErrPostTitleLength) ||
				errors.Is(err, models.ErrPostContentLength) ||
				errors.Is(err, models.ErrNotValidTagID) {
				err = htmlResponse(w, "create-post.html", CreatePostData{CategoryURL: categoryUrl, ErrorText: err.Error(), Tags: tags}, 400)
				if err != nil {
					http.Error(w, http.StatusText(500), 500)
				}
			} else if errors.Is(err, models.ErrWrongCategory) {
				h.errorHandler(w, 400, nil)
			} else {
				h.errorHandler(w, 500, nil)
			}
			return
		}
		postUrl := fmt.Sprintf("/post?category=%s&post_id=%d", categoryUrl, postID)
		http.Redirect(w, r, postUrl, 303)
	default:
		h.errorHandler(w, 405, nil)
		return
	}
}

func (h *Handler) PostView(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		user = nil
	}

	params := r.URL.Query()
	if strings.HasPrefix(params.Get("post_id"), "0") {
		h.errorHandler(w, 404, nil)
		return
	}
	postID, err := strconv.Atoi(params.Get("post_id"))
	if err != nil {
		h.errorHandler(w, 404, nil)
		return
	}
	post, err := h.services.Posts.GetByID(postID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, models.ErrNoPost) {
			h.errorHandler(w, 404, nil)
			return
		}
		if !errors.Is(err, models.ErrCommentsNotFound) {
			h.errorHandler(w, 500, nil)
			return
		}
		return
	}

	err = htmlResponse(w, "post.html", PostData{user, post}, 200)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func (h *Handler) GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		user = nil
	}
	if r.Method != http.MethodGet {
		h.errorHandler(w, 405, user)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	url := parts[2]
	tagsFromQuery := r.URL.Query()["tag"]
	posts, err := h.services.Posts.GetPostsByCategory(url, tagsFromQuery)
	if err != nil {
		if errors.Is(err, models.ErrNotValidTagID) {
			h.errorHandler(w, 400, user)
			return
		}
		h.errorHandler(w, 500, user)
		return
	}
	tags, err := h.services.Posts.GetTagsByCategory(url)
	if err != nil {
		log.Println(err)
		if errors.Is(err, models.ErrWrongCategory) {
			h.errorHandler(w, 404, nil)
		} else {
			h.errorHandler(w, 500, nil)
		}
		return
	}
	err = htmlResponse(w, "sub-forum.html", PostsPageData{User: user, Posts: posts, Tags: tags}, 200)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func (h *Handler) CreatedPost(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		h.errorHandler(w, 403, user)
		return
	}
	if r.Method != "GET" {
		h.errorHandler(w, 405, user)
		return
	}
	if r.URL.Path != "/profile/created-posts" {
		h.errorHandler(w, 404, user)
		return
	}

	posts, categories, err := h.services.Posts.GetCreatedPosts(user.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = htmlResponse(w, "created-posts.html", CreatedPostsData{User: user, Posts: posts, Categories: categories}, 200)
}

func (h *Handler) VotedPosts(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		h.errorHandler(w, 403, user)
		return
	}
	if r.Method != "GET" {
		h.errorHandler(w, 405, user)
		return
	}
	if r.URL.Path != "/profile/voted-posts" {
		h.errorHandler(w, 404, user)
		return
	}
	posts, categories, err := h.services.Posts.GetVotedPosts(user.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = htmlResponse(w, "voted-posts.html", CreatedPostsData{User: user, Posts: posts, Categories: categories}, 200)
}
