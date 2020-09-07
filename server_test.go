package dawn

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
)

var (
	m = mockModule{}
)

func Test_Server_New(t *testing.T) {
	app := fiber.New()

	s := New(
		App(app),
		Modulars(m),
	)

	assert.Equal(t, app, s.app)
	assert.Len(t, s.mods, 1)
	assert.Equal(t, "anonymous", s.mods[0].String())
}

func Test_Server_Default(t *testing.T) {
	s := Default()

	require.NotNil(t, s.app)
	assert.Len(t, s.app.Stack()[0], 1)
}

func Test_Server_AddModulars(t *testing.T) {
	s := New()

	s.AddModulars(m)
	assert.Len(t, s.mods, 1)
	assert.Equal(t, "anonymous", s.mods[0].String())
}

func Test_Server_Run(t *testing.T) {
	s := New(Modulars(m))

	go func() {
		time.Sleep(time.Millisecond * 100)
		assert.NoError(t, s.app.Shutdown())
	}()

	assert.NoError(t, s.Run(""))
}

func Test_Server_RunTls(t *testing.T) {
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

func Test_Server_Shutdown(t *testing.T) {
	require.NotNil(t, (&Server{}).Shutdown())
	require.NotNil(t, New().Shutdown())
}

func Test_Server_Router(t *testing.T) {
	require.Nil(t, (&Server{}).Router())
	require.NotNil(t, New().Router())
}
