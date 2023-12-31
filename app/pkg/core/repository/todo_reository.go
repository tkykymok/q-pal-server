package repository

import (
	"app/api/requests"
	"app/pkg/exmodels"
	"app/pkg/models"
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

type TodoRepository interface {
	ReadAllTodos(ctx context.Context) (models.TodoSlice, error)
	ReadTodosWithRelated(ctx context.Context, request *requests.GetTodosWithRelated) (*[]exmodels.TodoWithRelated, error)
	ReadTodoById(ctx context.Context, id int) (*models.Todo, error)
	CreateTodo(ctx context.Context, exec boil.ContextExecutor, todo *models.Todo) error
	UpdateTodo(ctx context.Context, exec boil.ContextExecutor, todo *models.Todo) error
}

type todoRepository struct {
}

func NewTodoRepo() TodoRepository {
	return &todoRepository{}
}

func (r todoRepository) ReadAllTodos(ctx context.Context) (models.TodoSlice, error) {
	return models.Todos().AllG(ctx)
}

func (r todoRepository) ReadTodosWithRelated(ctx context.Context, request *requests.GetTodosWithRelated) (*[]exmodels.TodoWithRelated, error) {
	// SELECTするカラム
	selectCols := []string{
		models.TodoTableColumns.ID,
		models.TodoTableColumns.Title,
		models.TodoTableColumns.Completed,
		models.UserTableColumns.Name,
		models.TodoTableColumns.CreatedAt,
	}

	// QueryModの生成
	mods := []qm.QueryMod{
		qm.Select(strings.Join(selectCols, ",")),
		qm.LeftOuterJoin("users on todos.userId = users.id"),
	}
	// WHERE句
	if request.ID != 0 {
		mods = append(mods, qm.Where("todos.id=?", request.ID))
	}
	if request.UserId != 0 {
		mods = append(mods, qm.Where("users.id=?", request.UserId))
	}

	var result []exmodels.TodoWithRelated
	err := models.Todos(mods...).BindG(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to read todos with related: %w", err)
	}

	return &result, nil
}

func (r todoRepository) ReadTodoById(ctx context.Context, id int) (*models.Todo, error) {
	return models.FindTodoG(ctx, id)
}

func (r todoRepository) CreateTodo(ctx context.Context, exec boil.ContextExecutor, todo *models.Todo) error {
	err := todo.Insert(ctx, exec, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}

func (r todoRepository) UpdateTodo(ctx context.Context, exec boil.ContextExecutor, todo *models.Todo) error {
	_, err := todo.Update(ctx, exec, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}
