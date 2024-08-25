package migrations

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// V004Migration represents the second migration, adds a jwt signing key to the settings table
type V004Migration struct {
	gorm.DB
}

// Up adds a jwt signing key to the settings table
func (m *V004Migration) Up(ctx context.Context, tx *sql.Tx) error {
	// Generate a random 32-byte string
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return err
	}
	encodedValue := base64.StdEncoding.EncodeToString(randomBytes)
	settingsModel := m.DB.Model(&v001setting{})
	err = settingsModel.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&v001setting{
		Key:   "jwt_signing_key",
		Value: encodedValue,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// Down removes the cookie encryption key from the settings table
func (m *V004Migration) Down(ctx context.Context, tx *sql.Tx) error {
	// Remove the row with the key cookie_encryption_key from the settings table using GORM
	settingsModel := m.DB.Model(&v001setting{})
	err := settingsModel.Where("key = ?", "jwt_signing_key").Delete(&v001setting{}).Error
	if err != nil {
		return err
	}

	return nil
}

// InitializeV004Migration initializes the V004Migration
func InitializeV004Migration(db gorm.DB) *V004Migration {
	migration := &V004Migration{DB: db}
	goose.AddMigrationContext(migration.Up, migration.Down)
	return migration
}
