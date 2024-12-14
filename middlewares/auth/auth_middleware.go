package auth_middleware

import (
	"errors"
	"os"
	"strings"
	"github.com/wiklapandu/todo-api/databases"
	"github.com/wiklapandu/todo-api/models"
	"github.com/wiklapandu/todo-api/types"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthCheck(c *fiber.Ctx) error {
	tokenHeader := c.Get("Authorization")

	if tokenHeader == "" {
		response := types.TypeAuthErrorResponseJson{}
		response.Message = "Token is required"
		response.Status = fiber.StatusUnauthorized
		return c.Status(response.Status).JSON(response)
	}

	parts := strings.Split(tokenHeader, " ")
	if len(parts) != 2 {
		response := types.TypeAuthErrorResponseJson{}
		response.Message = "Invalid Authorization header format"
		response.Status = fiber.StatusUnauthorized
		return c.Status(response.Status).JSON(response)
	}

	if parts[0] != "Bearer" {
		response := types.TypeAuthErrorResponseJson{}
		response.Message = "Invalid token type"
		response.Status = fiber.StatusUnauthorized
		return c.Status(response.Status).JSON(response)
	}

	tokenString := parts[1]

	token, err := jwt.ParseWithClaims(tokenString, &types.ClaimsAccessToken{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "Malformed token"})
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "Invalid signature"})
		case errors.Is(err, jwt.ErrTokenExpired):
			response := types.TypeAuthResponseJson{}
			response.Message = "Token is Expired! regenerate access token."
			response.Status = 402
			return c.Status(response.Status).JSON(response)
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "Invalid token"})
		}
	}

	claims := token.Claims.(*types.ClaimsAccessToken)

	if claims.Type != "access" {
		response := types.TypeAuthResponseJson{}
		response.Message = "Invalid token! this is not access token. please login again"
		response.Status = fiber.StatusUnauthorized
		return c.Status(response.Status).JSON(response)
	}

	db, _ := databases.DB()

	user := &models.User{}
	if result := db.First(user, "email = ?", claims.Email); result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "User is not availabel",
			"data":    nil,
		})
	}

	c.Locals("user", user)

	return c.Next()
}
