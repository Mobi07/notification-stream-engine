package errors

type AppError struct {
	Err       error
	Retryable bool
}

func (e AppError) Error() string {
	return e.Err.Error()
}

// Helper constructors (very useful)
func NewRetryable(err error) AppError {
	return AppError{
		Err:       err,
		Retryable: true,
	}
}

func NewNonRetryable(err error) AppError {
	return AppError{
		Err:       err,
		Retryable: false,
	}
}