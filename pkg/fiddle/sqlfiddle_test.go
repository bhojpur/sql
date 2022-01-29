package fiddle

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

func TestFiddle_CreateSchema4Mysql(t *testing.T) {
	fiddle := NewFiddle("")
	res, err := fiddle.CreateSchema(Mysql5_6, `create table person(id int not null auto_increment,
		name varchar(8),
		birthday datetime,
		constraint pk__person primary key(id));`)
	assert.NoError(t, err)
	fmt.Println(res)
	ret, err := fiddle.RunSQL(Mysql5_6, res.Code, "select * from person;")
	assert.NoError(t, err)
	fmt.Println(ret)
	ret, err = fiddle.RunSQL(Mysql5_6, res.Code, "select * from person1;")
	assert.Error(t, err)
	fmt.Println(err)
}
func TestFiddle_CreateSchema4Oracle(t *testing.T) {
	fiddle := NewFiddle("")
	res, err := fiddle.CreateSchema(Oracle11gR2, `create table table1(
       id number(9) not null primary key,
       a varchar2(40),
       b varchar2(40),
       c varchar2(40)
);`)
	assert.NoError(t, err)
	fmt.Println(res)
	ret, err := fiddle.RunSQL(Oracle11gR2, res.Code, "select * from table1;")
	assert.NoError(t, err)
	fmt.Println(ret)
}
func TestFiddle_CreateSchema4MssSQL(t *testing.T) {
	fiddle := NewFiddle("")
	res, err := fiddle.CreateSchema(Oracle11gR2, `create table table1(
       id int primary key,
       a varchar(40),
       b varchar(40),
       c varchar(40)
);`)
	assert.NoError(t, err)
	fmt.Println(res)
	ret, err := fiddle.RunSQL(Oracle11gR2, res.Code, "select * from table1;")
	assert.NoError(t, err)
	fmt.Println(ret)
}
