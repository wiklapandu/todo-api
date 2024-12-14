package todo_service

import (
	"fmt"
	"time"
	"github.com/wiklapandu/todo-api/databases"
	"github.com/wiklapandu/todo-api/models"
	"github.com/wiklapandu/todo-api/server/httputils"
	"github.com/wiklapandu/todo-api/types"

	"github.com/gofiber/fiber/v2"
)

func Index(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("user").(*models.User)
	db, _ := databases.DB()

	todos := []models.Todo{}
	if result := db.Preload("User").Preload("Author").Find(&todos, "user_id = ?", currentUser.ID); result.Error != nil || result.RowsAffected < 1 {
		response := types.TypeResponseTodoIndex{}
		response.Message = "Todo is not found"
		if result.Error != nil {
			fmt.Printf("[Todo Show] error: %s \n", result.Error)
		}

		fmt.Println("last query:", result.Statement.SQL.String())
		return ctx.Status(response.Status).JSON(response)
	}

	todosPublic := []types.TodoPublic{}

	for i := 0; i < len(todos); i++ {
		item := todos[i]
		todosPublic = append(todosPublic, types.TodoPublic{
			ID: item.ID,
			FromUser: &types.UserPublic{
				ID:    item.Author.ID,
				Name:  item.Author.Name,
				Email: item.Author.Email,
			},
			Title:       item.Title,
			Description: item.Description,
			Status:      item.Status,
			Color:       item.Color,
			DueTo:       item.DueTo,
			Published:   item.UpdatedAt.String(),
		})
	}

	response := types.TypeResponseTodoIndex{}
	response.Message = "Successfully getting todos"
	response.Data = &todosPublic
	response.Status = fiber.StatusOK

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func Show(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	db, _ := databases.DB()

	var todo = &models.Todo{}
	if result := db.First(todo, "id = ?", id); result.Error != nil || result.RowsAffected < 1 {
		response := types.TypeResponseTodoShow{}
		response.Message = "Todo is not found"
		response.Status = fiber.StatusBadRequest
		if result.Error != nil {
			fmt.Printf("[Todo Show] error: %s \n", result.Error)
		}

		fmt.Println("last query:", result.Statement.SQL.String())
		return ctx.Status(response.Status).JSON(response)
	}

	authorTodo := &models.User{}
	var fromUser *types.UserPublic

	if result := db.First(authorTodo, "id = ?", todo.FromUser); result.Error == nil {
		fromUser = &types.UserPublic{
			ID:    authorTodo.ID,
			Name:  authorTodo.Name,
			Email: authorTodo.Email,
		}
	}

	response := types.TypeResponseTodoShow{}
	response.Message = "Successfully getting data"
	response.Status = fiber.StatusOK
	response.Data = &types.TodoPublic{
		ID:          todo.ID,
		FromUser:    fromUser,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
		Color:       todo.Color,
		DueTo:       todo.DueTo,
		Published:   todo.CreatedAt.String(),
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func Store(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("user").(*models.User)
	var body types.RequestStoreTodo

	if err := ctx.BodyParser(&body); err != nil {
		fmt.Println("error: ", err.Error())
		return err
	}

	if err := httputils.Validate(body); len(err) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "invalid validation field!",
			"data":    err,
		})
	}

	data := &models.Todo{
		FromUser:    currentUser.ID,
		Title:       body.Title,
		Description: body.Description,
		Status:      body.Status,
		Color:       body.Color,
		UserId:      currentUser.ID,
	}

	if body.DueTo != "" {
		layout := "2006-01-02 15:04:05"
		parse, _ := time.Parse(layout, body.DueTo)
		data.DueTo = &parse
	}

	if currentUser.CheckPermission("Can Sign User") {
		if body.ForUser != nil {
			data.UserId = *body.ForUser
		}
	}

	db, _ := databases.DB()

	if result := db.Create(data); result.RowsAffected < 1 || result.Error != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed insert todo!",
			"data":    nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Successfully create todo",
		"data":    data,
	})
}

func Delete(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully delete todo",
		"data":    nil,
	})
}

func Update(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("user").(*models.User)
	id := ctx.Params("id")

	var body types.RequestStoreTodo

	if err := ctx.BodyParser(&body); err != nil {
		fmt.Println("error: ", err.Error())
		return err
	}

	if err := httputils.Validate(body); len(err) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "invalid validation field!",
			"data": fiber.Map{
				"errors": err,
			},
		})
	}

	db, _ := databases.DB()

	var todo *models.Todo
	if result := db.First(&todo, "id = ? AND user_id = ?", id, currentUser.ID); result.Error != nil || result.RowsAffected < 1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Failed! unknow todo",
			"data":    nil,
		})
	}

	todo.Title = body.Title
	todo.Color = body.Color
	todo.Status = body.Status

	if body.DueTo != "" {
		layout := "2006-01-02 15:04:05"
		parse, _ := time.Parse(layout, body.DueTo)
		todo.DueTo = &parse
	}

	if currentUser.CheckPermission("Can Sign User") {
		if body.ForUser != nil {
			todo.UserId = *body.ForUser
		}
	}

	// jsonPrint, _ := json.Marshal(todo)
	// fmt.Println("todo: ", string(jsonPrint))

	if body.DueTo != "" {
		layout := "2006-01-02 15:04:05"
		parse, _ := time.Parse(layout, body.DueTo)
		todo.DueTo = &parse
	}

	if result := db.Save(todo); result.RowsAffected < 1 || result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Failed! fail to create todo",
			"data": fiber.Map{
				"errorMessage": result.Error,
			},
		})
	}

	author := &models.User{}
	if result := db.First(author, "id = ?", todo.FromUser); result.Error != nil || result.RowsAffected < 1 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Failed! fail to create todo",
			"data": fiber.Map{
				"errorMessage": result.Error,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully update todo",
		"data": types.TodoPublic{
			ID:    todo.ID,
			Title: todo.Title,
			FromUser: &types.UserPublic{
				ID:    author.ID,
				Name:  author.Name,
				Email: author.Email,
			},
			Description: todo.Description,
			Status:      todo.Status,
			Color:       todo.Color,
			DueTo:       todo.DueTo,
			Published:   todo.UpdatedAt.String(),
		},
	})
}
