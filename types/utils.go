package types

import (
	//"bytes"
	//"encoding/binary"
	//"errors"
	//"fmt"
	"strconv"
)

func decimalStringToBinaryString(str string) string {
	var newString string
	if (str == "0") {
		newString = "0"
	} else {
		newString = ""
		for (str != "0") {
			lastChar := str[len(str) - 1]
			remainder := (int(lastChar) - int('0')) % 2
			newString = strconv.Itoa(remainder) + newString
			str = divideByTwo(str)
		}
	}
	return newString
}

func divideByTwo(str string) string {
    newString := ""
	newDigit := 0
	add := 0

	for _, ch := range str {
        newDigit = (int(ch) - int('0')) / 2 + add
        newString = newString + strconv.Itoa(newDigit)
		if ((int(ch) - int('0')) % 2 == 1) {
			add = 5
		} else {
			add = 0
		}
	}

    if (string(newString) != "0" && newString[0:1] == "0") {
		newString = newString[1:]
	}

	return newString
}

func compare(a *UInt128, b *UInt128) int {
	if (a.high < b.high) {
		return -1
	} else if (a.high > b.high) {
		return 1
	} else if (a.low < b.low) {
		return -1
	} else if (a.low > b.low) {
		return 1
	} else {
		return 0
	}
}

func isZero(val *UInt128) bool {
	return val.high == 0 && val.low == 0
}

func isNegative(val *Int128) bool {
	return (val.high & -0x8000000000000000) != 0
}


func (value *UInt128) not(val *UInt128) {
	val.high = (^val.high) >> 0
	val.low = (^val.high) >> 0
}

func (value *UInt128) or(a *UInt128, b *UInt128) {
	a.high = (a.high | b.high) >> 0
	a.low = (a.low | b.low) >> 0
}

func (value *UInt128) and(a *UInt128, b *UInt128) {
	a.high = (a.high & b.high) >> 0
	a.low = (a.low & b.low) >> 0
}

func (value *UInt128) xor(a *UInt128, b *UInt128) {
	a.high = (a.high ^ b.high) >> 0
	a.low = (a.low ^ b.low) >> 0
}

func (value *UInt128) rightShiftUnsigned(val *UInt128, bits uint) {
	if (bits >= 128) {
		val.high = 0
		val.low = 0
		return
	}

	if (bits < 64) {
		mask := value.genMask(bits);
		shifted := (val.high & mask) >> 0
		val.high = val.high >> bits
		val.low = ((val.low >> bits) | (shifted << (64 - bits))) >> 0
		return
	}

	bits = bits - 64;
	val.low = (val.high >> bits)
	val.high = 0;
}

func (value *UInt128) leftShift(val *UInt128, bits uint) {
	if (bits >= 128) {
		val.high = 0
		val.low = 0
		return
	}
	
	if ( bits < 64 ) {
		mask := (^value.genMask(64 - bits)) >> 0
		shifted := (val.low & mask) >> (64 - bits)
		val.low = (val.low << bits) >> 0
		val.high = (val.high << bits | shifted) >> 0
		return
	}
	
	bits = bits - 64
	val.high = (val.low << bits) >> 0
	val.low = 0;
}

func (value *UInt128) add(a *UInt128, b *UInt128) {
	var carry uint64 = 0
	low := a.low + b.low
	if (a.low > UINT64_MAX - b.low) {
		carry = 1
	}

	high := a.high + b.high + carry

	a.high = high
	a.low = low
}

func (value *UInt128) sub(a *UInt128, b *UInt128) {
	newB := UInt128 {
		high: b.high,
		low: b.low,
	}
	
	value.twosComplement(&newB)
	value.add(a, &newB)
}

func (value *UInt128) multiply(a *UInt128, b *UInt128) {
	ans := UInt128 {
		high: 0,
		low: 0,
	}

	bits := value.nbits(b)

	for i := 0; i < bits; i++ {
		if (b.low & 1 == 1) {
			value.add(&ans, a)
		}
		value.leftShift(a, 1)
		value.rightShiftUnsigned(b, 1)
	}

	a.high = ans.high
	a.low = ans.low
}

