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

func TestBuilderUpdate(t *testing.T) {
	sql, args, err := Update(Eq{"a": 2}).From("table1").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "UPDATE table1 SET a=? WHERE a=?", sql)
	assert.EqualValues(t, []interface{}{2, 1}, args)
	sql, args, err = Update(Eq{"a": 2, "b": 1}).From("table1").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "UPDATE table1 SET a=?,b=? WHERE a=?", sql)
	assert.EqualValues(t, []interface{}{2, 1, 1}, args)
	sql, args, err = Update(Eq{"a": 2}, Eq{"b": 1}).From("table1").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "UPDATE table1 SET a=?,b=? WHERE a=?", sql)
	assert.EqualValues(t, []interface{}{2, 1, 1}, args)
	sql, args, err = Update(Eq{"a": 2, "b": Incr(1)}).From("table2").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "UPDATE table2 SET a=?,b=b+? WHERE a=?", sql)
	assert.EqualValues(t, []interface{}{2, 1, 1}, args)
	sql, args, err = Update(Eq{"a": 2, "b": Incr(1), "c": Decr(1), "d": Expr("select count(*) from table2")}).From("table2").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "UPDATE table2 SET a=?,b=b+?,c=c-?,d=(select count(*) from table2) WHERE a=?", sql)
	assert.EqualValues(t, []interface{}{2, 1, 1, 1}, args)
	sql, args, err = Update(Eq{"a": 2}).Where(Eq{"a": 1}).ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrNoTableName, err)
	sql, args, err = Update(Eq{}).From("table1").Where(Eq{"a": 1}).ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrNoColumnToUpdate, err)
	var builder = Builder{cond: NewCond()}
	sql, args, err = builder.Update(Eq{"a": 2, "b": 1}).From("table1").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "UPDATE table1 SET a=?,b=? WHERE a=?", sql)
	assert.EqualValues(t, []interface{}{2, 1, 1}, args)
	sql, args, err = Update(Eq{"a": 1}, Expr("c = c+1")).From("table1").Where(Eq{"b": 2}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "UPDATE table1 SET a=?,c = c+1 WHERE b=?", sql)
	assert.EqualValues(t, []interface{}{1, 2}, args)
	sql, args, err = Update(Eq{"a": 2}).From("table1").ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "UPDATE table1 SET a=?", sql)
	assert.EqualValues(t, []interface{}{2}, args)
}
