default: buf-gen

BIN_NAME:=$(notdir $(shell pwd))
# 使用分支名作为version
VERSION := $(shell git branch | grep \* | cut -d ' ' -f2)
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
IMAGE_NAME := "flamefatex/${BIN_NAME}"

.PHONY: buf-update
buf-update:
	@docker run --volume "`pwd`/proto/src:/workspace" --workdir /workspace bufbuild/buf:1.16.0 mod update
.PHONY: buf-gen
buf-gen:
	@docker run --volume "`pwd`/proto:/workspace" --workdir /workspace bufbuild/buf:1.16.0 generate
.PHONY: local-buf-update
local-buf-update:
	@cd proto/src && buf mod update
.PHONY: local-buf-gen
local-buf-gen:
	@cd proto && buf generate
.PHONY: gorm-gen
gorm-gen:
	@rm -rf model/query
	@mkdir -p model/query
	@go run cmd/gorm-gen/gorm-gen.go --path model/query


.PHONY: build
build:
	@echo "building ${BIN_NAME} ${VERSION} ${GIT_COMMIT} ${GIT_DIRTY}"
	@go build -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.VersionPrerelease=DEV" -o bin/${BIN_NAME}
.PHONY: image
image:
	@echo "building image ${BIN_NAME} ${VERSION} ${GIT_COMMIT} ${GIT_DIRTY}"
	@docker build --build-arg APP_NAME=${BIN_NAME} --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=${GIT_COMMIT}${GIT_DIRTY} -t ${IMAGE_NAME}:latest .