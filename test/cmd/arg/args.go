package arg

import (
	"flag"
	"fmt"
	"strings"
)

type (
	Arg[t any] struct {
		Name, Tip string
		Value     t
	}
	Args struct{}
)

var (
	Lab    = Arg[FlagLab]{Name: "lab", Tip: "运行实验小节"}
	Run    = Arg[FlagRun]{Name: "run", Tip: "运行某个测试"}
	Count  = Arg[int]{Name: "count", Value: 1, Tip: "运行次数"}
	Batch  = Arg[bool]{Name: "batch", Tip: "使用`go test -failfast -count`批量一次执行"}
	Envs   = Arg[FlagEnvs]{Name: "env", Tip: "以key=value格式设置环境变量（可多次使用）"}
	Docker = Arg[bool]{Name: "docker", Tip: "在docker中运行"}
)

func (a Args) GetDockerCmd(dir string) string {
	var input []string
	flag.Visit(
		func(f *flag.Flag) {
			if f.Name == Docker.Name {
				return
			}
			switch value := f.Value.(type) {
			case interface{ IsBoolFlag() bool }:
				if value.IsBoolFlag() {
					input = append(input, "-"+f.Name)
					return
				}
			case *FlagEnvs:
				for _, env := range *value {
					input = append(input, "-"+f.Name, env)
				}
				return
			}
			input = append(input, "-"+f.Name, f.Value.String())
		},
	)
	return fmt.Sprintf("cd %s && go run . %s", dir, strings.Join(input, " "))
}
