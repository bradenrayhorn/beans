package beans

import (
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

type Amount struct {
	decimal apd.Decimal
	set     bool
}

func NewAmount(coefficient int64, exponent int32) Amount {
	return Amount{decimal: *apd.New(coefficient, exponent), set: true}
}

func NewEmptyAmount() Amount {
	return Amount{decimal: *apd.New(0, 0), set: false}
}

func NewAmountWithBigInt(coefficient *big.Int, exponent int32) Amount {
	bigInt := &apd.BigInt{}
	bigInt.SetMathBigInt(coefficient)
	return Amount{decimal: *apd.NewWithBigInt(bigInt, exponent), set: true}
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

// If the amount is empty returns a new zero amount.
func (a *Amount) OrZero() Amount {
	if a.Empty() {
		return NewAmount(0, 0)
	}

	return *a
}

func (a *Amount) String() string {
	if !a.set {
		return ""
	}
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
	if strings.TrimSpace(amountString.String()) == "" {
		return nil
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

func (a Amount) MarshalJSON() ([]byte, error) {
	type res struct {
		Coefficient *big.Int `json:"coefficient"`
		Exponent    int32    `json:"exponent"`
	}
	if !a.set {
		return json.Marshal(nil)
	}
	return json.Marshal(res{Coefficient: a.Coefficient(), Exponent: a.Exponent()})
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

// positive rule

type validatablePositive struct {
	Amount
}

func Positive(a Amount) validatablePositive {
	return validatablePositive{a}
}

func (v validatablePositive) Validate() error {
	if v.Amount.decimal.Negative {
		return errors.New(":field must be positive")
	}

	return nil
}

// non zero rule

type validatableNonZero struct {
	Amount
}

func NonZero(a Amount) validatableNonZero {
	return validatableNonZero{a}
}

func (v validatableNonZero) Validate() error {
	if v.Amount.set && v.Amount.decimal.IsZero() {
		return errors.New(":field must not be zero")
	}

	return nil
}