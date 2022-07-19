package beans_test

import (
	"errors"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/stretchr/testify/assert"
)

type StructPass struct{}

func (s StructPass) Validate() error {
	return nil
}

type StructError1 struct{}

func (s StructError1) Validate() error {
	return errors.New("error 1")
}

type StructError2 struct{}

func (s StructError2) Validate() error {
	return errors.New("error 2")
}

func TestValidateOnePassingField(t *testing.T) {
	result := beans.Validate(StructPass{})
	assert.Nil(t, result)
}

func TestValidateOneFailingField(t *testing.T) {
	result := beans.Validate(StructError1{})
	assert.Equal(t, "error 1.", result.Error())
}

func TestValidateOneEachField(t *testing.T) {
	result := beans.Validate(StructError1{}, StructPass{})
	assert.Equal(t, "error 1.", result.Error())
}

func TestValidateMultipleFails(t *testing.T) {
	result := beans.Validate(StructError1{}, StructPass{}, StructError2{})
	assert.Equal(t, "error 1. error 2.", result.Error())
}
