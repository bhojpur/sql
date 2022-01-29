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

func (b *Builder) limitWriteTo(w Writer) error {
	if strings.TrimSpace(b.dialect) == "" {
		return ErrDialectNotSetUp
	}
	if b.limitation != nil {
		limit := b.limitation
		if limit.offset < 0 || limit.limitN <= 0 {
			return ErrInvalidLimitation
		}
		// erase limit condition
		b.limitation = nil
		defer func() {
			b.limitation = limit
		}()
		ow := w.(*BytesWriter)
		switch strings.ToLower(strings.TrimSpace(b.dialect)) {
		case ORACLE:
			if len(b.selects) == 0 {
				b.selects = append(b.selects, "*")
			}
			var final *Builder
			selects := b.selects
			b.selects = append(selects, "ROWNUM RN")
			var wb *Builder
			if b.optype == setOpType {
				wb = Dialect(b.dialect).Select("at.*", "ROWNUM RN").
					From(b, "at")
			} else {
				wb = b
			}
			if limit.offset == 0 {
				final = Dialect(b.dialect).Select(selects...).From(wb, "at").
					Where(Lte{"at.RN": limit.limitN})
			} else {
				sub := Dialect(b.dialect).Select("*").
					From(b, "at").Where(Lte{"at.RN": limit.offset + limit.limitN})
				final = Dialect(b.dialect).Select(selects...).From(sub, "att").
					Where(Gt{"att.RN": limit.offset})
			}
			return final.WriteTo(ow)
		case SQLITE, MYSQL, POSTGRES:
			// if type UNION, we need to write previous content back to current writer
			if b.optype == setOpType {
				if err := b.WriteTo(ow); err != nil {
					return err
				}
			}
			if limit.offset == 0 {
				fmt.Fprint(ow, " LIMIT ", limit.limitN)
			} else {
				fmt.Fprintf(ow, " LIMIT %v OFFSET %v", limit.limitN, limit.offset)
			}
		case MSSQL:
			if len(b.selects) == 0 {
				b.selects = append(b.selects, "*")
			}
			var final *Builder
			selects := b.selects
			b.selects = append(append([]string{fmt.Sprintf("TOP %d %v", limit.limitN+limit.offset, b.selects[0])},
				b.selects[1:]...), "ROW_NUMBER() OVER (ORDER BY (SELECT 1)) AS RN")
			var wb *Builder
			if b.optype == setOpType {
				wb = Dialect(b.dialect).Select("*", "ROW_NUMBER() OVER (ORDER BY (SELECT 1)) AS RN").
					From(b, "at")
			} else {
				wb = b
			}
			if limit.offset == 0 {
				final = Dialect(b.dialect).Select(selects...).From(wb, "at")
			} else {
				final = Dialect(b.dialect).Select(selects...).From(wb, "at").Where(Gt{"at.RN": limit.offset})
			}
			return final.WriteTo(ow)
		default:
			return ErrNotSupportType
		}
	}
	return nil
}
