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

// Between implmentes between condition
type Between struct {
	Col     string
	LessVal interface{}
	MoreVal interface{}
}

var _ Cond = Between{}

// WriteTo write data to Writer
func (between Between) WriteTo(w Writer) error {
	if _, err := fmt.Fprintf(w, "%s BETWEEN ", between.Col); err != nil {
		return err
	}
	if lv, ok := between.LessVal.(expr); ok {
		if err := lv.WriteTo(w); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprint(w, "?"); err != nil {
			return err
		}
		w.Append(between.LessVal)
	}
	if _, err := fmt.Fprint(w, " AND "); err != nil {
		return err
	}
	if mv, ok := between.MoreVal.(expr); ok {
		if err := mv.WriteTo(w); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprint(w, "?"); err != nil {
			return err
		}
		w.Append(between.MoreVal)
	}
	return nil
}

// And implments And with other conditions
func (between Between) And(conds ...Cond) Cond {
	return And(between, And(conds...))
}

// Or implments Or with other conditions
func (between Between) Or(conds ...Cond) Cond {
	return Or(between, Or(conds...))
}

// IsValid tests if the condition is valid
func (between Between) IsValid() bool {
	return len(between.Col) > 0
}
