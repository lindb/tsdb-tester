package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/lindb/tsdb-tester/pkg"
)

var (
	url string
)

func init() {
	flag.StringVar(&url, "url", "http://127.0.0.1:9000/api/v1/exec", "The api url of LinDB broker.")
}

func main() {
	flag.Parse()
	var wait sync.WaitGroup
	tests, err := loadAllTests()
	if err != nil {
		log.Fatalf("load all tests err %v", err)
	}
	wait.Add(len(tests))
	cli := pkg.NewClient(url)

	totalSuccess := &atomic.Int32{}
	totalFail := &atomic.Int32{}

	for i := range tests {
		tt := tests[i]
		go func() {
			defer wait.Done()
			tr := pkg.NewTester(tt, cli)
			success, fail := tr.Run()
			totalSuccess.Add(success)
			totalFail.Add(fail)
			// print test result
			fmt.Println(tr.String())
		}()
	}
	wait.Wait()

	fmt.Printf("Summary: total %s, success %s, fail %s\n",
		color.BlueString("%d", totalFail.Load()+totalSuccess.Load()),
		color.GreenString("%d", totalSuccess.Load()),
		color.RedString("%d", totalFail.Load()))
}

func loadAllTests() ([]string, error) {
	tests := make([]string, 0)
	// tests must be in t folder or subdir in t folder
	err := filepath.Walk("./t/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".test") {
			name := strings.TrimPrefix(strings.TrimSuffix(path, ".test"), "t/")
			tests = append(tests, name)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tests, nil
}
