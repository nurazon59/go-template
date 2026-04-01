package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
	template "github.com/nurazon59/go-template"

	"github.com/adrg/xdg"
)

var CLI struct {
	Version struct{} `cmd:"" help:"Show version."`
	Config  struct{} `cmd:"" help:"Show config."`
}

func Run() {
	ctx := kong.Parse(&CLI)
	configPath, err := xdg.ConfigFile("go-template/config.yaml")
	if err != nil {
		template.Init("")
	}
	_ = template.Init(configPath)
	switch ctx.Command() {
	case "version":
		fmt.Println("0.1.0")
	default:
		panic(ctx.Command())
	}
}
