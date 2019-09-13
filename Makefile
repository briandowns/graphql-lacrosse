GO = go

BINDIR                 = bin
BINARY                 = graphql-lacrosse

UNAME_S = $(shell uname -s)

override BUILD_ENV     += CGO_ENABLED=0 GOARCH=amd64
override BUILD_FLAGS   += -ldflags '-w -s -extldflags "-static"'

$(BINDIR)/$(BINARY): $(BINDIR)
ifeq ($(UNAME_S),Linux)
	$(BUILD_ENV) GOOS=linux $(GO) build -v -o $@ $(BUILD_FLAGS) .
endif
ifeq ($(UNAME_S),Darwin)
	$(BUILD_ENV) GOOS=darwin $(GO) build -v -o $@ $(BUILD_FLAGS) .
endif
ifeq ($(UNAME_S),FreeBSD)
	$(BUILD_ENV) GOOS=freebsd $(GO) build -v -o $@ $(BUILD_FLAGS) .
endif

$(BINDIR):
	mkdir -p $@

.PHONY: test
test:
	$(GO) test -v -cover ./...

.PHONY: check
check:
	if [ -d vendor ]; then cp -r vendor/* ${GOPATH}/src/; fi
	GO111MODULE=off gosec -exclude G304 ./...

.PHONY: clean
clean:
	$(GO) clean
	rm -f bin/*

.PHONY: docs
docs:
	@godoc -http=:6060 2>/dev/null &
	@printf "To view common docs, point your browser to:\n"
	@printf "\n\thttp://127.0.0.1:6060/pkg/github.com/briandowns/$(pkg)\n\n"
	@sleep 1
	@open "http://127.0.0.1:6060/pkg/github.com/briandowns/$(pkg)"
