// Copyright 2022, Initialize All Once Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iteng

import (
	"testing"
)

type min_test struct {
	a   int
	b   int
	exp int
}

var min_tests = []min_test{
	{a: 1, b: 2, exp: 1},
	{a: 2, b: 1, exp: 1},
	{a: 3, b: 3, exp: 3},
}

func Test_min(t *testing.T) {
	for _, test := range min_tests {
		ret := min(test.a, test.b)
		if ret != test.exp {
			t.Errorf("min(%d,%d) = %d; expected %d", test.a, test.b, ret, test.exp)
		}
	}
}

type minf_test struct {
	a   float64
	b   float64
	exp float64
}

var minf_tests = []minf_test{
	{a: 1.0, b: 2.0, exp: 1.0},
	{a: 2.0, b: 1.0, exp: 1.0},
	{a: 3.0, b: 3.0, exp: 3.0},
}

func Test_minf(t *testing.T) {
	for _, test := range minf_tests {
		ret := minf(test.a, test.b)
		if ret != test.exp {
			t.Errorf("minf(%f,%f) = %f; expected %f", test.a, test.b, ret, test.exp)
		}
	}
}

var maxf_tests = []minf_test{
	{a: 1.0, b: 2.0, exp: 2.0},
	{a: 2.0, b: 1.0, exp: 2.0},
	{a: 3.0, b: 3.0, exp: 3.0},
}

func Test_maxf(t *testing.T) {
	for _, test := range maxf_tests {
		ret := maxf(test.a, test.b)
		if ret != test.exp {
			t.Errorf("maxf(%f,%f) = %f; expected %f", test.a, test.b, ret, test.exp)
		}
	}
}
