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

type condOr []Cond

var _ Cond = condOr{}

// Or sets OR conditions
func Or(conds ...Cond) Cond {
	var result = make(condOr, 0, len(conds))
	for _, cond := range conds {
		if cond == nil || !cond.IsValid() {
			continue
		}
		result = append(result, cond)
	}
	return result
}

// WriteTo implments Cond
func (o condOr) WriteTo(w Writer) error {
	for i, cond := range o {
		var needQuote bool
		switch cond.(type) {
		case condAnd, expr:
			needQuote = true
		case Eq:
			needQuote = (len(cond.(Eq)) > 1)
		case Neq:
			needQuote = (len(cond.(Neq)) > 1)
		}
		if needQuote {
			fmt.Fprint(w, "(")
		}
		err := cond.WriteTo(w)
		if err != nil {
			return err
		}
		if needQuote {
			fmt.Fprint(w, ")")
		}
		if i != len(o)-1 {
			fmt.Fprint(w, " OR ")
		}
	}
	return nil
}
func (o condOr) And(conds ...Cond) Cond {
	return And(o, And(conds...))
}
func (o condOr) Or(conds ...Cond) Cond {
	return Or(o, Or(conds...))
}
func (o condOr) IsValid() bool {
	return len(o) > 0
}
