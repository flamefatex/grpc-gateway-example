.PHONY: buf-update buf-gen local-buf-update local-buf-gen package clean

default: buf-gen


buf-update:
	@@docker run --volume "`pwd`/proto/src:/workspace" --workdir /workspace bufbuild/buf:1.16.0 mod update
buf-gen:
	@docker run --volume "`pwd`/proto:/workspace" --workdir /workspace bufbuild/buf:1.16.0 generate

local-buf-update:
	@cd proto/src && buf mod update
local-buf-gen:
	@cd proto && buf generate
