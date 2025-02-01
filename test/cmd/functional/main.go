package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	test "mit-6.2840-test"
	"mit-6.2840-test/cmd/docker"
)

func main() {
	args := getArgs()
	var tests []string

	done := make(chan struct{})
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan
		fmt.Println("\r exiting...")
		close(done)
	}()
	var name, wordPath string
	lab, runName := args.Lab.GetLab(), args.Run.String()
	if lab != nil {
		if len(runName) > 0 {
			tests = []string{runName}
		} else {
			tests = lab.Tests
		}
		name, wordPath = args.Lab.String(), lab.Path
		run(name, wordPath, tests, args, done)
	} else {
		for _, lab := range test.Labs {
			run(lab.Name, lab.Path, lab.Tests, args, done)
		}
	}
}

func run(name, workDir string, tests []string, args Args, done chan struct{}) {
	if args.Docker {
		command := args.GetDockerCmd("test/cmd/functional")
		volume := fmt.Sprintf("%s:%s", test.ReportDir, test.DockerReportDir)
		docker.Start(done, command, test.RootDir, volume)
		return
	}
	reportDir := test.GetReportDir("functional", args.Lab.String(), args.Run.String())
	err := os.Chdir(filepath.Join(test.RootDir, workDir))
	if err != nil {
		panic(err)
	}
	test.Exec(reportDir, name, tests, args.Count, args.Batch, done)
}
