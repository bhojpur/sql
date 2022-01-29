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
	"reflect"
	"strings"
)

type condNotIn condIn

var _ Cond = condNotIn{}

// NotIn generate NOT IN condition
func NotIn(col string, values ...interface{}) Cond {
	return condNotIn{col, values}
}
func (condNotIn condNotIn) handleBlank(w Writer) error {
	_, err := fmt.Fprint(w, "0=0")
	return err
}
func (condNotIn condNotIn) WriteTo(w Writer) error {
	if len(condNotIn.vals) <= 0 {
		return condNotIn.handleBlank(w)
	}
	switch condNotIn.vals[0].(type) {
	case []int8:
		vals := condNotIn.vals[0].([]int8)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []int16:
		vals := condNotIn.vals[0].([]int16)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []int:
		vals := condNotIn.vals[0].([]int)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []int32:
		vals := condNotIn.vals[0].([]int32)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []int64:
		vals := condNotIn.vals[0].([]int64)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []uint8:
		vals := condNotIn.vals[0].([]uint8)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []uint16:
		vals := condNotIn.vals[0].([]uint16)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []uint:
		vals := condNotIn.vals[0].([]uint)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []uint32:
		vals := condNotIn.vals[0].([]uint32)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []uint64:
		vals := condNotIn.vals[0].([]uint64)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []string:
		vals := condNotIn.vals[0].([]string)
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, val := range vals {
			w.Append(val)
		}
	case []interface{}:
		vals := condNotIn.vals[0].([]interface{})
		if len(vals) <= 0 {
			return condNotIn.handleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(vals))
		if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		w.Append(vals...)
	case expr:
		val := condNotIn.vals[0].(expr)
		if _, err := fmt.Fprintf(w, "%s NOT IN (", condNotIn.col); err != nil {
			return err
		}
		if err := val.WriteTo(w); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, ")"); err != nil {
			return err
		}
	case *Builder:
		val := condNotIn.vals[0].(*Builder)
		if _, err := fmt.Fprintf(w, "%s NOT IN (", condNotIn.col); err != nil {
			return err
		}
		if err := val.WriteTo(w); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, ")"); err != nil {
			return err
		}
	default:
		v := reflect.ValueOf(condNotIn.vals[0])
		if v.Kind() == reflect.Slice {
			l := v.Len()
			if l == 0 {
				return condNotIn.handleBlank(w)
			}
			questionMark := strings.Repeat("?,", l)
			if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
				return err
			}
			for i := 0; i < l; i++ {
				w.Append(v.Index(i).Interface())
			}
		} else {
			questionMark := strings.Repeat("?,", len(condNotIn.vals))
			if _, err := fmt.Fprintf(w, "%s NOT IN (%s)", condNotIn.col, questionMark[:len(questionMark)-1]); err != nil {
				return err
			}
			w.Append(condNotIn.vals...)
		}
	}
	return nil
}
func (condNotIn condNotIn) And(conds ...Cond) Cond {
	return And(condNotIn, And(conds...))
}
func (condNotIn condNotIn) Or(conds ...Cond) Cond {
	return Or(condNotIn, Or(conds...))
}
func (condNotIn condNotIn) IsValid() bool {
	return len(condNotIn.col) > 0 && len(condNotIn.vals) > 0
}
