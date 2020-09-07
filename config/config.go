package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type (
	Config struct {
		v   *viper.Viper
		mut sync.RWMutex
	}

	// A DecoderConfigOption can be passed to Unmarshal to configure
	// mapstructure.DecoderConfig options
	DecoderConfigOption func(*mapstructure.DecoderConfig)
)

var (
	global *Config
)

func init() {
	global = New()
}

func New() *Config {
	return &Config{v: viper.New()}
}

// Load config
func Load(configPath string, configName ...string) *Config {
	v := viper.New()

	name := "config"
	if len(configName) > 0 {
		name = configName[0]
	}

	v.SetConfigName(name)
	v.AddConfigPath(configPath)
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config: failed to read in %s: %w", name, err))
	}

	global = &Config{v: v}

	return global
}

// Load all config contents in the dir path
func LoadAll(configPath string) error {
	return filepath.Walk(configPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dir, filename := filepath.Split(path)
			name := strings.TrimSuffix(filename, filepath.Ext(filename))

			v := viper.New()

			v.SetConfigName(name)
			v.AddConfigPath(dir)
			if err := v.ReadInConfig(); err != nil {
				return errors.Wrapf(err, "failed to read in %s", path)
			}

			rel, _ := filepath.Rel(configPath, path)

			global.MergeConfigMap(configMap(getKeys(rel), v.AllSettings()))
		}
		return nil
	})
}

// configMap combines configuration recursively
func configMap(keys []string, value interface{}) map[string]interface{} {
	if len(keys) == 1 {
		return map[string]interface{}{keys[0]: value}
	}
	return map[string]interface{}{keys[0]: configMap(keys[1:], value)}
}

func getKeys(path string) (keys []string) {
	path = strings.TrimSuffix(path, filepath.Ext(path))

	b := new(bytes.Buffer)
	for i := 0; i < len(path); i++ {
		if path[i] == os.PathSeparator {
			//b.WriteByte('.')
			keys = append(keys, b.String())
			b.Reset()
		} else {
			b.WriteByte(path[i])
		}
	}

	keys = append(keys, b.String())

	return
}

func Get(key string, defaultValue ...interface{}) interface{} {
	return global.Get(key, defaultValue...)
}
func (c *Config) Get(key string, defaultValue ...interface{}) interface{} {
	c.mut.RLock()
	defer c.mut.RUnlock()

	return c.GetValue(key, defaultValue...)
}

func GetValue(key string, defaultValue ...interface{}) interface{} {
	return global.GetValue(key, defaultValue...)
}
func (c *Config) GetValue(key string, defaultValue ...interface{}) interface{} {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.Get(key)
}

func GetBool(key string, defaultValue ...bool) bool {
	return global.GetBool(key, defaultValue...)
}
func (c *Config) GetBool(key string, defaultValue ...bool) bool {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.GetBool(key)
}

func GetFloat64(key string, defaultValue ...float64) float64 {
	return global.GetFloat64(key, defaultValue...)
}
func (c *Config) GetFloat64(key string, defaultValue ...float64) float64 {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.GetFloat64(key)
}

func GetInt(key string, defaultValue ...int) int {
	return global.GetInt(key, defaultValue...)
}
func (c *Config) GetInt(key string, defaultValue ...int) int {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.GetInt(key)
}

func GetInt64(key string, defaultValue ...int64) int64 {
	return global.GetInt64(key, defaultValue...)
}
func (c *Config) GetInt64(key string, defaultValue ...int64) int64 {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.GetInt64(key)
}

func GetString(key string, defaultValue ...string) string {
	return global.GetString(key, defaultValue...)
}
func (c *Config) GetString(key string, defaultValue ...string) string {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.GetString(key)
}

func GetStringMap(key string, defaultValue ...map[string]interface{}) map[string]interface{} {
	return global.GetStringMap(key, defaultValue...)
}
func (c *Config) GetStringMap(key string, defaultValue ...map[string]interface{}) map[string]interface{} {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.GetStringMap(key)
}

