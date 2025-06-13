package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/onexstack/onexstack/pkg/db"
	"github.com/spf13/pflag"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// GORM Model自动生成工具
const helpText = `Usage: main [flags] arg [arg...]

This is a pflag example.

Flags:
`

// 定义数据库查询接口
type Querier interface {
	// FilterWithNameAndRole 按名称和角色查询记录
	FilterWithNameAndRole(name string) ([]gen.T, error)
}

// GenerateConfig 保存代码生成的配置
type GenerateConfig struct {
	ModelPackagePath string
	GenerateFunc     func(g *gen.Generator)
}

var generateConfigs = map[string]GenerateConfig{
	"mb": {ModelPackagePath: "../../internal/apiserver/model", GenerateFunc: GenerateMiniBlogModels},
}

// 命令行参数
var (
	addr       = pflag.StringP("addr", "a", "127.0.0.1:3306", "MySQL host address.")
	username   = pflag.StringP("username", "u", "miniblog", "Username to connect to the database.")
	password   = pflag.StringP("password", "p", "miniblog1234", "Password to use when connecting to the database.")
	database   = pflag.StringP("db", "d", "miniblog", "Database name to connect to.")
	modelPath  = pflag.String("model-pkg-path", "", "Generated model code's package name.")
	components = pflag.StringSlice("components", []string{"mb"}, "Generated model code's for specified component.")
	help       = pflag.BoolP("help", "h", false, "Show this help message.")
)

func main() {
	// 自定义使用说明
	pflag.Usage = func() {
		fmt.Printf("%s", helpText)
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if *help {
		pflag.Usage()
		return
	}
	// 初始化数据库连接
	dbInstance, err := initializeDatabase()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	for _, component := range *components {
		processComponent(component, dbInstance)
	}
}

func initializeDatabase() (*gorm.DB, error) {
	dbOptions := &db.MySQLOptions{
		Addr:     *addr,
		Username: *username,
		Password: *password,
		Database: *database,
	}
	return db.NewMySQL(dbOptions)
}

// 处理单个组件以生成代码
func processComponent(component string, dbInstance *gorm.DB) {
	config, ok := generateConfigs[component]
	if !ok {
		log.Printf("Component %s not found in configuration.Skipping.", component)
		return
	}
	modelPkgPath := resolveModelPackagePath(config.ModelPackagePath)

	// 创建生成器实例
	generator := createGenerator(modelPkgPath)
	generator.UseDB(dbInstance)

	// 应用自定义生成器选项
	applyGeneratorOptions(generator)

	// 使用指定函数生成模型
	config.GenerateFunc(generator)

	generator.Execute()
}

func resolveModelPackagePath(defaultPath string) string {
	if *modelPath != "" {
		return *modelPath
	}
	absPath, err := filepath.Abs(defaultPath)
	if err != nil {
		log.Printf("Error resolving path: %v", err)
		return defaultPath
	}
	return absPath
}

// 初始化并返回一个新的生成器实例
func createGenerator(packagePath string) *gen.Generator {
	return gen.NewGenerator(gen.Config{
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		ModelPkgPath:  packagePath,
		WithUnitTest:  true,
		FieldNullable: true,  // 对于数据库中可空的字段，使用指针类型
		FieldSignable: false, // 禁用无符号属性以提高兼容性
		// 不包含GORM的索引标签
		// Miniblog支持开启sqlite数据库的内存模式
		// 通过MariaDB表结构生成的GROM Model结构体的列标签并不能用来自动创建SQLite表结构
		// 因此关闭列标签的生成
		FieldWithIndexTag: false,
		FieldWithTypeTag:  false, // 不包含GORM的类型标签
	})
}

// 设置自定义生成器选项
// 默认生成的时间为default:current_timestamp() 该标签不能被SQLite识别
// 需要使用FieldGROMTag方法重命名为default:current_timestamp
func applyGeneratorOptions(g *gen.Generator) {
	g.WithOpts(
		gen.FieldGORMTag("createdAt", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
		gen.FieldGORMTag("updatedAt", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
	)
}

// 为miniblog组件生成模型
func GenerateMiniBlogModels(g *gen.Generator) {
	g.GenerateModelAs(
		"user",
		"UserM",
		gen.FieldIgnore("placeholder"),
		// 同理 uniqueIndex也需要重新指定
		gen.FieldGORMTag("username", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_username")
			return tag
		}),
		gen.FieldGORMTag("userID", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_userID")
			return tag
		}),
		gen.FieldGORMTag("phone", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_phone")
			return tag
		}),
	)
	g.GenerateModelAs(
		"post",
		"PostM",
		gen.FieldIgnore("placeholder"),
		gen.FieldGORMTag("postID", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_post_postID")
			return tag
		}),
	)
	g.GenerateModelAs(
		"casbin_rule",
		"CasbinRuleM",
		// 将casbin_rule表ptype字段重命名，符合Go开发规范
		gen.FieldRename("ptype", "PType"),
		// 忽略表中placeholder字段
		gen.FieldIgnore("placeholder"),
	)
}
