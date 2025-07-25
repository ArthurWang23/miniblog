# 将所有的 .mk 文件都在 all.mk 文件中导入，方便管理.
include scripts/make-rules/common.mk
include scripts/make-rules/tools.mk # include at second order
include scripts/make-rules/golang.mk
include scripts/make-rules/generate.mk
include scripts/make-rules/swagger.mk