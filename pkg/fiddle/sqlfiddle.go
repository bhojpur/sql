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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	Mysql5_6      = 9
	Oracle11gR2   = 4
	PostgreSQL96  = 17
	PostgreSQL93  = 15
	SQLite_WebSQL = 7
	SQLite_SQLjs  = 5
	MSSQL2017     = 18
)

/*
{
  "_id" : "9_ed3b0c",
  "short_code" : "ed3b0c",
  "schema_structure" : [ {
    "columns" : [ {
      "name" : "id",
      "type" : "INT(10)"
    }, {
      "name" : "name",
      "type" : "VARCHAR(8)"
    }, {
      "name" : "birthday",
      "type" : "DATETIME(19)"
    } ],
    "table_name" : "person",
    "table_type" : "TABLE"
  } ]
}
*/
const (
	defaultBaseURL = "http://sqlfiddle.com"
)

type Fiddle struct {
	baseURL string
}

func NewFiddle(baseURL string) *Fiddle {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Fiddle{baseURL}
}

type column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type schemaStructure struct {
	Columns   []column `json:"columns"`
	TableName string   `json:"table_name"`
	TableType string   `json:"table_type"`
}
type response struct {
	Id         string            `json:"_id"`
	Code       string            `json:"short_code"`
	Structures []schemaStructure `json:"schema_structure"`
	ErrorMsg   string            `json:"error"`
}

func (f *Fiddle) CreateSchema(dbType int, sql string) (*response, error) {
	payload := map[string]interface{}{"statement_separator": ";", "db_type_id": dbType, "ddl": sql}
	bs, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", f.baseURL+"/backend/createSchema?_action=create", bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	if res.ErrorMsg != "" {
		return &res, errors.New(res.ErrorMsg)
	}
	return &res, nil
}

/*
http://sqlfiddle.com/backend/executeQuery?_action=query
{"db_type_id":9,"schema_short_code":"ed3b0c","statement_separator":";","sql":"select * from person;"}
{
  "ID" : 2,
  "sets" : [ {
    "RESULTS" : {
      "COLUMNS" : [ "id", "name", "birthday" ],
      "DATA" : [ ]
    },
    "SUCCEEDED" : true,
    "STATEMENT" : "select * from person",
    "EXECUTIONTIME" : 2,
    "EXECUTIONPLANRAW" : {
      "COLUMNS" : [ "id", "select_type", "table", "type", "possible_keys", "key", "key_len", "ref", "rows", "filtered", "Extra" ],
      "DATA" : [ [ "1", "SIMPLE", "person", "ALL", null, null, null, null, "1", "100.00", null ] ]
    },
    "EXECUTIONPLAN" : {
      "COLUMNS" : [ "id", "select_type", "table", "type", "possible_keys", "key", "key_len", "ref", "rows", "filtered", "Extra" ],
      "DATA" : [ [ "1", "SIMPLE", "person", "ALL", null, null, null, null, "1", "100.00", null ] ]
    }
  } ]
}
*/
type sqlResult struct {
	Statement     string `json:"STATEMENT"`
	Succeeded     bool   `json:"SUCCEEDED"`
	ErrorMessage  string `json:"ERRORMESSAGE"`
	ExecutionTime int    `json:"EXECUTIONTIME"`
	Results       struct {
		Columns []string        `json:"COLUMNS"`
		Data    [][]interface{} `json:"DATA"`
	} `json:"RESULTS"`
	ExecutionPlanRaw struct {
		Columns []string        `json:"COLUMNS"`
		Data    [][]interface{} `json:"DATA"`
	} `json:"EXECUTIONPLANRAW"`
	ExecutionPlan struct {
		Columns []string        `json:"COLUMNS"`
		Data    [][]interface{} `json:"DATA"`
	} `json:"EXECUTIONPLAN"`
}
type sqlResponse struct {
	Id   int64 `json:"ID"`
	Sets []sqlResult
}

func (f *Fiddle) RunSQL(dbType int, code, sql string) (*sqlResponse, error) {
	payload := map[string]interface{}{"schema_short_code": code, "statement_separator": ";", "db_type_id": dbType, "sql": sql}
	bs, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", f.baseURL+"/backend/executeQuery?_action=query", bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res sqlResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	for e := range res.Sets {
		result := res.Sets[e]
		if result.ErrorMessage != "" {
			return &res, fmt.Errorf("something wrong with sql %q in pos %v, detail: %q",
				result.Statement, e+1, result.ErrorMessage)
		}
	}
	return &res, nil
}
