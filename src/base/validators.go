package base

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strings"
)

func ValidateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	pattern := `^[a-zA-Z0-9_-]{4,24}$`
	return regexp.MustCompile(pattern).MatchString(username)
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	pattern := `^[a-zA-Z0-9_!@#$%^&*]{8,24}$`
	return regexp.MustCompile(pattern).MatchString(password)
}

func ValidateFilename(fl validator.FieldLevel) bool {
	filename := fl.Field().String()

	pattern := `^[^<>:"/\\|?*]{1,200}$`
	return regexp.MustCompile(pattern).MatchString(filename)
}

func CreateValidator() *validator.Validate {
	schemaValidator := validator.New()

	schemaValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := schemaValidator.RegisterValidation("username", ValidateUsername); err != nil {
		panic(err)
	}
	if err := schemaValidator.RegisterValidation("password", ValidatePassword); err != nil {
		panic(err)
	}
	if err := schemaValidator.RegisterValidation("filename", ValidateFilename); err != nil {
		panic(err)
	}

	return schemaValidator
}

func getErrorMessageForTag(tag string) string {
	switch tag {
	case "username":
		return "Username can only contain latin symbols, " +
			"numbers, symbols '_-' with length 4-24"
	case "password":
		return "Password can only contain latin symbols, " +
			"numbers, symbols '_!@#$%^&*' with length 8-24"
	case "filename":
		return "File name should not contain symbols <>:\"\\/|?* " +
			"and should have length 1-200"
	case "required":
		return "Field required"
	case "gte":
		return "The field length is less than the specified length"
	}
	return ""
}
