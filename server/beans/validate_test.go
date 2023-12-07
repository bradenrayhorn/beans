package beans_test

import (
	"errors"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/assert"
)

type StructPass struct{}

func (s StructPass) Validate() error {
	return nil
}

type StructError1 struct{}

func (s StructError1) Validate() error {
	return errors.New(":field error 1")
}

type StructError2 struct{}

func (s StructError2) Validate() error {
	return errors.New(":field error 2")
}

func TestValidateField(t *testing.T) {
	t.Run("one passing", func(t *testing.T) {
		result := beans.ValidateFields(beans.Field("Field 1", StructPass{}))
		assert.Nil(t, result)
	})

	t.Run("one failing", func(t *testing.T) {
		result := beans.ValidateFields(beans.Field("Field 1", StructError1{}))
		_, msg := result.(beans.Error).BeansError()
		assert.Equal(t, "Field 1 error 1.", msg)
	})

	t.Run("one failing and one passing", func(t *testing.T) {
		result := beans.ValidateFields(beans.Field("Field 1", StructError1{}, StructPass{}))
		_, msg := result.(beans.Error).BeansError()
		assert.Equal(t, "Field 1 error 1.", msg)
	})

	t.Run("multiple failing", func(t *testing.T) {
		result := beans.ValidateFields(beans.Field("Field 1", StructError1{}, StructPass{}, StructError2{}))
		_, msg := result.(beans.Error).BeansError()
		assert.Equal(t, "Field 1 error 1, Field 1 error 2.", msg)
	})

	t.Run("multiple fields", func(t *testing.T) {
		result := beans.ValidateFields(
			beans.Field("Field 1", StructError1{}, StructError2{}),
			beans.Field("Field 2", StructError2{}),
		)
		_, msg := result.(beans.Error).BeansError()
		assert.Equal(t, "Field 1 error 1, Field 1 error 2. Field 2 error 2.", msg)
	})
}

type StructEmpty struct{ empty bool }

func (s StructEmpty) Empty() bool {
	return s.empty
}

func TestRequired(t *testing.T) {
	t.Run("fails if empty", func(t *testing.T) {
		err := beans.Required(StructEmpty{empty: true}).Validate()
		assert.Equal(t, ":field is required", err.Error())
	})

	t.Run("succeeds if not empty", func(t *testing.T) {
		err := beans.Required(StructEmpty{empty: false}).Validate()
		assert.Nil(t, err)
	})
}

type StructFieldError struct{}

func (s StructFieldError) Validate() error {
	return errors.New(":field is invalid")
}

func TestField(t *testing.T) {
	t.Run("replaces field variable with name", func(t *testing.T) {
		err := beans.Field("A Field", StructFieldError{}).Validate()
		assert.Equal(t, "A Field is invalid", err.Error())
	})
}
