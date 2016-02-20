// Copyright 2015-2016 Chao Yan. All rights reserved.

/*
GoTrade is a FIX protocol electronic trading and order management system
written in Golang, structured for typical multi-asset institutional use.

Dependencies:
  gogo/protobuf
    A fork of golang/protobuf with tweaks and extras
  quickfixgo/quickfix
    FIX Engine in Golang
  nats-io/nats
    Performant messaging bus in Golang, seems resemblance of TibRV

There are sub-packages within the gotrade package for various components:
  core/...:
    The low-level API that gives consumers all the knobs they need
  proto/...:
    The messaging protocol of various entities, in Protobuf format.
  services/...:
    Core services managing multi-asset order flow
  database/...: (not provided yet)
    SQL scripts and data layer APIs (PostgreSQL here but the idea is to
    support different storage engines without breaking everything)

Then on top of the core packages, we have:
  cmd/...:
    Command-line executables, a.k.a final products
  test/...:
    Integration tests, etc.

To avoid cyclic imports, imports should never pull in higher-level
APIs into a lower-level package.  For example, you could import all of
core and shell from cmd/... or test/..., but you couldn't import any
of shell from core/....

*/
package gotrade

import (
	_ "github.com/cyanly/gotrade/services/marketconnectors/simulator"
	_ "github.com/cyanly/gotrade/services/orderrouter"
)
