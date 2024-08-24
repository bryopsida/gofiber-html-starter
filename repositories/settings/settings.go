package settings

import (
	"errors"
	"strconv"

	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Setting represents a key-value pair in the settings table
type setting struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"uniqueIndex"`
	Value string
}

func (setting) TableName() string {
	return "settings"
}

// SettingsRepository handles database operations for settings
type settingsRepository struct {
	db *gorm.DB
}

// NewSettingsRepository initializes the repository with a database connection
func NewSettingsRepository(db *gorm.DB) interfaces.ISettingsRepository {
	return &settingsRepository{db: db}
}

// GetString retrieves a string value for a given key
func (r *settingsRepository) GetString(key string) (string, error) {
	var setting setting
	if err := r.db.Where("key = ?", key).First(&setting).Error; err != nil {
		return "", err
	}
	return setting.Value, nil
}

// GetInt retrieves an integer value for a given key
func (r *settingsRepository) GetInt(key string) (int, error) {
	value, err := r.GetString(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

// GetBool retrieves a boolean value for a given key
func (r *settingsRepository) GetBool(key string) (bool, error) {
	value, err := r.GetString(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(value)
}

// Set sets a value for a given key
func (r *settingsRepository) Set(key string, value interface{}) error {
	strValue := ""
	switch v := value.(type) {
	case string:
		strValue = v
	case int:
		strValue = strconv.Itoa(v)
	case bool:
		strValue = strconv.FormatBool(v)
	default:
		return errors.New("unsupported value type")
	}

	setting := setting{Key: key, Value: strValue}
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&setting).Error
}
