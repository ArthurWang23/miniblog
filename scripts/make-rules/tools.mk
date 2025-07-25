# ==============================================================================
# 工具相关的 Makefile
#

TOOLS ?= golangci-lint goimports protoc-plugins swagger addlicense protoc-go-inject-tag protolint

tools.verify: $(addprefix tools.verify., $(TOOLS))

tools.install: $(addprefix tools.install., $(TOOLS))

tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.2
	@golangci-lint completion bash > $(HOME)/.golangci-lint.bash
	@if ! grep -q .golangci-lint.bash $(HOME)/.bashrc; then echo "source \$$HOME/.golangci-lint.bash" >> $(HOME)/.bashrc; fi

install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

install.protoc-plugins:
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	@$(GO) install github.com/onexstack/protoc-gen-defaults@v0.0.2
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.24.0
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.24.0

install.swagger:
	@$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest

install.addlicense:
	@$(GO) install github.com/marmotedu/addlicense@latest

install.protoc-go-inject-tag:
	@$(GO) install github.com/favadi/protoc-go-inject-tag@latest

install.protolint:
	@$(GO) install github.com/yoheimuta/protolint/cmd/protolint@latest

# 伪目标（防止文件与目标名称冲突）
.PHONY: tools.verify tools.install tools.install.% tools.verify.% install.golangci-lint \
	install.goimports install.protoc-plugins install.swagger \
	install.addlicense install.protoc-go-inject-tag protolint