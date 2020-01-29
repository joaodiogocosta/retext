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
	RootPaths []string
}

func Parse() *Args {
	flag.Parse()
	rootPaths := flag.Args()

	fmt.Println(rootPaths)
	if len(rootPaths) == 0 {
		fmt.Println("Please provide a file or directory as argument")
		os.Exit(0)
	}

	return &Args{RootPaths: rootPaths}
}
