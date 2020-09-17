package fiberx

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestMobile(t *testing.T) {
	at := assert.New(t)

	v := validator.New()
	at.Nil(v.RegisterValidation("mobile", MobileRule))

	at.Nil(v.Var("13888888888", "mobile"))
	at.NotNil(v.Var("23888888888", "mobile"))
	at.NotNil(v.Var("1388888888", "mobile"))
	at.NotNil(v.Var(13888888888, "mobile"))
}

func BenchmarkName(b *testing.B) {
	ab := assert.New(b)

	v := validator.New()
	ab.Nil(v.RegisterValidation("mobile", MobileRule))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = v.Var("13888888888", "mobile")
	}
}
