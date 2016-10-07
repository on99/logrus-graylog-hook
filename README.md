# Graylog Hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:" />
Use this hook to send logs to multiple Graylog servers over UDP.

====

[![Build Status](https://travis-ci.org/on99/logrus-graylog-hook.svg?branch=master)](https://travis-ci.org/on99/logrus-graylog-hook)
[![Go Report Card](http://goreportcard.com/badge/on99/logrus-graylog-hook)](http://goreportcard.com/report/on99/logrus-graylog-hook)
[![codecov](https://codecov.io/gh/on99/logrus-graylog-hook/branch/master/graph/badge.svg)](https://codecov.io/gh/on99/logrus-graylog-hook)

## Motivation
There is no way to setup a UDP load balancer in front of a Graylog server cluster due to the [chunking feature](http://docs.graylog.org/en/2.1/pages/gelf.html#chunking), which requires all chunks of a GELF message being sent to the same Graylog server node.

Client side load balancing provided by this logrus hook could be a solution to this issue.

## Algorithm
[Weighted-random algorithm](http://eli.thegreenplace.net/2010/01/22/weighted-random-generation-in-python/) is used to select Graylog server node as log destination every time when the hook is fired by `logrus`.

## Install
This package depends on 

* [gelf](https://github.com/Graylog2/go-gelf/gelf)
* [logrus](https://github.com/Sirupsen/logrus)

`$ go get -u github.com/on99/logrus-graylog-hook`

## Usage
```go
package main

import (
	"io/ioutil"
	"time"

	"github.com/Sirupsen/logrus"
	graylog "github.com/on99/logrus-graylog-hook"
)

func main() {
	// graylog configuration
	config := graylog.NewConfig()
	config.Facility = "graylog_hook"
	config.HealthCheckInterval = 5 * time.Second // check graylog nodes health periodically
	config.StaticMeta = map[string]interface{}{  // static meta that always sent to graylog
		"go_version": "1.7.1",
	}

	// create hook
	hook := graylog.New(config)

	// set graylog nodes
	// 5/10 chances log will be sent to node-1
	// 3/10 chances log will be sent to node-2
	// 2/10 chances log will be sent to node-3
	hook.SetNodeConfigs(
		graylog.NodeConfig{
			UDPAddress:     "node-1.graylog:12201",
			HealthCheckURL: "node-1.graylog/api/system/lbstatus",
			Weight:         5,
		},
		graylog.NodeConfig{
			UDPAddress:     "node-2.graylog:12201",
			HealthCheckURL: "node-2.graylog/api/system/lbstatus",
			Weight:         3,
		},
		graylog.NodeConfig{
			UDPAddress:     "node-3.graylog:12201",
			HealthCheckURL: "node-3.graylog/api/system/lbstatus",
			Weight:         2,
		},
	)

	// start health check
	// all graylog nodes are alive by default
	hook.StartHealthCheck()

	logrus.AddHook(hook)             // add graylog hook to logrus
	logrus.SetOutput(ioutil.Discard) // discard logrus output

	// log something, enjoy a cup of coffee
	logrus.WithField("key", "value").Info("log message")
}
```

