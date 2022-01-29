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

// IsNull defines IS NULL condition
type IsNull [1]string

var _ Cond = IsNull{""}

// WriteTo write SQL to Writer
func (isNull IsNull) WriteTo(w Writer) error {
	_, err := fmt.Fprintf(w, "%s IS NULL", isNull[0])
	return err
}

// And implements And with other conditions
func (isNull IsNull) And(conds ...Cond) Cond {
	return And(isNull, And(conds...))
}

// Or implements Or with other conditions
func (isNull IsNull) Or(conds ...Cond) Cond {
	return Or(isNull, Or(conds...))
}

// IsValid tests if this condition is valid
func (isNull IsNull) IsValid() bool {
	return len(isNull[0]) > 0
}

// NotNull defines NOT NULL condition
type NotNull [1]string

var _ Cond = NotNull{""}

// WriteTo write SQL to Writer
func (notNull NotNull) WriteTo(w Writer) error {
	_, err := fmt.Fprintf(w, "%s IS NOT NULL", notNull[0])
	return err
}

// And implements And with other conditions
func (notNull NotNull) And(conds ...Cond) Cond {
	return And(notNull, And(conds...))
}

// Or implements Or with other conditions
func (notNull NotNull) Or(conds ...Cond) Cond {
	return Or(notNull, Or(conds...))
}

// IsValid tests if this condition is valid
func (notNull NotNull) IsValid() bool {
	return len(notNull[0]) > 0
}
