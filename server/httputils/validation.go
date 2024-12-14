package httputils

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidatorError map[string][]string

func Validate(data interface{}) ValidatorError {
	errors := ValidatorError{}

	file, err := os.Open("langs/validations/en.json")
	if err != nil {
		log.Println("Failed to open file:", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	var messages map[string]string
	if err := json.Unmarshal([]byte(byteValue), &messages); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	fileAttribute, errAttribute := os.Open("langs/attributes.json")
	if errAttribute != nil {
		log.Println("Failed to open file:", errAttribute)
	}
	defer fileAttribute.Close()

	byteAttribute, _ := io.ReadAll(fileAttribute)
	var attributes map[string]string
	if err := json.Unmarshal([]byte(byteAttribute), &attributes); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			field := err.Field()
			message := strings.Replace(strings.Replace(messages[err.Tag()], ":attribute", GetOrDefault(attributes, field, field), 1), ":tags", GetOrDefault(attributes, err.Param(), err.Param()), 1)

			attributeText := strings.ReplaceAll(strings.ToLower(GetOrDefault(attributes, field, field)), " ", "_")
			errors[attributeText] = append(errors[field], message)
		}
	}

	return errors
}

func GetOrDefault(data map[string]string, key, defaultValue string) string {
	if value, exists := data[key]; exists {
		return value
	}
	return defaultValue
}
