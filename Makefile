BIN = $(GOPATH)/bin
NODE_BIN = ./node_modules/.bin
PID = .pid
GO_FILES = $(filter-out bindata.go, $(shell find src/app -type f -name "*.go"))
BINDATA_FLAGS ?= -debug -pkg=static -prefix=src/app/server/static -ignore=src\\/app\\/server\\/static\\/bindata.go

clean:
	@echo cleaned

kill:
	@kill `cat $(PID)` || true

serve: clean
	# $(NODE_BIN)/webpack --progress --colors
	@make restart
	# @$(NODE_BIN)/webpack-dev-server --config webpack.hot.config.js $$! > $(WDSPID) &
	# @ANYBAR_WEBPACK=yep $(NODE_BIN)/webpack --progress --colors --watch $$! > $watch(WPPID) &
	@fswatch $(GO_FILES) | xargs -n1 -I{} make restart || make kill
	# @kill `cat $(WPPID)` || true
	# @kill `cat $(WDSPID)` || typerue

restart: src/app/server/static/bindata.go
	@make kill
	go install app
	$(BIN)/app & echo $$! > $(PID)

src/app/server/static/bindata.go:
	$(BIN)/go-bindata $(BINDATA_FLAGS) -o=$@ src/app/server/static/...
