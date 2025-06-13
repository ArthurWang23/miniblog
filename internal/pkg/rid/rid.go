package rid

import "github.com/onexstack/onexstack/pkg/id"

// 生成唯一标识符
// 需要为每一条REST资源生成唯一标识符uid（Unique Identifier UID）
// 例如更新资源、删除资源时，需要提供唯一ID
// 在生成uid时，需要注意
// 冲突问题：需要确保生成的标识符在业务范围内不冲突，在分布式环境中生成标识符时，需要特别注意跨节点冲突问题
// 大量请求同时生成 ID 可能会成为瓶颈。需要选择性能较高的算法
// 不应透露敏感信息
// 如果使用时间戳等生成 ID，需结合业务需求分析唯一性范围
// 应具备足够的扩展能力，例如未来新增资源或迁移到分布式环境时，仍然能生成唯一标识

// uid生成方法
// 使用数据库主键（Primary Key）：会暴露系统的数据规模，并且数据库 ID 是可预测的，攻击者可以轻松的基于当前的 ID，模拟一个存在的 ID，并尝试访问

// 36 位 UUID：能够在时间和空间上保证唯一性

// 雪花算法（Snowflake Algorithm）：分布式 ID 生成算法，通过时间戳、机器编号和自增序列号组合生成唯一标识符。生成的 ID 通常是 64 位整数，按时间递增排序
// 适合高性能、高并发的分布式场景。其优点是基于时间排序，生成的 ID 一定时间内是有序的。缺点是现相对复杂，可能受机器时钟漂移影响，需要引入额外的依赖（如机器编号）

// 数据库自增 ID 配合随机化：使用数据库的自增 ID（如递增主键）加上一个随机后缀（或编码）生成唯一标识符。
// 这是一种简单且实用的方法。在绝大部分项目中，全局使用一个数据库，数据库不存在主键冲突问题。
//其优点是可以生成短小的唯一 ID，例如：user-uvalgf。缺点是并非完全去中心化，需要依赖数据库生成初始 ID

// 基于时间戳的自定义生成：按时间生成标识符，用当前时间戳加一定后缀或随机值的方式实现。
// 该方法优点是简单直观，适合对时间敏感的业务。缺点是 ID 可预测，有可能出现冲突，需结合随机数后缀或机器编号

// 分布式 ID：在分布式系统或微服务中，数据库的自增主键仅在单个服务或数据库表中唯一。但多个服务或数据源需要统一标识一个记录时，自增主键可能会冲突。唯一标识符可以跨系统、跨服务维度保证唯一性

// 自定义的唯一标识符可以包含更丰富的信息，例如时间戳、数据中心编号、业务类型等
// 唯一标识符的生成一般不依赖于数据库，可以在程序中独立生成。这种去中心化的设计减少了对数据库的依赖，提升了系统的扩展能力，特别是在分布式环境中
// miniblog 通过rid自定义唯一标识符

const defaultABC = "abcedfghijklmnopqrstuvwxyz1234567890"

type ResourceID string

const (
	// 定义用户资源标识符
	UserID ResourceID = "user"
	// 定义帖子资源标识符
	PostID ResourceID = "post"
)

// 将资源标识符转换为字符串
func (rid ResourceID) String() string {
	return string(rid)
}

// 创建带前缀的唯一标识符
func (rid ResourceID) New(counter uint64) string {
	// 使用自定义选项生成唯一标识符
	uniqueStr := id.NewCode(
		counter,
		id.WithCodeChars([]rune(defaultABC)),
		id.WithCodeL(6),
		id.WithCodeSalt(Salt()),
	)
	return rid.String() + "-" + uniqueStr
}

// 分布式友好，不同机器id不同节点生成ID不冲突
// 避免数据库依赖
