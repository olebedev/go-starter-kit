BIN = $(GOPATH)/bin
NODE_BIN = ./node_modules/.bin
PID = .pid
GO_FILES = $(filter-out bindata.go, $(shell find src/app -type f -name "*.go"))
BINDATA = src/app/server/data/bindata.go
BINDATA_FLAGS ?= -debug -pkg=data -prefix=src/app/server/data -ignore=src\\/app\\/server\\/data\\/bindata.go

clean:
	@echo cleaned

kill:
	@kill `cat $(PID)` || true

serve: clean
	@$(NODE_BIN)/webpack --progress --colors
	@make restart
	@$(NODE_BIN)/webpack-dev-server --config webpack.hot.config.js $$! > $(PID)_wds &
	@ANYBAR_WEBPACK=yep $(NODE_BIN)/webpack --progress --colors --watch $$! > $(PID)_wp &
	@fswatch $(GO_FILES) | xargs -n1 -I{} make restart || make kill
	@kill `cat $(PID)_wp` || true
	@kill `cat $(PID)_wds` || true

restart: $(BINDATA)
	@make kill
	@go install app
	@$(BIN)/app run & echo $$! > $(PID)

$(BINDATA):
	$(BIN)/go-bindata $(BINDATA_FLAGS) -o=$@ src/app/server/data/...
