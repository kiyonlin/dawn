package sql

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kiyonlin/dawn"
	"github.com/kiyonlin/dawn/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	m        = &sqlModule{conns: make(map[string]*gorm.DB)}
	fallback = "testing"
)

type sqlModule struct {
	dawn.Module
	conns    map[string]*gorm.DB
	fallback string
}

// New gets the modular
func New() dawn.Modular {
	return m
}

// String is module name
func (*sqlModule) String() string {
	return "dawn:sql"
}

// Init does connection work to each database by config:
//  Default = "testing"
//  [Connections]
//  [Connections.testing]
//  Driver = "sqlite"
//  [Connections.mysql]
//  Driver = "mysql"
func (m *sqlModule) Init(wg *sync.WaitGroup) {
	defer wg.Done()

	// extract sql config
	c := config.Sub("sql")

	m.fallback = c.GetString("default", fallback)

	connsConfig := c.GetStringMap("connections")

	if len(connsConfig) == 0 {
		m.conns[m.fallback] = connect(m.fallback, config.New())
		return
	}

	// connect each db in config
	for name := range connsConfig {
		cfg := c.Sub("connections." + name)
		m.conns[name] = connect(name, cfg)
	}
}

func connect(name string, c *config.Config) (db *gorm.DB) {
	driver := c.GetString("driver", "sqlite")

	var err error
	switch strings.ToLower(driver) {
	case "sqlite":
		db, err = resolveSqlite(c)
	default:
		panic(fmt.Sprintf("dawn:sql unknown driver %s of %s", driver, name))
	}

	if err != nil || db == nil {
		panic(fmt.Sprintf("dawn:sql failed to connect %s(%s): %v", name, driver, err))
	}

	return
}

// Boot does general ping to each connections and
// sets cleanup to close each connections
func (m *sqlModule) Boot(wg *sync.WaitGroup, cleanup chan<- dawn.Cleanup) {
	defer wg.Done()

	cleanup <- func() {
		// close every connections
		for _, gdb := range m.conns {
			if db, err := gdb.DB(); err == nil {
				_ = db.Close()
			}
		}
	}
}

// Conn gets sql connection by specific name or fallback
func Conn(name ...string) *gorm.DB {
	n := m.fallback

	if len(name) > 0 {
		n = name[0]
	}

	return m.conns[n]
}

var l disabledLogger

type disabledLogger struct{}

func (disabledLogger) LogMode(logger.LogLevel) logger.Interface {
	return disabledLogger{}
}
func (disabledLogger) Info(context.Context, string, ...interface{})                    {}
func (disabledLogger) Warn(context.Context, string, ...interface{})                    {}
func (disabledLogger) Error(context.Context, string, ...interface{})                   {}
func (disabledLogger) Trace(context.Context, time.Time, func() (string, int64), error) {}
