package service

import (
	"errors"
	"fmt"
	"forum/internal/models"
	s "forum/internal/storage"
	"io"
	"os"
	"strconv"
	"strings"
	time2 "time"
)

type PostService struct {
	postStorage    s.Posts
	userStorage    s.Users
	commentStorage s.Comments
}

const (
	ContentSize  int64 = 20 * 1024 * 1024
	pathToImages       = "ui/static/images/post_images/"
)

func NewPostService(storage *s.Storages) *PostService {
	return &PostService{storage.Posts, storage.Users, storage.Comments}
}

func (p *PostService) Create(input PostCreateInput) (int, error) {
	// validating downloading file
	const op = "PostService.Create"
	var (
		flagTags bool
		imgPath  string
	)
	if input.HasImage {
		if input.Image.FileHeader.Size > ContentSize {
			return 0, models.ErrImageSize
		}
		if input.Image.FileHeader.Header["Content-Type"][0] != "image/jpeg" && input.Image.FileHeader.Header["Content-Type"][0] != "image/png" && input.Image.FileHeader.Header["Content-Type"][0] != "image/gif" {
			return 0, models.ErrImageExtension
		}

		// save file
		err := os.Mkdir(pathToImages, 0777)
		if err != nil {
			return 0, err
		}
		imageBytes, err := io.ReadAll(input.Image.File)
		if err != nil {
			return 0, fmt.Errorf("%s - read file error:%w", op, err)
		}
		t := time2.StampMilli
		imgPath = fmt.Sprintf(pathToImages + input.Title + input.Image.FileHeader.Filename + t)
		err = os.WriteFile(imgPath, imageBytes, 0666)
		if err != nil {
			return 0, fmt.Errorf("%s - create file error:%w", op, err)
		}
	}

	tagsFromCategory, err := p.GetTagsByCategory(input.Category)
	if err != nil {
		return 0, err
	}

	if len(input.Title) > 70 {
		return 0, models.ErrPostTitleLength
	} else if len(input.Content) > 500 {
		return 0, models.ErrPostContentLength
	}
	if len(input.Tags) == 0 {
		return 0, models.ErrNotValidTagID
	}
	tags := make([]int, len(input.Tags))
	for i, id := range input.Tags {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return 0, models.ErrNotValidTagID
		}
		tags[i] = idInt

		for i := range tagsFromCategory {
			if tagsFromCategory[i].ID == idInt {
				flagTags = true
			}
		}
	}

	if !flagTags {
		return 0, models.ErrNotValidTagID
	}

	category, err := p.postStorage.GetCategoryByUrl(input.Category)
	if err != nil {
		return 0, err
	}

	postModel := models.Post{
		Title:        input.Title,
		Content:      input.Content,
		Image:        strings.TrimPrefix(imgPath, "ui"),
		CategoryID:   category.ID,
		UserID:       input.UserId,
		Tags:         tags,
		LikeCount:    0,
		DislikeCount: 0,
	}

	postID, err := p.postStorage.Create(postModel)
	if err != nil {
		return 0, err
	}
	return postID, nil
}

func (p *PostService) GetCategoriesWithInfo() (Categories, error) {
	var c Categories
	cats, err := p.postStorage.GetAllCategories()
	if err != nil {
		return Categories{}, err
	}

	for _, cat := range cats {
		var category Category
		category.Title = cat.Name
		category.Description = cat.Description
		category.Url = cat.Url
		count, err := p.postStorage.PostsCountByCategory(cat.ID)
		if err != nil {
			return Categories{}, fmt.Errorf("PostService.GetCategoriesWithInfo - posts count get error: %w", err)
		}
		category.PostsCount = count

		latestPost, err := p.postStorage.GetLastPostByCategory(cat.ID)
		if err != nil {
			if errors.Is(err, models.ErrPostsByCatNotFound) {
				category.LatestPost = Post{}
			} else {
				return Categories{}, err
			}
		} else {
			author, err := p.userStorage.GetByID(latestPost.UserID)
			if err != nil {
				return Categories{}, err
			}
			category.LatestPost.Author = author.Username
			category.LatestPost.Title = latestPost.Title
			category.LatestPost.CreatedDate = latestPost.CreatedDate
		}
		c.Categories = append(c.Categories, category)
	}
	return c, nil
}

func (p *PostService) GetPostsByCategory(url string, tagsIDs []string) (*PostsByCategory, error) {
	posts := PostsByCategory{}
	category, err := p.postStorage.GetCategoryByUrl(url)
	if err != nil {
		return nil, err
	}
	posts.Category = category
	var errFromRepo error
	if len(tagsIDs) != 0 {
		posts.Posts, errFromRepo = p.GetByTags(tagsIDs)
	} else {
		posts.Posts, errFromRepo = p.postStorage.GetPostsByCategory(category.Name)
	}
	if errFromRepo != nil {
		if errors.Is(err, models.ErrPostsByCatNotFound) || errors.Is(err, models.ErrNotValidTagID) {
			posts.Posts = nil
		} else {
			return nil, err
		}
	} else {
		posts.CommentsCount, err = p.commentStorage.GetCountOnPost(category.ID)
		if err != nil {
			return nil, err
		}
		posts.Authors, err = p.userStorage.GetAuthorsOfPost(category.ID)
		if err != nil {
			return nil, err
		}
		posts.Tags, err = p.postStorage.GetTagsByCategoryID(category.ID)
		if err != nil {
			return nil, err
		}
	}

	return &posts, nil
}

func (p *PostService) GetByID(postID int) (*PostViewData, error) {
	var postView PostViewData

	post, err := p.postStorage.GetPostByID(postID)
	if err != nil {
		return nil, err
	}
	comments, err := p.commentStorage.GetByPostID(postID)
	if err != nil {
		return nil, err
	}

	authors, err := p.commentStorage.GetAuthors(postID)
	if err != nil {
		return nil, err
	}
	postView.Tags, err = p.postStorage.GetTagsByPostID(postID)
	if err != nil {
		return nil, err
	}
	postAuthor, err := p.userStorage.GetByID(post.UserID)
	if err != nil {
		return nil, err
	}
	postView.Post = post
	postView.Author = postAuthor.Username
	postView.Comments = comments
	postView.Authors = authors
	return &postView, err
}

func (p *PostService) GetCreatedPosts(userID int) ([]models.Post, map[int]string, error) {
	posts, err := p.postStorage.GetPostsByUserID(userID)
	if err != nil {
		return nil, nil, err
	}
	categories, err := p.postStorage.GetMapCategories()
	if err != nil {
		return nil, nil, err
	}
	return posts, categories, nil
}

func (p *PostService) GetVotedPosts(userID int) ([]models.Post, map[int]string, error) {
	categories, err := p.postStorage.GetMapCategories()
	if err != nil {
		return nil, nil, err
	}
	posts, err := p.postStorage.GetVotesPosts(userID)
	if err != nil {
		return nil, nil, err
	}
	return posts, categories, nil
}

func (p *PostService) GetTagsByCategory(categoryUrl string) ([]models.Tag, error) {
	category, err := p.postStorage.GetCategoryByUrl(categoryUrl)
	if err != nil {
		return nil, err
	}
	tags, err := p.postStorage.GetTags(category.ID)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (p *PostService) GetByTags(tagsIDs []string) ([]models.Post, error) {
	tags := make([]int, len(tagsIDs))
	for i, id := range tagsIDs {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return nil, models.ErrNotValidTagID
		}
		tags[i] = idInt
	}

	posts, err := p.postStorage.GetByTags(tags)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
