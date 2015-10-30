# go-starter-kit

[![Join the chat at https://gitter.im/olebedev/go-starter-kit](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/olebedev/go-starter-kit?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

This project contains a quick starter kit for **Facebook React** Single Page Apps with **Golang** server side render and also with a set of useful features for rapid development of efficient applications.

## What it contains?

* server side render via [go-duktape](https://github.com/olebedev/go-duktape)
* api requests between your react application and server side application directly  via [fetch polyfill](https://github.com/olebedev/go-duktape-fetch) for go-duktape at server side, and it is possible to process requests with user session as well
* title, Open Graph and other domain-specific meta tags render for each page at the server and at the client
* server side redirect
* embedding static files into artefact via bindata
* high performance [gin](https://github.com/gin-gonic/gin/) framework
* advanced cli via [cli](https://github.com/codegangsta/cli)
* Makefile based project
* one(!) terminal window process for development
* routing via [react-router](https://github.com/rackt/react-router)
* ES6 & JSX via [babel-loader](https://github.com/babel/babel-loader) with minimal runtime dependency footprint
* flux architecture via [Flummox](https://github.com/acdlite/flummox)
* stylus css styles without global namespace via [css-loader](https://github.com/webpack/css-loader) & css-module
* separate css file to avoid FOUC
* hot reload components & styles via [react-hot-loader](https://github.com/gaearon/react-hot-loader) & [HMR](http://webpack.github.io/docs/hot-module-replacement.html)
* webpack bundle builder
* eslint and golint rules for Makefile

## Dependencies

* [golang](http://golang.org/)
* [node.js](https://nodejs.org/) with [npm](https://www.npmjs.com/), only for javascript workflow
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
$ go get app/...
$ go get github.com/jteeuwen/go-bindata/...
```
Start dev server:

```
$ make serve
```
that's it. Open [http://localhost:5001/](http://localhost:5001/) at your browser. Now you ready to start coding your awesome project.

## Build

Install dependencies and just type `NODE_ENV=production make build`. This rule is producing webpack build and regular golang build after that. Result you can find at `$GOPATH/bin`.

## TODO

- [x] add react-transform-hmr instead react-hot-loader
- [ ] update react to 0.14.x
- [ ] update react-router to 1.x
- [ ] add redux and drop flummox
- [ ] add postcss and drop stylus (?)
- [ ] improve README and write an article to describe the project