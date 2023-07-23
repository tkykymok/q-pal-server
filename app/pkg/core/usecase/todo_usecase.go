package usecase

import (
	"app/api/presenter"
	"app/api/requests"
	"app/pkg/broadcast"
	"app/pkg/core/repository"
	"app/pkg/models"
	"app/pkg/outputs"
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TodoUsecase interface {
	FetchAllTodos(ctx context.Context) (*[]outputs.Todo, error)
	FetchTodosWithRelated(ctx context.Context, request *requests.GetTodosWithRelated) (*[]outputs.TodoWithRelated, error)
	FetchTodoById(ctx context.Context, id int) (*outputs.Todo, error)
	InsertTodo(ctx context.Context, todo *requests.AddTodo) error
	UpdateTodo(ctx context.Context, todo *requests.UpdateTodo) error

	BroadcastNewTodo(todo *models.Todo) error
}

type todoUsecase struct {
	repository repository.TodoRepository
}

func NewTodoUsecase(r repository.TodoRepository) TodoUsecase {
	return &todoUsecase{
		repository: r,
	}
}

func (u todoUsecase) FetchAllTodos(ctx context.Context) (*[]outputs.Todo, error) {
	todos := make([]outputs.Todo, 0)
	result, err := u.repository.ReadAllTodos(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read all todos: %w", err)
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

func (u todoUsecase) FetchTodosWithRelated(ctx context.Context, request *requests.GetTodosWithRelated) (*[]outputs.TodoWithRelated, error) {
	todos := make([]outputs.TodoWithRelated, 0)
	result, err := u.repository.ReadTodosWithRelated(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to read todos with related: %w", err)
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

func (u todoUsecase) FetchTodoById(ctx context.Context, id int) (*outputs.Todo, error) {
	result, err := u.repository.ReadTodoById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read todo by id: %w", err)
	}

	todo := outputs.Todo{
		ID:        result.ID,
		Title:     result.Title,
		Completed: result.Completed,
		CreatedAt: result.CreatedAt,
	}

	return &todo, nil
}

func (u todoUsecase) InsertTodo(ctx context.Context, todo *requests.AddTodo) error {
	// Start a new transaction
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	cTodo := models.Todo{
		Title: todo.Title,
	}
	// Pass the transaction to the repository method
	err = u.repository.CreateTodo(ctx, tx, &cTodo)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	// broadcast the new Todo
	err = u.BroadcastNewTodo(&cTodo)
	if err != nil {
		return fmt.Errorf("failed to broadcast new todo: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (u todoUsecase) UpdateTodo(ctx context.Context, todo *requests.UpdateTodo) error {
	uTodo := models.Todo{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: todo.Completed,
	}
	err := u.repository.UpdateTodo(ctx, &uTodo)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}

func (u todoUsecase) BroadcastNewTodo(todo *models.Todo) error {
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
