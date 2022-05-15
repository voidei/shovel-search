# scoop-search

[![](https://goreportcard.com/badge/github.com/shilangyu/scoop-search)](https://goreportcard.com/report/github.com/shilangyu/scoop-search)
[![](https://github.com/shilangyu/scoop-search/workflows/ci/badge.svg)](https://github.com/shilangyu/scoop-search/actions)

Fast `scoop search` replacement with some extra features.

## Hook

Instead of using `scoop-search.exe <term>` you can setup a hook that will run `scoop-search.exe` whenever you use native `scoop search`

Add this to your Powershell profile (usually located at `$PROFILE`)

```ps1
Invoke-Expression (&scoop-search --hook)
```

## Features

Returns identical results as `scoop search`, though it can also take advantage of the [ScoopSearch API](https://github.com/ScoopInstaller/ScoopInstaller.AzureFunctions) to search online buckets.

## Benchmarks

Done with [hyperfine](https://github.com/sharkdp/hyperfine). `scoop-search` is on average 50 times faster.

```sh
❯ hyperfine --warmup 1 'scoop-search google' 'scoop search google'
Benchmark #1: scoop-search google
  Time (mean ± σ):      76.1 ms ±   1.9 ms    [User: 0.8 ms, System: 5.8 ms]
  Range (min … max):    73.4 ms …  82.7 ms    37 runs

Benchmark #2: scoop search google
  Time (mean ± σ):      3.910 s ±  0.015 s    [User: 1.4 ms, System: 7.9 ms]
  Range (min … max):    3.888 s …  3.928 s    10 runs

Summary
  'scoop-search google' ran
   51.37 ± 1.31 times faster than 'scoop search google'
```

_ran on AMD Ryzen 5 3600 @ 3.6GHz_

## Building
To enable online lookup, you'll need to link an API key:
```sh
# How you find this API key is up to you
go build -ldflags="-X 'main.scoopSearchApiKey=something'"
```
