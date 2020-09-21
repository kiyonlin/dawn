package sql

import (
	"github.com/kiyonlin/dawn/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// resolveSqlite resolve sqlite connections with config:
// Driver = "sqlite"
// Database = "file:dawn?mode=memory&cache=shared&_fk=1"
// Prefix = "dawn_"
// Log = false
func resolveSqlite(c *config.Config) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.GetString("Prefix"),
			SingularTable: false,
		},
	}

	// disable logger
	if !c.GetBool("Log") {
		gormConfig.Logger = l
	}

	dbname := c.GetString("Database", "file:dawn?mode=memory&cache=shared&_fk=1")

	return gorm.Open(sqlite.Open(dbname), gormConfig)
}
