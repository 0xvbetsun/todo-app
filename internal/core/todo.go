// Package core represents domain's entities
package core

// Todolist it is an entity that represents user's list of todos
type Todolist struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// TodoItem it is an entity that represents user's single todo
type TodoItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// UpdateListData it is a DTO for passing data to the List service layer
type UpdateListData struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

// UpdateItemData it is a DTO for passing data to the Todo service layer
type UpdateItemData struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
