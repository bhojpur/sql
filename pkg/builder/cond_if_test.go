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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCond_If(t *testing.T) {
	cond1 := If(1 > 0, Eq{"a": 1}, Eq{"b": 1})
	sql, err := ToBoundSQL(cond1)
	assert.NoError(t, err)
	assert.EqualValues(t, "a=1", sql)
	cond2 := If(1 < 0, Eq{"a": 1}, Eq{"b": 1})
	sql, err = ToBoundSQL(cond2)
	assert.NoError(t, err)
	assert.EqualValues(t, "b=1", sql)
	cond3 := If(1 > 0, cond2, Eq{"c": 1})
	sql, err = ToBoundSQL(cond3)
	assert.NoError(t, err)
	assert.EqualValues(t, "b=1", sql)
	cond4 := If(2 < 0, Eq{"d": "a"})
	sql, err = ToBoundSQL(cond4)
	assert.NoError(t, err)
	assert.EqualValues(t, "", sql)
	cond5 := And(cond1, cond2, cond3, cond4)
	sql, err = ToBoundSQL(cond5)
	assert.NoError(t, err)
	assert.EqualValues(t, "a=1 AND b=1 AND b=1", sql)
}
