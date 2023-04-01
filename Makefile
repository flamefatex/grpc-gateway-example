.PHONY: buf-gen buf-update package clean

default: buf-gen


buf-update:
	@cd proto/src && buf mod update
buf-gen:
	@cd proto && buf generate


