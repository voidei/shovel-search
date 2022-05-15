package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/valyala/fastjson"
)

type (
	match struct {
		name, version, bin string
	}
	matchMap = map[string][]match
)

var (
	knownBuckets      map[string]string
	scoopSearchApiKey string
)

// resolves the path to scoop folder
func scoopHome() (res string) {
	if value, ok := os.LookupEnv("SCOOP"); ok {
		res = value
	} else {
		var err error
		res, err = os.UserHomeDir()
		checkWith(err, "Could not determine home dir")

		res += "\\scoop"
	}

	return
}

func scoopKnownRepos() (res map[string]string) {
	res = make(map[string]string)
	var parser fastjson.Parser

	raw, err := os.ReadFile(scoopHome() + "\\apps\\scoop\\current\\buckets.json")
	check(err)

	result, _ := parser.ParseBytes(raw)
	object, _ := result.Object()

	object.Visit(func(k []byte, v *fastjson.Value) {
		res[string(v.GetStringBytes())] = string(k)
	})
	return
}

func main() {
	args := parseArgs()
	knownBuckets = scoopKnownRepos()

	// print posh hook and exit if requested
	if args.hook {
		fmt.Println(poshHook)
		os.Exit(0)
	}

	// get buckets path
	bucketsPath := scoopHome() + "\\buckets"

	hasResults := printResults(scoopLocalSearch(bucketsPath, args.query), false)
	// print results and exit with status code
	if !hasResults && scoopSearchApiKey != "" {
		hasResults = printResults(scoopSearchAPI(args.query, args.officialOnly), true)
	}

	if !hasResults {
		fmt.Println("No results found.")
		os.Exit(1)
	}
}

func scoopLocalSearch(bucketsPath string, term string) matchMap {
	buckets, err := os.ReadDir(bucketsPath)
	checkWith(err, "Scoop folder does not exist")

	// start workers that will find matching manifests
	matches := struct {
		sync.Mutex
		data matchMap
	}{}
	matches.data = make(matchMap)
	var wg sync.WaitGroup

	for _, bucket := range buckets {
		if !bucket.IsDir() {
			continue
		}

		wg.Add(1)
		go func(file os.DirEntry) {
			// check if $bucketName/bucket exists, if not use $bucketName
			bucketPath := bucketsPath + "\\" + file.Name()
			if f, err := os.Stat(bucketPath + "\\bucket"); !os.IsNotExist(err) && f.IsDir() {
				bucketPath += "\\bucket"
			}

			res := matchingManifests(bucketPath, term)
			matches.Lock()
			matches.data[file.Name()] = res
			matches.Unlock()
			wg.Done()
		}(bucket)
	}
	wg.Wait()
	return matches.data
}

func scoopSearchAPI(term string, officialOnly bool) matchMap {
	var arena fastjson.Arena
	var parser fastjson.Parser

	body := arena.NewObject()
	body.Set("count", arena.NewTrue())
	if officialOnly {
		body.Set("filter", arena.NewString("Metadata/OfficialRepositoryNumber eq 1"))
	} else {
		body.Set("filter", arena.NewString(""))
	}
	body.Set("highlight", arena.NewString(""))
	body.Set("highlightPreTag", arena.NewString(""))
	body.Set("highlightPostTag", arena.NewString(""))
	body.Set("orderby", arena.NewString("search.score() desc, Metadata/OfficialRepositoryNumber desc, NameSortable asc"))
	body.Set("search", arena.NewString(term))
	body.Set("searchMode", arena.NewString("all"))
	body.Set("select", arena.NewString("Name,NamePartial,NameSuffix,Version,Metadata/Repository,Metadata/OfficialRepository"))
	body.Set("skip", arena.NewNumberInt(0))
	body.Set("top", arena.NewNumberInt(20))

	request, err := http.NewRequest(
		"POST",
		"https://scoopsearch.search.windows.net/indexes/apps/docs/search?api-version=2020-06-30",
		bytes.NewReader(body.MarshalTo(nil)),
	)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return make(matchMap)
		} else {
			check(err)
		}
	}
	request.Header.Set("api-key", scoopSearchApiKey)
	request.Header.Set("content-type", "application/json")

	response, err := http.DefaultClient.Do(request)
	check(err)
	defer response.Body.Close()

	raw, err := ioutil.ReadAll(response.Body)
	check(err)

	result, _ := parser.ParseBytes(raw)
	values := result.GetArray("value")

	var wg sync.WaitGroup
	matches := struct {
		sync.Mutex
		data matchMap
	}{}
	matches.data = make(matchMap)
	for _, searchResult := range values {
		wg.Add(1)
		go func(searchRes *fastjson.Value) {
			defer wg.Done()

			var bucket string
			if searchRes.Get("Metadata").GetBool("OfficialRepository") {
				bucket = knownBuckets[string(searchRes.Get("Metadata").GetStringBytes("Repository"))]
			} else {
				bucketSplit := strings.Split(string(searchRes.Get("Metadata").GetStringBytes("Repository")), "/")
				bucket = bucketSplit[len(bucketSplit)-2] + "_" + bucketSplit[len(bucketSplit)-1]
			}
			matches.Lock()
			matches.data[bucket] = append(matches.data[bucket], match{
				string(searchRes.GetStringBytes("Name")),
				string(searchRes.GetStringBytes("Version")),
				"",
			})
			matches.Unlock()
		}(searchResult)
	}
	wg.Wait()
	return matches.data
}

