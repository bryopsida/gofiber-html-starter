package migrations

import (
	"context"
	"database/sql"

	"github.com/bryopsida/gofiber-pug-starter/services/password"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// V003Migration represents the third migration, creates the default admin user
type V003Migration struct {
	gorm.DB
}

// Up creates the default admin user
func (m *V003Migration) Up(ctx context.Context, tx *sql.Tx) error {
	userModel := m.DB.Model(&v001user{})
	// Initialize the password service
	passwordService := password.NewPasswordService()

	// Hash the password "admin"
	passwordHash, err := passwordService.Hash("admin")
	if err != nil {
		return err
	}

	// Create the admin user
	adminUser := v001user{
		Username:     "admin",
		Email:        "admin@localhost",
		Role:         "admin",
		PasswordHash: passwordHash,
	}

	// Insert the admin user into the database
	if err := userModel.Create(&adminUser).Error; err != nil {
		return err
	}

	return nil
}

// Down removes the default admin user
func (m *V003Migration) Down(ctx context.Context, tx *sql.Tx) error {
	userModel := m.DB.Model(&v001user{})

	// Remove the admin user from the database
	if err := userModel.Where("username = ?", "admin").Delete(&v001user{}).Error; err != nil {
		return err
	}

	return nil
}

// InitializeV003Migration initializes the V003Migration
func InitializeV003Migration(db gorm.DB) *V003Migration {
	migration := &V003Migration{DB: db}
	goose.AddMigrationContext(migration.Up, migration.Down)
	return migration
}
