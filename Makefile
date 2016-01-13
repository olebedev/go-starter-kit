BIN = $(GOPATH)/bin
NODE_BIN = $(shell npm bin)
PID = .pid
GO_FILES = $(filter-out ./server/bindata.go, $(shell find ./server  -type f -name "*.go")) ./main.go
TEMPLATES = $(wildcard server/data/templates/*.html)
BINDATA = server/bindata.go
BINDATA_FLAGS = -pkg=server -prefix=server/data
BUNDLE = server/data/static/build/bundle.js
APP = $(shell find client -type f)
LDFLAGS = "-w -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.gittag=`git describe --tags || true` -X main.githash=`git rev-parse HEAD || true`" 
TARGET = $(BIN)/app

build: clean $(TARGET)

clean:
	@rm -rf server/data/static/build/*
	@rm -rf server/data/bundle.server.js
	@rm -rf $(BINDATA)
	@echo cleaned

$(BUNDLE): $(APP)
	@$(NODE_BIN)/webpack --progress --colors

$(TARGET): $(BUNDLE) $(BINDATA)
	@go build -ldflags $(LDFLAGS) -o $@

kill:
	@kill `cat $(PID)` || true

serve: clean $(BUNDLE)
	@make restart
	@BABEL_ENV=dev node hot.proxy &
	@$(NODE_BIN)/webpack --watch &
	@fswatch $(GO_FILES) $(TEMPLATES) | xargs -n1 -I{} make restart || make kill

restart: BINDATA_FLAGS += -debug
restart: $(BINDATA)
	@make kill
	@echo restart the app...
	@go build -o $(TARGET)
	@$(TARGET) run & echo $$! > $(PID)

$(BINDATA):
	$(BIN)/go-bindata $(BINDATA_FLAGS) -o=$@ server/data/...

lint:
	@eslint client || true
	@golint $(filter-out ./main.go, $(GO_FILES)) || true