func GetStringMapString(key string, defaultValue ...map[string]string) map[string]string {
	return global.GetStringMapString(key, defaultValue...)
}
func (c *Config) GetStringMapString(key string, defaultValue ...map[string]string) map[string]string {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.GetStringMapString(key)
}

func GetStringSlice(key string, defaultValue ...[]string) []string {
	return global.GetStringSlice(key, defaultValue...)
}
func (c *Config) GetStringSlice(key string, defaultValue ...[]string) []string {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.GetStringSlice(key)
}

func GetTime(key string, defaultValue ...time.Time) time.Time {
	return global.GetTime(key, defaultValue...)
}
func (c *Config) GetTime(key string, defaultValue ...time.Time) time.Time {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.GetTime(key)
}

func GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	return global.GetDuration(key, defaultValue...)
}
func (c *Config) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	c.mut.RLock()
	defer c.mut.RUnlock()

	if len(defaultValue) > 0 {
		c.v.SetDefault(key, defaultValue[0])
	}

	return c.v.GetDuration(key)
}

func AllSettings() map[string]interface{} {
	return global.AllSettings()
}
func (c *Config) AllSettings() map[string]interface{} {
	c.mut.RLock()
	defer c.mut.RUnlock()

	return c.v.AllSettings()
}

func Unmarshal(rawVal interface{}, opts ...DecoderConfigOption) error {
	return global.Unmarshal(rawVal, opts...)
}
func (c *Config) Unmarshal(rawVal interface{}, opts ...DecoderConfigOption) error {
	c.mut.RLock()
	defer c.mut.RUnlock()

	var args []viper.DecoderConfigOption
	for _, opt := range opts {
		args = append(args, viper.DecoderConfigOption(opt))
	}

	return c.v.Unmarshal(rawVal, args...)
}

func UnmarshalKey(key string, rawVal interface{}, opts ...DecoderConfigOption) error {
	return global.UnmarshalKey(key, rawVal, opts...)
}
func (c *Config) UnmarshalKey(key string, rawVal interface{}, opts ...DecoderConfigOption) error {
	c.mut.RLock()
	defer c.mut.RUnlock()

	var args []viper.DecoderConfigOption
	for _, opt := range opts {
		args = append(args, viper.DecoderConfigOption(opt))
	}

	return c.v.UnmarshalKey(key, rawVal, args...)
}

func MergeConfigMap(cfg map[string]interface{}) {
	global.MergeConfigMap(cfg)
}
func (c *Config) MergeConfigMap(cfg map[string]interface{}) {
	c.mut.Lock()
	defer c.mut.Unlock()

	_ = c.v.MergeConfigMap(cfg)
}

func Sub(key string) *Config {
	return global.Sub(key)
}

// Sub gets sub config
func (c *Config) Sub(key string) *Config {
	c.mut.RLock()
	defer c.mut.RUnlock()

	var newConf *Config
	if v := c.v.Sub(key); v != nil {
		newConf = &Config{v: v}
	} else {
		newConf = New()
	}

	return newConf
}

func Set(key string, value interface{}) {
	global.Set(key, value)
}
func (c *Config) Set(key string, value interface{}) {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.v.Set(key, value)
}

func Has(key string) bool {
	return global.Has(key)
}
func (c *Config) Has(key string) bool {
	c.mut.RLock()
	defer c.mut.RUnlock()

	return c.v.IsSet(key)
}

// LoadEnv loads env from .env or command line
// Use prefix to avoid conflicts with other env variables
// Same config key in env will override that in config file
func LoadEnv(prefix ...string) {
	global.LoadEnv(prefix...)
}
func (c *Config) LoadEnv(prefix ...string) {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.v.AutomaticEnv()
	if len(prefix) > 0 && prefix[0] != "" {
		c.v.SetEnvPrefix(prefix[0])
	}
	replacer := strings.NewReplacer(".", "_")
	c.v.SetEnvKeyReplacer(replacer)
}
