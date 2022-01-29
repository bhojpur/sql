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

// Like defines like condition
type Like [2]string

var _ Cond = Like{"", ""}

// WriteTo write SQL to Writer
func (like Like) WriteTo(w Writer) error {
	if _, err := fmt.Fprintf(w, "%s LIKE ?", like[0]); err != nil {
		return err
	}
	// FIXME: if use other regular express, this will be failed. but for compatible, keep this
	if like[1][0] == '%' || like[1][len(like[1])-1] == '%' {
		w.Append(like[1])
	} else {
		w.Append("%" + like[1] + "%")
	}
	return nil
}

// And implements And with other conditions
func (like Like) And(conds ...Cond) Cond {
	return And(like, And(conds...))
}

// Or implements Or with other conditions
func (like Like) Or(conds ...Cond) Cond {
	return Or(like, Or(conds...))
}

// IsValid tests if this condition is valid
func (like Like) IsValid() bool {
	return len(like[0]) > 0 && len(like[1]) > 0
}
