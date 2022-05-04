package service

type Deps struct {
	AuthStorage     AuthStorage
	TodoListStorage TodoListStorage
	TodoItemStorage TodoItemStorage
}

type Service struct {
	Auth     *AuthService
	TodoList *TodoListService
	TodoItem *TodoItemService
}

func NewService(deps Deps) *Service {
	return &Service{
		Auth:     NewAuthService(deps.AuthStorage),
		TodoList: NewTodoListService(deps.TodoListStorage),
		TodoItem: NewTodoItemService(deps.TodoItemStorage),
	}
}
