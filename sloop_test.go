package dawn

import (
	"testing"
	"time"

	"github.com/kiyonlin/dawn/config"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	m = mockModule{}
)

func Test_Sloop_New(t *testing.T) {
	app := fiber.New()

	s := New(
		App(app),
		Modulars(m),
	)

	assert.Equal(t, app, s.app)
	assert.Len(t, s.mods, 1)
	assert.Equal(t, "anonymous", s.mods[0].String())
}

func Test_Sloop_Default(t *testing.T) {
	config.Set("debug", true)
	s := Default(fiber.Config{})

	require.NotNil(t, s.app)
	assert.Len(t, s.app.Stack()[0], 1)
}

func Test_Sloop_AddModulars(t *testing.T) {
	s := New()

	s.AddModulars(m)
	assert.Len(t, s.mods, 1)
	assert.Equal(t, "anonymous", s.mods[0].String())
}

func Test_Sloop_Run(t *testing.T) {
	s := New(Modulars(m))

	go func() {
		time.Sleep(time.Millisecond * 100)
		assert.NoError(t, s.app.Shutdown())
	}()

	assert.NoError(t, s.Run(""))
}

func Test_Sloop_RunTls(t *testing.T) {
	s := New()

	t.Run("invalid addr", func(t *testing.T) {
		assert.NotNil(t, s.RunTls(":99999", "./.github/testdata/ssl.pem", "./.github/testdata/ssl.key"))
	})

	t.Run("invalid ssl info", func(t *testing.T) {
		assert.NotNil(t, s.RunTls("", "./.github/README.md", "./.github/README.md"))
	})

	t.Run("with ssl", func(t *testing.T) {
		go func() {
			time.Sleep(time.Millisecond * 100)
			assert.NoError(t, s.app.Shutdown())
		}()

		assert.NoError(t, s.RunTls("", "./.github/testdata/ssl.pem", "./.github/testdata/ssl.key"))
	})
}

func Test_Sloop_Shutdown(t *testing.T) {
	require.NotNil(t, (&Sloop{}).Shutdown())
	require.Nil(t, New().Shutdown())
}

func Test_Sloop_Router(t *testing.T) {
	require.Nil(t, (&Sloop{}).Router())
	require.NotNil(t, New().Router())
}
