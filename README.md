# scoop-search

[![](https://goreportcard.com/badge/github.com/shilangyu/scoop-search)](https://goreportcard.com/report/github.com/shilangyu/scoop-search)
[![](https://github.com/shilangyu/scoop-search/workflows/ci/badge.svg)](https://github.com/shilangyu/scoop-search/actions)

Fast `scoop search` drop-in replacement 🚀
Forked from [scoop-search](https://github.com/shilangyu/scoop-search) by [shilangyu](https://github.com/shilangyu)

## Installation

```sh
scoop install https://raw.githubusercontent.com/voidei/shovel-search/master/shovel-search.json
```

## Hook

Instead of using `shovel-search.exe <term>` you can setup a hook that will run `shovel-search.exe` whenever you use native `shovel search`

Add this to your Powershell profile (usually located at `$PROFILE`)

```ps1
Invoke-Expression (&shovel-search --hook)
```

## Features

Behaves just like `shovel search` and returns identical output. If any differences are found please open an issue.

**Non-goal**: any additional features unavailable in scoop search

## Benchmarks

Done with [hyperfine](https://github.com/sharkdp/hyperfine). `shovel-search` is on average 50 times faster.

# TODO - replace this w/ my own shovel-search results, so i'm not ripping off the original author's scoop-search results
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
