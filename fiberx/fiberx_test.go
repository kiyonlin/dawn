package fiberx

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/gofiber/fiber/v2/utils"

	"github.com/stretchr/testify/assert"

	"github.com/gofiber/fiber/v2"
)

func Test_Fiberx_ErrorHandler(t *testing.T) {
	at := assert.New(t)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrHandler,
	})

	t.Run("StatusUnprocessableEntity", func(t *testing.T) {
		app.Get("/422", func(c *fiber.Ctx) error {
			type User struct {
				Username string `validate:"required"`
				Field1   string `validate:"required,lt=10"`
				Field2   string `validate:"required,gt=1"`
			}

			user := User{
				Username: "kiyon",
				Field1:   "This field is always too long.",
				Field2:   "1",
			}

			return v.Struct(user)
		})

		resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/422", nil))
		at.Nil(err)
		at.Equal(fiber.StatusUnprocessableEntity, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		at.Nil(err)
		at.JSONEq(`{"code":422, "data":{"Field1":" must be less than 10 characters in length", "Field2":" must be greater than 1 character in length"}}`, string(body))
	})

	t.Run("normal error", func(t *testing.T) {
		app.Get("/", func(c *fiber.Ctx) error {
			return errors.New("hi, i'm an error")
		})

		resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
		at.Nil(err)
		at.Equal(fiber.StatusInternalServerError, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		at.Nil(err)
		at.JSONEq(`{"code":500, "message":"Internal Server Error"}`, string(body))
	})

	t.Run("fiber error", func(t *testing.T) {
		app.Get("/400", func(c *fiber.Ctx) error {
			return fiber.ErrBadRequest
		})

		resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/400", nil))
		at.Nil(err)
		at.Equal(fiber.StatusBadRequest, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		at.Nil(err)
		at.JSONEq(`{"code":400, "message":"Bad Request"}`, string(body))
	})
}

func Test_Fiberx_ValidateBody(t *testing.T) {
	at := assert.New(t)
	t.Run("success", func(t *testing.T) {

		app := fiber.New()

		app.Post("/", func(c *fiber.Ctx) error {
			type User struct {
				Username string `validate:"required" json:"username"`
			}

			var u User
			if err := ValidateBody(c, &u); err != nil {
				return err
			}

			return c.SendString(u.Username)
		})

		res := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewReader([]byte("username=kiyon")))
		res.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationForm)
		resp, err := app.Test(res)
		at.Nil(err)
		at.Equal(fiber.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		at.Nil(err)
		at.Equal("kiyon", string(body))
	})

	t.Run("error", func(t *testing.T) {
		c := fiber.New().AcquireCtx(&fasthttp.RequestCtx{})
		at.NotNil(ValidateBody(c, nil))
	})
}

func Test_Fiberx_ValidateQuery(t *testing.T) {
	at := assert.New(t)
	t.Run("success", func(t *testing.T) {

		app := fiber.New()

		app.Get("/", func(c *fiber.Ctx) error {
			type User struct {
				Username string `validate:"required" json:"username"`
			}

			var u User
			if err := ValidateQuery(c, &u); err != nil {
				return err
			}

			return c.SendString(u.Username)
		})

		resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/?username=kiyon", nil))
		at.Nil(err)
		at.Equal(fiber.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		at.Nil(err)
		at.Equal("kiyon", string(body))
	})

	t.Run("error", func(t *testing.T) {
		fctx := &fasthttp.RequestCtx{}
		fctx.Request.URI().SetQueryString("a=b")
		c := fiber.New().AcquireCtx(fctx)
		at.NotNil(ValidateQuery(c, nil))
	})
}

func Test_Fiberx_2xx(t *testing.T) {
	tt := []struct {
		code int
		fn   func(c *fiber.Ctx, msg ...string) error
	}{
		{fiber.StatusOK, RespOK},
		{fiber.StatusCreated, RespCreated},
		{fiber.StatusAccepted, RespAccepted},
		{fiber.StatusNonAuthoritativeInformation, RespNonAuthoritativeInformation},
		{fiber.StatusNoContent, RespNoContent},
		{fiber.StatusResetContent, RespResetContent},
		{fiber.StatusPartialContent, RespPartialContent},
		{fiber.StatusMultiStatus, RespMultiStatus},
		{fiber.StatusAlreadyReported, RespAlreadyReported},
	}

	for _, tc := range tt {
		t.Run(strconv.Itoa(tc.code), func(t *testing.T) {
			at := assert.New(t)
			fn := tc.fn
			app := fiber.New()
			app.Get("/", func(c *fiber.Ctx) error {
				if tc.code == fiber.StatusNoContent {
					return fn(c, "I will be removed")
				}
				return fn(c)
			})

			resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
			at.Nil(err)
			at.Equal(tc.code, resp.StatusCode)

			body, err := ioutil.ReadAll(resp.Body)
			at.Nil(err)
			expected := fmt.Sprintf("{\"code\":%d,\"message\":\"%s\"}", tc.code, utils.StatusMessage(tc.code))
			if tc.code == fiber.StatusNoContent {
				expected = ""
			}
			at.Equal(expected, string(body))
		})
	}
}

func Test_Fiberx_RequestID(t *testing.T) {
	at := assert.New(t)
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderXRequestID, "id")
		return RespOK(c)
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	at.Nil(err)
	at.Equal(fiber.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	at.Nil(err)
	at.Equal(`{"code":200,"message":"OK","request_id":"id"}`, string(body))
}

func Test_Fiberx_Logger(t *testing.T) {
	at := assert.New(t)

	app := fiber.New()
	app.Use(Logger())
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderXRequestID, "id")
		return RespOK(c)
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.ErrForbidden
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	at.Nil(err)
	at.Equal(fiber.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	at.Nil(err)
	at.Equal(`{"code":200,"message":"OK","request_id":"id"}`, string(body))

	resp, err = app.Test(httptest.NewRequest(fiber.MethodGet, "/error", nil))
	at.Nil(err)
	at.Equal(fiber.StatusForbidden, resp.StatusCode)

	body, err = ioutil.ReadAll(resp.Body)
	at.Nil(err)
	at.Equal("Forbidden", string(body))
}
