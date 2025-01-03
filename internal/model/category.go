package model

type Category struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Tasks    []Todo `json:"tasks,omitempty"`
}
