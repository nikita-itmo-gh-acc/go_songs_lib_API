package query

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

func getValidatorFunc(regex *regexp.Regexp) func(interface{}, interface{}) bool {
	return func(i interface{}, o interface{}) bool {
		param, ok := i.(string)
		if !ok {
			return false
		}
		return regex.MatchString(param)
	}
}

func SetQueryValidators() {
	dateRegex := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
	linkRegex := regexp.MustCompile(`^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$`)

	govalidator.CustomTypeTagMap.Set("date", govalidator.CustomTypeValidator(getValidatorFunc(dateRegex)))
	govalidator.CustomTypeTagMap.Set("link", govalidator.CustomTypeValidator(getValidatorFunc(linkRegex)))
}

