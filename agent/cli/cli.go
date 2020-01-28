package cli

import (
	"flag"
	"fmt"
	"os"
)

// TODO: Parse yml config file
// - https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64
// - https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64


type Args struct {
	RootPath string
}

func Parse() *Args {
	flag.Parse()
	rootPath := flag.Arg(0)

	if rootPath == "" {
		fmt.Println("Please provide a file or directory as argument")
		os.Exit(0)
	}

	return &Args{RootPath: rootPath}
}
