GO           ?= go
GOCACHE      ?= $(CURDIR)/.gocache
GOMODCACHE   ?= $(CURDIR)/.gomodcache

export GOCACHE
export GOMODCACHE

GO_PACKAGES := $(shell env GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GO) list ./... 2>/dev/null)

.PHONY: build run test race integration mocks tidy clean

build:
	@mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GO) build ./...

run:
	@mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GO) run ./main.go

test:
	@mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GO) test -count=1 ./...

race:
	@mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GO) test -race -count=1 ./...

integration:
	@mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GO) test -count=1 ./api/tests

mocks:
	@mkdir -p pkg/mocks
	$(GO) run github.com/golang/mock/mockgen@v1.6.0 \
		-destination=pkg/mocks/devices_mock.go \
		-package=mocks \
		-mock_names "Repository=MockRepository,KeyStore=MockKeyStore,KeyGenerator=MockKeyGenerator,SignerFactory=MockSignerFactory,Signer=MockSigner,SignatureStore=MockSignatureStore" \
		github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices \
		Repository,KeyStore,KeyGenerator,SignerFactory,Signer,SignatureStore
	$(GO) run github.com/golang/mock/mockgen@v1.6.0 \
		-destination=pkg/mocks/api_devices_service_mock.go \
		-package=mocks \
		-mock_names "Service=MockDevicesService" \
		github.com/fiskaly/coding-challenges/signing-service-challenge/api/v0/devices \
		Service

tidy:
	@mkdir -p $(GOCACHE) $(GOMODCACHE)
	$(GO) mod tidy

clean:
	rm -rf $(GOCACHE) $(GOMODCACHE)
	rm -rf pkg/mocks/*_mock.go
