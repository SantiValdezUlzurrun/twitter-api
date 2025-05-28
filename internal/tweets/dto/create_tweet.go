package dto

type CreateTweet struct {
	UserID  uint   `json:"user_id"`
	Content string `json:"content"`
}
