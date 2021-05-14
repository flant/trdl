GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

all: fmt build start

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/vault-plugin-secrets-trdl cmd/vault-plugin-secrets-trdl/main.go

start:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins -log-level trace

enable:
	VAULT_ADDR='http://127.0.0.1:8200' vault secrets enable -path=trdl vault-plugin-secrets-trdl

fixtures:
	curl -X POST "http://127.0.0.1:8200/v1/trdl/configure" -H  "accept: */*" -H  "Content-Type: application/json" -H  "X-Vault-Token: root" -d "{\"git_repo_url\":\"https://github.com/werf/trdl-test-project.git\",\"required_number_of_verified_signatures_on_commit\":0}"

clean:
	rm -f ./vault/plugins/vault-plugin-secrets-trdl

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt start enable
