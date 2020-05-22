GO := go

PROJECT := foot_event/cmd

# TODO: rm -mod=vendor once go 1.13 is unsupported
GO_BUILD := $(GO) build  -buildmode=pie $(EXTRA_FLAGS)

.DEFAULT: foot_event/cmd

foot_event/cmd:
	$(GO_BUILD) -o foot_event/cmd .

all: foot_event/cmd

release:

lint:
	$(GO) vet ./...
	$(GO) fmt ./...


clean:
	rm -f runc runc-*
	rm -f contrib/cmd/recvtty/recvtty
	rm -rf release
	rm -rf man/man8
