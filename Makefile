GO ?= go
ifeq (, $(shell which govvv))
$(shell go install github.com/ahmetb/govvv@latest)
endif

GIT_INFO := $(shell govvv -flags)
GO_BUILD = $(GO) build -ldflags "$(GIT_INFO)"

build:
	$(GO_BUILD) -o bin/turnip cmd/main.go

clean:
	@rm -rf bin
