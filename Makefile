SOURCES := main.go
BIN := github-webhook

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

$(BIN): $(SOURCES)
	CGO_ENABLED=0 go build .

build: $(BIN)

test: .test-done build

.test-done: $(BIN)
	@echo "Running tests..."
	go test ./...
	touch .test-done

.PHONY : $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
    GOOS=$(os) GOARCH=amd64 go build -o release/$(BIN)-v1.0.0-$(os)-amd64

.PHONY: release
release: windows linux darwin

clean:
	$(RM) $(call FixPath, $(BIN))

docker: $(BIN)
	docker build -t $(BIN):latest .

run: $(BIN) docker
	docker run -it -p 5000:5000 $(BIN):latest

.PHONY: all build test release clean docker run
