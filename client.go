package client

import (
	"encoding/json"
	"fmt"
	"json-client/models"
	"net/http"
)

// Client - HTTP клиент
type Client struct {
	url  string
	http *http.Client
}

// NewClient - создает клиента
func NewClient() *Client {
	return &Client{
		url:  "https://jsonplaceholder.typicode.com",
		http: &http.Client{},
	}
}

// GetPosts - получение постов пользователя
func (c *Client) GetPosts(userID int) ([]models.Post, error) {
	// запрос
	resp, err := c.http.Get(fmt.Sprintf("%s/posts?userId=%d", c.url, userID))
	if err != nil {
		return nil, fmt.Errorf("ошибка: %v", err)
	}
	defer resp.Body.Close()

	// статус
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("статус %d", resp.StatusCode)
	}

	// тело JSON
	var posts []models.Post
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, fmt.Errorf("ошибка: %v", err)
	}

	return posts, nil
}

// GetComments - получение комментариев к посту
func (c *Client) GetComments(postID int) ([]models.Comment, error) {
	// запрос
	resp, err := c.http.Get(fmt.Sprintf("%s/posts/%d/comments", c.url, postID))
	if err != nil {
		return nil, fmt.Errorf("ошибка: %v", err)
	}
	defer resp.Body.Close()

	// статус
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("статус 200")
	}

	// тело JSON
	var comments []models.Comment
	err = json.NewDecoder(resp.Body).Decode(&comments)
	if err != nil {
		return nil, fmt.Errorf("ошибка: %v", err)
	}

	return comments, nil
}
