package st1201

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"math"
)

func TestLength(t *testing.T) {
	fp, _ := NewFPEncoderWithPrecision(0.0, 10.0, 0.1)
	assert.Equal(t, fp.fieldLength, 1)

	fp2, _ := NewFPEncoderWithPrecision(0.0, 100.0, 0.1)
	assert.Equal(t, fp2.fieldLength, 2)

	fp3, _ := NewFPEncoderWithPrecision(0.0, 2000000000.0, 1.0)
	assert.Equal(t, fp3.fieldLength, 4)

	fp4, _ := NewFPEncoderWithLength(0.0, 15000.0, 8)
	assert.Equal(t, fp4.fieldLength, 8)
}

func TestEncode(t *testing.T) {
	fp, _ := NewFPEncoderWithLength(0.0, 1e9, 8)
	value := 3.14159
	expected := []byte{0x00, 0x00, 0x00, 0x06, 0x48, 0x7e, 0x7c, 0x06}
	encoded, _ := fp.Encode(value)
	assert.Equal(t, len(encoded), 8)
	assert.Equal(t, encoded, expected)

	value2 := math.Inf(1)
    encoded2, _ := fp.Encode(value2)
    assert.Equal(t, encoded2[0], byte(0xc8))

    value3 := math.Inf(-1)
    encoded3, _ := fp.Encode(value3)
    assert.Equal(t, encoded3[0], byte(0xe8))

    value4 := math.NaN()
    encoded4, _ := fp.Encode(value4)
    assert.Equal(t, encoded4[0], byte(0xd0))

}

func TestDecode(t *testing.T) {
	fp, _ := NewFPEncoderWithLength(0.0, 1e9, 8)
	
	encoded := []byte{0x00, 0x00, 0x00, 0x06, 0x48, 0x7e, 0x7c, 0x06}
	val, _ := fp.Decode(encoded)
	assert.InEpsilon(t, val, 3.14159, 1e-8)

	posInf := []byte{0xc8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	val2, _ := fp.Decode(posInf)
	assert.Equal(t, val2, math.Inf(1))

	negInf := []byte{0xe8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	val3, _ := fp.Decode(negInf)
	assert.Equal(t, val3, math.Inf(-1))

	nan := []byte{0xd0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	val4, _ := fp.Decode(nan)
	assert.True(t, math.IsNaN(val4))

}
