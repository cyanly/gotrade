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
$ go get -u github.com/cyanly/gotrade

$ cd $GOPATH/src/github.com/cyanly/gotrade
$ go get -u -t ./...
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

Pre-Requisites:
  - Go 1.4 or higher
  - ``` go get github.com/erikstmartin/go-testdb ```
  - ``` go get github.com/nats-io/gnatsd ```


The best way to see goTrade in action is to take a look at tests (see Benchmark section below):<br>
**OrderRouter** and **MarketConnector** test cases will mock a testdb and messaging bus for end-to-end, message to message test. 
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

Machine: `Intel Core i5 CPU @ 2.80GHz` + `Ubuntu 14.04 Desktop x86_64`

  - `test/benchmark/client2fix_test.go`
  - **CL &#8658; OR**:   <br>*Client send order protobuf to OrderRouter(OR)*
  - **OR &#8658; MC**:   <br>*OrderRouter process order and dispatch persisted order entity to target MarketConnector*
  - **MC &#8658; FIX**:  <br>*MarketConnector translate into NewOrderSingle FIX message based on the session with its counterparty*
  - **FIX &#8658; MC**:  <br>*MarketConnector received FIX message on its order, here Simulator sending a fully FILL execution*
  - **EXE &#8658; CL**:  <br>*MarketConnector publish processed and persisted Execution onto messaging bus, here our Client will listen to*

Included: 
  - from order to FIX to a fully fill execution message to execution protobuf published back
  - serialsing/deserialsing mock order into protobuf messages
  - Request/Publish and Response/Subscribe via NATS.io message bus
  - Time spent in the Linux TCP/IP stack
  - Decode FIX messages and reply by a simulated broker
  
Excluded:
  - Database transaction time (hard-wired to an inline mock DB driver) 

Result:   
  - **`0.176ms per op,  5670 order+fill pairs per sec`**
<p align="center">
  <img src="https://cdn.rawgit.com/cyanly/gotrade/gh-pages/benchmark.png" alt=""/>
</p>


## Limitations


## Contributing

**GoTrade** © 2016+, Chao Yan. Released under the [GNU] General Public License.<br>
Authored and maintained by Chao Yan with help from contributors ([list][contributors]). <br>
Contributions are welcome. 

> [cyan.ly](http://cyan.ly) &nbsp;&middot;&nbsp;
> GitHub [@cyanly](https://github.com/cyanly) &nbsp;&middot;&nbsp;

[GNU]: http://www.gnu.org/licenses/gpl-3.0.en.html
[contributors]: http://github.com/cyanly/gotrade/contributors
