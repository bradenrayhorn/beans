package beans

import "github.com/segmentio/ksuid"

type ID ksuid.KSUID

func (id ID) String() string {
	return ksuid.KSUID(id).String()
}

func BeansIDFromString(id string) (ID, error) {
	parsedID, err := ksuid.Parse(id)
	if err != nil {
		return ID(ksuid.Nil), err
	}
	return ID(parsedID), nil
}

func NewBeansID() ID {
	return ID(ksuid.New())
}
