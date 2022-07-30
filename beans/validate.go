package beans

import (
	"strings"
)

type Validatable interface {
	Validate() error
}

func Validate(objects ...Validatable) error {
	messages := []string{}
	for _, o := range objects {
		if err := o.Validate(); err != nil {
			messages = append(messages, err.Error()+".")
		}
	}

	if len(messages) == 0 {
		return nil
	}

	return NewError(EINVALID, strings.Join(messages, " "))
}
