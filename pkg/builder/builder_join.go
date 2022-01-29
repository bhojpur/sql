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

// InnerJoin sets inner join
func (b *Builder) InnerJoin(joinTable, joinCond interface{}) *Builder {
	return b.Join("INNER", joinTable, joinCond)
}

// LeftJoin sets left join SQL
func (b *Builder) LeftJoin(joinTable, joinCond interface{}) *Builder {
	return b.Join("LEFT", joinTable, joinCond)
}

// RightJoin sets right join SQL
func (b *Builder) RightJoin(joinTable, joinCond interface{}) *Builder {
	return b.Join("RIGHT", joinTable, joinCond)
}

// CrossJoin sets cross join SQL
func (b *Builder) CrossJoin(joinTable, joinCond interface{}) *Builder {
	return b.Join("CROSS", joinTable, joinCond)
}

// FullJoin sets full join SQL
func (b *Builder) FullJoin(joinTable, joinCond interface{}) *Builder {
	return b.Join("FULL", joinTable, joinCond)
}

// Join sets join table and conditions
func (b *Builder) Join(joinType string, joinTable, joinCond interface{}) *Builder {
	switch joinCond.(type) {
	case Cond:
		b.joins = append(b.joins, join{joinType, joinTable, joinCond.(Cond)})
	case string:
		b.joins = append(b.joins, join{joinType, joinTable, Expr(joinCond.(string))})
	}
	return b
}
