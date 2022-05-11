// Package core represents domain's entities
package core

// User it is an entity of end user for this application
type User struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
