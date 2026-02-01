package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Init initializes a GORM database connection.
// For simplicity in this common lib, we use a basic sqlite driver,
// but in a real production lib, this might be more generic or configurable.
func Init(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}
