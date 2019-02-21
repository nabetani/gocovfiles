NAME     := gocovfiles
VERSION  := v0.0.1
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell  find . -type f -name '*.go' | grep -v samplesrc)

LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: test
test:
	go test -cover -v

.PHONY: samplerun
samplerun:
	go build $(LDFLAGS) -o bin/$(NAME)
	go test ./samplesrc/* -cover -coverprofile=cover.out
	bin/$(NAME)
 