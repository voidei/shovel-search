package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type parsedArgs struct {
	query        string
	hook         bool
	officialOnly bool
}

func parseArgs() parsedArgs {
	var query string
	hook := flag.Bool("hook", false, "Prints a hook that redirects 'scoop search' here")
	officialOnly := *flag.Bool("known-only", false, "When searching online, only return known buckets") || *flag.Bool("official-only", false, "Equivalent of --known-only")
	help := flag.Bool("help", false, "Prints this help message")
	flag.Parse()

	if (flag.NArg() == 0 && !*hook) || *help {
		fmt.Printf("Usage: %s [OPTIONS] <query>\n", os.Args[0])
		fmt.Println("")
		fmt.Println("Performs search on all available buckets and online if local results are not found.")
		fmt.Println("")
		fmt.Println("    --hook      \tPrints 'scoop search' hook")
		fmt.Println("    --known-only\tWhen searching online, filter only known buckets")
		fmt.Println("    --help      \tPrints this help message")
		os.Exit(1)
	} else {
		query = strings.Join(flag.Args(), " ")
	}

	return parsedArgs{query, *hook, officialOnly}
}

const poshHook = `function scoop { if ($args[0] -eq "search") { scoop-search.exe @($args | Select-Object -Skip 1) } else { scoop.ps1 @args } }`
