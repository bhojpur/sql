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

func TestBuilderInsert(t *testing.T) {
	sql, err := Insert(Eq{"c": 1, "d": 2}).Into("table1").ToBoundSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "INSERT INTO table1 (c,d) Values (1,2)", sql)
	sql, err = Insert(Eq{"e": 3}, Eq{"c": 1}, Eq{"d": 2}).Into("table1").ToBoundSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "INSERT INTO table1 (c,d,e) Values (1,2,3)", sql)
	sql, err = Insert(Eq{"c": 1, "d": Expr("SELECT b FROM t WHERE d=? LIMIT 1", 2)}).Into("table1").ToBoundSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "INSERT INTO table1 (c,d) Values (1,(SELECT b FROM t WHERE d=2 LIMIT 1))", sql)
	sql, err = Insert(Eq{"c": 1, "d": 2}).ToBoundSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrNoTableName, err)
	assert.EqualValues(t, "", sql)
	sql, err = Insert(Eq{}).Into("table1").ToBoundSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrNoColumnToInsert, err)
	assert.EqualValues(t, "", sql)
	sql, err = Insert(Eq{`a`: nil}).Into(`table1`).ToBoundSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, `INSERT INTO table1 (a) Values (null)`, sql)
	sql, args, err := Insert(Eq{`a`: nil, `b`: `str`}).Into(`table1`).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, `INSERT INTO table1 (a,b) Values (null,?)`, sql)
	assert.EqualValues(t, []interface{}{`str`}, args)
}
func TestBuidlerInsert_Select(t *testing.T) {
	sql, err := Insert().Into("table1").Select().From("table2").ToBoundSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "INSERT INTO table1 SELECT * FROM table2", sql)
	sql, err = Insert("a, b").Into("table1").Select("b, c").From("table2").ToBoundSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "INSERT INTO table1 (a, b) SELECT b, c FROM table2", sql)
}
