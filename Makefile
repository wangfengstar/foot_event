GO := go

PREFIX := $(DESTDIR)/usr/local
BINDIR := $(PREFIX)/sbin
MANDIR := $(PREFIX)/share/man

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_BRANCH_CLEAN := $(shell echo $(GIT_BRANCH) | sed -e "s/[^[:alnum:]]/-/g")
RUNC_IMAGE := runc_dev$(if $(GIT_BRANCH_CLEAN),:$(GIT_BRANCH_CLEAN))
PROJECT := foot_event


# TODO: rm -mod=vendor once go 1.13 is unsupported
GO_BUILD := $(GO) build  -buildmode=pie $(EXTRA_FLAGS) -tags "$(BUILDTAGS)" \
	-ldflags "-X main.gitCommit=$(COMMIT) -X main.version=$(VERSION) $(EXTRA_LDFLAGS)"
GO_BUILD_STATIC := CGO_ENABLED=1 $(GO) build  $(EXTRA_FLAGS) -tags "$(BUILDTAGS)  " \
	-ldflags "-w -extldflags -static -X main.gitCommit=$(COMMIT) -X main.version=$(VERSION) $(EXTRA_LDFLAGS)"

.DEFAULT: foot_event

runc:
	$(GO_BUILD) -o foot_event .

all: foot_event


static:
	$(GO_BUILD_STATIC) -o foot_event .

release:

lint:
	$(GO) vet ./...
	$(GO) fmt ./...


clean:
	rm -f runc runc-*
	rm -f contrib/cmd/recvtty/recvtty
	rm -rf release
	rm -rf man/man8
