package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var (
	configPath = "./"
	configName = "foo"
	key        = "foo"
	value      = "bar"

	nonExistKey = "non"

	mergeCfg = map[string]interface{}{
		"merge": "cfg",
	}

	fu  forUnmarshal
	sub SubConfig
)

type (
	forUnmarshal struct {
		S string
		SubConfig
	}
	SubConfig struct {
		B bool
	}
)

func TestLoadReadError(t *testing.T) {
	nonConfigName := "non config name"
	err := Load(configPath, nonConfigName)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to read in")
	assert.Contains(t, err.Error(), nonConfigName)
}

func TestLoad(t *testing.T) {
	reset()

	err := Load(configPath, configName)
	assert.Nil(t, err)
	assert.Equal(t, value, GetString(key))

	defaultFileStorePath := "./data/dawn_store.db"
	assert.Equal(t, defaultFileStorePath, GetString("cache.file.path", defaultFileStorePath))
}

func TestAllGetFunctions(t *testing.T) {
	require.NoError(t, Load(configPath, configName))

	assert.Equal(t, "iface", Get("iface"))
	assert.Equal(t, "di", Get(nonExistKey, "di"))

	assert.Equal(t, "iface", GetValue("iface"))
	assert.Equal(t, "di", GetValue(nonExistKey, "di"))

	assert.Equal(t, "s", GetString("string"))
	assert.Equal(t, "ds", GetString(nonExistKey, "ds"))

	assert.Equal(t, true, GetBool("Bool"))
	assert.Equal(t, true, GetBool(nonExistKey, true))

	assert.Equal(t, time.Second, GetDuration("Duration"))
	assert.Equal(t, time.Minute, GetDuration(nonExistKey, time.Minute))

	Time, _ := time.Parse("2006-01-02 15:04:05", "2020-03-07 12:31:19")
	assert.Equal(t, Time, GetTime("Time"))
	now := time.Now()
	assert.Equal(t, now, GetTime(nonExistKey, now))

	assert.Equal(t, 1, GetInt("Int"))
	assert.Equal(t, 2, GetInt(nonExistKey, 2))

	assert.Equal(t, int64(1), GetInt64("Int"))
	assert.Equal(t, int64(2), GetInt64(nonExistKey, 2))

	assert.Equal(t, 1.1, GetFloat64("Float64"))
	assert.Equal(t, 2.2, GetFloat64(nonExistKey, 2.2))

	assert.Equal(t, map[string]interface{}{"string": "Map"}, GetStringMap("StringMap"))
	assert.Equal(t, map[string]interface{}{"k1": "v1"},
		GetStringMap(nonExistKey, map[string]interface{}{"K1": "v1"}))

	assert.Equal(t, map[string]string{"string": "String"},
		GetStringMapString("StringMapString"))
	assert.Equal(t, map[string]string{"K1": "v1"},
		GetStringMapString(nonExistKey, map[string]string{"K1": "v1"}))

	assert.Equal(t, []string{"s1", "s2"}, GetStringSlice("StringSlice"))
	assert.Equal(t, []string{"s3", "s4"}, GetStringSlice(nonExistKey, []string{"s3", "s4"}))
}

func TestAllSettings(t *testing.T) {
	reset()
	assert.Len(t, AllSettings(), 0)
}

func TestUnmarshal(t *testing.T) {
	reset()

	Set("S", value)
	err := Unmarshal(&fu)
	assert.Nil(t, err)
	assert.Equal(t, value, fu.S)
}

func TestUnmarshalKey(t *testing.T) {
	reset()

	Set("SubConfig.B", true)
	err := UnmarshalKey("SubConfig", &sub)
	assert.Nil(t, err)
	assert.True(t, sub.B)
}

func TestMergeConfigMap(t *testing.T) {
	reset()

	MergeConfigMap(mergeCfg)
	assert.Equal(t, mergeCfg["merge"], GetString("merge"))
}

func TestSub(t *testing.T) {
	reset()

	Set("SubConfig.B", true)
	c := Sub("SubConfig")
	assert.True(t, c.GetBool("B"))

	require.NoError(t, os.Setenv("SUBCONFIG_B", "false"))
	assert.False(t, c.GetBool("B"))
}

func TestHas(t *testing.T) {
	reset()

	Set("SubConfig.B", true)

	assert.True(t, Has("SubConfig"))
	assert.True(t, Has("SubConfig.B"))
	assert.False(t, Has("B"))
}

func reset() {
	global = New()
}

func TestLoadAll(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		assert.NotNil(t, LoadAll("./error"))
	})

	t.Run("success", func(t *testing.T) {
		assert.Nil(t, LoadAll("./all"))
		assert.True(t, global.Has("http"))
		assert.True(t, global.Has("auth"))
		assert.True(t, global.Has("others.1"))
	})

	t.Run("env", func(t *testing.T) {
		assert.Nil(t, LoadAll("./all"))
		assert.Equal(t, false, global.GetBool("http.accesslog"))

		LoadEnv("DAWN")

		require.NoError(t, os.Setenv("DAWN_HTTP_ACCESSLOG", "true"))

		assert.Equal(t, true, global.GetBool("http.accesslog"))
	})
}
