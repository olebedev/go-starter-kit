BIN           = $(GOPATH)/bin
ON            = $(BIN)/on
GO_BINDATA    = $(BIN)/go-bindata
NODE_BIN      = $(shell npm bin)
PID           = .pid
GO_FILES      = $(filter-out ./server/bindata.go, $(shell find ./server  -type f -name "*.go"))
TEMPLATES     = $(wildcard server/data/templates/*.html)
BINDATA       = server/bindata.go
BINDATA_FLAGS = -pkg=main -prefix=server/data
BUNDLE        = server/data/static/build/bundle.js
APP           = $(shell find client -type f)
IMPORT_PATH   = $(shell pwd | sed "s|^$(GOPATH)/src/||g")
APP_NAME      = $(shell pwd | sed 's:.*/::')
TARGET        = $(BIN)/$(APP_NAME)
GIT_HASH      = $(shell git rev-parse HEAD)
LDFLAGS       = -w -X main.commitHash=$(GIT_HASH)
GLIDE         := $(shell command -v glide 2> /dev/null)

build: $(ON) $(GO_BINDATA) clean $(TARGET)

clean:
	@rm -rf server/data/static/build/*
	@rm -rf server/data/bundle.server.js
	@rm -rf $(BINDATA)

$(ON):
	go install $(IMPORT_PATH)/vendor/github.com/olebedev/on

$(GO_BINDATA):
	go install $(IMPORT_PATH)/vendor/github.com/jteeuwen/go-bindata/...

$(BUNDLE): $(APP)
	@$(NODE_BIN)/webpack --progress --colors --bail

$(TARGET): $(BUNDLE) $(BINDATA)
	@go build -ldflags '$(LDFLAGS)' -o $@ $(IMPORT_PATH)/server

kill:
	@kill `cat $(PID)` || true

serve: $(ON) $(GO_BINDATA) clean $(BUNDLE) restart
	@BABEL_ENV=dev node hot.proxy &
	@$(NODE_BIN)/webpack --watch &
	@$(ON) -m 2 $(GO_FILES) $(TEMPLATES) | xargs -n1 -I{} make restart || make kill

restart: BINDATA_FLAGS += -debug
restart: LDFLAGS += -X main.debug=true
restart: $(BINDATA) kill $(TARGET)
	@echo restart the app...
	@$(TARGET) run & echo $$! > $(PID)

$(BINDATA):
	$(GO_BINDATA) $(BINDATA_FLAGS) -o=$@ server/data/...

lint:
	@yarn run eslint || true
	@golint $(GO_FILES) || true

install:
	@yarn install

ifdef GLIDE
	@glide install
else
	$(warning "Skipping installation of Go dependencies: glide is not installed")
endif