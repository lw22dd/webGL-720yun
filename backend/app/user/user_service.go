package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"webGL-720yun/app/models"
	"webGL-720yun/pkg/database"
	"webGL-720yun/pkg/services/jwt"
	"webGL-720yun/pkg/services/redis"
	"webGL-720yun/pkg/utils"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	jwtService   *jwt.JWTService
	redisService *redis.RedisService
}

// NewUserService 创建用户服务
func NewUserService(jwtService *jwt.JWTService, redisService *redis.RedisService) *UserService {
	return &UserService{
		jwtService:   jwtService,
		redisService: redisService,
	}
}

// Register 用户注册（管理员创建学生账户）
func (s *UserService) Register(req *models.RegisterRequest) (*models.User, error) {
	// 检查角色是否存在
	var role models.Role
	if err := database.DB.First(&role, req.RoleID).Error; err != nil {
		return nil, errors.New("角色不存在")
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已存在")
		}
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		if err := database.DB.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
			return nil, errors.New("手机号已存在")
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		RoleID:   req.RoleID,
		Status:   models.UserStatusActive,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("用户创建失败: %v", err)
	}

	// 加载角色信息
	database.DB.Preload("Role").First(user, user.ID)

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// 查找用户
	var user models.User
	if err := database.DB.Preload("Role").Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, fmt.Errorf("数据库查询失败: %v", err)
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		return nil, errors.New("账户已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成访问令牌
	accessToken, err := s.jwtService.GenerateAccessToken(&user)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %v", err)
	}

	// 生成刷新令牌
	refreshToken, err := s.jwtService.GenerateRefreshToken(&user)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %v", err)
	}

	// 保存会话到Redis
	expiresIn := time.Duration(24) * time.Hour // 24小时
	if err := s.redisService.SaveUserSession(user.ID, accessToken, refreshToken, expiresIn); err != nil {
		return nil, fmt.Errorf("保存会话失败: %v", err)
	}

	// 设置用户在线状态
	onlineExpiresIn := time.Duration(30) * time.Minute // 30分钟
	if err := s.redisService.SetUserOnline(user.ID, onlineExpiresIn); err != nil {
		return nil, fmt.Errorf("设置在线状态失败: %v", err)
	}

	// 缓存用户信息
	if err := s.redisService.CacheUserInfo(user.ID, user, time.Duration(1)*time.Hour); err != nil {
		return nil, fmt.Errorf("缓存用户信息失败: %v", err)
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(time.Duration(24) * time.Hour),
		User:         &user,
	}, nil
}

// Logout 用户登出
func (s *UserService) Logout(userID uint, accessToken string) error {
	// 将访问令牌加入黑名单
	if err := s.redisService.AddToBlacklist(accessToken, time.Duration(24)*time.Hour); err != nil {
		return fmt.Errorf("添加令牌到黑名单失败: %v", err)
	}

	// 删除用户会话
	if err := s.redisService.DeleteUserSession(userID); err != nil {
		return fmt.Errorf("删除用户会话失败: %v", err)
	}

	// 设置用户离线
	if err := s.redisService.SetUserOffline(userID); err != nil {
		return fmt.Errorf("设置用户离线失败: %v", err)
	}

	// 删除缓存的用户信息
	if err := s.redisService.DeleteCachedUserInfo(userID); err != nil {
		return fmt.Errorf("删除缓存用户信息失败: %v", err)
	}

	return nil
}

// RefreshToken 刷新访问令牌
func (s *UserService) RefreshToken(refreshToken string) (*models.RefreshTokenResponse, error) {
	// 验证刷新令牌
	claims, err := s.jwtService.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 获取用户ID
	userID := claims.UserID

	// 检查用户会话是否存在
	session, err := s.redisService.GetUserSession(userID)
	if err != nil {
		return nil, errors.New("会话不存在或已过期")
	}

	// 验证刷新令牌是否匹配
	if session["refresh_token"] != refreshToken {
		return nil, errors.New("刷新令牌不匹配")
	}

	// 查找用户
	var user models.User
	if err := database.DB.Preload("Role").First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		return nil, errors.New("账户已被禁用")
	}

	// 生成新的访问令牌
	newAccessToken, err := s.jwtService.GenerateAccessToken(&user)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %v", err)
	}

	// 更新会话中的访问令牌
	if err := s.redisService.SaveUserSession(userID, newAccessToken, refreshToken, time.Duration(24)*time.Hour); err != nil {
		return nil, fmt.Errorf("更新会话失败: %v", err)
	}

	return &models.RefreshTokenResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Duration(24) * time.Hour),
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	// 先尝试从缓存获取
	if cachedUser, err := s.redisService.GetCachedUserInfo(userID); err == nil {
		var user models.User
		if userData, err := json.Marshal(cachedUser); err == nil {
			if err := json.Unmarshal(userData, &user); err == nil {
				return &user, nil
			}
		}
	}

	// 从数据库获取
	var user models.User
	if err := database.DB.Preload("Role").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 缓存用户信息
	s.redisService.CacheUserInfo(userID, user, time.Duration(1)*time.Hour)

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID uint, req *models.UpdateUserRequest) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 更新字段
	if req.Email != "" {
		// 检查邮箱是否已被使用
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已被使用")
		}
		user.Email = req.Email
	}

	if req.Phone != "" {
		// 检查手机号是否已被使用
		var existingUser models.User
		if err := database.DB.Where("phone = ? AND id != ?", req.Phone, userID).First(&existingUser).Error; err == nil {
			return nil, errors.New("手机号已被使用")
		}
		user.Phone = req.Phone
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if req.Status >= 0 {
		user.Status = req.Status
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("更新用户信息失败: %v", err)
	}

	// 更新缓存
	s.redisService.CacheUserInfo(userID, user, time.Duration(1)*time.Hour)

	// 重新加载角色信息
	database.DB.Preload("Role").First(&user, userID)

	return &user, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新密码
	user.Password = hashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}

	return nil
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(req *models.UserListRequest) (*models.UserListResponse, error) {
	var users []*models.User
	var total int64

	query := database.DB.Model(&models.User{}).Preload("Role")

	// 条件查询
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Email != "" {
		query = query.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.RoleID > 0 {
		query = query.Where("role_id = ?", req.RoleID)
	}
	if req.Status >= 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	return &models.UserListResponse{
		Total: total,
		Users: users,
	}, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 删除用户
	if err := database.DB.Delete(&user).Error; err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}

	// 删除相关数据
	database.DB.Where("user_id = ?", userID).Delete(&models.UserSession{})

	// 删除缓存
	s.redisService.DeleteCachedUserInfo(userID)
	s.redisService.DeleteUserSession(userID)
	s.redisService.SetUserOffline(userID)

	return nil
}
