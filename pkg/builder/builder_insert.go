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
	"bytes"
	"fmt"
	"sort"
)

// Insert creates an insert Builder
func Insert(eq ...interface{}) *Builder {
	builder := &Builder{cond: NewCond()}
	return builder.Insert(eq...)
}
func (b *Builder) insertSelectWriteTo(w Writer) error {
	if _, err := fmt.Fprintf(w, "INSERT INTO %s ", b.into); err != nil {
		return err
	}
	if len(b.insertCols) > 0 {
		fmt.Fprintf(w, "(")
		for _, col := range b.insertCols {
			fmt.Fprintf(w, col)
		}
		fmt.Fprintf(w, ") ")
	}
	return b.selectWriteTo(w)
}
func (b *Builder) insertWriteTo(w Writer) error {
	if len(b.into) <= 0 {
		return ErrNoTableName
	}
	if len(b.insertCols) <= 0 && b.from == "" {
		return ErrNoColumnToInsert
	}
	if b.into != "" && b.from != "" {
		return b.insertSelectWriteTo(w)
	}
	if _, err := fmt.Fprintf(w, "INSERT INTO %s (", b.into); err != nil {
		return err
	}
	var args = make([]interface{}, 0)
	var bs []byte
	var valBuffer = bytes.NewBuffer(bs)
	for i, col := range b.insertCols {
		value := b.insertVals[i]
		fmt.Fprint(w, col)
		if e, ok := value.(expr); ok {
			fmt.Fprintf(valBuffer, "(%s)", e.sql)
			args = append(args, e.args...)
		} else if value == nil {
			fmt.Fprintf(valBuffer, `null`)
		} else {
			fmt.Fprint(valBuffer, "?")
			args = append(args, value)
		}
		if i != len(b.insertCols)-1 {
			if _, err := fmt.Fprint(w, ","); err != nil {
				return err
			}
			if _, err := fmt.Fprint(valBuffer, ","); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprint(w, ") Values ("); err != nil {
		return err
	}
	if _, err := w.Write(valBuffer.Bytes()); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, ")"); err != nil {
		return err
	}
	w.Append(args...)
	return nil
}

type insertColsSorter struct {
	cols []string
	vals []interface{}
}

func (s insertColsSorter) Len() int {
	return len(s.cols)
}
func (s insertColsSorter) Swap(i, j int) {
	s.cols[i], s.cols[j] = s.cols[j], s.cols[i]
	s.vals[i], s.vals[j] = s.vals[j], s.vals[i]
}
func (s insertColsSorter) Less(i, j int) bool {
	return s.cols[i] < s.cols[j]
}

// Insert sets insert SQL
func (b *Builder) Insert(eq ...interface{}) *Builder {
	if len(eq) > 0 {
		var paramType = -1
		for _, e := range eq {
			switch t := e.(type) {
			case Eq:
				if paramType == -1 {
					paramType = 0
				}
				if paramType != 0 {
					break
				}
				for k, v := range t {
					b.insertCols = append(b.insertCols, k)
					b.insertVals = append(b.insertVals, v)
				}
			case string:
				if paramType == -1 {
					paramType = 1
				}
				if paramType != 1 {
					break
				}
				b.insertCols = append(b.insertCols, t)
			}
		}
	}
	if len(b.insertCols) == len(b.insertVals) {
		sort.Sort(insertColsSorter{
			cols: b.insertCols,
			vals: b.insertVals,
		})
	}
	b.optype = insertType
	return b
}
