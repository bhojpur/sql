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
	sql2 "database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	sqlfiddle "github.com/bhojpur/sql/pkg/fiddle"
	"github.com/stretchr/testify/assert"
)

const placeholderConverterSQL = "SELECT a, b FROM table_a WHERE b_id=(SELECT id FROM table_b WHERE b=?) AND id=? AND c=? AND d=? AND e=? AND f=?"
const placeholderConvertedSQL = "SELECT a, b FROM table_a WHERE b_id=(SELECT id FROM table_b WHERE b=$1) AND id=$2 AND c=$3 AND d=$4 AND e=$5 AND f=$6"
const placeholderBoundSQL = "SELECT a, b FROM table_a WHERE b_id=(SELECT id FROM table_b WHERE b=1) AND id=2.1 AND c='3' AND d=4 AND e='5' AND f=true"

func TestNoSQLQuoteNeeded(t *testing.T) {
	assert.False(t, noSQLQuoteNeeded(nil))
}
func TestPlaceholderConverter(t *testing.T) {
	var convertCases = []struct {
		before, after string
		mark          string
	}{
		{
			before: placeholderConverterSQL,
			after:  placeholderConvertedSQL,
			mark:   "$",
		},
		{
			before: "SELECT a, b, 'a?b' FROM table_a WHERE id=?",
			after:  "SELECT a, b, 'a?b' FROM table_a WHERE id=:1",
			mark:   ":",
		},
		{
			before: "SELECT a, b, 'a\\'?b' FROM table_a WHERE id=?",
			after:  "SELECT a, b, 'a\\'?b' FROM table_a WHERE id=$1",
			mark:   "$",
		},
		{
			before: "SELECT a, b, 'a\\'b' FROM table_a WHERE id=?",
			after:  "SELECT a, b, 'a\\'b' FROM table_a WHERE id=$1",
			mark:   "$",
		},
	}
	for _, kase := range convertCases {
		t.Run(kase.before, func(t *testing.T) {
			newSQL, err := ConvertPlaceholder(kase.before, kase.mark)
			assert.NoError(t, err)
			assert.EqualValues(t, kase.after, newSQL)
		})
	}
}
func BenchmarkPlaceholderConverter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConvertPlaceholder(placeholderConverterSQL, "$")
	}
}
func TestBoundSQLConverter(t *testing.T) {
	newSQL, err := ConvertToBoundSQL(placeholderConverterSQL, []interface{}{1, 2.1, "3", uint(4), "5", true})
	assert.NoError(t, err)
	assert.EqualValues(t, placeholderBoundSQL, newSQL)
	newSQL, err = ConvertToBoundSQL(placeholderConverterSQL, []interface{}{1, 2.1, sql2.Named("any", "3"), uint(4), "5", true})
	assert.NoError(t, err)
	assert.EqualValues(t, placeholderBoundSQL, newSQL)
	newSQL, err = ConvertToBoundSQL(placeholderConverterSQL, []interface{}{1, 2.1, "3", 4, "5"})
	assert.Error(t, err)
	assert.EqualValues(t, ErrNeedMoreArguments, err)
	newSQL, err = ToBoundSQL(Select("id").From("table").Where(In("a", 1, 2)))
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT id FROM table WHERE a IN (1,2)", newSQL)
	newSQL, err = ToBoundSQL(Eq{"a": 1})
	assert.NoError(t, err)
	assert.EqualValues(t, "a=1", newSQL)
	newSQL, err = ToBoundSQL(1)
	assert.Error(t, err)
	assert.EqualValues(t, ErrNotSupportType, err)
}
func TestSQL(t *testing.T) {
	newSQL, args, err := ToSQL(In("a", 1, 2))
	assert.NoError(t, err)
	assert.EqualValues(t, "a IN (?,?)", newSQL)
	assert.EqualValues(t, []interface{}{1, 2}, args)
	newSQL, args, err = ToSQL(Select("id").From("table").Where(In("a", 1, 2)))
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT id FROM table WHERE a IN (?,?)", newSQL)
	assert.EqualValues(t, []interface{}{1, 2}, args)
	newSQL, args, err = ToSQL(1)
	assert.Error(t, err)
	assert.EqualValues(t, ErrNotSupportType, err)
}

type fiddler struct {
	sessionCode string
	dbType      int
	f           *sqlfiddle.Fiddle
}

