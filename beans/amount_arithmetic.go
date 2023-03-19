package beans

import "github.com/cockroachdb/apd/v3"

type amountArithmetic struct{}

var Arithmetic amountArithmetic

func (a amountArithmetic) Add(amounts ...Amount) (Amount, error) {
	accumulator := apd.New(0, 0)

	for _, b := range amounts {
		_, err := apd.BaseContext.Add(accumulator, accumulator, &b.decimal)
		if err != nil {
			return NewEmptyAmount(), err
		}

	}

	bigInt := accumulator.Coeff.MathBigInt()
	if accumulator.Negative {
		bigInt.Neg(bigInt)
	}

	return NewAmountWithBigInt(bigInt, accumulator.Exponent), nil
}

func (a amountArithmetic) Negate(amount Amount) Amount {
	if amount.Empty() || amount.decimal.IsZero() {
		return NewAmount(0, 0)
	}

	newDecimal := amount.decimal
	newDecimal.Neg(&newDecimal)
	return Amount{set: true, decimal: newDecimal}
}
