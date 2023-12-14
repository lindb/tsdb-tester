package pkg

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Tester struct {
	database string
	testcase string
	cli      *Client

	testResults []TestResult
	fail        int32
	success     int32
}

func NewTester(testcase string, cli *Client) *Tester {
	return &Tester{
		testcase: testcase,
		cli:      cli,
	}
}

func (t *Tester) Run() (success, fail int32) {
	staetments := t.loadQueries()
	results := t.loadResults()
	if len(results) != len(staetments) {
		fmt.Println(staetments)
		panic(fmt.Sprintf("expect results(%d) not match tests(%d): %s", len(results), len(staetments), t.testcase))
	}

	for i := range staetments {
		stmt := staetments[i]
		start := time.Now()
		data, err := t.cli.Put(map[string]any{
			"db":  t.database,
			"sql": stmt.sql,
		})
		end := time.Now()
		var msg string
		if err != nil {
			fail++
			msg = err.Error()
		} else {
			ok, errMsg := equal(results[i].data, string(data))
			if ok {
				success++
			} else {
				msg = errMsg
				fail++
			}
		}
		t.testResults = append(t.testResults, TestResult{
			stmt: stmt,
			cost: end.Sub(start),
			msg:  msg,
		})
	}
	t.fail = fail
	t.success = success

	return
}

func (t *Tester) loadResults() (results []Result) {
	data, err := os.ReadFile(fmt.Sprintf("./t/%s.result", t.testcase))
	if err != nil {
		panic(err)
	}
	seps := bytes.Split(data, []byte("\n"))
	newStmt := true
	var stmt string
	for _, v := range seps {
		v = bytes.TrimSpace(v)
		s := string(v)
		switch {
		case len(s) == 0:
			// skip empty line
			continue
		case strings.HasPrefix(s, "#"):
			// skip # comment
			newStmt = true
			continue
		}

		if newStmt {
			stmt = s
		} else {
			stmt = fmt.Sprintf("%s\n%s", stmt, s)
		}

		// if the line has a ; in the end, we will treat new line as the new statement.
		newStmt = strings.HasSuffix(s, ";")
		if newStmt {
			results = append(results, Result{
				data: trim(stmt, ";"),
			})
			stmt = ""
		}
	}
	return
}

func (t *Tester) setting(stmt string) {
	if strings.HasPrefix(stmt, "@use") {
		t.database = trim(strings.TrimPrefix(stmt, "@use"), ";")
	}
}

func (t *Tester) loadQueries() (statments []Statement) {
	data, err := os.ReadFile(fmt.Sprintf("./t/%s.test", t.testcase))
	if err != nil {
		panic(err)
	}

	seps := bytes.Split(data, []byte("\n"))
	newStmt := true
	var stmt string
	var desc string
	for _, v := range seps {
		v = bytes.TrimSpace(v)
		s := string(v)
		switch {
		case len(s) == 0:
			// skip empty line
			continue
		case strings.HasPrefix(s, "#"):
			desc = strings.TrimPrefix(s, "#")
			// skip # comment
			newStmt = true
			continue
		case strings.HasPrefix(s, "@"):
			// setting statement
			t.setting(s)
			newStmt = true
			desc = ""
			continue
		}
		if newStmt {
			stmt = s
		} else {
			stmt = fmt.Sprintf("%s\n%s", stmt, s)
		}
		// if the line has a ; in the end, we will treat new line as the new statement.
		newStmt = strings.HasSuffix(s, ";")
		if newStmt {
			statments = append(statments, Statement{
				desc: strings.TrimSpace(desc),
				sql:  trim(stmt, ";"),
			})
			desc = ""
			stmt = ""
		}
	}
	return
}

func (t *Tester) String() string {
	var msg []string

	msg = append(msg, fmt.Sprintf("=== RUN   %s", t.testcase))
	totalCost := float64(0)
	for _, rs := range t.testResults {
		msg = append(msg, (&rs).RunMsg(t.testcase))
		totalCost += rs.cost.Seconds()
	}
	if t.fail > 0 {
		msg = append(msg, color.RedString("--- FAIL: %s (%.2f sec)", t.testcase, totalCost))
	} else {
		msg = append(msg, color.GreenString("--- PASS: %s (%.2f sec)", t.testcase, totalCost))
	}
	for _, rs := range t.testResults {
		msg = append(msg, (&rs).ResultMsg(t.testcase))
	}
	if t.fail > 0 {
		msg = append(msg, color.RedString("fail %s ", fmt.Sprintf("./t/%s.test", t.testcase)))
	} else {
		msg = append(msg, color.GreenString("ok %s ", fmt.Sprintf("./t/%s.test", t.testcase)))
	}
	return strings.Join(msg, "\n")
}
