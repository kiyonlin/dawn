package sql

import (
	"sync"
	"testing"

	"github.com/kiyonlin/dawn"
	"github.com/kiyonlin/dawn/config"

	"github.com/stretchr/testify/assert"
)

func Test_Sql_New(t *testing.T) {
	modular := New()
	_, ok := modular.(*sqlModule)
	assert.True(t, ok)
}

func Test_Sql_Module_Name(t *testing.T) {
	assert.Equal(t, "dawn:sql", m.String())
}

func Test_Sql_Init(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		at := assert.New(t)

		var wg sync.WaitGroup
		wg.Add(1)
		m.Init(&wg)
		wg.Wait()

		at.Equal(fallback, m.fallback)
		at.Len(m.conns, 1)
	})

	t.Run("with config", func(t *testing.T) {
		config.Set("sql.default", "sqlite")
		config.Set("sql.connections.sqlite", map[string]string{})

		at := assert.New(t)

		var wg sync.WaitGroup
		wg.Add(1)
		m.Init(&wg)
		wg.Wait()

		at.Equal("sqlite", m.fallback)
		at.Len(m.conns, 1)
	})
}

func Test_Sql_Boot(t *testing.T) {
	at := assert.New(t)

	m.conns[fallback] = connect(fallback, config.New())

	var (
		wg      sync.WaitGroup
		cleanup = make(chan dawn.Cleanup, 1)
	)
	wg.Add(1)
	m.Boot(&wg, cleanup)
	wg.Wait()

	(<-cleanup)()

	db, err := m.conns[fallback].DB()
	at.Nil(err)
	at.NotNil(db.Ping())
}

func Test_Sql_Conn(t *testing.T) {
	assert.Nil(t, Conn("non"))
}

func Test_Sql_connect(t *testing.T) {
	t.Run("unknow driver", func(t *testing.T) {
		defer func() {
			assert.Equal(t, "dawn:sql unknown driver test of name", recover())
		}()
		c := config.New()
		c.Set("driver", "test")
		connect("name", c)
	})

	t.Run("error", func(t *testing.T) {
		defer func() {
			assert.Equal(t, "dawn:sql failed to connect name(sqlite): unable to open database file: is a directory", recover())
		}()
		c := config.New()
		c.Set("database", "./")
		connect("name", c)
	})
}
