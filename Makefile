# ==============================================================================
# 定义全局 Makefile 变量方便后面引用
 
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
PROJ_ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR := $(PROJ_ROOT_DIR)/_output
# Protobuf文件存放路径
APIROOT = $(PROJ_ROOT_DIR)/pkg/api
 
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
all: tidy format lint build add-copyright
 
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

.PHONY: protoc
protoc: # 编译protobuf文件
	@echo "===========> Generate protobuf files"
	@mkdir -p $(PROJ_ROOT_DIR)/api/openapi
	@# --grpc-gateway_out用来在pkg/api/apiserver/v1目录下胜澈功能反向服务器代码apiserver.pb.gw.go
	@# --openapiv2_out用来在api/openapi/apiserver/v1目录下生成Swagger V2接口文档
	@# restfulapi中delete通常用于删除资源，按规范只应携带资源标识符，不建议附带请求体
	@# 然而某些场景（删除选项或多个待删除的资源列表）可设置allow_delete_body=true放宽delete请求限制
	@protoc                                              \
		--proto_path=$(APIROOT)                          \
		--proto_path=$(PROJ_ROOT_DIR)/third_party/protobuf    \
		--go_out=paths=source_relative:$(APIROOT)        \
		--go-grpc_out=paths=source_relative:$(APIROOT)   \
		--grpc-gateway_out=allow_delete_body=true,paths=source_relative:$(APIROOT) \
		--openapiv2_out=$(PROJ_ROOT_DIR)/api/openapi     \
		--openapiv2_opt=allow_delete_body=true,logtostderr=true \
		--defaults_out=paths=source_relative:$(APIROOT) \
		$(shell find $(APIROOT) -name *.proto)
	@find . -name "*.pb.go" -exec protoc-go-inject-tag -input={} \;

.PHONY: ca
ca: # 生成CA文件
	@mkdir -p $(OUTPUT_DIR)/cert
	@# 1.生成根证书私钥 (CA Key)
	@openssl genrsa -out $(OUTPUT_DIR)/cert/ca.key 4096
	@# 2.使用根私钥生成证书签名请求 (CA CSR)
	@openssl req -new -nodes -key $(OUTPUT_DIR)/cert/ca.key -sha256 -out $(OUTPUT_DIR)/cert/ca.csr \
		-subj "/C=CN/ST=Shanghai/L=Shanghai/O=miniblog/OU=it/CN=127.0.0.1/emailAddress=2826979176@qq.com"
	@# 3.使用根私钥签发根证书 (CA CRT) 使其自签名 -req指定输入文件为证书请求 -signkey指定用于自签名的私钥文件
	@openssl x509 -req -days 365 -in $(OUTPUT_DIR)/cert/ca.csr -signkey $(OUTPUT_DIR)/cert/ca.key -out $(OUTPUT_DIR)/cert/ca.crt
	@# 4.生成服务端私钥
	@openssl genrsa -out $(OUTPUT_DIR)/cert/server.key 2048
	@# 5.使用服务端私钥生成服务端证书签名系统
	@openssl req -new -key $(OUTPUT_DIR)/cert/server.key -out $(OUTPUT_DIR)/cert/server.csr \
		-subj "/C=CN/ST=Shanghai/L=Shanghai/O=serverdevops/OU=serverit/CN=localhost/emailAddress=2826979176@qq.com" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
	@# 6.使用根证书(CA)签发服务端证书
	@echo "[v3_req]" > $(OUTPUT_DIR)/cert/v3.ext
	@echo "subjectAltName=DNS:localhost,IP:127.0.0.1" >> $(OUTPUT_DIR)/cert/v3.ext
	@openssl x509 -days 356 -sha256 -req -CA $(OUTPUT_DIR)/cert/ca.crt -CAkey $(OUTPUT_DIR)/cert/ca.key \
		-CAcreateserial -in $(OUTPUT_DIR)/cert/server.csr -out $(OUTPUT_DIR)/cert/server.crt -extensions v3_req \
		-extfile $(OUTPUT_DIR)/cert/v3.ext
	@rm -f $(OUTPUT_DIR)/cert/v3.ext


.PHONY: test
test: # 执行单元测试.
	@echo "===========> Running unit tests"
	@mkdir -p $(OUTPUT_DIR)
	@go test -race -cover \
		-coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short \
		-v `go list ./...|egrep -v 'tools|vendor|third_party'`

.PHONY: cover
cover: test ## 执行单元测试，并校验覆盖率阈值.
	@echo "===========> Running code coverage tests"
	@go tool cover -func=$(OUTPUT_DIR)/coverage.out | awk -v target=$(COVERAGE) -f $(PROJ_ROOT_DIR)/scripts/coverage.awk


.PHONY:lint
lint: # 执行lint检查.
	@echo "===========> Running golangci to lint source codes"
	@golangci-lint run -c $(PROJ_ROOT_DIR)/.golangci.yaml $(PROJ_ROOT_DIR)/...