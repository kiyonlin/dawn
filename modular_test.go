package dawn

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockModule struct {
	Module
}

// go test -run Test_Modular_Embed_Empty_Module -race
func Test_Modular_Embed_Empty_Module(t *testing.T) {
	t.Parallel()

	module := mockModule{Module{}}

	assert.Implements(t, (*Modular)(nil), module)

	assert.Equal(t, "anonymous", module.String())

	wg := &sync.WaitGroup{}
	wg.Add(2)

	module.Init(wg)

	cleanup := make(chan Cleanup, 1)
	module.Boot(wg, cleanup)
	(<-cleanup)()

	wg.Wait()

	module.RegisterRoutes(nil)
}
