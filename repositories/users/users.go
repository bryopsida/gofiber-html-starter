package users

import (
	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"gorm.io/gorm"
)

type user struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	Role         string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

// Optionally, set a custom table name
func (user) TableName() string {
	return "users"
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new userRepository instance
func NewUserRepository(db *gorm.DB) interfaces.IUserRepository {
	return &userRepository{db: db}
}

func (userRepository) FromDTO(userDTO interfaces.User) user {
	return user{
		ID:           userDTO.ID,
		Username:     userDTO.Username,
		Email:        userDTO.Email,
		Role:         userDTO.Role,
		PasswordHash: userDTO.PasswordHash,
	}
}

func (userRepository) ToDTO(user user) interfaces.User {
	return interfaces.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		PasswordHash: user.PasswordHash,
	}
}

func (r *userRepository) CreateUser(user *interfaces.User) error {
	userDb := r.FromDTO(*user)
	return r.db.Create(userDb).Error
}

func (r *userRepository) GetUserByID(id uint) (*interfaces.User, error) {
	var user user
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	var retUser = r.ToDTO(user)
	return &retUser, nil
}

func (r *userRepository) GetUserByUsername(username string) (*interfaces.User, error) {
	var user user
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	var retUser = r.ToDTO(user)
	return &retUser, nil
}

func (r *userRepository) UpdateUser(user *interfaces.User) error {
	dbuser := r.FromDTO(*user)
	return r.db.Save(dbuser).Error
}

func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Delete(&user{}, id).Error
}
