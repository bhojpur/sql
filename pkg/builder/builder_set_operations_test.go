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

func TestBuilder_Union(t *testing.T) {
	sql, args, err := Select("*").From("t1").Where(Eq{"status": "1"}).
		Union("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Union("distinct", Select("*").From("t2").Where(Eq{"status": "3"})).
		Union("", Select("*").From("t2").Where(Eq{"status": "3"})).
		ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "(SELECT * FROM t1 WHERE status=?) UNION ALL (SELECT * FROM t2 WHERE status=?) UNION DISTINCT (SELECT * FROM t2 WHERE status=?) UNION (SELECT * FROM t2 WHERE status=?)", sql)
	assert.EqualValues(t, []interface{}{"1", "2", "3", "3"}, args)
	// sub-query will inherit dialect from the main one
	sql, args, err = MySQL().Select("*").From("t1").Where(Eq{"status": "1"}).
		Union("all", Select("*").From("t2").Where(Eq{"status": "2"}).Limit(10)).
		Union("", Select("*").From("t2").Where(Eq{"status": "3"})).
		ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "(SELECT * FROM t1 WHERE status=?) UNION ALL (SELECT * FROM t2 WHERE status=? LIMIT 10) UNION (SELECT * FROM t2 WHERE status=?)", sql)
	assert.EqualValues(t, []interface{}{"1", "2", "3"}, args)
	// will raise error
	_, _, err = MySQL().Select("*").From("t1").Where(Eq{"status": "1"}).
		Union("all", Oracle().Select("*").From("t2").Where(Eq{"status": "2"}).Limit(10)).
		ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrInconsistentDialect, err)
	// will raise error
	_, _, err = Select("*").From("table1").Where(Eq{"a": "1"}).
		Union("all", Select("*").From("table2").Where(Eq{"a": "2"})).
		Where(Eq{"a": 2}).Limit(5, 10).
		ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrNotUnexpectedUnionConditions, err)
	// will raise error
	_, _, err = Delete(Eq{"a": 1}).From("t1").
		Union("all", Select("*").From("t2").Where(Eq{"status": "2"})).ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrUnsupportedUnionMembers, err)
	// will be overwrote by SELECT op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Union("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Select("*").From("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
	// will be overwrote by DELETE op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Union("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Delete(Eq{"status": "1"}).From("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
	// will be overwrote by INSERT op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Union("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Insert(Eq{"status": "1"}).Into("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
}
func TestBuilder_Intersect(t *testing.T) {
	sql, args, err := Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Intersect("distinct", Select("*").From("t2").Where(Eq{"status": "3"})).
		Intersect("", Select("*").From("t2").Where(Eq{"status": "3"})).
		ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "(SELECT * FROM t1 WHERE status=?) INTERSECT ALL (SELECT * FROM t2 WHERE status=?) INTERSECT DISTINCT (SELECT * FROM t2 WHERE status=?) INTERSECT (SELECT * FROM t2 WHERE status=?)", sql)
	assert.EqualValues(t, []interface{}{"1", "2", "3", "3"}, args)
	// sub-query will inherit dialect from the main one
	sql, args, err = MySQL().Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"}).Limit(10)).
		Intersect("", Select("*").From("t2").Where(Eq{"status": "3"})).
		ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "(SELECT * FROM t1 WHERE status=?) INTERSECT ALL (SELECT * FROM t2 WHERE status=? LIMIT 10) INTERSECT (SELECT * FROM t2 WHERE status=?)", sql)
	assert.EqualValues(t, []interface{}{"1", "2", "3"}, args)
	// will raise error
	_, _, err = MySQL().Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Oracle().Select("*").From("t2").Where(Eq{"status": "2"}).Limit(10)).
		ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrInconsistentDialect, err)
	// will raise error
	_, _, err = Select("*").From("table1").Where(Eq{"a": "1"}).
		Intersect("all", Select("*").From("table2").Where(Eq{"a": "2"})).
		Where(Eq{"a": 2}).Limit(5, 10).
		ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrNotUnexpectedUnionConditions, err)
	// will raise error
	_, _, err = Delete(Eq{"a": 1}).From("t1").
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"})).ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrUnsupportedUnionMembers, err)
	// will be overwrote by SELECT op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Select("*").From("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
	// will be overwrote by DELETE op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Delete(Eq{"status": "1"}).From("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
	// will be overwrote by INSERT op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Insert(Eq{"status": "1"}).Into("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
}
func TestBuilder_Except(t *testing.T) {
	sql, args, err := Select("*").From("t1").Where(Eq{"status": "1"}).
		Except("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Except("distinct", Select("*").From("t2").Where(Eq{"status": "3"})).
		Except("", Select("*").From("t2").Where(Eq{"status": "3"})).
		ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "(SELECT * FROM t1 WHERE status=?) EXCEPT ALL (SELECT * FROM t2 WHERE status=?) EXCEPT DISTINCT (SELECT * FROM t2 WHERE status=?) EXCEPT (SELECT * FROM t2 WHERE status=?)", sql)
	assert.EqualValues(t, []interface{}{"1", "2", "3", "3"}, args)
	// sub-query will inherit dialect from the main one
	sql, args, err = MySQL().Select("*").From("t1").Where(Eq{"status": "1"}).
		Except("all", Select("*").From("t2").Where(Eq{"status": "2"}).Limit(10)).
		Except("", Select("*").From("t2").Where(Eq{"status": "3"})).
		ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "(SELECT * FROM t1 WHERE status=?) EXCEPT ALL (SELECT * FROM t2 WHERE status=? LIMIT 10) EXCEPT (SELECT * FROM t2 WHERE status=?)", sql)
	assert.EqualValues(t, []interface{}{"1", "2", "3"}, args)
	// will raise error
	_, _, err = MySQL().Select("*").From("t1").Where(Eq{"status": "1"}).
		Except("all", Oracle().Select("*").From("t2").Where(Eq{"status": "2"}).Limit(10)).
		ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrInconsistentDialect, err)
	// will raise error
	_, _, err = Select("*").From("table1").Where(Eq{"a": "1"}).
		Except("all", Select("*").From("table2").Where(Eq{"a": "2"})).
		Where(Eq{"a": 2}).Limit(5, 10).
		ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrNotUnexpectedUnionConditions, err)
	// will raise error
	_, _, err = Delete(Eq{"a": 1}).From("t1").
		Except("all", Select("*").From("t2").Where(Eq{"status": "2"})).ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrUnsupportedUnionMembers, err)
	// will be overwrote by SELECT op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Except("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Select("*").From("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
	// will be overwrote by DELETE op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Except("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Delete(Eq{"status": "1"}).From("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
	// will be overwrote by INSERT op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Except("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Insert(Eq{"status": "1"}).Into("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
}
func TestBuilder_SetOperations(t *testing.T) {
	sql, args, err := Select("*").From("t1").Where(Eq{"status": "1"}).
		Union("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Intersect("distinct", Select("*").From("t2").Where(Eq{"status": "3"})).
		Except("", Select("*").From("t2").Where(Eq{"status": "3"})).
		ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "(SELECT * FROM t1 WHERE status=?) UNION ALL (SELECT * FROM t2 WHERE status=?) INTERSECT DISTINCT (SELECT * FROM t2 WHERE status=?) EXCEPT (SELECT * FROM t2 WHERE status=?)", sql)
	assert.EqualValues(t, []interface{}{"1", "2", "3", "3"}, args)
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Union("distinct", Select("*").From("t2").Where(Eq{"status": "3"})).
		Except("", Select("*").From("t2").Where(Eq{"status": "3"})).
		ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "(SELECT * FROM t1 WHERE status=?) INTERSECT ALL (SELECT * FROM t2 WHERE status=?) UNION DISTINCT (SELECT * FROM t2 WHERE status=?) EXCEPT (SELECT * FROM t2 WHERE status=?)", sql)
	assert.EqualValues(t, []interface{}{"1", "2", "3", "3"}, args)
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Except("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Intersect("distinct", Select("*").From("t2").Where(Eq{"status": "3"})).
		Union("", Select("*").From("t2").Where(Eq{"status": "3"})).
		ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "(SELECT * FROM t1 WHERE status=?) EXCEPT ALL (SELECT * FROM t2 WHERE status=?) INTERSECT DISTINCT (SELECT * FROM t2 WHERE status=?) UNION (SELECT * FROM t2 WHERE status=?)", sql)
	assert.EqualValues(t, []interface{}{"1", "2", "3", "3"}, args)
	// sub-query will inherit dialect from the main one
	sql, args, err = MySQL().Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"}).Limit(10)).
		Intersect("", Select("*").From("t2").Where(Eq{"status": "3"})).
		ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "(SELECT * FROM t1 WHERE status=?) INTERSECT ALL (SELECT * FROM t2 WHERE status=? LIMIT 10) INTERSECT (SELECT * FROM t2 WHERE status=?)", sql)
	assert.EqualValues(t, []interface{}{"1", "2", "3"}, args)
	// will raise error
	_, _, err = MySQL().Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Oracle().Select("*").From("t2").Where(Eq{"status": "2"}).Limit(10)).
		ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrInconsistentDialect, err)
	// will raise error
	_, _, err = Select("*").From("table1").Where(Eq{"a": "1"}).
		Intersect("all", Select("*").From("table2").Where(Eq{"a": "2"})).
		Where(Eq{"a": 2}).Limit(5, 10).
		ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrNotUnexpectedUnionConditions, err)
	// will raise error
	_, _, err = Delete(Eq{"a": 1}).From("t1").
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"})).ToSQL()
	assert.Error(t, err)
	assert.EqualValues(t, ErrUnsupportedUnionMembers, err)
	// will be overwrote by SELECT op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Select("*").From("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
	// will be overwrote by DELETE op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Delete(Eq{"status": "1"}).From("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
	// will be overwrote by INSERT op
	sql, args, err = Select("*").From("t1").Where(Eq{"status": "1"}).
		Intersect("all", Select("*").From("t2").Where(Eq{"status": "2"})).
		Insert(Eq{"status": "1"}).Into("t2").ToSQL()
	assert.NoError(t, err)
	fmt.Println(sql, args)
}
