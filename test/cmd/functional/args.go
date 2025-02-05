package main

import (
	"flag"
	"mit-6.824-test/cmd/arg"
)

type Args struct {
	arg.Args
	Lab    arg.FlagLab
	Run    arg.FlagRun
	Count  int
	Batch  bool
	Docker bool
	Envs   arg.FlagEnvs
}

func getArgs() (args Args) {
	flag.Var(&args.Lab, arg.Lab.Name, arg.Lab.Tip)
	args.Run = arg.NewFlagRun(args.Lab.String())
	flag.Var(&args.Run, arg.Run.Name, arg.Run.Tip)
	flag.Var(&args.Envs, arg.Envs.Name, arg.Envs.Tip)

	flag.IntVar(&args.Count, arg.Count.Name, arg.Count.Value, arg.Lab.Tip)
	flag.BoolVar(&args.Batch, arg.Batch.Name, arg.Batch.Value, arg.Batch.Tip)
	flag.BoolVar(&args.Docker, arg.Docker.Name, arg.Docker.Value, arg.Docker.Tip)

	flag.Parse()
	// 如果没输入lab 但是输入了run 则设置lab为run所在的lab
	if len(args.Lab.String()) == 0 && len(args.Run.String()) > 0 {
		_ = args.Lab.Set(args.Run.GetLab().Name)
	}
	return
}
