package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	test "mit-6.824-test"
	"mit-6.824-test/cmd/docker"
)

func main() {
	args := getArgs()
	var tests []string
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan
		fmt.Println("\r exiting...")
		cancel()
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
		run(name, wordPath, tests, args, ctx)
	} else {
		for _, lab := range test.Labs {
			run(lab.Name, lab.Path, lab.Tests, args, ctx)
		}
	}
}

func run(name, workDir string, tests []string, args Args, ctx context.Context) {
	if args.Docker {
		command := args.GetDockerCmd("test/cmd/functional")
		volume := fmt.Sprintf("%s:%s", test.ReportDir, test.DockerReportDir)
		docker.Start(ctx, command, test.RootDir, volume)
		return
	}
	reportDir := test.GetReportDir("functional", args.Lab.String(), args.Run.String())
	err := os.Chdir(filepath.Join(test.RootDir, workDir))
	if err != nil {
		panic(err)
	}
	test.Exec(reportDir, name, tests, args.Count, args.Batch, ctx)
}
