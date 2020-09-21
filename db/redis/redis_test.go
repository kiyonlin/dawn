package redis

import (
	"sync"
	"testing"

	"github.com/kiyonlin/dawn/config"

	"github.com/go-redis/redis/v8"
	"github.com/kiyonlin/dawn"
	"github.com/stretchr/testify/assert"
)

func Test_Redis_New(t *testing.T) {
	modular := New()
	_, ok := modular.(*redisModule)
	assert.True(t, ok)
}

func Test_Redis_Module_Name(t *testing.T) {
	assert.Equal(t, "dawn:redis", m.String())
}

func Test_Redis_Init(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		m.Init(&wg)
		wg.Wait()
	})

	t.Run("error", func(t *testing.T) {
		defer func() {
			assert.Contains(t, recover(), "dawn:redis failed to ping")
		}()
		config.Load("./", "redis")
		config.Set("Redis.Connections.Default.Addr", "127.0.0.1:99999")

		var wg sync.WaitGroup
		wg.Add(1)
		m.Init(&wg)
		wg.Wait()
	})
}

func Test_Redis_Boot(t *testing.T) {
	m.conns = map[string]*redis.Client{
		fallback: redis.NewClient(&redis.Options{}),
	}

	var (
		wg      sync.WaitGroup
		cleanup = make(chan dawn.Cleanup, 1)
	)
	wg.Add(1)
	m.Boot(&wg, cleanup)
	wg.Wait()

	(<-cleanup)()
}

func Test_Redis_Conn(t *testing.T) {
	assert.Nil(t, Conn("non"))
}
