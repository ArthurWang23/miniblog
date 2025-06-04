# ==============================================================================
# 定义全局 Makefile 变量方便后面引用
 
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
PROJ_ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR := $(PROJ_ROOT_DIR)/_output
 
# ==============================================================================
# 定义版本相关变量 编译时自动注入版本信息


## 指定应用使用的version包，通过'-ldflags -X'向包中指定的变量注入值
VERSION_PACKAGE := github.com/ArthurWang23/miniblog/pkg/version
## 定义VERSION语义化版本号 --tags使用所有标签   --always如果仓库没有可用的标签，使用提交ID缩写替代
ifeq ($(origin VERSION),undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

## 检查代码仓库是否是dirty(默认dirty)
GIT_TREE_STATE := "dirty"
ifeq (,$(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE := "clean"
endif
## 获取构建时提交ID
GIT_COMMIT := $(shell git rev-parse HEAD)
## 通过-ldflags参数向version包注入了gitVersion,gitCommit,gitTreeState,buildDate，Info中另外三个信息通过runtime动态获取
GO_LDFLAGS += \
	-X $(VERSION_PACKAGE).gitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).gitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).gitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \


# ==============================================================================
# 定义默认目标为 all
.DEFAULT_GOAL := all
 
# 定义 Makefile all 伪目标，执行 `make` 时，会默认会执行 all 伪目标
.PHONY: all
all: tidy format build add-copyright
 
# ==============================================================================
# 定义其他需要的伪目标
 
.PHONY: build
build: tidy # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	@go build -v -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/mb-apiserver $(PROJ_ROOT_DIR)/cmd/mb-apiserver/main.go
 
.PHONY: format
format: # 格式化 Go 源码.
	@gofmt -s -w ./
 
.PHONY: add-copyright
add-copyright: # 添加版权头信息.
	@addlicense -v -f $(PROJ_ROOT_DIR)/scripts/boilerplate.txt $(PROJ_ROOT_DIR) --skip-dirs=third_party,vendor,$(OUTPUT_DIR)
 
.PHONY: tidy
tidy: # 自动添加/移除依赖包.
	@go mod tidy
 
.PHONY: clean
clean: # 清理构建产物、临时文件等.
	@-rm -vrf $(OUTPUT_DIR)