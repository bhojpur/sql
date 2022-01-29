package builder

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import "fmt"

type condAnd []Cond

var _ Cond = condAnd{}

// And generates AND conditions
func And(conds ...Cond) Cond {
	var result = make(condAnd, 0, len(conds))
	for _, cond := range conds {
		if cond == nil || !cond.IsValid() {
			continue
		}
		result = append(result, cond)
	}
	return result
}
func (and condAnd) WriteTo(w Writer) error {
	for i, cond := range and {
		_, isOr := cond.(condOr)
		_, isExpr := cond.(expr)
		wrap := isOr || isExpr
		if wrap {
			fmt.Fprint(w, "(")
		}
		err := cond.WriteTo(w)
		if err != nil {
			return err
		}
		if wrap {
			fmt.Fprint(w, ")")
		}
		if i != len(and)-1 {
			fmt.Fprint(w, " AND ")
		}
	}
	return nil
}
func (and condAnd) And(conds ...Cond) Cond {
	return And(and, And(conds...))
}
func (and condAnd) Or(conds ...Cond) Cond {
	return Or(and, Or(conds...))
}
func (and condAnd) IsValid() bool {
	return len(and) > 0
}
