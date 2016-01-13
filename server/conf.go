package server

import "fmt"

var Debug, CommitHash, BuildTime string

// Most easiest way to configure
// an application is define config as
// yaml string and then parse it into
// map.
// How it works see here:
//     https://github.com/olebedev/config
func confString() string {
	return fmt.Sprintf(`
debug: %s
commit: %s
buildTime: %s
port: 5000
title: Go Starter Kit
api:
  prefix: /api
duktape:
  path: static/build/bundle.js
`,
		Debug,
		CommitHash,
		BuildTime,
	)
}
