package beans

const (
	EINTERNAL     = "internal"
	EINVALID      = "invalid"
	ENOTFOUND     = "not_found"
	EUNAUTHORIZED = "unauthorized"
)

var (
	ErrorInternal     = &beansError{code: EINTERNAL, msg: "Internal error"}
	ErrorInvalid      = &beansError{code: EINVALID, msg: "Invalid data provided"}
	ErrorNotFound     = &beansError{code: ENOTFOUND, msg: "Not found"}
	ErrorUnauthorized = &beansError{code: EUNAUTHORIZED, msg: "Not authenticated"}
)

var codeToError = map[string]*beansError{
	EINTERNAL:     ErrorInternal,
	EINVALID:      ErrorInvalid,
	ENOTFOUND:     ErrorNotFound,
	EUNAUTHORIZED: ErrorUnauthorized,
}

func NewError(code string, msg string) Error {
	err := codeToError[code]
	return WrapError(err, &beansError{code, msg})
}

type Error interface {
	error
	BeansError() (string, string)
}

type beansError struct {
	code string
	msg  string
}

func (e beansError) Error() string {
	return e.msg
}

func (e beansError) BeansError() (string, string) {
	return e.code, e.msg
}

type wrappedBeansError struct {
	error
	beansError Error
}

func (e wrappedBeansError) Is(err error) bool {
	return e.beansError == err
}

func (e wrappedBeansError) Unwrap() error {
	return e.error
}

func (e wrappedBeansError) BeansError() (string, string) {
	return e.beansError.BeansError()
}

func WrapError(err error, parent Error) Error {
	return wrappedBeansError{error: err, beansError: parent}
}
