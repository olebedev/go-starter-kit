BIN           = $(GOPATH)/bin
NODE_BIN      = $(shell npm bin)
PID           = .pid
GO_FILES      = $(filter-out ./server/bindata.go, $(shell find ./server  -type f -name "*.go")) ./main.go
TEMPLATES     = $(wildcard server/data/templates/*.html)
BINDATA       = server/bindata.go
BINDATA_FLAGS = -pkg=main -prefix=server/data
BUNDLE        = server/data/static/build/bundle.js
APP           = $(shell find client -type f)
IMPORT_PATH   = $(shell echo `pwd` | sed "s|^$(GOPATH)/src/||g")
APP_NAME      = $(shell echo $(IMPORT_PATH) | sed 's:.*/::')
TARGET        = $(BIN)/$(APP_NAME)
GIT_HASH      = $(shell git rev-parse HEAD)
LDFLAGS       = -w -X main.commitHash=$(GIT_HASH)

build: clean $(TARGET)

clean:
	@rm -rf server/data/static/build/*
	@rm -rf server/data/bundle.server.js
	@rm -rf $(BINDATA)

$(BUNDLE): $(APP)
	@$(NODE_BIN)/webpack --progress --colors --bail

$(TARGET): $(BUNDLE) $(BINDATA)
	@go build -ldflags '$(LDFLAGS)' -o $@ $(IMPORT_PATH)/server

kill:
	@kill `cat $(PID)` || true

serve: clean $(BUNDLE) restart
	@BABEL_ENV=dev node hot.proxy &
	@$(NODE_BIN)/webpack --watch &
	@fswatch --event Updated $(GO_FILES) $(TEMPLATES) | xargs -n1 -I{} make restart || make kill

restart: BINDATA_FLAGS += -debug
restart: LDFLAGS += -X main.debug=true
restart: $(BINDATA) kill $(TARGET)
	@echo restart the app...
	@$(TARGET) run & echo $$! > $(PID)

$(BINDATA):
	$(BIN)/go-bindata $(BINDATA_FLAGS) -o=$@ server/data/...

lint:
	@eslint client || true
	@golint $(GO_FILES) || true
