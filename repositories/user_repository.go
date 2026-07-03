package repositories

import (
	"gin-mysql-demo/database"
	"gin-mysql-demo/models"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// 用户字段映射
var userFieldMap = map[string]string{
	"id":         "id",
	"username":   "username",
	"email":      "email",
	"nickname":   "nickname",
	"age":        "age",
	"gender":     "gender",
	"phone":      "phone",
	"status":     "status",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

// FindAll 使用QueryBuilder智能查询
func (r *UserRepository) FindAll(params map[string]string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 1. 构建基础查询
	query := r.db.Model(&models.User{})

	// 2. 使用自动过滤器
	// 支持的操作符: __like, __gt, __gte, __lt, __lte, __ne, __in, __between
	// 例如: ?username__like=张&age__gte=18&status=1
	query = AutoFilter(query, params, userFieldMap)

	// 3. 处理分页参数
	page := 1
	pageSize := 10
	if p, ok := params["page"]; ok {
		if p != "" {
			// 转换逻辑...
			page = 1 // 简化示例
		}
	}
	if ps, ok := params["page_size"]; ok {
		if ps != "" {
			// 转换逻辑...
		}
	}

	// 4. 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 5. 排序
	if orderBy, ok := params["order_by"]; ok && orderBy != "" {
		orderDir := "DESC"
		if dir, ok := params["order_dir"]; ok && dir == "ASC" {
			orderDir = "ASC"
		}
		query = query.Order(orderBy + " " + orderDir)
	} else {
		query = query.Order("created_at DESC")
	}

	// 6. 分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// 7. 预加载（可选）
	if preload, ok := params["preload"]; ok && preload == "orders" {
		query = query.Preload("Orders")
	}

	// 8. 执行查询
	err := query.Find(&users).Error
	return users, total, err
}

// FindWithBuilder 使用QueryBuilder的完整示例
func (r *UserRepository) FindWithBuilder(filters map[string]interface{}) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	qb := NewQueryBuilder()

	// 手动添加条件
	for key, value := range filters {
		switch v := value.(type) {
		case string:
			if v != "" {
				qb.AddEqual(key, v)
			}
		case int, int64, float64:
			if v != 0 {
				qb.AddEqual(key, v)
			}
		case []string:
			if len(v) > 0 {
				qb.AddIn(key, v)
			}
		case map[string]interface{}:
			// 处理范围查询
			if start, ok := v["start"]; ok {
				if end, ok := v["end"]; ok {
					qb.AddBetween(key, start, end)
				}
			}
		}
	}

	// 设置分页
	qb.SetPagination(1, 10)

	// 设置排序
	qb.SetOrder("created_at", "DESC")

	// 执行查询
	query := qb.Build(r.db.Model(&models.User{}))

	// 获取总数
	countQuery := qb.BuildCount(r.db.Model(&models.User{}))
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Find(&users).Error
	return users, total, err
}

// 高级条件查询示例
func (r *UserRepository) FindAdvanced(params map[string]string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 使用QueryBuilder
	qb := NewQueryBuilder()

	// 用户名模糊查询
	if username, ok := params["username"]; ok && username != "" {
		qb.AddLike("username", username)
	}

	// 年龄范围
	if ageStart, ok := params["age_start"]; ok && ageStart != "" {
		qb.AddCondition("age", "gte", ageStart)
	}
	if ageEnd, ok := params["age_end"]; ok && ageEnd != "" {
		qb.AddCondition("age", "lte", ageEnd)
	}

	// 状态精确匹配
	if status, ok := params["status"]; ok && status != "" {
		qb.AddEqual("status", status)
	}

	// 性别
	if gender, ok := params["gender"]; ok && gender != "" {
		qb.AddEqual("gender", gender)
	}

	// 创建日期范围
	if startDate, ok := params["start_date"]; ok {
		if endDate, ok := params["end_date"]; ok {
			qb.AddDateRange("created_at", startDate, endDate)
		}
	}

	// 多状态IN查询
	if statuses, ok := params["statuses"]; ok && statuses != "" {
		statusList := strings.Split(statuses, ",")
		qb.AddIn("status", statusList)
	}

	// 用户ID列表
	if ids, ok := params["ids"]; ok && ids != "" {
		idList := strings.Split(ids, ",")
		qb.AddIn("id", idList)
	}

	// 分页
	page := 1
	pageSize := 10
	if p, ok := params["page"]; ok && p != "" {
		page, _ = strconv.Atoi(p)
	}
	if ps, ok := params["page_size"]; ok && ps != "" {
		pageSize, _ = strconv.Atoi(ps)
	}
	qb.SetPagination(page, pageSize)

	// 排序
	if orderBy, ok := params["order_by"]; ok && orderBy != "" {
		orderDir := "DESC"
		if dir, ok := params["order_dir"]; ok && dir == "ASC" {
			orderDir = "ASC"
		}
		qb.SetOrder(orderBy, orderDir)
	}

	// 执行查询
	query := qb.Build(r.db.Model(&models.User{}))
	countQuery := qb.BuildCount(r.db.Model(&models.User{}))

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Find(&users).Error
	return users, total, err
}

// 其他基础方法
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) UpdateLastLogin(id uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).
		Update("last_login", gorm.Expr("NOW()")).Error
}
