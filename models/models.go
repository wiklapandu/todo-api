package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	ID          uuid.UUID  `gorm:"type:char(36);primary_key;column:id"`
	UserId      uuid.UUID  `json:"userId" gorm:"type:char(36);key;column:user_id"`
	User        User       `gorm:"foreignKey:UserId;references:ID"`
	Author      User       `gorm:"foreignKey:FromUser;references:ID"`
	FromUser    uuid.UUID  `json:"fromUserId" gorm:"type:char(36);key;column:from_user_id"`
	Title       string     `json:"title" gorm:"column:title"`
	Description string     `json:"description" gorm:"column:description"`
	Status      string     `json:"status" gorm:"column:status"`
	Color       string     `json:"color" gorm:"column:color"`
	DueTo       *time.Time `json:"dueTo" gorm:"column:due_to"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName overrides the default table name
func (Todo) TableName() string {
	return "todo_lists"
}

func (todo *Todo) BeforeCreate(tx *gorm.DB) (err error) {
	todo.ID = uuid.New()
	return
}

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key;column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	Email     string    `json:"email" gorm:"column:email"`
	Password  string    `json:"password" gorm:"column:password"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

type UserRole struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key;column:id"`
}

func (UserRole) TableName() string {
	return "user_role"
}

func (user *User) TableName() string {
	return "users"
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}

func (user *User) CheckPermission(permission string) bool {
	fmt.Println(user.Email)
	fmt.Println(permission)

	return true
}
