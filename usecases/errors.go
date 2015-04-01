package usecases

import "fmt"

// ErrorDomain defines domain of an error
type ErrorDomain string

const (
	// ErrorDomainDataAccess represents data layer domain
	ErrorDomainDataAccess ErrorDomain = "DataAccess"
	// ErrorDomainUseCase represents use-cases/interactors layer domain
	ErrorDomainUseCase ErrorDomain = "UseCase"
	// ErrorDomainApplication represents handler/controller/cli/app layer domain
	ErrorDomainApplication ErrorDomain = "Application"
)

// ErrorType defines type of an error
type ErrorType string

const (
	// ErrorTypeNotFound represents entity not found error
	ErrorTypeNotFound ErrorType = "NOT FOUND"
	// ErrorTypeConflict represents constrains violations / duplication cases
	ErrorTypeConflict ErrorType = "CONFLICT"
	// ErrorTypeOperational is usually an internal server error (unexpected)
	ErrorTypeOperational ErrorType = "OPERATIONAL"
)

// Error represents general error at use-case level
type Error struct {
	Domain ErrorDomain
	Type   ErrorType
	Msg    string
	Cause  error
}

// NewError constructs a new error structure
func NewError(errDomain ErrorDomain, errType ErrorType, msg string, cause error) *Error {
	return &Error{
		Domain: errDomain,
		Type:   errType,
		Msg:    msg,
		Cause:  cause,
	}
}

// NewDataAccessError constructs error of type ErrorTypeNotFound
func NewDataAccessError(errType ErrorType, msg string, cause error) *Error {
	return NewError(ErrorDomainDataAccess, errType, msg, cause)
}

// NewUseCaseError constructs error of type ErrorTypeConflict
func NewUseCaseError(errType ErrorType, msg string, cause error) *Error {
	return NewError(ErrorDomainUseCase, errType, msg, cause)
}

// NewApplicationError constructs error of type ErrorTypeOperational
func NewApplicationError(errType ErrorType, msg string, cause error) *Error {
	return NewError(ErrorDomainApplication, errType, msg, cause)
}

func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	if err.Cause != nil {
		return fmt.Sprintf("%v::%v ERROR: %v | CAUSED BY %v", err.Domain, err.Type, err.Msg, err.Cause)
	}
	return fmt.Sprintf("%v::%v ERROR: %v", err.Domain, err.Type, err.Msg)
}
