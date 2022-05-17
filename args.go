package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type parsedArgs struct {
	query                 string
	hook                  bool
	officialOnly          bool
	popularCommunityAndUp bool
	topResults            int
}

func parseArgs() parsedArgs {
	var query string
	hook := flag.Bool("hook", false, "Prints a hook that redirects 'scoop search' here")
	officialOnly := *flag.Bool("known-only", false, "When searching online, only return known buckets") || *flag.Bool("official-only", false, "Equivalent of --known-only")
	help := flag.Bool("help", false, "Prints this help message")
	topResults := flag.Int("top", 100, "Returns the top N results when searching online")
	popularCommunity := flag.Bool("popular", false, "When searching online, only returns popular buckets (50 stars and up)")
	flag.Parse()

	if *help {
		fmt.Printf("Usage: %s [OPTIONS] [<query>|/<query>/]\n", os.Args[0])
		fmt.Println("")
		fmt.Println("Performs search on all available buckets and online if local results are not found.")
		fmt.Println("If query is flanked by slashes, it is interpreted as a Golang regular expression")
		fmt.Println("locally but as plain text online.")
		fmt.Println("")
		fmt.Println("General:")
		fmt.Println("    --hook   \tPrints 'scoop search' hook")
		fmt.Println("    --help   \tPrints this help message")
		fmt.Println("")
		fmt.Println("Online searching:")
		fmt.Println("    --known  \tOnly show resuls from known buckets            \t[default: false]")
		fmt.Println("    --popular\tShow results from known buckets and            \t[default: false]")
		fmt.Println("             \tpopular community repositories (>50 stars)")
		fmt.Println("    --top N  \tReturns the top N results when searching online\t[default: 100]")
		os.Exit(1)
	} else {
		query = strings.Join(flag.Args(), " ")
	}

	return parsedArgs{query, *hook, officialOnly, *popularCommunity, *topResults}
}

const poshHook = `function scoop { if ($args[0] -eq "search") { scoop-search.exe @($args | Select-Object -Skip 1) } else { scoop.ps1 @args } }`
