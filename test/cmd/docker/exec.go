// Package docker
// docker指令操作docker
// 因为docker API需要在linux环境下使用，window上使用不便所以使用docker指令操作docker
package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const WorkDir = "/mit6.824"

type Done = chan struct{}

func Build(ctx context.Context, imgName, dockerfilePath string) {
	execCmd(ctx, "docker", "build", "-t", imgName, dockerfilePath)
}

func Run(ctx context.Context, name, volume, imgName, command string) {
	execCmd(
		ctx, "docker", "run", "--rm", "--name", name, "-v", volume,
		"--entrypoint",
		"/bin/sh", imgName, "-c", command,
	)
}

func execCmd(ctx context.Context, name string, arg ...string) {
	fmt.Println(append([]string{name}, arg...))
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	_ = cmd.Run()
}

func Start(ctx context.Context, command string, file string, volume string) {
	imgName := "mit6.824-test"
	Build(ctx, imgName, file)
	name := fmt.Sprintf("mit6.824-%d", time.Now().Unix())
	defer exec.Command("docker", "rm", "-f", name).Run()
	Run(ctx, name, volume, imgName, command)
}
