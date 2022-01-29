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

// WriteMap writes conditions' SQL to Writer, op could be =, <>, >, <, <=, >= and etc.
func WriteMap(w Writer, data map[string]interface{}, op string) error {
	var args = make([]interface{}, 0, len(data))
	var i = 0
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	for _, k := range keys {
		v := data[k]
		switch v.(type) {
		case expr:
			if _, err := fmt.Fprintf(w, "%s%s(", k, op); err != nil {
				return err
			}
			if err := v.(expr).WriteTo(w); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, ")"); err != nil {
				return err
			}
		case *Builder:
			if _, err := fmt.Fprintf(w, "%s%s(", k, op); err != nil {
				return err
			}
			if err := v.(*Builder).WriteTo(w); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, ")"); err != nil {
				return err
			}
		default:
			if _, err := fmt.Fprintf(w, "%s%s?", k, op); err != nil {
				return err
			}
			args = append(args, v)
		}
		if i != len(data)-1 {
			if _, err := fmt.Fprint(w, " AND "); err != nil {
				return err
			}
		}
		i = i + 1
	}
	w.Append(args...)
	return nil
}

// Lt defines < condition
type Lt map[string]interface{}

var _ Cond = Lt{}

// WriteTo write SQL to Writer
func (lt Lt) WriteTo(w Writer) error {
	return WriteMap(w, lt, "<")
}

// And implements And with other conditions
func (lt Lt) And(conds ...Cond) Cond {
	return condAnd{lt, And(conds...)}
}

// Or implements Or with other conditions
func (lt Lt) Or(conds ...Cond) Cond {
	return condOr{lt, Or(conds...)}
}

// IsValid tests if this Eq is valid
func (lt Lt) IsValid() bool {
	return len(lt) > 0
}

// Lte defines <= condition
type Lte map[string]interface{}

var _ Cond = Lte{}

// WriteTo write SQL to Writer
func (lte Lte) WriteTo(w Writer) error {
	return WriteMap(w, lte, "<=")
}

// And implements And with other conditions
func (lte Lte) And(conds ...Cond) Cond {
	return And(lte, And(conds...))
}

// Or implements Or with other conditions
func (lte Lte) Or(conds ...Cond) Cond {
	return Or(lte, Or(conds...))
}

// IsValid tests if this Eq is valid
func (lte Lte) IsValid() bool {
	return len(lte) > 0
}

// Gt defines > condition
type Gt map[string]interface{}

var _ Cond = Gt{}

// WriteTo write SQL to Writer
func (gt Gt) WriteTo(w Writer) error {
	return WriteMap(w, gt, ">")
}

// And implements And with other conditions
func (gt Gt) And(conds ...Cond) Cond {
	return And(gt, And(conds...))
}

// Or implements Or with other conditions
func (gt Gt) Or(conds ...Cond) Cond {
	return Or(gt, Or(conds...))
}

// IsValid tests if this Eq is valid
func (gt Gt) IsValid() bool {
	return len(gt) > 0
}

// Gte defines >= condition
type Gte map[string]interface{}

var _ Cond = Gte{}

// WriteTo write SQL to Writer
func (gte Gte) WriteTo(w Writer) error {
	return WriteMap(w, gte, ">=")
}

// And implements And with other conditions
func (gte Gte) And(conds ...Cond) Cond {
	return And(gte, And(conds...))
}

// Or implements Or with other conditions
func (gte Gte) Or(conds ...Cond) Cond {
	return Or(gte, Or(conds...))
}

// IsValid tests if this Eq is valid
func (gte Gte) IsValid() bool {
	return len(gte) > 0
}
