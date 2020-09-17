package fiberx

import (
	"strings"

	"github.com/go-playground/locales/en"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ent "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

var (
	v     *validator.Validate
	uni   *ut.UniversalTranslator
	trans ut.Translator
)

func init() {
	v = validator.New()

	uni = ut.New(en.New())
	trans, _ = uni.GetTranslator("en")

	if err := ent.RegisterDefaultTranslations(v, trans); err != nil {
		panic(err)
	}
}

// ErrHandler is Dawn's error handler
var ErrHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	res := Response{
		Code:    code,
		Message: utils.StatusMessage(code),
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		res.Code = fiber.StatusUnprocessableEntity
		res.Message = ""
		res.Data = removeTopStruct(errs.Translate(trans))
	} else if e, ok := err.(*fiber.Error); ok {
		res.Code = e.Code
		res.Message = e.Message
	}

	return Resp(c, res.Code, res)
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, msg := range fields {
		stripStruct := field[strings.Index(field, ".")+1:]
		res[stripStruct] = strings.TrimLeft(msg, stripStruct)
	}
	return res
}

// ValidateBody accepts a obj holds results from BodyParser
// and then do the validation by a validator
func ValidateBody(c *fiber.Ctx, obj interface{}) (err error) {
	if err := c.BodyParser(obj); err != nil {
		return err
	}

	return v.Struct(obj)
}

// ValidateQuery accepts a obj holds results from QueryParser
// and then do the validation by a validator
func ValidateQuery(c *fiber.Ctx, obj interface{}) (err error) {
	if err := c.QueryParser(obj); err != nil {
		return err
	}

	return v.Struct(obj)
}

// Response is a unified format for api results
type Response struct {
	// Code is the status code by default, but also can be
	// a custom code
	Code int `json:"code,omitempty"`
	// Message shows detail thing back to caller
	Message string `json:"message,omitempty"`
	// RequestID needs to be used with middleware
	RequestID string `json:"request_id,omitempty"`
	// Data accepts any thing as the response data
	Data interface{} `json:"data,omitempty"`
}

// Resp returns the
func Resp(c *fiber.Ctx, statusCode int, res ...Response) error {
	r := Response{}
	if len(res) > 0 {
		r = res[0]
	}

	if r.Code == 0 {
		r.Code = statusCode
	}

	if id := c.Response().Header.Peek(fiber.HeaderXRequestID); len(id) > 0 {
		r.RequestID = utils.GetString(id)
	}

	return c.Status(statusCode).JSON(r)
}

// respCommon
func respCommon(c *fiber.Ctx, code int, msg ...string) error {
	res := Response{
		Message: utils.StatusMessage(code),
	}

	if len(msg) > 0 {
		res.Message = msg[0]
	}
	return Resp(c, code, res)
}

// RespOK responses with status code 200 RFC 7231, 6.3.1
func RespOK(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusOK, msg...)
}

// RespCreated responses with status code 201 RFC 7231, 6.3.2
func RespCreated(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusCreated, msg...)
}

// RespAccepted responses with status code 202 RFC 7231, 6.3.3
func RespAccepted(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusAccepted, msg...)
}

// RespNonAuthoritativeInformation responses with status code 203 RFC 7231, 6.3.4
func RespNonAuthoritativeInformation(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusNonAuthoritativeInformation, msg...)
}

// RespNoContent responses with status code 204 RFC 7231, 6.3.5
func RespNoContent(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusNoContent, msg...)
}

// RespResetContent responses with status code 205 RFC 7231, 6.3.6
func RespResetContent(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusResetContent, msg...)
}

// RespPartialContent responses with status code 206 RFC 7233, 4.1
func RespPartialContent(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusPartialContent, msg...)
}

// RespMultiStatus responses with status code 207 RFC 4918, 11.1
func RespMultiStatus(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusMultiStatus, msg...)
}

// RespAlreadyReported responses with status code 208 RFC 5842, 7.1
func RespAlreadyReported(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusAlreadyReported, msg...)
}
