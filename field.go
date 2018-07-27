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

// Package gf256 implements arithmetic over the finite field GF[2⁸] as well as
// over the polynomial ring with coefficients in GF[2⁸].
package gf256

import "fmt"

// Num is a bit-vector representation of the polynomial used to represent
// numbers in GF[2⁸]. Concretely, values of Num will be unsigned integers
// between 0 and 255.
type Num uint

// Irreducible is a bit-vector representation of the irreducible polynomial
// used to define GF[2⁸]. This will be an unsigned integer with the ninth bit
// set and no higher-order bits set.
type Irreducible uint

// Field represents an instantiation of GF[2⁸].
type Field struct {
	// poly is a bit-vector representation of the irreducible
	// polynomial in Z₂[x] which define the irreducible congruence
	// used to define the field. Bit i in the bit-vector
	// represents the coefficient of term x^i in the irreducible
	// polynomial. A commonly used representation of GF[2⁸] is
	// Z₂[x]/(x⁸+x⁴+x³+x²+1) which corresponds to a bit-vector
	// 100011101, or 0x11d in hexadecimal, or 285 in decimal.
	poly Irreducible
	// g is the generator used for multiplication and division.
	// A common choice is x, which corresponds to the bit-vector
	// 10, or 2 in decimal.
	g Num
	// expTable[i] == g^i is built in NewField.
	expTable [255]Num
	// logtable[i] == log_g i is built in NewField; logtable[g^i] == i.
	logTable [256]int
}

// Zero returns the additive zero of the field f.
func (f *Field) Zero() Num {
	return Num(0)
}

// One returns the multiplicative unit of the field f.
func (f *Field) One() Num {
	return Num(1)
}

// Generator returns the generator used when defining the field f.
func (f *Field) Generator() Num {
	return f.g
}

// Polynomial returns the irreducible polynomial used when defining the field f.
func (f *Field) Polynomial() Irreducible {
	return f.poly
}

// Exp returns the generator of the field f raised to the power x.
func (f *Field) Exp(x int) Num {
	x = x % 255
	if x < 0 {
		x = x + 255
	}
	return f.expTable[x%255]
}

// Log returns the logarithm of x with respect to the generator of the
// field f, or an error if x==0.
func (f *Field) Log(x Num) (int, error) {
	if x == f.Zero() {
		return 0, fmt.Errorf("Taking log of zero.")
	}
	return f.logTable[x], nil
}

// Inv returns the multiplicative inverse of x, or an error if x==0.
func (f *Field) Inv(x Num) (Num, error) {
	if x == f.Zero() {
		return f.Zero(), fmt.Errorf("Taking inverse of zero.")
	}
	logX, _ := f.Log(x)
	return f.Exp(-logX), nil
}

// Add(x, y) returns the sum of x and y in the field f.
func (f *Field) Add(x, y Num) Num {
	return x ^ y
}

// Mul returns the product of x and y in the field f.
func (f *Field) Mul(x, y Num) Num {
	if x == f.Zero() || y == f.Zero() {
		return f.Zero()
	}
	logX, _ := f.Log(x)
	logY, _ := f.Log(y)
	return f.Exp(logX + logY)
}

// String returns a readable string representation of the number n in GF[2⁸].
func (n Num) String() string {
	return fmt.Sprintf("%b", uint(n))
}

// String returns a readable string representation of the irreducible
// polynomial p.
func (p Irreducible) String() string {
	return bitmaskToString(uint(p))
}

// NewField creates a new version of GF[2⁸] using the supplied
// irreducible polynomial and generator.
func NewField(polynomial Irreducible, generator Num) (*Field, error) {
	if polynomial|0x1FF != 0x1FF {
		return nil, fmt.Errorf("%v has too high degree.", polynomial)
	}
	if polynomial&0x100 == 0 {
		return nil, fmt.Errorf("%v has too low degree.", polynomial)
	}
	if generator == 0 || generator == 1 {
		return nil, fmt.Errorf("%v is not a generator.", generator)
	}
	f := &Field{
		poly: polynomial,
		g:    generator,
	}
	// Build expTable and logTable.
	for n := 0; n < 256; n++ {
		// Fill with zeroes to have know values everywhere.
		f.logTable[n] = 0
	}
	product := Num(0x01) // The number 1.
	for i := 0; i < 255; i++ {
		if i != 0 && product == 1 {
			return nil, fmt.Errorf("%v is not a generator.", f.g)
		}
		f.expTable[i] = product
		f.logTable[product] = i
		product = multiply(product, f.g, f.poly)
	}
	// Double-check that the generator has generated all of GF[2⁸]
	// by checking that every number other then zero and one has
	// non-zero logarithm.
	for n := 2; n < 256; n++ {
		if f.logTable[n] == 0 {
			return nil, fmt.Errorf("%v is not a generator.", f.g)
		}
	}
	return f, nil
}

func multiply(x, y Num, poly Irreducible) Num {
	// Repeated squaring; optimize for small y.
	product := Num(0)
	for y != 0 {
		if y&0x01 != 0 {
			product = product ^ x
		}
		x = x << 1
		y = y >> 1
	}
	// Reduce modulo the irreducible polynomial.
	// Casting poly to Num is fine since both Num
	// and Irreducible are represented as uint.
	poly_msb := msb(uint(poly))
	for product >= 256 {
		product_msb := msb(uint(product))
		product = product ^ (Num(poly) << (product_msb - poly_msb))
	}
	return product
}

func msb(n uint) uint {
	i := uint(0)
	for {
		n = n >> 1
		if n == 0 {
			return i
		}
		i = i + 1
	}
}

func bitmaskToString(n uint) string {
	terms := []string{"1", "x", "x²", "x³", "x⁴", "x⁵", "x⁶", "x⁷", "x⁸"}
	if n == 0 {
		return "0"
	}
	i := 0
	s := ""
	for n != 0 {
		if n&0x01 != 0 {
			nextTerm := fmt.Sprintf("x^%d", i)
			if i < 9 {
				nextTerm = terms[i]
			}
			if s != "" {
				s = nextTerm + "+" + s
			} else {
				s = nextTerm
			}
		}
		n = n >> 1
		i = i + 1
	}
	return s
}
