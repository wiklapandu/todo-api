package types

import (
	"time"

	"github.com/google/uuid"
)

type UserPublic struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type TodoPublic struct {
	ID          uuid.UUID   `json:"id"`
	FromUser    *UserPublic `json:"fromUser"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	Color       string      `json:"color"`
	DueTo       *time.Time  `json:"dueTo"`
	Published   string      `json:"published"`
}

type TypeResponseTodoIndex struct {
	TypeResponse
	Data *[]TodoPublic `json:"data"`
}

type TypeResponseTodoShow struct {
	TypeResponse
	Data *TodoPublic `json:"data"`
}

type RequestStoreTodo struct {
	Title       string     `validate:"required" json:"title" xml:"title" form:"title"`
	Description string     `validate:"required" json:"description" xml:"description" form:"description"`
	Status      string     `validate:"required" json:"status" xml:"status" form:"status"`
	Color       string     `validate:"required" json:"color" xml:"color" form:"color"`
	DueTo       string     `validate:"" json:"dueTo" xml:"dueTo" form:"dueTo"`
	ForUser     *uuid.UUID `validate:"" json:"forUser" xml:"forUser" form:"dueTo"`
}

type RequestUpdateTodo struct {
	Title       string     `validate:"required" json:"title" xml:"title" form:"title"`
	Description string     `validate:"required" json:"description" xml:"description" form:"description"`
	Status      string     `validate:"required" json:"status" xml:"status" form:"status"`
	Color       string     `validate:"required" json:"color" xml:"color" form:"color"`
	DueTo       string     `validate:"" json:"dueTo" xml:"dueTo" form:"dueTo"`
	ForUser     *uuid.UUID `validate:"" json:"forUser" xml:"forUser" form:"dueTo"`
}
