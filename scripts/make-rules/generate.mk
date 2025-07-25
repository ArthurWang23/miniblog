# ==============================================================================
# 用来进行代码生成的 Makefile
#

gen.add-copyright: tools.verify.addlicense ## 添加版权头信息.
	@addlicense -v -f $(PROJ_ROOT_DIR)/scripts/boilerplate.txt $(PROJ_ROOT_DIR) --skip-dirs=third_party,vendor,$(OUTPUT_DIR)

gen.ca: ## 生成 CA 文件.
	@mkdir -p $(OUTPUT_DIR)/cert
	@# 1. 生成根证书私钥 (CA Key)
	@openssl genrsa -out $(OUTPUT_DIR)/cert/ca.key 4096
	@# 2. 使用根私钥生成证书签名请求 (CA CSR)，有效期为 10 年
	@openssl req -new -nodes -key $(OUTPUT_DIR)/cert/ca.key -sha256 -days 3650 -out $(OUTPUT_DIR)/cert/ca.csr \
		-subj "/C=CN/ST=Guangdong/L=Shenzhen/O=onexstack/OU=it/CN=127.0.0.1/emailAddress=colin404@foxmail.com"
	@# 3. 使用根私钥签发根证书 (CA CRT)，使其自签名
	@openssl x509 -req -days 365 -in $(OUTPUT_DIR)/cert/ca.csr -signkey $(OUTPUT_DIR)/cert/ca.key -out $(OUTPUT_DIR)/cert/ca.crt
	@# 4. 生成服务端私钥
	@openssl genrsa -out $(OUTPUT_DIR)/cert/server.key 2048
	@# 5. 使用服务端私钥生成服务端的证书签名请求 (Server CSR)
	@openssl req -new -key $(OUTPUT_DIR)/cert/server.key -out $(OUTPUT_DIR)/cert/server.csr \
		-subj "/C=CN/ST=Guangdong/L=Shenzhen/O=serverdevops/OU=serverit/CN=localhost/emailAddress=colin404@foxmail.com" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
	@# 6. 使用根证书 (CA) 签发服务端证书 (Server CRT)
	@openssl x509 -days 365 -sha256 -req -CA $(OUTPUT_DIR)/cert/ca.crt -CAkey $(OUTPUT_DIR)/cert/ca.key \
		-CAcreateserial -in $(OUTPUT_DIR)/cert/server.csr -out $(OUTPUT_DIR)/cert/server.crt -extensions v3_req \
		-extfile <(printf "[v3_req]\nsubjectAltName=DNS:localhost,IP:127.0.0.1")

gen.protoc: tools.verify.protoc-gen-go ## 编译 protobuf 文件.
	@echo "===========> Generate protobuf files"
	@mkdir -p $(PROJ_ROOT_DIR)/api/openapi
	@# --grpc-gateway_out 用来在 pkg/api/apiserver/v1/ 目录下生成反向服务器代码 apiserver.pb.gw.go
	@# --openapiv2_out 用来在 api/openapi/apiserver/v1/ 目录下生成 Swagger V2 接口文档
	@protoc                                            \
		--proto_path=$(APIROOT)                          \
		--proto_path=$(PROJ_ROOT_DIR)/third_party/protobuf             \
		--go_out=paths=source_relative:$(APIROOT)        \
		--go-grpc_out=paths=source_relative:$(APIROOT)   \
		--grpc-gateway_out=allow_delete_body=true,paths=source_relative:$(APIROOT) \
		--openapiv2_out=$(PROJ_ROOT_DIR)/api/openapi \
		--openapiv2_opt=allow_delete_body=true,logtostderr=true \
		--defaults_out=paths=source_relative:$(APIROOT) \
		$(shell find $(APIROOT) -name *.proto)
	@find . -name "*.pb.go" -exec protoc-go-inject-tag -input={} \;

gen.generate:
	@GOWORK=off go generate ./...

# 伪目标（防止文件与目标名称冲突）
.PHONY: gen.add-copyright gen.ca gen.protoc gen.generate