func readPreparationSQLFromFile(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func newFiddler(fiddleServerAddr, dbDialect, preparationSQL string) (*fiddler, error) {
	var dbType int
	switch dbDialect {
	case MYSQL:
		dbType = sqlfiddle.Mysql5_6
	case MSSQL:
		dbType = sqlfiddle.MSSQL2017
	case POSTGRES:
		dbType = sqlfiddle.PostgreSQL96
	case ORACLE:
		dbType = sqlfiddle.Oracle11gR2
	case SQLITE:
		dbType = sqlfiddle.SQLite_WebSQL
	default:
		return nil, ErrNotSupportDialectType
	}
	f := sqlfiddle.NewFiddle(fiddleServerAddr)
	response, err := f.CreateSchema(dbType, preparationSQL)
	if err != nil {
		return nil, err
	}
	return &fiddler{sessionCode: response.Code, f: f, dbType: dbType}, nil
}
func (f *fiddler) executableCheck(obj interface{}) error {
	var sql string
	var err error
	switch obj.(type) {
	case *Builder:
		sql, err = obj.(*Builder).ToBoundSQL()
		if err != nil {
			return err
		}
	case string:
		sql = obj.(string)
	}
	_, err = f.f.RunSQL(f.dbType, f.sessionCode, sql)
	if err != nil {
		return err
	}
	return nil
}
func TestReadPreparationSQLFromFile(t *testing.T) {
	sqlFromFile, err := readPreparationSQLFromFile("testdata/mysql_fiddle_data.sql")
	assert.NoError(t, err)
	fmt.Println(sqlFromFile)
}
func TestNewFiddler(t *testing.T) {
	sqlFromFile, err := readPreparationSQLFromFile("testdata/mysql_fiddle_data.sql")
	assert.NoError(t, err)
	f, err := newFiddler("", MYSQL, sqlFromFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, f.sessionCode)
}
func TestExecutableCheck(t *testing.T) {
	sqlFromFile, err := readPreparationSQLFromFile("testdata/mysql_fiddle_data.sql")
	assert.NoError(t, err)
	f, err := newFiddler("", MYSQL, sqlFromFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, f.sessionCode)
	assert.NoError(t, f.executableCheck("SELECT * FROM table1"))
	err = f.executableCheck("SELECT * FROM table3")
	assert.Error(t, err)
}
func TestToSQLInDifferentDialects(t *testing.T) {
	sql, args, err := Postgres().Select().From("table1").Where(Eq{"a": "1"}.And(Neq{"b": "100"})).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM table1 WHERE a=$1 AND b<>$2", sql)
	assert.EqualValues(t, []interface{}{"1", "100"}, args)
	sql, args, err = MySQL().Select().From("table1").Where(Eq{"a": "1"}.And(Neq{"b": "100"})).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM table1 WHERE a=? AND b<>?", sql)
	assert.EqualValues(t, []interface{}{"1", "100"}, args)
	sql, args, err = MsSQL().Select().From("table1").Where(Eq{"a": "1"}.And(Neq{"b": "100"})).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM table1 WHERE a=@p1 AND b<>@p2", sql)
	assert.EqualValues(t, []interface{}{sql2.Named("p1", "1"), sql2.Named("p2", "100")}, args)
	// test sql.NamedArg in cond
	sql, args, err = MsSQL().Select().From("table1").Where(Eq{"a": sql2.NamedArg{Name: "param", Value: "1"}}.And(Neq{"b": "100"})).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM table1 WHERE a=@p1 AND b<>@p2", sql)
	assert.EqualValues(t, []interface{}{sql2.Named("p1", "1"), sql2.Named("p2", "100")}, args)
	sql, args, err = Oracle().Select().From("table1").Where(Eq{"a": "1"}.And(Neq{"b": "100"})).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM table1 WHERE a=:p1 AND b<>:p2", sql)
	assert.EqualValues(t, []interface{}{sql2.Named("p1", "1"), sql2.Named("p2", "100")}, args)
	// test sql.NamedArg in cond
	sql, args, err = Oracle().Select().From("table1").Where(Eq{"a": sql2.Named("a", "1")}.And(Neq{"b": "100"})).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM table1 WHERE a=:p1 AND b<>:p2", sql)
	assert.EqualValues(t, []interface{}{sql2.Named("p1", "1"), sql2.Named("p2", "100")}, args)
	sql, args, err = SQLite().Select().From("table1").Where(Eq{"a": "1"}.And(Neq{"b": "100"})).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM table1 WHERE a=? AND b<>?", sql)
	assert.EqualValues(t, []interface{}{"1", "100"}, args)
}
func TestToSQLInjectionHarmlessDisposal(t *testing.T) {
	sql, err := MySQL().Select("*").From("table1").Where(Cond(Eq{"name": "cat';truncate table table1;"})).ToBoundSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT * FROM table1 WHERE name='cat'';truncate table table1;'", sql)
	sql, err = MySQL().Update(Eq{`a`: 1, `b`: nil}).From(`table1`).ToBoundSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "UPDATE table1 SET a=1,b=null", sql)
	sql, args, err := MySQL().Update(Eq{`a`: 1, `b`: nil}).From(`table1`).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "UPDATE table1 SET a=?,b=null", sql)
	assert.EqualValues(t, []interface{}{1}, args)
}
