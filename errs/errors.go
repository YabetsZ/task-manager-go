package errs

type AppError struct {
	Code int
	Msg  string
	Err  error
}

func (e *AppError) Error() string {
	return e.Msg
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, msg string, err error) *AppError {
	return &AppError{code, msg, err}
}

var ()
