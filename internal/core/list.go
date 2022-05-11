// Package core represents domain's entities
package core

// ListItem it is an entity that represents M:M relation of Lists to Todos
type ListItem struct {
	ListID int
	ItemID int
}

// UsersList it is an entity that represents M:M relation of Lists to Users
type UsersList struct {
	UserID int
	ListID int
}
