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

import "errors"

var (
	// ErrNotSupportType not supported SQL type error
	ErrNotSupportType = errors.New("Not supported SQL type")
	// ErrNoNotInConditions no NOT IN params error
	ErrNoNotInConditions = errors.New("No NOT IN conditions")
	// ErrNoInConditions no IN params error
	ErrNoInConditions = errors.New("No IN conditions")
	// ErrNeedMoreArguments need more arguments
	ErrNeedMoreArguments = errors.New("Need more sql arguments")
	// ErrNoTableName no table name
	ErrNoTableName = errors.New("No table indicated")
	// ErrNoColumnToUpdate no column to update
	ErrNoColumnToUpdate = errors.New("No column(s) to update")
	// ErrNoColumnToInsert no column to insert
	ErrNoColumnToInsert = errors.New("No column(s) to insert")
	// ErrNotSupportDialectType not supported dialect type error
	ErrNotSupportDialectType = errors.New("Not supported dialect type")
	// ErrNotUnexpectedUnionConditions using union in a wrong way
	ErrNotUnexpectedUnionConditions = errors.New("Unexpected conditional fields in UNION query")
	// ErrUnsupportedUnionMembers unexpected members in UNION query
	ErrUnsupportedUnionMembers = errors.New("Unexpected members in UNION query")
	// ErrUnexpectedSubQuery Unexpected sub-query in SELECT query
	ErrUnexpectedSubQuery = errors.New("Unexpected sub-query in SELECT query")
	// ErrDialectNotSetUp dialect is not setup yet
	ErrDialectNotSetUp = errors.New("Dialect is not setup yet, try to use `Dialect(dbType)` at first")
	// ErrInvalidLimitation offset or limit is not correct
	ErrInvalidLimitation = errors.New("Offset or limit is not correct")
	// ErrUnnamedDerivedTable Every derived table must have its own alias
	ErrUnnamedDerivedTable = errors.New("Every derived table must have its own alias")
	// ErrInconsistentDialect Inconsistent dialect in same builder
	ErrInconsistentDialect = errors.New("Inconsistent dialect in same builder")
)
