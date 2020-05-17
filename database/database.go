package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "rafiki/settings"
)

func CreateConnection() (*gorm.DB, error) {

	// Get database details from environment variables
	host := GetDatabaseHost()
	user := GetDatabaseUser()
	name := GetDatabaseName()
	password := GetDatabasePassword()

	return gorm.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			user, password, host, name,
		),
	)
}
