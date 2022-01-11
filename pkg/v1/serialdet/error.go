package serialdet

// Error is a constant-like string error.
type Error string

func (err Error) Error() string {
	return string(err)
}

// Possible package errors.
const (
	ErrPermissionDenied         Error = "permission denied"
	ErrInvalidInformationHeader Error = "invalid information header"
	ErrInvalidRow               Error = "invalid row"
	ErrDriverNotDefined         Error = "driver is not defined"
)
