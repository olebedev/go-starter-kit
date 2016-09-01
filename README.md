# go-starter-kit [![wercker status](https://app.wercker.com/status/cd5a782c425b1feb06844dcc701e528c/s/master "wercker status")](https://app.wercker.com/project/bykey/cd5a782c425b1feb06844dcc701e528c) [![Join the chat at https://gitter.im/olebedev/go-starter-kit](https://img.shields.io/gitter/room/nwjs/nw.js.svg?maxAge=2592000&style=plastic)](https://gitter.im/olebedev/go-starter-kit?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

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
* routing via [react-router](https://github.com/reactjs/react-router)
* ES6 & JSX via [babel-loader](https://github.com/babel/babel-loader) with minimal runtime dependency footprint
* [redux](https://rackt.org/redux/) as state container
* [redux-devtools](https://github.com/gaearon/redux-devtools)
* css styles without global namespace via PostCSS, [css-loader](https://github.com/webpack/css-loader) & css-modules
* separate css file to avoid FOUC
* hot reloading via [react-transform](https://github.com/gaearon/babel-plugin-react-transform) & [HMR](http://webpack.github.io/docs/hot-module-replacement.html)
* webpack bundle builder
* eslint and golint rules for Makefile

## Dependencies

* [golang](https://golang.org/)
* [node.js](https://nodejs.org/) with [npm](https://www.npmjs.com/), only to build the application bundle at compile time
* [GNU make](https://www.gnu.org/software/make/)

Note that probably not works at windows.

## Project structure

##### The server's entry point
```
$ tree server
server
├── api.go
├── app.go
├── bindata.go <-- this file is gitignored, it will appear at compile time
├── conf.go
├── data
│   └── templates
│       └── react.html
├── main.go <-- main function declared here
├── react.go
└── utils.go
```

The `./server/` is flat golang package.

##### The client's entry point

It's simple React application

```
$ tree client
client
├── actions.js
├── components
│   ├── app
│   │   ├── favicon.ico
│   │   ├── index.js
│   │   └── styles.css
│   ├── homepage
│   │   ├── index.js
│   │   └── styles.css
│   ├── not-found
│   │   ├── index.js
│   │   └── styles.css
│   └── usage
│       ├── index.js
│       └── styles.css
├── css
│   ├── funcs.js
│   ├── global.css
│   ├── index.js
│   └── vars.js
├── index.js <-- main function declared here
├── reducers.js
├── router
│   ├── index.js
│   ├── routes.js
│   └── toString.js
└── store.js
```

The client app will be compiled into `server/data/static/build/`.  Then it will be embedded into go package via _go-bindata_. After that the package will be compiled into binary.

**Convention**: javascript app should declare [_main_](https://github.com/olebedev/go-starter-kit/blob/master/client/index.js#L4) function right in the global namespace. It will used to render the app at the server side.

## Install

Clone the repo:

```
$ git clone git@github.com:olebedev/go-starter-kit.git $GOPATH/src/github.com/<username>/<project>
$ cd $GOPATH/src/github.com/<username>/<project>
```
Install JavaScript dependencies:

```
$ npm i
```

Install Golang dependencies via revision locking tool - [srlt](https://github.com/olebedev/srlt). Make sure that you have srlt installed, environment variable `GO15VENDOREXPERIMENT=1` and _Golang_ >= 1.5.0.

```
$ srlt restore
```

This command will install dependencies into `./vendor/` folder located in root.

You can also install all dependencies at once by running:

```
$ make install
```

## Run development

Start dev server:

```
$ make serve
```

that's it. Open [http://localhost:5001/](http://localhost:5001/)(if you use default port) at your browser. Now you ready to start coding your awesome project.

## Build

Install dependencies and type `NODE_ENV=production make build`. This rule is producing webpack build and regular golang build after that. Result you can find at `$GOPATH/bin`. Note that the binary will be named **as the current project directory**.

## License
MIT
