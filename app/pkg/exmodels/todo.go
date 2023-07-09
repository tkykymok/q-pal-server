package exmodels

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type TodoWithRelated struct {
	ID        int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Title     string      `boil:"title" json:"title,omitempty" toml:"title" yaml:"title,omitempty"`
	Completed bool        `boil:"completed" json:"completed" toml:"completed" yaml:"completed"`
	Name      null.String `boil:"name" json:"name" toml:"name" yaml:"name"`
	CreatedAt time.Time   `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
}
