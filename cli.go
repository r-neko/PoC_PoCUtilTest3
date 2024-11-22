package main

import (
	"flag"
	"fmt"
	"os"
)

func parseCLIArgs() (string, string) {
	// コマンドライン引数の解析
	name := flag.String("name", "", "Repository name (required)")
	license := flag.String("license", "MIT", "License type: MIT or Unlicense (default: MIT)")
	flag.Parse()

	if *name == "" {
		fmt.Println("Error: --name is required")
		os.Exit(1)
	}

	return "PoC_" + *name, *license
}
