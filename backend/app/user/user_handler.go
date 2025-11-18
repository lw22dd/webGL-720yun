package user

import (
	"webGL-720yun/app/models"
	"webGL-720yun/pkg/middleware"
	"webGL-720yun/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	
	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	service *UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// RegisterRoutes 注册用户相关路由
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	// 公开路由
	public := router.Group("/auth")
	{
		public.POST("/login", h.Login)
		public.POST("/refresh", h.RefreshToken)
	}

	// 需要认证的路由
	protected := router.Group("")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.POST("/logout", h.Logout)
		protected.GET("/profile", h.GetProfile)
		protected.PUT("/profile", h.UpdateProfile)
		protected.PUT("/password", h.ChangePassword)
	}

	// 管理员路由
	admin := router.Group("/admin")
	admin.Use(authMiddleware.RequireAdmin())
	{
		admin.POST("/users", h.CreateUser)
		admin.GET("/users", h.GetUserList)
		admin.GET("/users/:id", h.GetUserByID)
		admin.PUT("/users/:id", h.UpdateUser)
		admin.DELETE("/users/:id", h.DeleteUser)
	}
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c.Writer, "请求参数错误")
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		utils.Unauthorized(c.Writer, err.Error())
		return
	}

	utils.Success(c.Writer, response)
}

// Register 用户注册（管理员创建学生账户）
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c.Writer, "请求参数错误")
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		utils.Error(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c.Writer, user)
}

// Logout 用户登出
func (h *UserHandler) Logout(c *gin.Context) {
	userID, _, _ := middleware.GetCurrentUser(c)
	
	// 获取访问令牌
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.BadRequest(c.Writer, "认证令牌格式错误")
		return
	}

	if err := h.service.Logout(userID, parts[1]); err != nil {
		utils.InternalServerError(c.Writer, "登出失败")
		return
	}

	utils.Success(c.Writer, gin.H{"message": "登出成功"})
}

// RefreshToken 刷新访问令牌
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c.Writer, "请求参数错误")
		return
	}

	response, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.Unauthorized(c.Writer, err.Error())
		return
	}

	utils.Success(c.Writer, response)
}

// GetProfile 获取个人资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _, _ := middleware.GetCurrentUser(c)
	
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c.Writer, "用户不存在")
		return
	}

	utils.Success(c.Writer, user)
}

// UpdateProfile 更新个人资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _, _ := middleware.GetCurrentUser(c)
	
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c.Writer, "请求参数错误")
		return
	}

	user, err := h.service.UpdateUser(userID, &req)
	if err != nil {
		utils.Error(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c.Writer, user)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _, _ := middleware.GetCurrentUser(c)
	
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c.Writer, "请求参数错误")
		return
	}

	if err := h.service.ChangePassword(userID, &req); err != nil {
		if err.Error() == "旧密码错误" {
			utils.Error(c.Writer, http.StatusBadRequest, err.Error())
		} else {
			utils.InternalServerError(c.Writer, err.Error())
		}
		return
	}

	utils.Success(c.Writer, gin.H{"message": "密码修改成功"})
}

// GetUserByID 根据ID获取用户
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c.Writer, "用户ID格式错误")
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		utils.NotFound(c.Writer, "用户不存在")
		return
	}

	utils.Success(c.Writer, user)
}

// GetUserList 获取用户列表
func (h *UserHandler) GetUserList(c *gin.Context) {
	var req models.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c.Writer, "请求参数错误")
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	response, err := h.service.GetUserList(&req)
	if err != nil {
		utils.InternalServerError(c.Writer, "获取用户列表失败")
		return
	}

	utils.Success(c.Writer, gin.H{
		"total": response.Total,
		"users": response.Users,
		"page":  req.Page,
		"page_size": req.PageSize,
	})
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c.Writer, "用户ID格式错误")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c.Writer, "请求参数错误")
		return
	}

	user, err := h.service.UpdateUser(uint(id), &req)
	if err != nil {
		utils.Error(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c.Writer, user)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c.Writer, "用户ID格式错误")
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		utils.InternalServerError(c.Writer, "删除用户失败")
		return
	}

	utils.Success(c.Writer, gin.H{"message": "删除成功"})
}