package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"regexp"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.RegisterValidation("phone", validatePhoneNumber)
	if err != nil {
		return err
	}
	if err := cv.Validator.Struct(i); err != nil {
		// 在这里自定义错误消息
		return echo.NewHTTPError(400, "Validation failed: "+err.Error())
	}
	return nil
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()

	// 自定义电话号码的正则表达式
	phoneRegex := `^\+?\d{1,3}[-.\s]?\(?\d{1,3}\)?[-.\s]?\d{1,4}[-.\s]?\d{1,4}$`

	r, err := regexp.Compile(phoneRegex)
	if err != nil {
		return false
	}
	// 使用正则表达式匹配电话号码
	return r.MatchString(phoneNumber)

}
