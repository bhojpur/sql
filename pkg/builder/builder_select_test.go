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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_Select(t *testing.T) {
	sql, args, err := Select("c, d").From("table1").ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT c, d FROM table1", sql)
	assert.EqualValues(t, []interface{}(nil), args)
	sql, args, err = Select("c, d").From("table1").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT c, d FROM table1 WHERE a=?", sql)
	assert.EqualValues(t, []interface{}{1}, args)
	_, _, err = Select("c, d").ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrNoTableName, err)
}
func TestBuilderSelectGroupBy(t *testing.T) {
	sql, args, err := Select("c").From("table1").GroupBy("c").Having("count(c)=1").ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT c FROM table1 GROUP BY c HAVING count(c)=1", sql)
	assert.EqualValues(t, 0, len(args))
	fmt.Println(sql, args)
}
func TestBuilderSelectOrderBy(t *testing.T) {
	sql, args, err := Select("c").From("table1").OrderBy("c DESC").ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT c FROM table1 ORDER BY c DESC", sql)
	assert.EqualValues(t, 0, len(args))
	fmt.Println(sql, args)
}
func TestBuilder_From(t *testing.T) {
	// simple one
	sql, args, err := Select("c").From("table1").ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT c FROM table1", sql)
	assert.EqualValues(t, 0, len(args))
	// from sub with alias
	sql, args, err = Select("sub.id").From(Select("id").From("table1").Where(Eq{"a": 1}),
		"sub").Where(Eq{"b": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT sub.id FROM (SELECT id FROM table1 WHERE a=?) sub WHERE b=?", sql)
	assert.EqualValues(t, []interface{}{1, 1}, args)
	// from sub without alias and with conditions
	sql, args, err = Select("sub.id").From(Select("id").From("table1").Where(Eq{"a": 1})).Where(Eq{"b": 1}).ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrUnnamedDerivedTable, err)
	// from sub without alias and conditions
	sql, args, err = Select("sub.id").From(Select("id").From("table1").Where(Eq{"a": 1})).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT sub.id FROM (SELECT id FROM table1 WHERE a=?)", sql)
	assert.EqualValues(t, []interface{}{1}, args)
	// from union with alias
	sql, args, err = Select("sub.id").From(
		Select("id").From("table1").Where(Eq{"a": 1}).Union(
			"all", Select("id").From("table1").Where(Eq{"a": 2})), "sub").Where(Eq{"b": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT sub.id FROM ((SELECT id FROM table1 WHERE a=?) UNION ALL (SELECT id FROM table1 WHERE a=?)) sub WHERE b=?", sql)
	assert.EqualValues(t, []interface{}{1, 2, 1}, args)
	// from union without alias
	_, _, err = Select("sub.id").From(
		Select("id").From("table1").Where(Eq{"a": 1}).Union(
			"all", Select("id").From("table1").Where(Eq{"a": 2}))).Where(Eq{"b": 1}).ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrUnnamedDerivedTable, err)
	// will raise error
	_, _, err = Select("c").From(Insert(Eq{"a": 1}).From("table1"), "table1").ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrUnexpectedSubQuery, err)
	// will raise error
	_, _, err = Select("c").From(Delete(Eq{"a": 1}).From("table1"), "table1").ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrUnexpectedSubQuery, err)
	// from a sub-query in different dialect
	_, _, err = MySQL().Select("sub.id").From(
		Oracle().Select("id").From("table1").Where(Eq{"a": 1}), "sub").Where(Eq{"b": 1}).ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrInconsistentDialect, err)
	// from a sub-query (dialect set up)
	sql, args, err = MySQL().Select("sub.id").From(
		MySQL().Select("id").From("table1").Where(Eq{"a": 1}), "sub").Where(Eq{"b": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT sub.id FROM (SELECT id FROM table1 WHERE a=?) sub WHERE b=?", sql)
	assert.EqualValues(t, []interface{}{1, 1}, args)
	// from a sub-query (dialect not set up)
	sql, args, err = MySQL().Select("sub.id").From(
		Select("id").From("table1").Where(Eq{"a": 1}), "sub").Where(Eq{"b": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT sub.id FROM (SELECT id FROM table1 WHERE a=?) sub WHERE b=?", sql)
	assert.EqualValues(t, []interface{}{1, 1}, args)
}
