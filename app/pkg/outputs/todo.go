package outputs

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title,omitempty"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type TodoWithRelated struct {
	ID        int         `json:"id"`
	Title     string      `json:"title,omitempty"`
	Completed bool        `json:"completed"`
	Name      null.String `json:"name"`
	CreatedAt time.Time   `json:"created_at,omitempty"`
}
