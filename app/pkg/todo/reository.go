package todo

import (
	"app/api/requests"
	"app/pkg/exmodels"
	"app/pkg/models"
	"context"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"strings"
)

type Repository interface {
	ReadAllTodos(ctx context.Context) (models.TodoSlice, error)
	ReadTodosWithRelated(ctx context.Context, request *requests.GetTodosWithRelated) (*[]exmodels.TodoWithRelated, error)
	ReadTodoById(ctx context.Context, id int) (*models.Todo, error)
	CreateTodo(ctx context.Context, todo *models.Todo) error
	UpdateTodo(ctx context.Context, todo *models.Todo) error
}

type repository struct {
}

func NewRepo() Repository {
	return &repository{}
}

func (r repository) ReadAllTodos(ctx context.Context) (models.TodoSlice, error) {
	return models.Todos().AllG(ctx)
}

func (r repository) ReadTodosWithRelated(ctx context.Context, request *requests.GetTodosWithRelated) (*[]exmodels.TodoWithRelated, error) {
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
		log.Fatal(err)
	}

	return &result, nil
}

func (r repository) ReadTodoById(ctx context.Context, id int) (*models.Todo, error) {
	return models.FindTodoG(ctx, id)
}

func (r repository) CreateTodo(ctx context.Context, todo *models.Todo) error {
	err := todo.InsertG(ctx, boil.Infer())
	if err != nil {
		return nil
	}
	return err
}

func (r repository) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	_, err := todo.UpdateG(ctx, boil.Infer())
	if err != nil {
		return nil
	}
	return err
}
