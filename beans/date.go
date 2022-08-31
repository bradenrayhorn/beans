package beans

import "time"

type Date struct {
	time.Time
	set bool
}

func NewDate(date time.Time) Date {
	return Date{Time: date, set: true}
}

func (t *Date) UnmarshalJSON(b []byte) error {
	date, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
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
