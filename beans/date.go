package beans

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"
)

type Date struct {
	time.Time
	set bool
}

func NewDate(date time.Time) Date {
	return Date{Time: time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC), set: true}
}

func (t *Date) UnmarshalJSON(b []byte) error {
	date, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		var parseError *time.ParseError
		if errors.As(err, &parseError) {
			return &json.UnmarshalTypeError{
				Value:  string(b),
				Offset: 0,
				Type:   reflect.TypeOf(t),
			}
		}

		return err
	}
	t.Time = date
	t.set = true
	return nil
}

func (d Date) Empty() bool {
	return !d.set
}

func (d Date) String() string {
	return d.Format("2006-01-02")
}
