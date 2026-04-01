package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
	template "github.com/nurazon59/go-template"
)

var CLI struct {
	Version struct{} `cmd:"" help:"Show version."`
	Config  struct{} `cmd:"" help:"Show config."`
}

func Run() {
	ctx := kong.Parse(&CLI)
	_ = template.Init("./testdata/config.yaml")
	switch ctx.Command() {
	case "version":
		fmt.Println("0.1.0")
	default:
		panic(ctx.Command())
	}
}
