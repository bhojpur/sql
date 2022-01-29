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
	"strings"
)

func (b *Builder) setOpWriteTo(w Writer) error {
	if b.limitation != nil || b.cond.IsValid() ||
		b.orderBy != "" || b.having != "" || b.groupBy != "" {
		return ErrNotUnexpectedUnionConditions
	}
	for idx, o := range b.setOps {
		current := o.builder
		if current.optype != selectType {
			return ErrUnsupportedUnionMembers
		}
		if len(b.setOps) == 1 {
			if err := current.selectWriteTo(w); err != nil {
				return err
			}
		} else {
			if b.dialect != "" && b.dialect != current.dialect {
				return ErrInconsistentDialect
			}
			if idx != 0 {
				if o.distinctType == "" {
					fmt.Fprint(w, fmt.Sprintf(" %s ", strings.ToUpper(o.opType)))
				} else {
					fmt.Fprint(w, fmt.Sprintf(" %s %s ", strings.ToUpper(o.opType), strings.ToUpper(o.distinctType)))
				}
			}
			fmt.Fprint(w, "(")
			if err := current.selectWriteTo(w); err != nil {
				return err
			}
			fmt.Fprint(w, ")")
		}
	}
	return nil
}
