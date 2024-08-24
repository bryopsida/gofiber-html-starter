package interfaces

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
