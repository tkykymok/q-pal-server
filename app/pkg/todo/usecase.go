package todo

import (
	"app/api/presenter"
	"app/api/requests"
	"app/pkg/broadcast"
	"app/pkg/models"
	"app/pkg/outputs"
	"context"
)

type Usecase interface {
	FetchAllTodos(ctx context.Context) (*[]outputs.Todo, error)
	FetchTodosWithRelated(ctx context.Context, request *requests.GetTodosWithRelated) (*[]outputs.TodoWithRelated, error)
	FetchTodoById(ctx context.Context, id int) (*outputs.Todo, error)
	InsertTodo(ctx context.Context, todo *requests.AddTodo) error
	UpdateTodo(ctx context.Context, todo *requests.UpdateTodo) error

	BroadcastNewTodo(todo *models.Todo) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Usecase {
	return &service{
		repository: r,
	}
}

func (s service) FetchAllTodos(ctx context.Context) (*[]outputs.Todo, error) {
	todos := make([]outputs.Todo, 0)
	result, err := s.repository.ReadAllTodos(ctx)
	if err != nil {
		return nil, err
	}

	for _, t := range result {
		todo := outputs.Todo{
			ID:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		}
		todos = append(todos, todo)
	}

	return &todos, nil
}

func (s service) FetchTodosWithRelated(ctx context.Context, request *requests.GetTodosWithRelated) (*[]outputs.TodoWithRelated, error) {
	todos := make([]outputs.TodoWithRelated, 0)
	result, err := s.repository.ReadTodosWithRelated(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, t := range *result {
		todo := outputs.TodoWithRelated{
			ID:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
			Name:      t.Name,
			CreatedAt: t.CreatedAt,
		}
		todos = append(todos, todo)
	}

	return &todos, nil
}

func (s service) FetchTodoById(ctx context.Context, id int) (*outputs.Todo, error) {
	result, err := s.repository.ReadTodoById(ctx, id)
	if err != nil {
		return nil, err
	}

	todo := outputs.Todo{
		ID:        result.ID,
		Title:     result.Title,
		Completed: result.Completed,
		CreatedAt: result.CreatedAt,
	}

	return &todo, nil
}

func (s service) InsertTodo(ctx context.Context, todo *requests.AddTodo) error {
	cTodo := models.Todo{
		Title: todo.Title,
	}
	err := s.repository.CreateTodo(ctx, &cTodo)
	if err != nil {
		return err
	}

	// broadcast the new Todo
	err = s.BroadcastNewTodo(&cTodo)
	if err != nil {
		return err
	}
	return nil
}

func (s service) UpdateTodo(ctx context.Context, todo *requests.UpdateTodo) error {
	uTodo := models.Todo{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
	}
	return s.repository.UpdateTodo(ctx, &uTodo)
}

func (s service) BroadcastNewTodo(todo *models.Todo) error {
	// Convert the Todo model to presenter.Todo
	pTodo := presenter.Todo{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
		CreatedAt: todo.CreatedAt.Format("2023-01-01 00:00:00"),
	}

	broadcast.TodoClient.SendNewTodo(pTodo)
	return nil
}