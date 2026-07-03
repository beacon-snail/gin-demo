package repositories

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// QueryConditon 查询条件
type QueryCondition struct {
	Field    string      // 字段名
	Operator string      // 操作符，如 =, <, >, LIKE 等
	Value    interface{} // 值
}

// QueryBuilder 查询构建器

type QueryBuilder struct {
	conditions []QueryCondition // 查询条件
	orderBy    string           // 排序字段
	orderDir   string
	page       int      // 页码
	pageSize   int      // 每页数量
	preloads   []string // 预加载的关联表
	selects    []string // 选择的字段
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		conditions: []QueryCondition{},
		preloads:   []string{},
		selects:    []string{},
		page:       1,
		pageSize:   10,
		orderDir:   "DESC",
	}
}

// AddCondition 添加查询条件
func (qb *QueryBuilder) AddCondition(field string, operator string, value interface{}) *QueryBuilder {
	qb.conditions = append(qb.conditions, QueryCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return qb
}

// AddEqual 等于
func (qb *QueryBuilder) AddEqual(field string, value interface{}) *QueryBuilder {
	return qb.AddCondition(field, "eq", value)
}

// AddLike 模糊查询
func (qb *QueryBuilder) AddLike(field string, value string) *QueryBuilder {
	return qb.AddCondition(field, "like", value)
}

// AddIn 查询包含
func (qb *QueryBuilder) AddIn(field string, values interface{}) *QueryBuilder {
	return qb.AddCondition(field, "in", values)
}

// AddBetween 查询范围
func (qb *QueryBuilder) AddBetween(field string, start interface{}, end interface{}) *QueryBuilder {
	return qb.AddCondition(field, "between", []interface{}{start, end})
}

// AddDateRange 日期范围查询（自动处理日期格式）
func (qb *QueryBuilder) AddDateRange(field string, startDate, endDate string) *QueryBuilder {
	if startDate != "" && endDate != "" {
		qb.conditions = append(qb.conditions, QueryCondition{
			Field:    field,
			Operator: "between",
			Value:    []interface{}{startDate + " 00:00:00", endDate + " 23:59:59"},
		})
	} else if startDate != "" {
		qb.conditions = append(qb.conditions, QueryCondition{
			Field:    field,
			Operator: "gte",
			Value:    startDate + " 00:00:00",
		})
	} else if endDate != "" {
		qb.conditions = append(qb.conditions, QueryCondition{
			Field:    field,
			Operator: "lte",
			Value:    endDate + " 23:59:59",
		})
	}
	return qb
}

// SetPagination 设置分页参数
func (qb *QueryBuilder) SetPagination(page int, pageSize int) *QueryBuilder {
	if page > 0 {
		qb.page = page
	}
	if pageSize > 0 {
		qb.pageSize = pageSize
	}
	return qb
}

// SetOrder 设置排序参数
func (qb *QueryBuilder) SetOrder(field string, dir string) *QueryBuilder {
	qb.orderBy = field

	if dir == "asc" || dir == "desc" {
		qb.orderDir = dir
	}
	return qb
}

// AddPreload 添加预加载的关联表
func (qb *QueryBuilder) AddPreload(fields ...string) *QueryBuilder {
	qb.preloads = append(qb.preloads, fields...)
	return qb
}

// SetSelect 设置选择的字段
func (qb *QueryBuilder) SetSelect(fields ...string) *QueryBuilder {
	qb.selects = append(qb.selects, fields...)
	return qb
}

// Build 构建查询条件

func (qb *QueryBuilder) Build(db *gorm.DB) *gorm.DB {
	query := db

	// 应用选择字段
	if len(qb.selects) > 0 {
		query = query.Select(qb.selects)
	}

	// 应用预加载
	for _, preload := range qb.preloads {
		query = query.Preload(preload)
	}

	// 应用查询条件
	for _, cond := range qb.conditions {
		query = qb.applyCondition(query, cond)
	}

	// 应用排序
	if qb.orderBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", qb.orderBy, qb.orderDir))
	}

	// 应用分页
	offset := (qb.page - 1) * qb.pageSize
	query = query.Offset(offset).Limit(qb.pageSize)

	return query
}

// applyCondition 应用单个条件
func (qb *QueryBuilder) applyCondition(db *gorm.DB, cond QueryCondition) *gorm.DB {
	switch cond.Operator {
	case "eq":
		return db.Where(cond.Field+" = ?", cond.Value)
	case "like":
		return db.Where(cond.Field+" LIKE ?", "%"+cond.Value.(string)+"%")
	case "in":
		return db.Where(cond.Field+" IN ?", cond.Value)
	case "between":
		return db.Where(cond.Field+" BETWEEN ? AND ?", cond.Value.([]interface{})[0], cond.Value.([]interface{})[1])
	case "gt":
		return db.Where(cond.Field+" > ?", cond.Value)
	case "lt":
		return db.Where(cond.Field+" < ?", cond.Value)
	case "gte":
		return db.Where(cond.Field+" >= ?", cond.Value)
	case "lte":
		return db.Where(cond.Field+" <= ?", cond.Value)
	case "is_null":
		return db.Where(cond.Field + " IS NULL")
	case "is_not_null":
		return db.Where(cond.Field + " IS NOT NULL")
	default:
		return db
	}
}

// BuildCount 构建计数查询
func (qb *QueryBuilder) BuildCount(db *gorm.DB) *gorm.DB {
	query := db.Model(&struct{}{})

	// 应用查询条件（不包含分页和排序）
	for _, cond := range qb.conditions {
		query = qb.applyCondition(query, cond)
	}

	return query
}

// ParseQueryParams 从map解析查询参数（自动映射）
func ParseQueryParams(params map[string]string, fieldMapping map[string]string) *QueryBuilder {
	qb := NewQueryBuilder()

	for key, value := range params {
		if value == "" {
			continue
		}

		// 映射字段名
		field := key
		if mapped, ok := fieldMapping[key]; ok {
			field = mapped
		}

		// 解析操作符
		parts := strings.Split(key, "__")
		operator := "eq"
		if len(parts) > 1 {
			operator = parts[len(parts)-1]
			// 去除操作符部分，获取真实字段名
			field = strings.Join(parts[:len(parts)-1], "__")
			if mapped, ok := fieldMapping[field]; ok {
				field = mapped
			}
		}

		// 根据操作符处理值
		switch operator {
		case "like":
			qb.AddLike(field, value)
		case "in":
			values := strings.Split(value, ",")
			qb.AddIn(field, values)
		case "between":
			// 格式: start,end
			parts := strings.Split(value, ",")
			if len(parts) == 2 {
				qb.AddBetween(field, parts[0], parts[1])
			}
		case "gt", "gte", "lt", "lte", "ne":
			qb.AddCondition(field, operator, value)
		case "is_null":
			qb.AddCondition(field, "is_null", nil)
		case "is_not_null":
			qb.AddCondition(field, "is_not_null", nil)
		default:
			qb.AddEqual(field, value)
		}
	}

	return qb
}

// 字段映射配置
type FieldMapping map[string]string

// 自动映射查询参数
func AutoFilter(query *gorm.DB, params map[string]string, fieldMap map[string]string) *gorm.DB {
	qb := ParseQueryParams(params, fieldMap)
	return qb.Build(query)
}
