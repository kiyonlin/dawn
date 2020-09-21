package sql

import (
	"context"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/kiyonlin/dawn"
	"github.com/kiyonlin/dawn/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func Test_Sql_New(t *testing.T) {
	t.Parallel()

	modular := New()
	_, ok := modular.(*sqlModule)
	assert.True(t, ok)
}

func Test_Sql_Module_Name(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "dawn:sql", m.String())
}

func Test_Sql_Init(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		m.conns = map[string]*gorm.DB{}

		at := assert.New(t)

		var wg sync.WaitGroup
		wg.Add(1)
		m.Init(&wg)
		wg.Wait()

		at.Equal(fallback, m.fallback)
		at.Len(m.conns, 1)
	})

	t.Run("with config", func(t *testing.T) {
		m.conns = map[string]*gorm.DB{}

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
	m.conns = map[string]*gorm.DB{}

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
	t.Run("unknown driver", func(t *testing.T) {
		defer func() {
			assert.Equal(t, "dawn:sql unknown driver test of name", recover())
		}()
		c := config.New()
		c.Set("driver", "test")
		connect("name", c)
	})

	t.Run("sqlite", func(t *testing.T) {
		c := config.New()
		gdb := connect("name", c)

		gdb.Logger.LogMode(logger.Info)
		ctx := context.Background()
		gdb.Logger.Info(ctx, "info")
		gdb.Logger.Warn(ctx, "warn")
		gdb.Logger.Trace(ctx, time.Now(), func() (string, int64) { return "", 0 }, nil)
	})

	t.Run("mysql", func(t *testing.T) {
		defer func() {
			assert.Contains(t, recover(),
				"dawn:sql failed to connect name(mysql):")
		}()

		c := config.New()
		c.Set("Driver", "mysql")
		c.Set("ParseTime", false)
		connect("name", c)
	})

	t.Run("postgres", func(t *testing.T) {
		defer func() {
			assert.Contains(t, recover(),
				"dawn:sql failed to connect name(postgres):")
		}()

		c := config.New()
		c.Set("Driver", "postgres")
		connect("name", c)
	})
}
