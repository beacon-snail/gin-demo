package controllers

import (
	"strconv"

	"gin-mysql-demo/dto"
	"gin-mysql-demo/services"
	"gin-mysql-demo/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		service: services.NewUserService(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	user, err := c.service.Register(&req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.Success(ctx, user)
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	resp, err := c.service.Login(&req)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.Success(ctx, resp)
}

// GetUser 获取用户详情
func (c *UserController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "invalid user id")
		return
	}

	user, err := c.service.GetUser(uint(id))
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}

	utils.Success(ctx, user)
}

// GetUsers 获取用户列表（支持智能过滤）
func (c *UserController) GetUsers(ctx *gin.Context) {
	// 获取查询参数
	params := make(map[string]string)

	// 自动提取所有查询参数
	for key, values := range ctx.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	users, total, err := c.service.GetUsers(params)
	if err != nil {
		utils.InternalError(ctx, err.Error())
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(params["page"])
	if page == 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(params["page_size"])
	if pageSize == 0 {
		pageSize = 10
	}

	utils.Success(ctx, dto.UserListResponse{
		Total: total,
		Page:  page,
		Size:  pageSize,
		Data:  users,
	})
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "invalid user id")
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	user, err := c.service.UpdateUser(uint(id), &req)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}

	utils.Success(ctx, user)
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "invalid user id")
		return
	}

	if err := c.service.DeleteUser(uint(id)); err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "user deleted successfully", nil)
}
