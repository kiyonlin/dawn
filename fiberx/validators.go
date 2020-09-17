package fiberx

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/utils"
)

var (
	mobileRegexp *regexp.Regexp
)

func init() {
	mobileRegexp = regexp.MustCompile(`^1[3456789]\d{9}$`)
}

// MobileRule validates mobile phone number
func MobileRule(fl validator.FieldLevel) bool {
	return mobileRegexp.Match(utils.GetBytes(fl.Field().String()))
}
