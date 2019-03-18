package types

import (
	"fmt"
)

type Int128 struct {
	high int64
	low int64
}

func NewInt128(s string) *Int128 {
	fmt.Println("....")
	// Empty string bad.
	if len(s) == 0 {
		return nil
	}

	// Pick off leading sign.
	neg := false
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	// Convert unsigned.
	un := NewUInt128(s)

	newValue := Int128 {
		high: 0,
		low: 0,
	}

	if neg {
		un.twosComplement(un)

		high := int64(un.high)
		low := int64(un.low)
		newValue.high = int64(high)
		newValue.low = int64(low)
	} else {
		high := int64(un.high)
		low := int64(un.low)
		newValue.high = int64(high)
		newValue.low = int64(low)
	}

	return &newValue
}

func (value *Int128) IsNegative() bool {
	return isNegative(value)
}
