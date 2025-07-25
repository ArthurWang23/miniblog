# 重构makefile结构

# 定义默认目标为 all
.DEFAULT_GOAL := all

# 定义 Makefile all 伪目标，执行 `make` 时，会默认会执行 all 伪目标
all: tidy format build cover add-copyright

# ==============================================================================
# Includes

# 确保 `include common.mk` 位于第一行，common.mk 中定义了一些变量，后面的子 makefile 有依赖
include scripts/make-rules/common.mk
include scripts/make-rules/all.mk

# ==============================================================================

# Usage

define USAGE_OPTIONS

选项:
  BINS             要构建的二进制文件。默认值为cmd中的所有文件。
                   此选项可用于以下命令：make build
                   示例：make build BINS="miniblog"
  VERSION          编译到二进制文件中的版本信息。
  V                设置为1以启用详细的构建信息输出。默认值为0。
endef
export USAGE_OPTIONS

## --------------------------------------
## Binaries
## --------------------------------------

##@ build:

build: go.tidy  ## 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	@$(MAKE) go.build

## --------------------------------------
## Testing
## --------------------------------------

##@ test:

test: ## 执行单元测试.
	@$(MAKE) go.test

cover: ## 执行单元测试，并校验覆盖率阈值.
	@$(MAKE) go.cover

## --------------------------------------
## Cleanup
## --------------------------------------

##@ clean:

clean: ## 清理构建产物、临时文件等. 例如 _output 目录.
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## --------------------------------------
## Lint / Verification
## --------------------------------------

##@ lint and verify:

lint: ## 执行静态代码检查.
	@$(MAKE) go.lint

tidy: ## 自动添加/移除依赖包.
	@$(MAKE) go.tidy

format: tools.verify.protolint ## 格式化 Go 源码及 protobuf 文件.
	@$(MAKE) go.format
	@protolint -fix -config_path ${PROJ_ROOT_DIR}/.protolint.yaml $(shell find $(APIROOT) -name *.proto)

## --------------------------------------
## Generate / Manifests
## --------------------------------------

##@ generate

ca: ## 生成 CA 文件.
	@$(MAKE) gen.ca

protoc: ## 编译 protobuf 文件.
	@$(MAKE) gen.protoc

## --------------------------------------
## Hack / Tools
## --------------------------------------

##@ hack/tools:

swagger: ## 聚合 swagger 文档到一个 openapi.yaml 文件中.
	@$(MAKE) swagger.run

serve-swagger: ## 运行 Swagger 文档服务器.
	@$(MAKE) swagger.serve

add-copyright: ## 添加版权头信息.
	@$(MAKE) gen.add-copyright

help: Makefile ## 打印 Makefile help 信息.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<TARGETS> <OPTIONS>\033[0m\n\n\033[35mTargets:\033[0m\n"} /^[0-9A-Za-z._-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-45s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' Makefile #$(MAKEFILE_LIST)
	@echo -e "$$USAGE_OPTIONS"

# 伪目标（防止文件与目标名称冲突）
.PHONY: all build test cover clean lint tidy format ca protoc swagger serve-swagger add-copyright help