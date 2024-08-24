package database

import (
	"gorm.io/gorm"
)

var (
	// DBConn is the database connection, can be referenced in case its not passed but passing it is recommended
	DBConn *gorm.DB
)
