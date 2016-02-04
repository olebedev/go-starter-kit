package main

var (
	// Debug var to switch mode from outside
	debug string
	// CommitHash exported to assign it from main.go
	commitHash string
)

// Most easiest way to configure
// an application is define config as
// yaml string and then parse it into
// map.
// How it works see here:
//     https://github.com/olebedev/config
var confString = `
debug: false
commit: 0
port: 5000
title: Go Starter Kit
api:
  prefix: /api
duktape:
  path: static/build/bundle.js
`
