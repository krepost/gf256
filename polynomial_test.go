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

func ExamplePolynomial() {
	f, _ := NewField(0x11d, 0x2)
	p := Polynomial{0xff, 0x01, 0x00, 0x17, 0x02, 0x01}
	fmt.Println(p)
	fmt.Println(f.ToString(p))
	// Output:
	// x^5 + 10 x^4 + 10111 x^3 + x + 11111111
	// x^5 + α x^4 + α^129 x^3 + x + α^175
}

func ExampleComputeWithPolynomials() {
	f, _ := NewField(0x11d, 0x2)
	p1 := Polynomial{0xff, 0x01, 0x00, 0x17, 0x02, 0x01}
	p2 := Polynomial{0x01, 0x00, 0x01}
	fmt.Println(p1)
	fmt.Println(p2)
	fmt.Println(f.EvaluatePolynomial(p1, 0x02))
	fmt.Println(f.EvaluatePolynomial(p2, 0x02))
	fmt.Println(f.AddPolynomials(p1, p2))
	fmt.Println(f.AddPolynomials(p2, p1))
	fmt.Println(f.MultiplyPolynomials(p1, p2))
	fmt.Println(f.MultiplyPolynomials(p2, p1))
	// Output:
	// x^5 + 10 x^4 + 10111 x^3 + x + 11111111
	// x^2 + 1
	// 1000101
	// 101
	// x^5 + 10 x^4 + 10111 x^3 + x^2 + x + 11111110
	// x^5 + 10 x^4 + 10111 x^3 + x^2 + x + 11111110
	// x^7 + 10 x^6 + 10110 x^5 + 10 x^4 + 10110 x^3 + 11111111 x^2 + x + 11111111
	// x^7 + 10 x^6 + 10110 x^5 + 10 x^4 + 10110 x^3 + 11111111 x^2 + x + 11111111
}

func ExampleLongDivision() {
	f, _ := NewField(0x11d, 0x2)
	nominator := Polynomial{0xff, 0x01, 0x00, 0x17, 0x02, 0x01}
	denominator := Polynomial{0x01, 0x00, 0x01}
	quotient, remainder, _ := f.DividePolynomials(nominator, denominator)
	fmt.Println(quotient)
	fmt.Println(remainder)
	// Output:
	// x^3 + 10 x^2 + 10110 x + 10
	// 10111 x + 11111101
}

func ExampleLongDivisionZeroQuotient() {
	f, _ := NewField(0x11d, 0x2)
	nominator := Polynomial{0x01, 0x00, 0x01}
	denominator := Polynomial{0xff, 0x01, 0x00, 0x17, 0x02, 0x01}
	quotient, remainder, _ := f.DividePolynomials(nominator, denominator)
	fmt.Println(quotient)
	fmt.Println(remainder)
	// Output:
	// 0
	// x^2 + 1
}

func ExampleLongDivisionSameDegree() {
	f, _ := NewField(0x11d, 0x2)
	nominator := Polynomial{0x17, 0x01, 0x02}
	denominator := Polynomial{0x01, 0x00, 0x04}
	quotient, remainder, _ := f.DividePolynomials(nominator, denominator)
	fmt.Println(quotient)
	fmt.Println(remainder)
	// Output:
	// 10001110
	// x + 10011001
}

func ExampleLongDivisionIgnoreLeadingZeros() {
	f, _ := NewField(0x11d, 0x2)
	nominator := Polynomial{0x17, 0x01, 0x02}
	denominator := Polynomial{0x04, 0x00, 0x00}
	quotient, remainder, _ := f.DividePolynomials(nominator, denominator)
	fmt.Println(quotient)
	fmt.Println(remainder)
	// Output:
	// 10001110 x^2 + 1000111 x + 11001100
	// 0
}

func ExampleLongDivisionZeroDenominator() {
	f, _ := NewField(0x11d, 0x2)
	nominator := Polynomial{0x17, 0x01, 0x02}
	denominator := Polynomial{0x00, 0x00, 0x00}
	_, _, err := f.DividePolynomials(nominator, denominator)
	fmt.Println(err)
	// Output:
	// Division by zero polynomial: 10 x^2 + x + 10111.
}
