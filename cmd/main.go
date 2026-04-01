package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	template "github.com/nurazon59/go-template"
)

var CLI struct {
	Version struct{} `cmd:"" help:"Show version."`
	Config  struct{} `cmd:"" help:"Show config."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "version":
		fmt.Println("0.1.0")
	case "config":
		template.Init()
	default:
		panic(ctx.Command())
	}
}
