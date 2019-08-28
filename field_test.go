// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gf256

import "fmt"
import "testing"

func ExampleNum() {
	n := Num(0x17)
	fmt.Println(n)
	// Output: 10111
}

func ExampleField() {
	f, _ := NewField(0x11d, 0x2)
	fmt.Println(f.Polynomial())
	fmt.Println(f.Generator())
	fmt.Println(f.Zero())
	fmt.Println(f.One())
	// Output:
	// x⁸+x⁴+x³+x²+1
	// 10
	// 0
	// 1
}

func ExampleField_Exp() {
	f, _ := NewField(0x11d, 0x2)
	fmt.Println(f.Exp(0))
	fmt.Println(f.Exp(1))
	fmt.Println(f.Exp(17))
	fmt.Println(f.Exp(51))
	fmt.Println(f.Exp(255))
	// Output:
	// 1
	// 10
	// 10011000
	// 1010
	// 1
}

func ExampleField_Log() {
	f, _ := NewField(0x11d, 0x2)
	n := Num(0x0a)
	l, _ := f.Log(n)
	fmt.Println(n)
	fmt.Println(l)
	// Output:
	// 1010
	// 51
}

func ExampleField_Inv() {
	f, _ := NewField(0x11d, 0x2)
	n, _ := f.Inv(Num(0x0a))
	fmt.Println(n)
	// Output:
	// 11011101
}

func ExampleField_Add() {
	f, _ := NewField(0x11d, 0x2)
	x, y := Num(0x0a), Num(0x1f)
	fmt.Println(x, y, f.Add(x, y))
	// Output:
	// 1010 11111 10101
}

func ExampleField_Mul() {
	f, _ := NewField(0x11d, 0x2)
	x, y := Num(0x0a), Num(0x1f)
	fmt.Println(x, y, f.Mul(x, y))
	// Output:
	// 1010 11111 11000110
}

func TestToString(t *testing.T) {
	testData := []struct {
		coefficients uint
		toString     string
	}{
		{0x00, "0"},
		{0x01, "1"},
		{0x02, "10"},
		{0x03, "11"},
		{0x17, "10111"},
		{0x11d, "100011101"},
		{0x310, "1100010000"},
	}
	for _, data := range testData {
		if s := Num(data.coefficients).String(); s != data.toString {
			t.Errorf("For Num(%v): expected %s, got %s",
				data.coefficients, data.toString, s)
		}
	}
}

