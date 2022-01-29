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

import (
	"fmt"
)

// UpdateCond defines an interface that cond could be used with update
type UpdateCond interface {
	IsValid() bool
	OpWriteTo(op string, w Writer) error
}

// Update creates an update Builder
func Update(updates ...Cond) *Builder {
	builder := &Builder{cond: NewCond()}
	return builder.Update(updates...)
}
func (b *Builder) updateWriteTo(w Writer) error {
	if len(b.from) <= 0 {
		return ErrNoTableName
	}
	if len(b.updates) <= 0 {
		return ErrNoColumnToUpdate
	}
	if _, err := fmt.Fprintf(w, "UPDATE %s SET ", b.from); err != nil {
		return err
	}
	for i, s := range b.updates {
		if err := s.OpWriteTo(",", w); err != nil {
			return err
		}
		if i != len(b.updates)-1 {
			if _, err := fmt.Fprint(w, ","); err != nil {
				return err
			}
		}
	}
	if !b.cond.IsValid() {
		return nil
	}
	if _, err := fmt.Fprint(w, " WHERE "); err != nil {
		return err
	}
	return b.cond.WriteTo(w)
}
