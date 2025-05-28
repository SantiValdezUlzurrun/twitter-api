package dto

type User struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}