func TestNewFieldWithZeroGenerator(t *testing.T) {
	_, err := NewField(0x11d, 0x0)
	if err == nil {
		t.Errorf("Expected error return value from NewField().")
	}
	if err.Error() != "0 is not a generator." {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestNewFieldWithUnitGenerator(t *testing.T) {
	_, err := NewField(0x11d, 0x1)
	if err == nil {
		t.Errorf("Expected error return value from NewField().")
	}
	if err.Error() != "1 is not a generator." {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestNewFieldWithBadGenerator(t *testing.T) {
	_, err := NewField(0x11d, 0x20)
	if err == nil {
		t.Errorf("Expected error return value from NewField().")
	}
	if err.Error() != "100000 is not a generator." {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestNewFieldWithLowDegreePolynomial(t *testing.T) {
	_, err := NewField(0x3, 0x2)
	if err == nil {
		t.Errorf("Expected error return value from NewField().")
	}
	if err.Error() != "x+1 has too low degree." {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestNewFieldWithHighDegreePolynomial(t *testing.T) {
	_, err := NewField(0x200, 0x2)
	if err == nil {
		t.Errorf("Expected error return value from NewField().")
	}
	if err.Error() != "x^9 has too high degree." {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestNewFieldWithReduciblePolynomial(t *testing.T) {
	_, err := NewField(0x101, 0x02) // x⁸+1 == (x+1)⁸ in Z₂[x].
	if err == nil {
		t.Errorf("Expected error return value from NewField().")
	}
	if err.Error() != "10 is not a generator." {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestNewFieldWithCorrectParameters(t *testing.T) {
	f, err := NewField(0x11d, 0x02)
	if err != nil {
		t.Errorf("Could not create GF[2⁸]: %v.", err)
	}
	if f == nil {
		t.Error("Unexpected nil return value from NewField().")
	}
}

func TestFieldAccessors(t *testing.T) {
	f, err := NewField(0x11d, 0x02)
	if err != nil {
		t.Errorf("Could not create GF[2⁸]: %v.", err)
		return // Avoid crashing due to dereferencing nil below.
	}
	zero := f.Zero().String()
	one := f.One().String()
	generator := f.Generator().String()
	polynomial := f.Polynomial().String()
	if zero != "0" {
		t.Errorf("Unexpected zero: expected 0, got %v.", zero)
	}
	if one != "1" {
		t.Errorf("Unexpected one: expected 1, got %v.", one)
	}
	if generator != "10" {
		t.Errorf("Unexpected generator: expected 10, got %v.", generator)
	}
	if polynomial != "x⁸+x⁴+x³+x²+1" {
		t.Errorf("Unexpected irreducible polynomial: expected x⁸+x⁴+x³+x²+1, got %v.", polynomial)
	}
}

func TestArithmeticOperators(t *testing.T) {
	f, err := NewField(0x11d, 0x02)
	if err != nil {
		t.Errorf("Could not create GF[2⁸]: %v.", err)
		return // Avoid crashing due to dereferencing nil below.
	}
	for i := uint(0); i < 256; i++ {
		x := Num(i)
		if y := f.Add(x, f.Zero()); x != y {
			t.Errorf("Error adding with zero: expected %v, got %v.", x, y)
		}
		if y := f.Mul(x, f.Zero()); f.Zero() != y {
			t.Errorf("Error multiplying by one: expected %v, got %v.", f.Zero(), y)
		}
		if y := f.Mul(x, f.One()); x != y {
			t.Errorf("Error multiplying by one: expected %v, got %v.", x, y)
		}
		if x != f.Zero() {
			inv, err := f.Inv(x)
			if err != nil {
				t.Errorf("Error computing inverse of %v: %v", x, err)
			}
			if y := f.Mul(x, inv); f.One() != y {
				t.Errorf("Error multiplying %v by its inverse: expected 1, got %v.", x, y)
			}
		}
	}
	for i := uint(1); i < 256; i++ {
		for j := uint(1); j < 256; j++ {
			x := Num(i)
			y := Num(j)
			logX, errX := f.Log(x)
			logY, errY := f.Log(y)
			if errX != nil {
				t.Errorf("Error computing logarithm of %v: %v", x, errX)
			}
			if errY != nil {
				t.Errorf("Error computing logarithm of %v: %v", y, errY)
			}
			expected := f.Exp(logX + logY)
			actual := f.Mul(x, y)
			if expected != actual {
				t.Errorf("%v × %v: expected %v, got %v.", x, y, expected, actual)
			}
		}
	}
}

func TestAddition(t *testing.T) {
	f, err := NewField(0x11d, 0x02)
	if err != nil {
		t.Errorf("Could not create GF[2⁸]: %v.", err)
		return // Avoid crashing due to dereferencing nil below.
	}
	testData := []struct {
		term1       Num
		term2       Num
		expectedSum Num
	}{
		{0x02, 0x04, 0x06}, // x + x² == x²+x.
		{0x05, 0x11, 0x14}, // x²+1 + x⁴+1 == x⁴+x².
		{0x80, 0x80, 0x00}, // x⁷ + x⁷ == 0.
		{0x7f, 0x1f, 0x60}, // x⁶+x⁵+x⁴+x³+x²+x+1 + x⁴+x³+x²+x+1 == x⁶+x⁵.
	}
	for _, data := range testData {
		actualSum := f.Add(data.term1, data.term2)
		if data.expectedSum != actualSum {
			t.Errorf("%v + %v: expected %v, actual %v.", data.term1, data.term2, data.expectedSum, actualSum)
		}
	}
}

func TestInverse(t *testing.T) {
	f, err := NewField(0x11d, 0x02)
	if err != nil {
		t.Errorf("Could not create GF[2⁸]: %v.", err)
		return // Avoid crashing due to dereferencing nil below.
	}
	testData := []struct {
		number          Num
		expectedInverse Num
	}{
		{0x02, 0x8e}, // 1 / x == x⁷+x³+x²+x.
		{0x05, 0xa7}, // 1 / x²+1 == x⁷+x⁵+x²+x+1.
		{0xba, 0x07}, // 1 / x⁷+x⁵+x⁴+x³+x == x²+x+1.
		{0x80, 0x1b}, // 1 / x⁷ == x⁴+x³+x+1.
		{0xff, 0xfd}, // 1 / x⁷+x⁶+x⁵+x⁴+x³+x²+x+1 == x⁷+x⁶+x⁵+x⁴+x³+x²+1.
	}
	for _, data := range testData {
		actualInverse, err := f.Inv(data.number)
		if err != nil {
			t.Errorf("1 / %v: got error %v", data.number, err)
		}
		if data.expectedInverse != actualInverse {
			t.Errorf("1 / %v: expected %v, actual %v.", data.number, data.expectedInverse, actualInverse)
		}
	}
}

func TestMultiplication(t *testing.T) {
	f, err := NewField(0x11d, 0x02)
	if err != nil {
		t.Errorf("Could not create GF[2⁸]: %v.", err)
		return // Avoid crashing due to dereferencing nil below.
	}
	testData := []struct {
		factor1         Num
		factor2         Num
		expectedProduct Num
	}{
		{0x02, 0x04, 0x08}, // x × x² == x³.
		{0x05, 0x11, 0x55}, // x²+1 × x⁴+1 == x⁶+x⁴+x²+1.
		{0x80, 0x80, 0x13}, // x⁷ × x⁷ == x⁴+x+1.
		{0x7f, 0x19, 0x03}, // x⁶+x⁵+x⁴+x³+x²+x+1 × x⁴+x³+x²+x+1 == x+1.
		{0xff, 0xff, 0xe2}, // x⁷+x⁶+x⁵+x⁴+x³+x²+x+1 × x⁷+x⁶+x⁵+x⁴+x³+x²+x+1 == x⁷+x⁶+x⁵+x.
	}
	for _, data := range testData {
		actualProduct := f.Mul(data.factor1, data.factor2)
		if data.expectedProduct != actualProduct {
			t.Errorf("%v × %v: expected %v, actual %v.", data.factor1, data.factor2, data.expectedProduct, actualProduct)
		}
	}
}
