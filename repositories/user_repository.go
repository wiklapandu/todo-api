package repositories

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/wiklapandu/todo-api/models"
	"github.com/wiklapandu/todo-api/server/httputils"
	"github.com/wiklapandu/todo-api/types"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(user *models.User, db *gorm.DB, password string) error {
	result := db.Where("email = ?", user.Email).First(&user)
	if result.Error != nil {
		return fmt.Errorf("user is not registered")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println(err)
		return fmt.Errorf("email and password is not matches")
	}

	return nil
}

func Register() {}

func Update() {}

func Delete() {}

func ForgetPassword() {}

func GenerateToken(user *models.User, db *gorm.DB) (string, string, error) {
	session := map[string]interface{}{}

	if err := db.Table("sessions").Where("user_id = ?", user.ID).Take(&session).Error; err != nil {
		timeNow := time.Now().UTC().Format("2006-01-02 15:04:05")
		session["id"] = uuid.New()
		session["user_id"] = user.ID
		session["created_at"] = timeNow
		session["updated_at"] = timeNow

		if err := db.Table("sessions").Create(session).Error; err != nil {
			fmt.Println("Error creating session:", err)

			return "", "", fiber.ErrInternalServerError
		} else {
			fmt.Println("Session created successfully.")
		}
	}

	updateAt, err := httputils.ParseTime(session["updated_at"].(string))
	if err != nil {
		return "", "", fmt.Errorf(err.Error())
	}

	if time.Since(updateAt) > (time.Hour * 24) {
		session["updated_at"] = time.Now().UTC().Format("2006-01-02 15:04:05")
		db.Table("sessions").Where("id = ?", session["id"]).Updates(session)
	}

	timeAt, err := httputils.ParseTime(session["updated_at"].(string))
	if err != nil {
		fmt.Printf("error 0: %v \n", err)
	}

	claimsRefreshToken := types.ClaimsRefreshToken{}
	claimsRefreshToken.Subject = user.Email
	claimsRefreshToken.Type = "refresh"
	claimsRefreshToken.ExpiresAt = jwt.NewNumericDate(timeAt.Add((time.Hour * 24) * 7))

	claimsAccessToken := types.ClaimsAccessToken{}
	claimsAccessToken.IssuedAt = jwt.NewNumericDate(time.Now())
	claimsAccessToken.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 2))
	claimsAccessToken.Subject = user.ID.String()
	claimsAccessToken.Email = user.Email
	claimsAccessToken.Name = user.Name
	claimsAccessToken.Type = "access"

	var secretKey = []byte(os.Getenv("JWT_SECRET"))

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefreshToken)
	signedRefreshToken, err := refreshToken.SignedString(secretKey)
	if err != nil {
		log.Fatalf("Error signing refresh token: %v", err)
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccessToken)
	signedAccessToken, err := accessToken.SignedString(secretKey)
	if err != nil {
		log.Fatalf("Error signing access token: %v", err)
	}

	return signedAccessToken, signedRefreshToken, nil
}
