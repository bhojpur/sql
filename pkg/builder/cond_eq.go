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

// Incr implements a type used by Eq
type Incr int

// Decr implements a type used by Eq
type Decr int

// Eq defines equals conditions
type Eq map[string]interface{}

var _ Cond = Eq{}

// OpWriteTo writes conditions with special operator
func (eq Eq) OpWriteTo(op string, w Writer) error {
	var i = 0
	for _, k := range eq.sortedKeys() {
		v := eq[k]
		switch v.(type) {
		case []int, []int64, []string, []int32, []int16, []int8, []uint, []uint64, []uint32, []uint16, []interface{}:
			if err := In(k, v).WriteTo(w); err != nil {
				return err
			}
		case expr:
			if _, err := fmt.Fprintf(w, "%s=(", k); err != nil {
				return err
			}
			if err := v.(expr).WriteTo(w); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, ")"); err != nil {
				return err
			}
		case *Builder:
			if _, err := fmt.Fprintf(w, "%s=(", k); err != nil {
				return err
			}
			if err := v.(*Builder).WriteTo(w); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, ")"); err != nil {
				return err
			}
		case Incr:
			if _, err := fmt.Fprintf(w, "%s=%s+?", k, k); err != nil {
				return err
			}
			w.Append(int(v.(Incr)))
		case Decr:
			if _, err := fmt.Fprintf(w, "%s=%s-?", k, k); err != nil {
				return err
			}
			w.Append(int(v.(Decr)))
		case nil:
			if _, err := fmt.Fprintf(w, "%s=null", k); err != nil {
				return err
			}
		default:
			if _, err := fmt.Fprintf(w, "%s=?", k); err != nil {
				return err
			}
			w.Append(v)
		}
		if i != len(eq)-1 {
			if _, err := fmt.Fprint(w, op); err != nil {
				return err
			}
		}
		i = i + 1
	}
	return nil
}

// WriteTo writes SQL to Writer
func (eq Eq) WriteTo(w Writer) error {
	return eq.OpWriteTo(" AND ", w)
}

// And implements And with other conditions
func (eq Eq) And(conds ...Cond) Cond {
	return And(eq, And(conds...))
}

// Or implements Or with other conditions
func (eq Eq) Or(conds ...Cond) Cond {
	return Or(eq, Or(conds...))
}

// IsValid tests if this Eq is valid
func (eq Eq) IsValid() bool {
	return len(eq) > 0
}

// sortedKeys returns all keys of this Eq sorted with sort.Strings.
// It is used internally for consistent ordering when generating
func (eq Eq) sortedKeys() []string {
	keys := make([]string, 0, len(eq))
	for key := range eq {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
