package cli

import (
	"flag"
	"github.com/joaodiogocosta/retext/client"
	"fmt"
	"os"
)

// TODO: Parse yml config file
// - https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64
// - https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64


type Args struct {
	RootPaths []string
	ConnectionAdapter int
}

func Parse() *Args {
	var dry bool
	flag.BoolVar(&dry, "dry", false, "print updates to the console")

	flag.Parse()
	rootPaths := flag.Args()

	if len(rootPaths) == 0 {
		fmt.Println("Please provide a file or directory as argument")
		os.Exit(0)
	}

	args := &Args{RootPaths: rootPaths}

	if dry {
		args.ConnectionAdapter = client.DryAdapter
	}

	return args
}
