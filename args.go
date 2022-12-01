package main

import "os"

type parsedArgs struct {
	query string
	hook  bool
}

func parseArgs() parsedArgs {
	var hook bool
	var query string

	if len(os.Args) == 1 {
		// pass
	} else if os.Args[1] == "--hook" {
		hook = true
	} else {
		query = os.Args[1]
	}

	return parsedArgs{query, hook}
}

const poshHook = `function shovel { if ($args[0] -eq "search") { shovel-search.exe @($args | Select-Object -Skip 1) } else { shovel.ps1 @args } }`
