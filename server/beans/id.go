package beans

import (
	"encoding/json"

	"github.com/segmentio/ksuid"
)

type ID ksuid.KSUID

func (id ID) String() string {
	return ksuid.KSUID(id).String()
}

func (id ID) Empty() bool {
	return ksuid.KSUID(id).IsNil()
}

func (id *ID) UnmarshalJSON(b []byte) error {
	var idString string
	if err := json.Unmarshal(b, &idString); err != nil {
		return err
	}

	parsedID, err := BeansIDFromString(idString)
	if err != nil {
		return err
	}
	*id = parsedID
	return nil
}

func (id ID) MarshalJSON() ([]byte, error) {
	if id.Empty() {
		return json.Marshal(nil)
	}
	return json.Marshal(id.String())
}

func BeansIDFromString(id string) (ID, error) {
	if id == "" {
		return ID(ksuid.Nil), nil
	}

	parsedID, err := ksuid.Parse(id)
	if err != nil {
		return ID(ksuid.Nil), err
	}
	return ID(parsedID), nil
}

func NewBeansID() ID {
	return ID(ksuid.New())
}

func EmptyID() ID {
	return ID(ksuid.Nil)
}
