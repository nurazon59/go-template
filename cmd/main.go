package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Version struct{} `cmd:"" help:"Show version."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "version":
		fmt.Println("0.1.0")
	default:
		panic(ctx.Command())
	}
}
