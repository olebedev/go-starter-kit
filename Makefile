BIN = $(GOPATH)/bin
NODE_BIN = $(shell npm bin)
PID = .pid
GO_FILES = $(filter-out src/app/server/bindata.go, $(shell find src/app -type f -name "*.go"))
BINDATA = src/app/server/bindata.go
BINDATA_FLAGS = -pkg=server -prefix=src/app/server/data
BUNDLE = src/app/server/data/static/build/bundle.js
APP = $(shell find src/app/client -type f)

build: clean $(BIN)/app

clean:
	@rm -rf src/app/server/data/static/build/*
	@rm -rf src/app/server/data/bundle.server.js
	@rm -rf $(BINDATA)
	@echo cleaned

$(BUNDLE): $(APP)
	@$(NODE_BIN)/webpack --progress --colors

$(BIN)/app: $(BUNDLE) $(BINDATA)
	@go install -ldflags "-w -X main.buildstamp `date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.gittag `git describe --tags || true` -X main.githash `git rev-parse HEAD || true`" app

kill:
	@kill `cat $(PID)` || true

serve: clean $(BUNDLE)
	@make restart
	@$(NODE_BIN)/webpack-dev-server --config webpack.hot.config.js $$! > $(PID)_wds &
	@ANYBAR_WEBPACK=yep  NODE_ENV=disable-hmr-plugin-at-server $(NODE_BIN)/webpack --progress --colors --watch $$! > $(PID)_wp &
	@fswatch $(GO_FILES) | xargs -n1 -I{} make restart || make kill
	@kill `cat $(PID)_wp` || true
	@kill `cat $(PID)_wds` || true

restart: BINDATA_FLAGS += -debug
restart: $(BINDATA)
	@make kill
	@go install app
	@$(BIN)/app run & echo $$! > $(PID)

$(BINDATA):
	$(BIN)/go-bindata $(BINDATA_FLAGS) -o=$@ src/app/server/data/...

lint:
	@eslint src/app/client || true
	@golint $(filter-out src/app/main.go, $(GO_FILES)) || true
	@golint -min_confidence=1 app
