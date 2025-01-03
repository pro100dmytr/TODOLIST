package model

type Todo struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Completed  bool   `json:"completed"`
	CategoryID *int   `json:"category_id,omitempty"`
	UserID     int    `json:"user_id"`
}