func matchingManifests(path string, term string) (res []match) {
	term = strings.ToLower(term)
	files, err := os.ReadDir(path)
	check(err)

	var parser fastjson.Parser

	for _, file := range files {
		name := file.Name()

		// its not a manifest, skip
		if !strings.HasSuffix(name, ".json") {
			continue
		}

		// parse relevant data from manifest
		raw, err := os.ReadFile(path + "\\" + name)
		check(err)
		result, _ := parser.ParseBytes(raw)

		version := string(result.GetStringBytes("version"))

		stem := name[:len(name)-5]

		if strings.Contains(strings.ToLower(stem), term) {
			// the name matches
			res = append(res, match{stem, version, ""})
		} else {
			// the name did not match, lets see if any binary files do
			var bins []string
			bin := result.Get("bin") // can be: nil, string, [](string | []string)

			if bin == nil {
				// no binaries
				continue
			}

			const badManifestErrMsg = `Cannot parse "bin" attribute in a manifest. This should not happen. Please open an issue about it with steps to reproduce`

			switch bin.Type() {
			case fastjson.TypeString:
				bins = append(bins, string(bin.GetStringBytes()))
			case fastjson.TypeArray:
				for _, stringOrArray := range bin.GetArray() {
					switch stringOrArray.Type() {
					case fastjson.TypeString:
						bins = append(bins, string(stringOrArray.GetStringBytes()))
					case fastjson.TypeArray:
						// check only first two, the rest are command flags
						stringArray := stringOrArray.GetArray()
						bins = append(bins, string(stringArray[0].GetStringBytes()))
						if len(stringArray) > 1 {
							bins = append(bins, string(stringArray[1].GetStringBytes()))
						}
					default:
						log.Fatalln(badManifestErrMsg)
					}
				}
			default:
				log.Fatalln(badManifestErrMsg)
			}

			for _, bin := range bins {
				bin = filepath.Base(bin)
				if strings.Contains(strings.ToLower(strings.TrimSuffix(bin, filepath.Ext(bin))), term) {
					res = append(res, match{stem, version, bin})
					break
				}
			}
		}
	}

	sort.SliceStable(res, func(i, j int) bool {
		// case insensitive comparison where hyphens are ignored
		return strings.ToLower(strings.ReplaceAll(res[i].name, "-", "")) <= strings.ToLower(strings.ReplaceAll(res[j].name, "-", ""))
	})

	return
}

func printResults(data matchMap, fromScoopSearch bool) (anyMatches bool) {
	// sort by bucket names
	entries := 0
	sortedKeys := make([]string, 0, len(data))
	for k := range data {
		entries += len(data[k])
		sortedKeys = append(sortedKeys, k)
	}

	if fromScoopSearch {
		// Hoisting known buckets down to the bottom
		sort.SliceStable(sortedKeys, func(i, j int) bool {
			isIOfficial := !strings.Contains(sortedKeys[i], "_")
			isJOfficial := !strings.Contains(sortedKeys[j], "_")
			if isIOfficial && !isJOfficial {
				return false
			} else if !isIOfficial && isJOfficial {
				return true
			} else {
				return sortedKeys[i] > sortedKeys[j]
			}
		})
	} else {
		sort.Strings(sortedKeys)
	}

	for _, k := range sortedKeys {
		v := data[k]

		if len(v) > 0 {
			anyMatches = true
			break
		}
	}

	// reserve additional space assuming each variable string has length 1. Will save time on initial allocations
	var display strings.Builder
	display.Grow((len(sortedKeys)*12 + entries*11))

	if anyMatches && fromScoopSearch {
		fmt.Println("Results from other buckets...")
	}

	for _, k := range sortedKeys {
		v := data[k]

		if len(v) > 0 {
			anyMatches = true
			display.WriteString("'")
			display.WriteString(k)
			display.WriteString("' bucket")
			if fromScoopSearch {
				if strings.Contains(k, "_") {
					display.WriteString(" (https://github.com/")
					display.WriteString(strings.Replace(k, "_", "/", 1))
					display.WriteString("):\n")
				} else {
					display.WriteString(" (install using 'scoop install ")
					display.WriteString(k)
					display.WriteString("/<app>'):\n")
				}
			} else {
				display.WriteString(":\n")
			}
			for _, m := range v {
				display.WriteString("    ")
				display.WriteString(m.name)
				display.WriteString(" (")
				display.WriteString(m.version)
				display.WriteString(")")
				if m.bin != "" {
					display.WriteString(" --> includes '")
					display.WriteString(m.bin)
					display.WriteString("'")
				}
				display.WriteString("\n")
			}
			display.WriteString("\n")
		}
	}

	os.Stdout.WriteString(display.String())
	return
}
