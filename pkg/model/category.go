package model

// Category is the mapping to the category database table.
type Category struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
}
