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
	"sort"
)

// Neq defines not equal conditions
type Neq map[string]interface{}

var _ Cond = Neq{}

// WriteTo writes SQL to Writer
func (neq Neq) WriteTo(w Writer) error {
	var args = make([]interface{}, 0, len(neq))
	var i = 0
	for _, k := range neq.sortedKeys() {
		v := neq[k]
		switch v.(type) {
		case []int, []int64, []string, []int32, []int16, []int8:
			if err := NotIn(k, v).WriteTo(w); err != nil {
				return err
			}
		case expr:
			if _, err := fmt.Fprintf(w, "%s<>(", k); err != nil {
				return err
			}
			if err := v.(expr).WriteTo(w); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, ")"); err != nil {
				return err
			}
		case *Builder:
			if _, err := fmt.Fprintf(w, "%s<>(", k); err != nil {
				return err
			}
			if err := v.(*Builder).WriteTo(w); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, ")"); err != nil {
				return err
			}
		default:
			if _, err := fmt.Fprintf(w, "%s<>?", k); err != nil {
				return err
			}
			args = append(args, v)
		}
		if i != len(neq)-1 {
			if _, err := fmt.Fprint(w, " AND "); err != nil {
				return err
			}
		}
		i = i + 1
	}
	w.Append(args...)
	return nil
}

// And implements And with other conditions
func (neq Neq) And(conds ...Cond) Cond {
	return And(neq, And(conds...))
}

// Or implements Or with other conditions
func (neq Neq) Or(conds ...Cond) Cond {
	return Or(neq, Or(conds...))
}

// IsValid tests if this condition is valid
func (neq Neq) IsValid() bool {
	return len(neq) > 0
}

// sortedKeys returns all keys of this Neq sorted with sort.Strings.
// It is used internally for consistent ordering when generating
func (neq Neq) sortedKeys() []string {
	keys := make([]string, 0, len(neq))
	for key := range neq {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
