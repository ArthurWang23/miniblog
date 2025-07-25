# ==============================================================================
# Makefile helper functions for swagger
#

swagger.run: tools.verify.swagger
	@echo "===========> Generating swagger API docs"
	@swagger mixin `find $(PROJ_ROOT_DIR)/api/openapi -name "*.swagger.json"` \
		-q                                                    \
		--keep-spec-order                                     \
		--format=yaml                                         \
		--ignore-conflicts                                    \
		-o $(PROJ_ROOT_DIR)/api/openapi/apiserver/v1/openapi.yaml
	@echo "Generated at: $(PROJ_ROOT_DIR)/api/openapi/apiserver/v1/openapi.yaml"

swagger.serve: tools.verify.swagger
	@swagger serve -F=redoc --no-open --port 65534 $(PROJ_ROOT_DIR)/api/openapi/apiserver/v1/openapi.yaml

# 伪目标（防止文件与目标名称冲突）
.PHONY: swagger.run swagger.serve