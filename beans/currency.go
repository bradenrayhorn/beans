package beans

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/cockroachdb/apd/v3"
)

type Amount struct {
	decimal apd.Decimal
	set     bool
}

func NewAmount(coefficient int64, exponent int32) Amount {
	return Amount{decimal: *apd.New(coefficient, exponent), set: true}
}

func (a *Amount) Exponent() int32 {
	return a.decimal.Exponent
}

func (a *Amount) Coefficient() *big.Int {
	bigInt := a.decimal.Coeff.MathBigInt()
	if a.decimal.Negative {
		bigInt.Neg(bigInt)
	}

	return bigInt
}

func (a *Amount) String() string {
	return a.decimal.Text('f')
}

func (a *Amount) Empty() bool {
	return !a.set
}

func (a *Amount) UnmarshalJSON(b []byte) error {
	var amountString json.Number
	if err := json.Unmarshal(b, &amountString); err != nil {
		return err
	}

	dec, condition, err := apd.NewFromString(amountString.String())
	if err != nil {
		return err
	}
	if condition != 0 {
		return errors.New("invalid amount")
	}

	a.decimal = *dec
	a.set = true
	return nil
}

// max precision rule

type validatableMaxPrecision struct {
	Amount
}

func MaxPrecision(a Amount) validatableMaxPrecision {
	return validatableMaxPrecision{a}
}

func (m validatableMaxPrecision) Validate() error {
	if m.Amount.Exponent() < -2 {
		return errors.New(":field must have at most 2 decimal points")
	}

	return nil
}
