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

// Polynomial represents a polynomial with coefficients in GF[2⁸].
// The representation is an array slice of Num values: position i
// in the array slice holds the coefficient for x^i.
type Polynomial []Num

// IsIdenticalZero returns true is p is the zero polynomial.
func (f *Field) IsIdenticalZero(p Polynomial) bool {
	for _, coefficient := range p {
		if coefficient != f.Zero() {
			return false
		}
	}
	return true
}

// Normalize removes redundant initial zero coefficients from p.
func (f *Field) Normalize(p Polynomial) Polynomial {
	i := len(p) - 1
	for ; i > 0; i-- {
		if p[i] != f.Zero() {
			break
		}
	}
	// At this point, i==0 or p[i]!=f.Zero(), or both.
	// We want to keep elements up to and including position i.
	return p[:i+1]
}

// EvaluatePolynomial evaluates the polynomial p at point x.
func (f *Field) EvaluatePolynomial(p Polynomial, x Num) Num {
	result := f.Zero()
	power := f.One()
	for _, coefficient := range p {
		result = f.Add(result, f.Mul(coefficient, power))
		power = f.Mul(power, x)
	}
	return result
}

// AddPolynomials returns p1+p2.
func (f *Field) AddPolynomials(p1, p2 Polynomial) (sum Polynomial) {
	length := 0
	if length < len(p1) {
		length = len(p1)
	}
	if length < len(p2) {
		length = len(p2)
	}
	sum = make([]Num, length)
	for i := range sum {
		sum[i] = f.Zero()
		if i < len(p1) {
			sum[i] = f.Add(sum[i], p1[i])
		}
		if i < len(p2) {
			sum[i] = f.Add(sum[i], p2[i])
		}
	}
	return sum
}

// MultiplyPolynomials returns p1×p2.
func (f *Field) MultiplyPolynomials(p1, p2 Polynomial) (product Polynomial) {
	// The code below implements long multiplication using addition and multiplication
	// in the Galois field used for the polynomial coefficients.
	product = make([]Num, len(p1)+len(p2)-1)
	for i, _ := range product {
		product[i] = f.Zero()
	}
	for i1, n1 := range p1 {
		for i2, n2 := range p2 {
			product[i1+i2] = f.Add(product[i1+i2], f.Mul(n1, n2))
		}
	}
	return product
}

// DividePolynomials returns the quotient and remainder when dividing
// nom by den, or an error if den is the zero polynomial.
func (f *Field) DividePolynomials(nom, den Polynomial) (quot, rem Polynomial, err error) {
	if f.IsIdenticalZero(den) {
		return nil, nil, fmt.Errorf("Division by zero polynomial: %v.", nom)
	}
	den = f.Normalize(den) // Ensure non-zero highest-order coefficient.
	if len(nom) < len(den) {
		return Polynomial{f.Zero()}, nom, nil
	}
	// The code below implements long division using addition and multiplication
	// in the Galois field used for the polynomial coefficients.
	rem = Polynomial(make([]Num, len(nom)))
	for i, n := range nom {
		rem[i] = n
	}
	degreeDiff := len(nom) - len(den)
	quot = Polynomial(make([]Num, degreeDiff+1))
	dInv, _ := f.Inv(den[len(den)-1])
	for i := len(quot) - 1; i >= 0; i-- {
		quot[i] = f.Mul(rem[i+len(den)-1], dInv)
		for j, n := range den {
			rem[i+j] = f.Add(rem[i+j], f.Mul(quot[i], n))
		}
	}
	return quot, f.Normalize(rem), nil
}

// ToString returns a human-readable string representation of the polynomial.
// Each coefficient is expressed in terms of the field generator.
func (f *Field) ToString(p Polynomial) string {
	var s string
	for power := len(p) - 1; power >= 0; power-- {
		n := p[power]
		if n == f.Zero() {
			continue
		}
		log, _ := f.Log(n)
		coeff := fmt.Sprintf("α^%d", log)
		switch log {
		case 0:
			coeff = "1"
		case 1:
			coeff = "α"
		}
		monomial := fmt.Sprintf("x^%d", power)
		switch power {
		case 0:
			monomial = "1"
		case 1:
			monomial = "x"
		}
		term := coeff + " " + monomial
		if log == 0 {
			term = monomial
		} else {
			if power == 0 {
				term = coeff
			}
		}
		if s == "" {
			s = term
		} else {
			s = s + " + " + term
		}
	}
	if s == "" {
		s = "0"
	}
	return s
}

func (p Polynomial) String() string {
	var s string
	for power := len(p) - 1; power >= 0; power-- {
		n := p[power]
		if n == 0 {
			continue
		}
		coeff := fmt.Sprintf("%b", n)
		monomial := fmt.Sprintf("x^%d", power)
		switch power {
		case 0:
			monomial = "1"
		case 1:
			monomial = "x"
		}
		term := coeff + " " + monomial
		if coeff == "1" {
			term = monomial
		} else {
			if power == 0 {
				term = coeff
			}
		}
		if s == "" {
			s = term
		} else {
			s = s + " + " + term
		}
	}
	if s == "" {
		s = "0"
	}
	return s
}
