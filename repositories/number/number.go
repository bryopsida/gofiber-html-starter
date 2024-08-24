package number

import (
	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"gorm.io/gorm"
)

// Number represents the number entity
type number struct {
	ID    string `gorm:"primaryKey"`
	Value uint64
}

// TableName returns the table name for the number entity
func (number) TableName() string {
	return "numbers"
}

// gormNumberRepository handles database operations for numbers
type gormNumberRepository struct {
	db *gorm.DB
}

// NewNumberRepository creates a new gormNumberRepository instance
func NewNumberRepository(db *gorm.DB) interfaces.INumberRepository {
	return &gormNumberRepository{db: db}
}

// Save saves a number
// - number: the number to save
// Returns an error if the save operation fails
func (r *gormNumberRepository) Save(incomingNumb interfaces.Number) error {
	num := number{
		ID:    incomingNumb.ID,
		Value: uint64(incomingNumb.Number),
	}
	return r.db.Save(&num).Error
}

// FindByID finds a number by its ID
// - id: the ID of the number to find
// Returns the number if found, otherwise returns an error
func (r *gormNumberRepository) FindByID(id string) (*interfaces.Number, error) {
	var num number
	if err := r.db.First(&num, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &interfaces.Number{
		ID:     num.ID,
		Number: num.Value,
	}, nil
}

// DeleteByID deletes a number by its ID
// - id: the ID of the number to delete
// Returns an error if the delete operation fails
func (r *gormNumberRepository) DeleteByID(id string) error {
	return r.db.Delete(&number{}, "id = ?", id).Error
}
