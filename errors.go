package postalcode

import (
	"errors"
	"fmt"
)

// Common errors
var (
	// ErrNotFound is returned when postal code is not found
	ErrNotFound = errors.New("postal code not found")

	// ErrInvalidZipCode is returned when zip code format is invalid
	ErrInvalidZipCode = errors.New("invalid zip code format")

	// ErrInvalidZipPrefix is returned when zip prefix format is invalid
	ErrInvalidZipPrefix = errors.New("invalid zip prefix format: must be 3 digits")

	// ErrEmptyFile is returned when import file is empty
	ErrEmptyFile = errors.New("import file is empty")

	// ErrInvalidFileFormat is returned when file format is incorrect
	ErrInvalidFileFormat = errors.New("invalid file format")

	// ErrDatabaseConnection is returned when database connection fails
	ErrDatabaseConnection = errors.New("database connection failed")

	// ErrInvalidSearchParams is returned when search parameters are invalid
	ErrInvalidSearchParams = errors.New("invalid search parameters")
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ImportError represents an import operation error
type ImportError struct {
	Line    int
	Message string
	Err     error
}

func (e *ImportError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("import error at line %d: %s - %v", e.Line, e.Message, e.Err)
	}
	return fmt.Sprintf("import error at line %d: %s", e.Line, e.Message)
}

func (e *ImportError) Unwrap() error {
	return e.Err
}

// NewImportError creates a new import error
func NewImportError(line int, message string, err error) error {
	return &ImportError{
		Line:    line,
		Message: message,
		Err:     err,
	}
}
