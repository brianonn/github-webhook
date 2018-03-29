PKGS := $(shell go list ./ ... | grep -v /vendor)
SOURCES := main.go
DEPSDIR := ./Godeps
DEPSFILE= $(DEPSDIR)/Godeps.json
VENDIR  := ./vendor
BIN := github-webhook

GOPATHBIN := $(GOPATH)/bin
GODEP := $(GOPATHBIN)/godep

PLATFORMS := windows linux darwin
os = $(word 1, $@)

ifeq ($(OS), WINDOWS_NT)
	RM = del /S /Q
	FixPath = $(subst /,\,$1)
else
    RM = rm -fr
    FixPath = $1
endif

all: build test
	@echo "Done."

$(BIN): $(SOURCES) $(DEPSFILE)
	$(GODEP) go build

$(GODEP):
	go get -u github.com/tools/godep

install: all
	$(GODEP) go install

$(DEPSFILE): $(SOURCES)
	$(GODEP) save

build: $(BIN)

test: .test-done build

.test-done: $(BIN)
	@echo "Running tests..."
	$(GODEP) go test
	touch .test-done

.PHONY : $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
    GOOS=$(os) GOARCH=amd64 go build -o release/$(BIN)-v1.0.0-$(os)-amd64

.PHONY: release
release: windows linux darwin

clean:
	$(RM) $(call FixPath, $(BIN))

realclean: clean
	$(RM) $(call FixPath, $(DEPSDIR))
	$(RM) $(call FixPath, $(VENDIR))

docker: $(BIN)
	docker build -t $(BIN):latest .

run: $(BIN) docker
	docker run -it -p 5000:5000 $(BIN):latest

.PHONY: all clean realclean prepare install deps build test docker run
