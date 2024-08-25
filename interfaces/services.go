package interfaces

import (
	"github.com/golang-jwt/jwt/v5"
)

// IIncrementService is an interface for incrementing numbers
type IIncrementService interface {
	// Increment increments a number
	// - id: the ID of the number to increment
	// Returns the incremented number if successful, otherwise returns an error
	Increment(id string) (uint64, error)
}

// IPasswordService is an interface for password hashing and verification
type IPasswordService interface {
	// Hash hashes a plaintext password
	Hash(plaintext string) (string, error)
	// Verify verifies a plaintext password against an encoded hash
	Verify(plaintext, encodedHash string) (bool, error)
}

// ISettingsService is an interface fetching and setting settings
type ISettingsService interface {
	// GetString gets a string setting by key
	// - key: the key of the setting to get
	// Returns the setting value if found, otherwise returns an error
	GetString(key string) (string, error)
	// GetInt gets an integer setting by key
	// - key: the key of the setting to get
	// Returns the setting value if found, otherwise returns an error
	GetInt(key string) (int, error)
	// GetBool gets a boolean setting by key
	// - key: the key of the setting to get
	// Returns the setting value if found, otherwise returns an error
	GetBool(key string) (bool, error)
	// SetString sets a string setting by key
	// - key: the key of the setting to set
	// - value: the value to set
	// Returns an error if the set operation fails
	SetString(key string, value string) error
	// SetInt sets an integer setting by key
	// - key: the key of the setting to set
	// - value: the value to set
	// Returns an error if the set operation fails
	SetInt(key string, value int) error
	// SetBool sets a boolean setting by key
	// - key: the key of the setting to set
	// - value: the value to set
	// Returns an error if the set operation fails
	SetBool(key string, value bool) error
}

// IUsersService is an interface for user operations
type IUsersService interface {
	// CreateUser creates a new user
	// - user: the user to create
	// Returns an error if the create operation fails
	CreateUser(user *User) error
	// GetUserByID gets a user by ID
	// - id: the ID of the user to get
	// Returns the user if found, otherwise returns an error
	GetUserByID(id uint) (*User, error)
	// GetUserByUsername gets a user by username
	// - username: the username of the user to get
	// Returns the user if found, otherwise returns an error
	GetUserByUsername(username string) (*User, error)
	// UpdateUser updates a user
	// - user: the user to update
	// Returns an error if the update operation fails
	UpdateUser(user *User) error
	// DeleteUser deletes a user by ID
	// - id: the ID of the user to delete
	// Returns an error if the delete operation fails
	DeleteUser(id uint) error
}

// IJWTService is an interface for JWT operations
type IJWTService interface {
	// Generate generates a JWT token for a user
	// - user: the user to generate a token for
	// Returns the generated token if successful, otherwise returns an error
	Generate(user *User) (string, error)
	// Validate validates a JWT token
	// - token: the token to validate
	// Returns the token if valid, otherwise returns an error
	Validate(token string) (*jwt.Token, error)

	UserFromClaims(ctx IRequestContext) (*User, error)
}
