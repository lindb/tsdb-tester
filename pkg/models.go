package pkg

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Statement struct {
	desc string
	sql  string
}

// String returns statement desc, if desc is mepty returns sql.
func (stmt *Statement) String() string {
	if stmt.desc != "" {
		return stmt.desc
	}
	strs := strings.Split(stmt.sql, "\n")
	if len(strs) > 1 {
		return strs[0] + " ..."
	}
	return strs[0]
}

type Result struct {
	data string
}

type TestResult struct {
	stmt Statement
	cost time.Duration
	msg  string
}

func (r *TestResult) RunMsg(parent string) string {
	return fmt.Sprintf("=== RUN   %s/%s", parent, r.stmt.String())
}

func (r *TestResult) ResultMsg(parent string) string {
	var msg []string
	if r.msg != "" {
		msg = append(msg, color.RedString("    --- FAIL: %s/%s (%.2f sec)", parent, r.stmt.String(), r.cost.Seconds()))
		strs := strings.Split(r.stmt.sql, "\n")
		msg = append(msg, fmt.Sprintf("\tStatement: %s", strings.Join(strs, "\n\t\t   ")))
		msg = append(msg, fmt.Sprintf("\t%s", r.msg))
	} else {
		msg = append(msg, color.GreenString("    --- PASS: %s/%s (%.2f sec)", parent, r.stmt.String(), r.cost.Seconds()))
	}
	return strings.Join(msg, "\n")
}
