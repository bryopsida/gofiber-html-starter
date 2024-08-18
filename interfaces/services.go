package interfaces

// IIncrementService is an interface for incrementing numbers
type IIncrementService interface {
	// Increment increments a number
	// - id: the ID of the number to increment
	// Returns the incremented number if successful, otherwise returns an error
	Increment(id string) (uint64, error)
}
