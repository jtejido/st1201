// Encoding and decoding of floating point values per MISB ST 1201.3 (5 October 2017)
// http://www.gwg.nga.mil/misb/st_pubs.html
//
// Scope:
// This Standard (ST) describes the method for mapping floating-point values to integer values and
// the reverse, mapping integer values back to their original floating-point value to within an
// acceptable precision. There are many ways of optimizing the transmission of floating-point
// values from one system to another; the purpose of this ST is to provide a single method which
// can be used for all floating-point ranges and valid precisions. This ST provides a method for a
// forward and reverse linear mapping of a specified range of floating-point values to a specified
// integer range of values based on the number of bytes desired to be used for the integer value.
// Additionally, it provides a set of special values which can be used to transmit non-numerical
// “signals” to a receiving system.
package st1201

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

type FPEncoder struct {
	a, b, bPow, dPow, sF, sR, zOffset float64
	fieldLength                       int
}

var logOf2 float64 = math.Log(2.0)

// Construct an encoder with the desired field length
// min - The minimum floating point value to be encoded
// max - The maximum floating point value to be encoded
// length - The field length, in bytes (1, 2, 4, or 8)
func NewFPEncoderWithLength(min, max float64, length int) (fpe *FPEncoder, err error) {
	fpe = new(FPEncoder)

	if length == 1 || length == 2 || length == 4 || length == 8 {
		fpe.preCompute(min, max, length)
		return
	}

	return nil, fmt.Errorf("Only 1, 2, 4, and 8 are valid field lengths")
}

// Construct an encoder with the desired precision, automatically selecting field length
// min - The minimum floating point value to be encoded
// max - The maximum floating point value to be encoded
// precision - The required precision
func NewFPEncoderWithPrecision(min, max, precision float64) (fpe *FPEncoder, err error) {
	fpe = new(FPEncoder)

	bits := math.Ceil(log2((max-min)/precision) + 1)

	length := math.Ceil(bits / 8)

	if length <= 2 {
		fpe.preCompute(min, max, int(length))
		return
	} else if length <= 4 {
		fpe.preCompute(min, max, 4)
		return
	} else if length <= 8 {
		fpe.preCompute(min, max, 8)
		return
	}

	return nil, fmt.Errorf("The specified range and precision cannot be represented using a 64-bit integer")
}

// Encode a floating point value as a byte array
// Note: Positive and negative infinity and NaN will be encoded by setting special flags defined by the ST.
func (fpe *FPEncoder) Encode(val float64) (encoded []byte, err error) {

	encoded = make([]byte, fpe.fieldLength)

	if val == math.Inf(0) {
		encoded[0] = 0xc8
	} else if val == math.Inf(-1) {
		encoded[0] = 0xe8
	} else if math.IsNaN(val) {
		encoded[0] = 0xd0
	} else if val < fpe.a || val > fpe.b {
		return nil, fmt.Errorf("Value must be in range [ %v, %v]", fpe.a, fpe.b)
	} else {
		b := new(bytes.Buffer)
		d := math.Floor(fpe.sF*(val-fpe.a) + fpe.zOffset)
		switch fpe.fieldLength {
		case 1:
			b.WriteByte(byte(d))
			break
		case 2:
			binary.Write(b, binary.BigEndian, uint16(d))
			break
		case 4:
			binary.Write(b, binary.BigEndian, uint32(d))
			break
		case 8:
			binary.Write(b, binary.BigEndian, uint64(d))
			break
		}

		encoded = b.Bytes()
	}
	return
}

// Decode a byte array containing an encoded floating point value
// bytes - The encoded array
func (fpe *FPEncoder) Decode(bytes []byte) (val float64, err error) {

	if len(bytes) != fpe.fieldLength {
		err = fmt.Errorf("Array length does not match expected field length")
		return
	} else if bytes[0] == 0xc8 {
		val = math.Inf(0)
	} else if bytes[0] == 0xe8 {
		val = math.Inf(-1)
	} else if bytes[0] == 0xd0 {
		val = math.NaN()
	} else {
		var l float64
		switch fpe.fieldLength {
		case 1:
			l = float64(int(bytes[0]))
			break
		case 2:
			l = float64(binary.BigEndian.Uint16(bytes))
			break
		case 4:
			l = float64(binary.BigEndian.Uint32(bytes))
			break
		case 8:
			l = float64(binary.BigEndian.Uint64(bytes))
			break
		}

		val = fpe.sR*(l-fpe.zOffset) + fpe.a

		if val < fpe.a || val > fpe.b {
			return 0.0, fmt.Errorf("Error decoding floating point value; out of range")
		}
	}

	return
}

// Compute constants used for encoding and decoding
// min - The minimum floating point value to be encoded
// max - The maximum floating point value to be encoded
// length - The field length, in bytes
// Section 8.9
func (fpe *FPEncoder) preCompute(min, max float64, length int) {
	fpe.fieldLength = length
	fpe.a = min
	fpe.b = max
	fpe.bPow = math.Ceil(log2(fpe.b - fpe.a))
	fpe.dPow = float64(8*fpe.fieldLength - 1)
	fpe.sF = math.Pow(2, fpe.dPow-fpe.bPow)
	fpe.sR = math.Pow(2, fpe.bPow-fpe.dPow)

	if fpe.a < 0 && fpe.b > 0 {
		fpe.zOffset = fpe.sF*fpe.a - math.Floor(fpe.sF*fpe.a)
	}
}

func log2(val float64) float64 {
	return math.Log(val) / logOf2
}
