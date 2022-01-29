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

type condIf struct {
	condition bool
	condTrue  Cond
	condFalse Cond
}

var _ Cond = condIf{}

// If returns Cond via condition
func If(condition bool, condTrue Cond, condFalse ...Cond) Cond {
	var c = condIf{
		condition: condition,
		condTrue:  condTrue,
	}
	if len(condFalse) > 0 {
		c.condFalse = condFalse[0]
	}
	return c
}
func (condIf condIf) WriteTo(w Writer) error {
	if condIf.condition {
		return condIf.condTrue.WriteTo(w)
	} else if condIf.condFalse != nil {
		return condIf.condFalse.WriteTo(w)
	}
	return nil
}
func (condIf condIf) And(conds ...Cond) Cond {
	return And(condIf, And(conds...))
}
func (condIf condIf) Or(conds ...Cond) Cond {
	return Or(condIf, Or(conds...))
}
func (condIf condIf) IsValid() bool {
	if condIf.condition {
		return condIf.condTrue != nil
	}
	return condIf.condFalse != nil
}
