package test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

func Exec(reportDir string, name string, tests []string, count int, batch bool, ctx context.Context) {
	test := newTest(reportDir, name, tests, count, batch, ctx)

	group, ctx := errgroup.WithContext(ctx)
	group.Go(test.checkAndPrint)
	group.Go(test.collectReport)

	for _, t := range tests {
		testName := t
		group.Go(
			func() error {
				test.execTest(testName, count)
				return nil
			},
		)
	}
	if err := group.Wait(); err != nil {
		panic(err)
	}
}

type Test struct {
	report       *Report
	testResultCh chan Result
	storagePath  string
	finish       atomic.Bool
	count        int
	batch        bool

	ctx context.Context
}

func newTest(reportPath, name string, tests []string, count int, batch bool, ctx context.Context) *Test {
	t := &Test{
		testResultCh: make(chan Result, len(tests)*10),
		storagePath:  reportPath,
		count:        count,
		batch:        batch,
		ctx:          ctx,
	}
	t.report = NewReport(name, t.getExecCount(), tests)
	return t
}

func (t *Test) checkAndPrint() error {
	execCount := t.getExecCount()
	for !t.finish.Load() {

		select {
		case <-t.ctx.Done():
			t.finish.Store(true)
			return nil
		default:
			t.printReport()
			finish := true
			for _, test := range t.report.tests {
				finish = finish && (t.report.success.LoadOrZero(test)+t.report.fail.LoadOrZero(test)) == execCount
			}
			if finish {
				t.printReport()
				t.finish.Store(true)
				close(t.testResultCh)
				return nil
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
	return nil
}

func (t *Test) collectReport() error {
	for response := range t.testResultCh {
		t.report.isUpdate.Store(true)
		if response.success {
			fmt.Print("\r", " Success ", response.name)
			t.report.success.Add(response.name, 1)
		} else {
			fmt.Print("\r", " Fail ", response.name)
			t.report.fail.Add(response.name, 1)
		}
	}
	return nil
}

func (t *Test) printReport() {
	if !t.report.isUpdate.CompareAndSwap(true, false) {
		return
	}
	t.report.makeReport()
	fmt.Print(
		fmt.Sprintf(
			"\u001B[25l\u001B[H\u001B[2J\u001B[1m----------------------%s test report --------------------------\u001B[0m\n Commed: go %s\n Fail report: %s\n%s\n",
			t.report.name,
			strings.Join(t.buildCommandArg("preview", t.count), " "),
			t.storagePath,
			strings.Join(t.report.content, "\n"),
		),
	)
}

type Result struct {
	name    string
	success bool
}

func (t *Test) execTest(name string, count int) {
	response, arg := Result{name: name}, t.buildCommandArg(name, count)
	for i := 0; i < t.getExecCount(); i++ {
		select {
		case <-t.ctx.Done():
			return
		default:
			cmd := exec.Command("go", arg...)
			log, err := cmd.CombinedOutput()
			if err != nil {
				response.success = false
				t.logErr(name, log, i)
			} else {
				response.success = true
			}
			t.testResultCh <- response
		}
	}
}

func (t *Test) buildCommandArg(name string, count int) []string {
	if t.batch {
		return []string{"test", "-run", name, "-failfast", "-count", strconv.Itoa(t.count)}
	}
	return []string{"test", "-run", name}
}

func (t *Test) getExecCount() int {
	if t.batch {
		return 1
	}
	return t.count
}

func (t *Test) logErr(name string, log []byte, number int) {
	dir := filepath.Join(t.storagePath, name)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(filepath.Join(dir, strconv.Itoa(number)+".ansi"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, _ = file.Write([]byte("----------------------------------" + name + " fail----------------------------------\n"))
	_, _ = file.Write(log)
}

type Report struct {
	name          string
	count         int
	tests         []string
	success, fail *RWMutexMap[string, int]
	isUpdate      atomic.Bool
	content       []string
}

func NewReport(name string, count int, tests []string) *Report {
	r := &Report{
		name:    name,
		count:   count,
		tests:   tests,
		content: make([]string, len(tests), len(tests)),
		success: NewRWMutexMap[string, int](),
		fail:    NewRWMutexMap[string, int](),
	}
	r.isUpdate.Store(true)
	for _, test := range r.tests {
		r.success.Store(test, 0)
		r.fail.Store(test, 0)
	}
	return r
}

const progressLen = 25

// makeReport
// title | progress | result
// %-25s | %-25s    | %10s
func (r *Report) makeReport() {
	for i, test := range r.tests {
		switch true {
		case r.success.LoadOrZero(test) == r.count:
			r.content[i] = fmt.Sprintf(
				" \u001B[34m%-25s \u001B[92m|%-25s\u001B[0m %10s",
				test,
				strings.Repeat("#", progressLen),
				"All pass",
			)
		case r.success.LoadOrZero(test)+r.fail.LoadOrZero(test) == r.count:
			if r.count == 0 {
				r.content[i] = fmt.Sprintf(
					" \u001B[34m%-25s \u001B[31m%s\u001B[0m", test, "Fail",
				)
			} else {
				r.content[i] = fmt.Sprintf(
					" \u001B[34m%-25s \u001B[31m%s%.2f\u001B[0m", test, "通过率：",
					float32(r.success.LoadOrZero(test))/float32(r.success.LoadOrZero(test)+r.fail.LoadOrZero(test)),
				)
			}

		default:
			progress := strings.Repeat("#", (r.success.LoadOrZero(test)+r.fail.LoadOrZero(test))*progressLen/r.count)
			log := fmt.Sprintf(
				" \u001B[34m%-25s \u001B[92m|%-25s\u001B[0m %10s", test,
				progress,
				fmt.Sprintf("%d/%d", r.success.LoadOrZero(test)+r.fail.LoadOrZero(test), r.count),
			)
			if r.fail.LoadOrZero(test) > 0 {
				log += fmt.Sprintf(" \u001B[31mfail:%d\u001B[0m", r.fail.LoadOrZero(test))
			}
			r.content[i] = log
		}
	}
}
