SHELL := /bin/bash

generate: protoc
	go generate ./...

protoc: protoc-gen protoc-gen-go protoc-gen-go-grpc

build: generate deps build_client build_server

build_client:
	go build -o bin/client ./cmd/client

build_server:
	go build -o bin/server ./cmd/server

deps:
	go mod vendor
	go mod tidy

protoc-gen:
ifeq (, $(shell which protoc))
	@{ \
	echo protoc must be installed for the code to work ;\
	exit 1 ;\
	}
endif

protoc-gen-go:
ifeq (, $(shell which protoc-gen-go))
	@{ \
	set -e ;\
	PROTOC_GO_TMP_DIR=$$(mktemp -d) ;\
	cd $$PROTOC_GO_TMP_DIR ;\
	go mod init tmp ;\
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26 ;\
	rm -rf $$PROTOC_GO_TMP_DIR ;\
	}
PROTOC_GEN_GO=$(GOBIN)/protoc-gen-go
else
PROTOC_GEN_GO=$(shell which protoc-gen-go)
endif

protoc-gen-go-grpc:
ifeq (, $(shell which protoc-gen-go-grpc))
	@{ \
	set -e ;\
	PROTOC_GO_GRPC_TMP_DIR=$$(mktemp -d) ;\
	cd $$PROTOC_GO_GRPC_TMP_DIR ;\
	go mod init tmp ;\
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1 ;\
	rm -rf $$PROTOC_GO_GRPC_TMP_DIR ;\
	}
PROTOC_GEN_GO_GRPC=$(GOBIN)/protoc-gen-go-grpc
else
PROTOC_GEN_GO_GRPC=$(shell which protoc-gen-go-grpc)
endif

certs:
	@{ \
  	set -e; \
	openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365 -subj '/CN=localhost' -nodes \
	-extensions san -config \
	<(cat /etc/ssl/openssl.cnf \
	<(printf "\n[san]\nsubjectAltName=DNS:localhost,IP:127.0.0.1\n")); \
	}