# GoTrade

> GoTrade is a **FIX** protocol electronic trading and order management system written in Golang, structured for typical multi-asset instituional use

<p align="center">
  <img src="https://cdn.rawgit.com/cyanly/gotrade/gh-pages/orderrouting.svg" alt=""/>
</p>

[![GoDoc](https://godoc.org/github.com/cyanly/gotrade?status.png)](https://godoc.org/github.com/cyanly/gotrade) [![Build Status](https://travis-ci.org/quickfixgo/quickfix.svg?branch=master)](https://travis-ci.org/quickfixgo/quickfix)

## Status
This project is currently more of a proof of concept. It is no where near in completeness of a commerical product. This public repo serves as mostly for the purpose of experimenting and share of ideas.

## Getting Started
```
$ go get github.com/cyanly/gotrade
```

## Features

- [x] Trade in real-time via FIX through the broker-neutral API.
- [x] Normalized FIX order flow behavior across multiple FIX versions and asset classes.
- [x] Pure Go.
  - [x] Platform neutral: write once, build for any operating systems and arch (Linux/Windows/OSX etc).
  - [x] Native code performance.
  - [x] Ease of deployment.
  - [x] Lack of OOP verbosity, works for small and big teams.
- [x] Protobuf.
  - [x] Binary encoding format, efficient yet extensible.
  - [x] Easy Language Interoperability (C++, Python, Java, C#, Javascript, etc).
  - [x] Protocol backward compatibility.

## Design
```
└─ gotrade/
   ├─ core/                 -> The low-level API that gives consumers all the knobs they need
   │  ├─ order/
   │  │  └─ execution/
   │  ├─ service/
   └─ proto/...             -> Protobuf messaging protocol of various entities
   └─ services/             -> Core services managing multi-asset order flow
   │  ├─ orderrouter/       -> Centralized management facility for multi-asset global order flow
   │  ├─ marketconnectors/  -> Managing FIX connection to each trading venue, also performs pre-trade risk checks
   │
   └─cmd/...                -> Command-line executables, binary build targets
   
```


## Examples
**OrderRouter** and **MarketConnector** test cases will mock a testdb and messaging bus for end-to-end, message to message test. 

Pre-Requisites:
  - Go 1.3 or higher
  - ``` go get github.com/erikstmartin/go-testdb ```

Run test cases in services:
```
$ cd $GOPATH/src/github.com/cyanly/gotrade/services/orderrouter
$ go test -v 

$ cd $GOPATH/src/github.com/cyanly/gotrade/services/marketconnectors/simulator
$ go test -v 
```

<p align="center">
  <img src="https://cdn.rawgit.com/cyanly/gotrade/gh-pages/servicestest.png" alt=""/>
</p>

## Benchmark


## Limitations


## Thanks

**GoTrade** © 2016+, Chao Yan. Released under the [GNU] General Public License.<br>
Authored and maintained by Chao Yan with help from contributors ([list][contributors]).

> [cyan.ly](http://cyan.ly) &nbsp;&middot;&nbsp;
> GitHub [@cyanly](https://github.com/cyanly) &nbsp;&middot;&nbsp;

[MIT]: http://www.gnu.org/licenses/gpl-3.0.en.html
[contributors]: http://github.com/cyanly/gotrade/contributors
