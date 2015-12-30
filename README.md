# go-starter-kit

[![Join the chat at https://gitter.im/olebedev/go-starter-kit](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/olebedev/go-starter-kit?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

This project contains a quick starter kit for **Facebook React** Single Page Apps with **Golang** server side render and also with a set of useful features for rapid development of efficient applications.

## What it contains?

* server side render via [go-duktape](https://github.com/olebedev/go-duktape)
* api requests between your react application and server side application directly  via [fetch polyfill](https://github.com/olebedev/go-duktape-fetch) for go-duktape at server side, and it is possible to process requests with user session as well
* title, Open Graph and other domain-specific meta tags render for each page at the server and at the client
* server side redirect
* embedding static files into artefact via bindata
* high performance [echo](https://github.com/labstack/echo) framework
* advanced cli via [cli](https://github.com/codegangsta/cli)
* Makefile based project
* one(!) terminal window process for development
* routing via [react-router](https://github.com/rackt/react-router)
* ES6 & JSX via [babel-loader](https://github.com/babel/babel-loader) with minimal runtime dependency footprint
* [redux](http://rackt.org/redux/) as state container
* [redux-devtools](https://github.com/gaearon/redux-devtools)
* stylus css styles without global namespace via [css-loader](https://github.com/webpack/css-loader) & css-modules
* separate css file to avoid FOUC
* hot reloading via [react-transform](https://github.com/gaearon/babel-plugin-react-transform) & [HMR](http://webpack.github.io/docs/hot-module-replacement.html)
* webpack bundle builder
* eslint and golint rules for Makefile

## Dependencies

* [golang](http://golang.org/)
* [node.js](https://nodejs.org/) with [npm](https://www.npmjs.com/), only to build the application bundle at compile time
* [GNU make](https://www.gnu.org/software/make/)
* [fswatch](https://github.com/emcrisostomo/fswatch/)

Note that probably not works at windows.

## Install

Clone the repo:

```
$ git clone git@github.com:olebedev/go-starter-kit.git && cd go-starter-kit
```
Install javscript dependencies:

```
$ npm i
```
Install Golang dependencies:

```
$ export GOPATH=`pwd` # the most important step, ensure that you do it
$ export GOBIN=$GOPATH/bin # optional, redefine, if it already was defined
$ go get app
$ go get github.com/jteeuwen/go-bindata/...
```
You will get this output after `go get app`, at the first time:

```bash
src/app/server/app.go:64: undefined: Asset
src/app/server/app.go:65: undefined: AssetDir
src/app/server/react.go:191: undefined: Asset
```

don't worry about this, see [this](https://github.com/olebedev/go-starter-kit/issues/5#issuecomment-142585756) comment.

Start dev server:

```
$ make serve
```
that's it. Open [http://localhost:5001/](http://localhost:5001/)(if you use default port) at your browser. Now you ready to start coding your awesome project.

## Build

Install dependencies and just type `NODE_ENV=production make build`. This rule is producing webpack build and regular golang build after that. Result you can find at `$GOPATH/bin`.

## TODO

- [x] migrate from react-hot-loader to react-transform-hmr
- [x] update react to 0.14.x
- [x] update react-router to 1.x
- [x] render final HTML markup at Golang side
- [x] migrate from Flummox to Redux
- [ ] migrate from Stylus to PostCSS
- [x] migrate from Gin to Echo
- [ ] improve README and write an article to describe the project
