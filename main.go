package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzuse/core"
)

func main() {
	app := cli.NewSingleProgram(&cli.SingleProgramConfig{
		Name:    "gzuse",
		Usage:   "gzuse is a portable usage",
		Version: Version,
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:    "memory-percent",
				Usage:   "used memory percent",
				EnvVars: []string{"MEMORY_PERCENT"},
			},
			&cli.UintFlag{
				Name:    "cpu-percent",
				Usage:   "used cpu percent",
				EnvVars: []string{"CPU_PERCENT"},
			},
			&cli.StringFlag{
				Name:    "memory-size",
				Usage:   "used memory size",
				EnvVars: []string{"MEMORY"},
			},
			&cli.UintFlag{
				Name:    "cpu-core",
				Usage:   "used cpu core",
				EnvVars: []string{"CPU_CORE"},
			},
		},
	})

	app.Command(func(ctx *cli.Context) error {
		return core.Run(&core.Config{
			MemoryPercent: ctx.Uint("memory-percent"),
			CPUPercent:    ctx.Uint("cpu-percent"),
			MemorySize:    ctx.String("memory-size"),
			CPUCore:       ctx.Uint("cpu-core"),
		})
	})

	app.Run()
}
