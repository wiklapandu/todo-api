package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/wiklapandu/todo-api/databases"
	"github.com/wiklapandu/todo-api/models"
	"github.com/wiklapandu/todo-api/repositories"
	"github.com/wiklapandu/todo-api/server/httputils"
	"github.com/wiklapandu/todo-api/types"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {
	body := &types.RequestLogin{}

	if err := c.BodyParser(body); err != nil {
		return err
	}

	if errors := httputils.Validate(body); len(errors) >= 1 {
		response := types.TypeAuthErrorResponseJson{}
		response.Status = fiber.StatusBadRequest
		response.Message = "Invalid validation"
		response.Data = errors
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	db, err := databases.DB()

	if err != nil {
		response := types.TypeAuthErrorResponseJson{}
		response.Status = fiber.StatusInternalServerError
		response.Message = "Internal Server Error"
		return c.Status(response.Status).JSON(response)
	}

	var user = &models.User{Email: body.Email}

	if err := repositories.Login(user, db, body.Pass); err != nil {
		response := types.TypeAuthErrorResponseJson{}
		response.Status = fiber.StatusBadRequest
		response.Message = err.Error()
		return c.Status(response.Status).JSON(response)
	}

	signedAccessToken, signedRefreshToken, err := repositories.GenerateToken(user, db)

	if err != nil {
		response := types.TypeAuthResponseJson{}
		response.Status = fiber.StatusInternalServerError
		response.Message = err.Error()
		return c.Status(response.Status).JSON(response)
	}

	response := types.TypeAuthResponseJson{}
	response.Status = fiber.StatusOK
	response.Message = "Success fully login"
	response.Data = &types.TypeDataToken{
		RefreshToken: signedRefreshToken,
		AccessToken:  signedAccessToken,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func Register(c *fiber.Ctx) error {
	body := &types.RequestRegister{}

	if err := c.BodyParser(body); err != nil {
		return err
	}

	if err := httputils.Validate(body); err != nil {
		response := types.TypeAuthErrorResponseJson{}
		response.Status = fiber.StatusBadRequest
		response.Message = "Invalid validation"
		response.Data = err
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	db, _ := databases.DB()

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Pass), bcrypt.DefaultCost)
	body.Pass = string(hash)

	db.Create(&models.User{Name: body.Name, Email: body.Email, Password: body.Pass})

	response := types.TypeAuthResponseJson{}
	response.Status = fiber.StatusOK
	response.Message = "Successfully Register"

	return c.Status(fiber.StatusOK).JSON(response)
}

func ForgetPassword() {

}

func GenerateAccessToken(c *fiber.Ctx) error {
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
			response := types.TypeAuthExpiredAccessToken{}
			response.Message = "Token is Expired! please login again."
			response.Status = 401
			response.RedirectLogin = true
			return c.Status(response.Status).JSON(response)
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": fiber.StatusUnauthorized, "message": "Invalid token"})
		}
	}

	claims := token.Claims.(*types.ClaimsRefreshToken)

	if claims.Type != "refresh" {
		response := types.TypeAuthExpiredAccessToken{}
		response.Message = "Invalid token! this is not refresh token. please login again"
		response.Status = fiber.StatusUnauthorized
		response.RedirectLogin = true
		return c.Status(response.Status).JSON(response)
	}

	db, _ := databases.DB()

	var user = models.User{Email: claims.Subject}
	if result := db.Find(&user); result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		response := types.TypeAuthErrorResponseJson{}
		response.Status = fiber.StatusBadRequest
		response.Message = "Token is invalid. please login again"
		fmt.Printf("[Generate Access Token] error: %s", result.Error)
		return c.Status(response.Status).JSON(response)
	}

	var secretKey = []byte(os.Getenv("JWT_SECRET"))

	claimsAccessToken := types.ClaimsAccessToken{}
	claimsAccessToken.IssuedAt = jwt.NewNumericDate(time.Now())
	claimsAccessToken.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 2))
	claimsAccessToken.Subject = user.ID.String()
	claimsAccessToken.Email = user.Email
	claimsAccessToken.Name = user.Name
	claimsAccessToken.Type = "access"

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccessToken)
	signedAccessToken, err := accessToken.SignedString(secretKey)
	if err != nil {
		log.Fatalf("Error signing access token: %v", err)
	}

	response := types.TypeAuthResponseJson{}
	response.Message = "Successfully generate access token!"
	response.Status = fiber.StatusOK
	response.Data.AccessToken = signedAccessToken
	response.Data.RefreshToken = tokenString
	return c.Status(response.Status).JSON(response)
}
