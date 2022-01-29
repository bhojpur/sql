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

// Not defines NOT condition
type Not [1]Cond

var _ Cond = Not{}

// WriteTo writes SQL to Writer
func (not Not) WriteTo(w Writer) error {
	if _, err := fmt.Fprint(w, "NOT "); err != nil {
		return err
	}
	switch not[0].(type) {
	case condAnd, condOr:
		if _, err := fmt.Fprint(w, "("); err != nil {
			return err
		}
	case Eq:
		if len(not[0].(Eq)) > 1 {
			if _, err := fmt.Fprint(w, "("); err != nil {
				return err
			}
		}
	case Neq:
		if len(not[0].(Neq)) > 1 {
			if _, err := fmt.Fprint(w, "("); err != nil {
				return err
			}
		}
	}
	if err := not[0].WriteTo(w); err != nil {
		return err
	}
	switch not[0].(type) {
	case condAnd, condOr:
		if _, err := fmt.Fprint(w, ")"); err != nil {
			return err
		}
	case Eq:
		if len(not[0].(Eq)) > 1 {
			if _, err := fmt.Fprint(w, ")"); err != nil {
				return err
			}
		}
	case Neq:
		if len(not[0].(Neq)) > 1 {
			if _, err := fmt.Fprint(w, ")"); err != nil {
				return err
			}
		}
	}
	return nil
}

// And implements And with other conditions
func (not Not) And(conds ...Cond) Cond {
	return And(not, And(conds...))
}

// Or implements Or with other conditions
func (not Not) Or(conds ...Cond) Cond {
	return Or(not, Or(conds...))
}

// IsValid tests if this condition is valid
func (not Not) IsValid() bool {
	return not[0] != nil && not[0].IsValid()
}
