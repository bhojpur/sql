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

type expr struct {
	sql  string
	args []interface{}
}

var _ Cond = expr{}

// Expr generate customerize SQL
func Expr(sql string, args ...interface{}) Cond {
	return expr{sql, args}
}
func (expr expr) OpWriteTo(op string, w Writer) error {
	return expr.WriteTo(w)
}
func (expr expr) WriteTo(w Writer) error {
	if _, err := fmt.Fprint(w, expr.sql); err != nil {
		return err
	}
	w.Append(expr.args...)
	return nil
}
func (expr expr) And(conds ...Cond) Cond {
	return And(expr, And(conds...))
}
func (expr expr) Or(conds ...Cond) Cond {
	return Or(expr, Or(conds...))
}
func (expr expr) IsValid() bool {
	return len(expr.sql) > 0
}
