P 			= $(shell pwd | sed 's:.*/::')
SOURCE  = $(wildcard *.go)
TAG     = $(shell git describe --tags)
GOBUILD = go build -ldflags '-w'

# $(tag) here will contain either `-1.0-` or just `-`
ALL = \
	$(foreach arch,64,\
    $(foreach tag,-$(TAG)-,\
	$(foreach suffix,win.exe linux osx,\
		build/$P$(tag)$(arch)-$(suffix))))

build: $(ALL)

# os is determined as thus: if variable of suffix exists, it's taken, if not, then
# suffix itself is taken
win.exe = windows
osx = darwin
build/$P-$(TAG)-64-%: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(firstword $($*) $*) GOARCH=amd64 $(GOBUILD) -o $@
	@cd $(@D) && tar cvzf $(@F).tar.gz $(@F)

build/$P-$(TAG)-32-%: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(firstword $($*) $*) GOARCH=386 $(GOBUILD) -o $@
	@cd $(@D) && tar cvzf $(@F).tar.gz $(@F)

build/$P-%: build/$P-$(TAG)-%
	@mkdir -p $(@D)
	cd $(@D) && ln -sf $(<F) $(@F)
