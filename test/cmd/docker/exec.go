// Package docker
// docker指令操作docker
// 因为docker API需要在linux环境下使用，window上使用不便所以使用docker指令操作docker
package docker

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const WorkDir = "/mit6.824"

type Done = chan struct{}

func Build(done Done, imgName, dockerfilePath string) {
	execCmd(done, "docker", "build", "-t", imgName, dockerfilePath)
}

func Run(done Done, name, volume, imgName, command string) {
	execCmd(
		done, "docker", "run", "--rm", "--name", name, "-v", volume,
		"--entrypoint",
		"/bin/sh", imgName, "-c", command,
	)
}

func execCmd(done Done, name string, arg ...string) {
	fmt.Println(append([]string{name}, arg...))
	cmd := exec.Command(name, arg...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	if done != nil {
		go func() {
			<-done
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
		}()
	}
	_ = cmd.Run()
}

func Start(done Done, command string, file string, volume string) {
	imgName := "mit6.824-test"
	Build(done, imgName, file)
	name := fmt.Sprintf("mit6.824-%d", time.Now().Unix())
	defer exec.Command("docker", "rm", "-f", name)
	Run(done, name, volume, imgName, command)
}
