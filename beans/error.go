package beans

const (
	EINTERNAL = "internal"
	EINVALID  = "invalid"
)

var (
	ErrorInternal = &beansError{code: EINTERNAL, msg: "internal error"}
	ErrorInvalid  = &beansError{code: EINVALID, msg: "invalid data provided"}
)

type Error interface {
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
	beansError *beansError
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

func WrapError(err error, beansError *beansError) error {
	return wrappedBeansError{error: err, beansError: beansError}
}
