package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"mit-6.2840-test/cmd/docker"
)

var ReportDir = filepath.Join(RootDir, "test", "report")
var DockerReportDir = fmt.Sprintf("%s/test/report", docker.WorkDir)

func GetReportDir(t, lab, run string) (path string) {
	defer func() {
		path = filepath.Join(path, time.Now().Format("20060102-150405"))
		err := os.MkdirAll(path, 0777)
		if err != nil {
			panic(err)
		}
	}()
	if len(lab) > 0 {
		if len(run) > 0 {
			path = filepath.Join(ReportDir, lab, run)
		} else {
			path = filepath.Join(ReportDir, lab, "all")
		}
	}
	if len(t) > 0 {
		return filepath.Join(path, t)
	}
	return ReportDir
}

func ReportSave(path string, log []byte, number int) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(path, strconv.Itoa(number)+".ansi")
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, _ = file.Write(log)
}
