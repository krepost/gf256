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

// Generates the two GF[2⁸] tables found on pages 191–193 of “W. H. Bussey,
// Tables of Galois fields of order less than 1,000. Bulletin of the American
// Mathematical Society, 16(4):188–206, 1910”.
package main

import (
	"fmt"
	"github.com/krepost/gf256"
	"sort"
)

type output struct {
	lambda int
	binary string
}

type byBinaryString []output

func (v byBinaryString) Len() int      { return len(v) }
func (v byBinaryString) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v byBinaryString) Less(i, j int) bool {
	if len(v[i].binary) < len(v[j].binary) {
		return true
	}
	if len(v[i].binary) > len(v[j].binary) {
		return false
	}
	return v[i].binary < v[j].binary
}

func main() {
	f, _ := gf256.NewField(0x11d, 0x2)
	table1 := make([]output, 256)
	table2 := make([]output, 256)
	for i := 0; i < 256; i++ {
		table1[i] = output{
			lambda: i,
			binary: f.Exp(i).String(),
		}
		table2[i] = table1[i]
	}
	sort.Sort(byBinaryString(table2))
	fmt.Println("λ,αβγδεζηθ,λ,αβγδεζηθ")
	for i := 1; i < 256; i++ {
		fmt.Printf("%d,%s,%d,%s\n",
			table1[i].lambda, table1[i].binary,
			table2[i].lambda, table2[i].binary)
	}
}
