package beans

import "errors"

type Validatable interface {
	Validate() error
}

func Validate(objects ...Validatable) error {
	message := ""
	for _, o := range objects {
		if err := o.Validate(); err != nil {
			message += err.Error()
		}
	}

	if message == "" {
		return nil
	}

	return errors.New(message)
}
