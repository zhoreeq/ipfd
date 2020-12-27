GOARCH := $(GOARCH)
GOOS := $(GOOS)
BINARY_NAME := ipfd
RELEASE_DIR := ipfd_release
PACKAGE_PATH := github.com/zhoreeq/ipfd/internal/app/ipfd
GITBRANCH := $$(git symbolic-ref -q HEAD --short)
GITVERSION := $$(git rev-parse --short HEAD)
FLAGS := -ldflags "-s -w -X $(PACKAGE_PATH).Version=$(GITBRANCH)-$(GITVERSION)"

all:
	GOARCH=$$GOARCH GOOS=$$GOOS go build $(FLAGS) ./cmd/$(BINARY_NAME)

clean:
	$(RM) $(BINARY_NAME) 
	$(RM) -rf $(RELEASE_DIR) $(RELEASE_DIR).tgz

release:
	$(RM) -rf $(RELEASE_DIR)
	mkdir $(RELEASE_DIR)
	cp $(BINARY_NAME) $(RELEASE_DIR)
	cp -r static $(RELEASE_DIR)
	cp -r templates $(RELEASE_DIR)
	cp -r migrations $(RELEASE_DIR)
	cp -r config.example $(RELEASE_DIR)
	tar czf $(RELEASE_DIR).tgz $(RELEASE_DIR)

test:
	go test -v ./...

.PHONY: all clean release test
