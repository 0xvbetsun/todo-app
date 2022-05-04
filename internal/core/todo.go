package core

type Todolist struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type UpdateListData struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type UpdateItemData struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
