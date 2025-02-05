package main

import (
	"context"
	"fmt"
	"io"
	test "mit-6.824-test"
	"mit-6.824-test/cmd/docker"
	"mit-6.824-test/util"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	args := getArgs()
	done := make(chan struct{})
	os.Environ()
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan
		fmt.Println("\r exiting...")
		close(done)
	}()
	start(args, done)
}

func start(args Args, done chan struct{}) {
	var title, runDir, command string
	lab, runName := args.Lab.GetLab(), args.Run.String()
	if lab != nil {
		if len(runName) > 0 {
			title = fmt.Sprintf("%s %s", lab.Name, lab.Name)
			command = fmt.Sprintf("time go test -run %s", args.Run)
		} else {
			title = lab.Name
			command = fmt.Sprintf("time go test -run %s", lab.Short)
		}
		runDir = lab.Path
	} else {
		panic("not run")
	}
	run(args, title, runDir, command, done)
}

func run(args Args, title, workDir, command string, done chan struct{}) {
	if args.Docker {
		command = args.GetDockerCmd("test/cmd/performance")
		volume := fmt.Sprintf("%s:%s", test.ReportDir, test.DockerReportDir)
		docker.Start(done, command, test.RootDir, volume)
		return
	}
	reportDir := test.GetReportDir("performance", args.Lab.String(), args.Run.String())
	err := os.Chdir(filepath.Join(test.RootDir, workDir))
	if err != nil {
		panic(err)
	}
	title = strings.ToTitle(fmt.Sprintf("%s performance test", title))
	fmt.Printf(
		"\u001B[2J\u001B[H\u001B[35m%s \u001B[0m\n"+
			"\u001B[34mRun: %s\u001B[0m\nReport: %s\n", title, command, util.ClickablePath(reportDir),
	)
	reportDir = filepath.Join(reportDir, time.Now().Format("20060102_150405"))
	err = execCommand(reportDir, "bash", "-c", command)
	if err != nil {
		panic(err)
	}
}

func execCommand(outputFile string, command string, args ...string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, command, args...)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nstopping ...")
		cancel()
	}()

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)
	cmd.Stdout, cmd.Stderr = multiWriter, multiWriter
	return cmd.Run()
}
