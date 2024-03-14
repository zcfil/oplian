ldflags=-X=build.CurrentCommit=zc3.14
GOFLAGS+=-ldflags="$(ldflags)"
GOCC?=go
.PHONY: all
all: oplian oplian-gateway oplian-op oplian-op-c2

oplian:
	rm -f oplian
	$(GOCC) build $(GOFLAGS) ./cmd/oplian/
.PHONY: oplian

oplian-gateway:
	rm -f oplian-gateway
	$(GOCC) build $(GOFLAGS) ./cmd/oplian-gateway/
.PHONY: oplian-gateway

oplian-op:
	rm -f oplian-op
	$(GOCC) build $(GOFLAGS) ./cmd/oplian-op/
.PHONY: oplian-op

oplian-op-c2:
	rm -f oplian-op-c2
	$(GOCC) build $(GOFLAGS) ./cmd/oplian-op-c2/
.PHONY: oplian-op-c2

.PHONY: clean
clean:
	rm -f oplian oplian-gateway oplian-op oplian-op-c2
