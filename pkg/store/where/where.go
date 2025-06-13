package where

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 通过where包定制化查询条件
var (
	defaultLimit = -1
)

// Tenant represents a tenant with a key and a function to retrieve its value
// 租户机制
// 在同一个应用实例和同一个数据库（集群）中为多个独立客户提供服务
// 确保数据在逻辑上或物理上隔离
type Tenant struct {
	Key       string
	ValueFunc func(ctx context.Context) string // retrieve the tenant's value based on the context
}

// interface for types that can modify GORM database queries
type Where interface {
	Where(db *gorm.DB) *gorm.DB
}

// represents a database query with its arguments
// contains the query condition and any associated parameters
type Query struct {
	// holds the condition to be used in the Gorm query
	// can be a string, a map, or a struct
	Query interface{}
	// arguments
	// will be used to replace placeholders in the query
	Args []interface{}
}

// a function that modifies options
type Option func(*Options)

// Options represents the configuration for a database query
// Options结构体中字段最后会通过以下方式来为*gorm.DB类型的实例添加查询条件
//
//	func (whr *Options) Where(db *gorm.DB) *gorm.DB{
//	  return db.Where(whr.Filters).Clauses(whr.Clauses...).Offset(whr.Offset).Limit(whr.Limit)
//	}
type Options struct {
	// 分页的起始位置
	Offset int `json:"offset"`
	// maximum number of results to return
	Limit int `json:"limit"`
	// filters contains key-value pairs for filtering records
	Filters map[any]any
	// contains custom clauses to be applied to the query
	// 查询子句
	Clauses []clause.Expression
	// queries to be executed
	Queries []Query
}

var registeredTenant Tenant

// 三种方式配置Options
// 通过NewWhere
// 通过便捷函数直接创建
// 通过链式调用

func WithOffset(offset int64) Option {
	return func(whr *Options) {
		if offset < 0 {
			offset = 0
		}
		whr.Offset = int(offset)
	}
}

func WithLimit(limit int64) Option {
	return func(whr *Options) {
		if limit <= 0 {
			limit = int64(defaultLimit)
		}
		whr.Limit = int(limit)
	}
}

func WithPage(page int, pageSize int) Option {
	return func(whr *Options) {
		if page == 0 {
			page = 1
		}
		if pageSize == 0 {
			pageSize = defaultLimit
		}
		whr.Offset = (page - 1) * pageSize
		whr.Limit = pageSize
	}
}

func WithFilters(filter map[any]any) Option {
	return func(whr *Options) {
		whr.Filters = filter
	}
}

func WithClauses(conds ...clause.Expression) Option {
	return func(whr *Options) {
		whr.Clauses = append(whr.Clauses, conds...)
	}
}

func WithQueries(queries interface{}, args ...interface{}) Option {
	return func(whr *Options) {
		whr.Queries = append(whr.Queries, Query{
			Query: queries,
			Args:  args,
		})
	}
}

// 通过函数选项模式配置*Options结构体
func NewWhere(opts ...Option) *Options {
	whr := &Options{
		Offset:  0,
		Limit:   defaultLimit,
		Filters: map[any]any{},
		Clauses: make([]clause.Expression, 0),
	}
	for _, opt := range opts {
		opt(whr)
	}
	return whr
}

func (whr *Options) O(offset int) *Options {
	if offset < 0 {
		offset = 0
	}
	whr.Offset = offset
	return whr
}

func (whr *Options) L(limit int) *Options {
	if limit <= 0 {
		limit = defaultLimit
	}
	whr.Limit = limit
	return whr
}

// set the pagination based on the page number and page size
func (whr *Options) P(page int, pageSize int) *Options {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultLimit
	}
	whr.Offset = (page - 1) * pageSize
	whr.Limit = pageSize
	return whr
}

func (whr *Options) C(conds ...clause.Expression) *Options {
	whr.Clauses = append(whr.Clauses, conds...)
	return whr
}

func (whr *Options) Q(query interface{}, args ...interface{}) *Options {
	whr.Queries = append(whr.Queries, Query{
		Query: query,
		Args:  args,
	})
	return whr
}

// retrieve the value associated with the registered tenant using the provided context
func (whr *Options) T(ctx context.Context) *Options {
	if registeredTenant.Key != "" && registeredTenant.ValueFunc != nil {
		whr.F(registeredTenant.Key, registeredTenant.ValueFunc(ctx))
	}
	return whr
}

// add filters
func (whr *Options) F(kvs ...any) *Options {
	if len(kvs)%2 != 0 {
		return whr
	}

	for i := 0; i < len(kvs); i += 2 {
		key := kvs[i]
		value := kvs[i+1]
		whr.Filters[key] = value
	}
	return whr
}

func (whr *Options) Where(db *gorm.DB) *gorm.DB {
	for _, query := range whr.Queries {
		conds := db.Statement.BuildCondition(query.Query, query.Args...)
		whr.Clauses = append(whr.Clauses, conds...)
	}
	return db.Where(whr.Filters).Clauses(whr.Clauses...).Offset(whr.Offset).Limit(whr.Limit)
}

// 提供一些便捷函数，用来快速创建一个指定了某类查询条件的*Options结构体实例
func O(offset int) *Options {
	return NewWhere().O(offset)
}

func L(limit int) *Options {
	return NewWhere().L(limit)
}

func P(page int, pageSize int) *Options {
	return NewWhere().P(page, pageSize)
}

func C(conds ...clause.Expression) *Options {
	return NewWhere().C(conds...)
}

func T(ctx context.Context) *Options {
	return NewWhere().F(registeredTenant.Key, registeredTenant.ValueFunc(ctx))
}

func F(kvs ...any) *Options {
	return NewWhere().F(kvs...)
}

// register a new tenant with the specified key and value function
func RegisterTenant(key string, valueFunc func(ctx context.Context) string) {
	registeredTenant = Tenant{
		Key:       key,
		ValueFunc: valueFunc,
	}
}

// 链式调用
// opts := NewWhere().O(10).L(20).F("name","John","status","active").
// C(clause.OrderBy{Columns:[]clause.OrderByColumn{
// 	{
// 		Column: clause.Column{Name: "created_at"},
// 		Desc:   true,
// 	},
// }}).P(2,10) 分页第2页每页10条数据
