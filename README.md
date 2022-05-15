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

Returns identical results as `scoop search`, though it doesn't behave the same as `scoop search`, so if that's what you want, please check upstream.

## Deviations from scoop search
### Regex searching:
```ps1
❯ scoop-search /balls?/
'extras' bucket:
    monero (0.17.3.2) --> includes 'monero-blockchain-blackball'

'games' bucket:
    neverball (1.6.0)
    paintball2 (build045)

'main' bucket:
    monero-cli (0.17.3.2) --> includes 'monero-blockchain-blackball'

'nonportable' bucket:
    virtualbox52-np (5.2.44) --> includes 'VBoxBalloonCtrl.exe'
    virtualbox-np (6.1.34) --> includes 'VBoxBalloonCtrl.exe'
    virtualbox-with-extension-pack-np (6.1.34) --> includes 'VBoxBalloonCtrl.exe'

'retools' bucket:
    dynamorio (9.0.1) --> includes 'balloon.exe'
```

### Online searching
```
❯ scoop-search noto
Results from other buckets...
'dodorz_scoop' bucket (https://github.com/dodorz/scoop):
    Noto-Font (2017.10.25)
    Noto-Mono-Font (2017.10.25)
    Noto-Serif-Font (2017.10.25)
    Noto-Sans-Font (2017.10.25)

'ShuguangSun_sgs-scoop-bucket' bucket (https://github.com/ShuguangSun/sgs-scoop-bucket):
    Noto-CJK-Mega-OTC (20190603)

'KnotUntied_scoop-fonts' bucket (https://github.com/KnotUntied/scoop-fonts):
    notomusic (2.000)
    notosansmath (2.001)
    notosansmono (2.006)

'nerd-fonts' bucket (install using 'scoop install nerd-fonts/<app>'):
    Noto-CJK-Mega-OTC (20190603)
    Noto-NF (2.1.0)
    Noto-NF-Mono (2.1.0)
    Source-Han-Noto-CJK-Ultra-OTC (20190603)
```

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