func (value *UInt128) divide(a *UInt128, b *UInt128) *UInt128 {
	quotient := UInt128 {
		high: 0,
		low: 0,
	}
	remainder := UInt128 {
		high: a.high,
		low: a.low,
	}
	divider := UInt128 {
		high: b.high,
		low: b.low,
	}

	if (isZero(b)) {
		return nil
	}
	if (compare(a, b) < 0) {
		a.high = 0
		a.low = 0
		return &quotient
	}

	var mask uint64 = 0x8000000000000000
	var dPadding uint = 0
	var rPadding uint = 0
	var count uint = 128

	for (count > 0) {
		if ((remainder.high & mask) != 0) {
			break
		}

		value.leftShift(&remainder, 1)
		rPadding++
		count--
	}
	
	remainder.high = a.high
	remainder.low = a.low

	count = 128
	for (count > 0) {
		if ( (divider.high & mask) != 0 ) {
			break
		}
		
		value.leftShift(&divider, 1)
		dPadding++
		count--
	}
	value.rightShiftUnsigned(&divider, rPadding)

	count = dPadding - rPadding + 1
	for (count > 0) {
		count--

		if (compare(&remainder, &divider) >= 0) {
			value.sub(&remainder, &divider)
			quotient.low = quotient.low | 1
		}
		if (count > 0) {
			value.leftShift(&quotient, 1)
			value.rightShiftUnsigned(&divider, 1)
		}
	}

	a.high = quotient.high
	a.low = quotient.low

	return &remainder
}

func (value *UInt128) nbits(val *UInt128) int {
	bits := 0
	high := val.high
	low := val.low

	if (high == 0) {
		for (low > 0) {
			low = low >> 1
			bits++
		}
		return bits
	}

	for (high > 0) {
		high = high >> 1
		bits++
	}
	return bits + 64
}

func (value *UInt128) twosComplement(val *UInt128) {
	val.high = (^val.high) >> 0
	val.low = (^val.low) >> 0

	var carry uint64 = 0
	low := val.low + 1
	if (val.low > UINT64_MAX - 1) {
		carry = 1
	}

	high := val.high + carry

	val.high = high
	val.low = low
}

func (value *UInt128) toBinaryString(val *UInt128) string {
	strLow := strconv.FormatUint(val.low, 2)
	if (val.high == 0) {
		return strLow
	}

	strHigh := strconv.FormatUint(val.high, 2)
	str := strHigh + value.paddingZero(strLow, 64);

	return str
}

func (value *UInt128) toHexString(val *UInt128) string {
	strLow := strconv.FormatUint(val.low, 16)
	if (val.high == 0) {
		return strLow
	}

	strHigh := strconv.FormatUint(val.high, 16)
	str := strHigh + value.paddingZero(strLow, 16);

	return str
}

// TODO: uncomplete
func (value *UInt128) ToDecimalString(val *UInt128) string {
	var output []string

	quotient := UInt128 {
		high: val.high,
		low: val.low,
	}
	
	var slc []string
	div := NewUInt128("1844674407370955161")

	for (!isZero(&quotient)) {
		slc = append(slc, strconv.FormatUint(quotient.low, 10))
		output = append(slc, output...)
		value.divide(&quotient, div)
	}
	
	if ( len(output) == 0 ) {
		return "0"
	} else {
		x := ""
		for _, comp := range output {
			x = x + value.paddingZero(comp, 9)
		}
		
		return x
	}
	return "0"
}

func (value *UInt128) paddingZero(data string, length int) string {
	zeros := length - len(data)
	padded := ""
	for (zeros > 0) {
		padded = padded + "0";
		zeros--
	}

	return padded + data;
}

func (value *UInt128) genMask(bits uint) uint64 {
	if (bits > 64) {
		bits = 64
	}
	if (bits < 0) {
		bits = 0
	}

	var val uint64 = 0
	for (bits > 0) {
		val = ((val << 1) | 1) >> 0
		bits--
	}
	return val;
}
