package src

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func TopicUrl(field validator.FieldLevel) bool {
	if str := field.Field().String(); len(str) != 0 {
		fmt.Println(str)
		if matched, _ := regexp.MatchString(`^\w{4,10}$`, str); matched {
			return true
		}
	}
	return false
}
