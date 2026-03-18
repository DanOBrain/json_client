package models

// Post - пост
type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Comment - комментарий
type Comment struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

// GetPost - получить пост
func GetPost(userID int, id int, title string, body string) Post {
	return Post{
		UserID: userID,
		ID:     id,
		Title:  title,
		Body:   body,
	}
}

// GetComment - получить комментарий
func GetComment(postID, id int, name string, email string, body string) Comment {
	return Comment{
		PostID: postID,
		ID:     id,
		Name:   name,
		Email:  email,
		Body:   body,
	}
}
