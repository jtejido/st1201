package main

import (
	"fmt"
	"github.com/jtejido/st1201"
	"math"
)

func testWithError(val float64) {
	fpEncoder, _ := st1201.NewFPEncoderWithLength(0.0, 1e5, 4)
	encoded, _ := fpEncoder.Encode(val)
	decoded, _ := fpEncoder.Decode(encoded)

	fmt.Printf("Encoded float with value %v and got %v, Amount of error at : %.8f", val, decoded, math.Abs(val-decoded))
}

func main() {
	pi := 3.14159
	testWithError(pi)
}
