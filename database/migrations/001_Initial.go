package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// capture point in time structure for this migration
// to ensure consistent behavior in cases when the model changes in the future
type v001setting struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"uniqueIndex"`
	Value string
}

func (v001setting) TableName() string {
	return "settings"
}

type v001number struct {
	ID    string `gorm:"primaryKey"`
	Value uint64
}

func (v001number) TableName() string {
	return "numbers"
}

type v001user struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	Role         string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

func (v001user) TableName() string {
	return "users"
}

// V001Migration represents the first migration, creates the user, settings and number repoisitories
type V001Migration struct {
	gorm.DB
}

// Up creates the user, settings and number tables
func (m *V001Migration) Up(ctx context.Context, tx *sql.Tx) error {
	mig := m.DB.Migrator()
	// needs to be consistent even when the model changes in the future
	err := mig.CreateTable(&v001setting{})
	if err != nil {
		return err
	}
	err = mig.CreateTable(&v001number{})
	if err != nil {
		return err
	}
	err = mig.CreateTable(&v001user{})
	if err != nil {
		return err
	}
	return nil
}

// Down drops the user, settings and number tables
func (m *V001Migration) Down(ctx context.Context, tx *sql.Tx) error {
	mig := m.DB.Migrator()
	err := mig.DropTable(&v001setting{})
	if err != nil {
		return err
	}

	err = mig.DropTable(&v001number{})
	if err != nil {
		return err
	}

	err = mig.DropTable(&v001user{})
	if err != nil {
		return err
	}
	return nil
}

// InitializeV001Migration initializes the first migration
func InitializeV001Migration(db gorm.DB) *V001Migration {
	migration := &V001Migration{DB: db}
	goose.AddMigrationContext(migration.Up, migration.Down)
	return migration
}
