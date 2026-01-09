package user

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"
	"webGL-720yun/pkg/middleware"
	"webGL-720yun/pkg/services/redis"
	"webGL-720yun/pkg/utils"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// 用户服务
type UserService struct {
	db           *gorm.DB
	jwtService   *middleware.JWTService
	redisService *redis.RedisService
}

func NewUserService(db *gorm.DB, jwtService *middleware.JWTService, redisService *redis.RedisService) *UserService {
	return &UserService{
		db:           db,
		jwtService:   jwtService,
		redisService: redisService,
	}
}

// 用户注册（管理员创建学生账户）
func (s *UserService) Register(req *RegisterRequest) (*User, error) {
	// 检查用户名是否已存在
	var existingUser User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已存在")
		}
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		if err := s.db.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
			return nil, errors.New("手机号已存在")
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户
	user := &User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Nickname: req.Nickname,
		RoleID:   req.RoleID,
		Status:   1, // 1表示活跃状态
	}

	// 处理Phone字段，仅在非空时设置，避免唯一索引冲突
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("用户创建失败: %v", err)
	}

	// 如果是学生角色，添加学生信息和关联
	// 先获取角色信息
	var role Role
	if err := tx.First(&role, req.RoleID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("获取角色信息失败: %v", err)
	}

	if role.Name == RoleStudent {
		// 检查学号是否已存在
		if req.StudentID != "" {
			var existingStudent Student
			if err := tx.Where("student_id = ?", req.StudentID).First(&existingStudent).Error; err == nil {
				tx.Rollback()
				return nil, errors.New("学号已存在")
			}
		}

		// 更新学生信息
		student := &Student{
			User:      *user,
			StudentID: req.StudentID,
			ClassID:   req.ClassID,
		}

		if err := tx.Save(student).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("学生信息保存失败: %v", err)
		}

		// 关联教师
		if len(req.TeacherIDs) > 0 {
			for _, teacherID := range req.TeacherIDs {
				studentTeacher := &StudentTeacher{
					StudentID: user.ID,
					TeacherID: teacherID,
				}
				if err := tx.Create(studentTeacher).Error; err != nil {
					tx.Rollback()
					return nil, fmt.Errorf("关联教师失败: %v", err)
				}
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("事务提交失败: %v", err)
	}

	return user, nil
}

// 用户登录
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户，支持用户名或邮箱登录，并预加载角色
	var user User
	if err := s.db.Preload("Role").Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, fmt.Errorf("数据库查询失败: %v", err)
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账户已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成访问令牌
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Username, user.Role.Name)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %v", err)
	}

	// 生成刷新令牌
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
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

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		User:         &user,
	}, nil
}

// 用户登出
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

// 刷新访问令牌
func (s *UserService) RefreshToken(refreshToken string) (*RefreshTokenResponse, error) {
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
	var user User
	if err := s.db.Preload("Role").First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户状态
	if user.Status != UserStatusActive {
		return nil, errors.New("账户已被禁用")
	}

	// 生成新的访问令牌
	newAccessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Username, user.Role.Name)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %v", err)
	}

	// 更新会话中的访问令牌
	if err := s.redisService.SaveUserSession(userID, newAccessToken, refreshToken, time.Duration(24)*time.Hour); err != nil {
		return nil, fmt.Errorf("更新会话失败: %v", err)
	}

	return &RefreshTokenResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
	}, nil
}

// 根据ID获取用户
func (s *UserService) GetUserByID(userID uint) (*User, error) {
	// 先尝试从缓存获取
	cachedUser, err := s.redisService.GetCachedUserInfo(userID)
	if err == nil && cachedUser != nil {
		var user User
		if userData, err := json.Marshal(cachedUser); err == nil {
			if err := json.Unmarshal(userData, &user); err == nil {
				return &user, nil
			}
		}
	}

	// 缓存穿透：检查用户是否不存在
	if cachedUser == nil && err == nil {
		// 缓存中标记为不存在，直接返回
		return nil, errors.New("用户不存在")
	}

	// 缓存击穿：使用分布式锁确保只有一个请求访问数据库
	lockKey := fmt.Sprintf("user:%d", userID)
	expiresIn := time.Duration(5) * time.Second // 锁过期时间

	// 尝试获取锁
	if s.redisService.AcquireLock(lockKey, expiresIn) {
		defer s.redisService.ReleaseLock(lockKey) // 释放锁

		// 再次尝试从缓存获取（防止锁等待期间其他请求已更新缓存）
		if cachedUser, err := s.redisService.GetCachedUserInfo(userID); err == nil && cachedUser != nil {
			var user User
			if userData, err := json.Marshal(cachedUser); err == nil {
				if err := json.Unmarshal(userData, &user); err == nil {
					return &user, nil
				}
			}
		}

		// 从数据库获取
		var user User
		if err := s.db.Preload("Role").First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 缓存空对象防止缓存穿透
				s.redisService.CacheEmptyUserInfo(userID, time.Duration(30)*time.Minute)
				return nil, errors.New("用户不存在")
			}
			return nil, err
		}

		// 缓存用户信息
		s.redisService.CacheUserInfo(userID, user, time.Duration(1)*time.Hour)

		return &user, nil
	} else {
		// 获取锁失败，等待一段时间后重试
		time.Sleep(100 * time.Millisecond)
		// 递归重试，最多重试3次
		for i := 0; i < 3; i++ {
			if cachedUser, err := s.redisService.GetCachedUserInfo(userID); err == nil && cachedUser != nil {
				var user User
				if userData, err := json.Marshal(cachedUser); err == nil {
					if err := json.Unmarshal(userData, &user); err == nil {
						return &user, nil
					}
				}
			}
			time.Sleep(200 * time.Millisecond)
		}

		// 重试后仍失败，直接访问数据库
		var user User
		if err := s.db.Preload("Role").First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 缓存空对象防止缓存穿透
				s.redisService.CacheEmptyUserInfo(userID, time.Duration(30)*time.Minute)
				return nil, errors.New("用户不存在")
			}
			return nil, err
		}

		return &user, nil
	}
}

// 更新用户信息
func (s *UserService) UpdateUser(userID uint, req *UpdateUserRequest) (*User, error) {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 更新字段
	if req.Email != "" {
		// 检查邮箱是否已被使用
		var existingUser User
		if err := tx.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			tx.Rollback()
			return nil, errors.New("邮箱已被使用")
		}
		user.Email = req.Email
	}

	if req.Phone != "" {
		// 检查手机号是否已被使用
		var existingUser User
		if err := tx.Where("phone = ? AND id != ?", req.Phone, userID).First(&existingUser).Error; err == nil {
			tx.Rollback()
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

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新用户信息失败: %v", err)
	}

	// 检查是否是学生角色
	var role Role
	if err := tx.First(&role, user.RoleID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("获取角色信息失败: %v", err)
	}

	if role.Name == RoleStudent {
		// 更新学生信息
		var student Student
		if err := tx.First(&student, userID).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("获取学生信息失败: %v", err)
		}

		// 更新班级信息
		if req.ClassID > 0 && req.ClassID != student.ClassID {
			student.ClassID = req.ClassID
			if err := tx.Save(&student).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("更新班级信息失败: %v", err)
			}
		}

		// 更新教师关联
		if len(req.TeacherIDs) > 0 {
			// 删除旧的关联
			if err := tx.Where("student_id = ?", userID).Delete(&StudentTeacher{}).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("删除旧教师关联失败: %v", err)
			}

			// 添加新的关联
			for _, teacherID := range req.TeacherIDs {
				studentTeacher := &StudentTeacher{
					StudentID: userID,
					TeacherID: teacherID,
				}
				if err := tx.Create(studentTeacher).Error; err != nil {
					tx.Rollback()
					return nil, fmt.Errorf("添加新教师关联失败: %v", err)
				}
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("事务提交失败: %v", err)
	}

	// 更新缓存
	s.redisService.CacheUserInfo(userID, user, time.Duration(1)*time.Hour)

	// 重新加载角色信息
	s.db.Preload("Role").First(&user, userID)

	return &user, nil
}

// 修改密码
func (s *UserService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 调用用户对象的ChangePassword方法
	if err := user.ChangePassword(req.OldPassword, req.NewPassword); err != nil {
		tx.Rollback()
		return err
	}

	// 保存更新后的用户
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新密码失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交失败: %v", err)
	}

	// 更新缓存
	s.redisService.DeleteCachedUserInfo(userID)

	return nil
}

// 获取用户列表
func (s *UserService) GetUserList(req *UserListRequest) (*UserListResponse, error) {
	var users []*User
	var total int64

	query := s.db.Model(&User{}).Preload("Role")

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
	// 只有当Status被显式设置时才添加条件（使用指针类型或特殊标记来判断是否被设置）
	// 对于int类型，我们无法直接判断是否被显式设置，所以需要修改逻辑
	// 移除默认的状态过滤，让所有状态的用户都能被查询到
	// 如果需要过滤特定状态，会通过前端显式传递status参数

	// 按班级ID查询
	if req.ClassID > 0 {
		query = query.Joins("JOIN students ON students.id = users.id").Where("students.class_id = ?", req.ClassID)
	}
	// 关键词搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", keyword, keyword, keyword)
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
	// 返回响应
	return &UserListResponse{
		PageInfo: utils.PageInfo{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
			Pages:    int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
		},
		Users: users,
	}, nil
}

// 重置密码
func (s *UserService) ResetPassword(req *ResetPasswordRequest) (string, error) {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找用户
	var user User
	if err := tx.Where("username = ?", req.Username).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户不存在")
		}
		return "", fmt.Errorf("数据库查询失败: %v", err)
	}

	// 调用用户对象的ResetPassword方法
	tempPassword, err := user.ResetPassword()
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("密码重置失败: %v", err)
	}

	// 保存更新后的用户
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("保存密码失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("事务提交失败: %v", err)
	}

	// 更新缓存
	s.redisService.DeleteCachedUserInfo(user.ID)

	return tempPassword, nil
}

// 删除用户
func (s *UserService) DeleteUser(userID uint) error {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除用户关联的学生-教师关系
	if err := tx.Where("student_id = ?", userID).Delete(&StudentTeacher{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除学生-教师关系失败: %v", err)
	}

	// 删除用户
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除用户失败: %v", err)
	}

	// 删除相关数据
	if err := tx.Where("user_id = ?", userID).Delete(&UserSession{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除用户会话失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交失败: %v", err)
	}

	// 删除缓存
	s.redisService.DeleteCachedUserInfo(userID)
	s.redisService.DeleteUserSession(userID)
	s.redisService.SetUserOffline(userID)

	return nil
}

// 解析Excel文件，提取学生数据
func (s *UserService) ParseExcel(file multipart.File) ([]StudentExcelData, error) {
	// 解析Excel文件，提取学生数据
	return s.parseExcelFile(file)
}

// 解析CSV文件，提取学生数据
func (s *UserService) ParseCSV(file multipart.File) ([]StudentExcelData, error) {
	// 解析CSV文件，提取学生数据
	return s.parseCSVFile(file)
}

// 解析文件，根据文件名后缀自动选择解析方式
func (s *UserService) ParseFile(file multipart.File, filename string) ([]StudentExcelData, error) {
	// 根据文件名后缀选择解析方式
	if strings.HasSuffix(strings.ToLower(filename), ".xlsx") || strings.HasSuffix(strings.ToLower(filename), ".xls") {
		return s.parseExcelFile(file)
	} else if strings.HasSuffix(strings.ToLower(filename), ".csv") {
		return s.parseCSVFile(file)
	}
	return nil, errors.New("不支持的文件格式，仅支持.xlsx, .xls和.csv文件")
}

// 内部方法：解析Excel文件
func (s *UserService) parseExcelFile(file multipart.File) ([]StudentExcelData, error) {
	// 创建excelize文件
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, errors.New("Excel文件中没有工作表")
	}

	// 获取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取Excel行失败: %v", err)
	}

	if len(rows) < 2 {
		return nil, errors.New("Excel文件中至少需要包含表头和一行数据")
	}

	// 解析表头，自动识别关键字段
	head := rows[0]
	nameIndex := -1
	studentIDIndex := -1
	emailIndex := -1
	phoneIndex := -1
	classNameIndex := -1

	// 自动识别字段
	for i, field := range head {
		field = strings.TrimSpace(field)
		fieldLower := strings.ToLower(field)

		// 识别姓名
		if nameIndex == -1 && (strings.Contains(fieldLower, "姓名") || strings.Contains(fieldLower, "name")) {
			nameIndex = i
		}

		// 识别学号
		if studentIDIndex == -1 && (strings.Contains(fieldLower, "学号") || strings.Contains(fieldLower, "student") || strings.Contains(fieldLower, "id")) {
			studentIDIndex = i
		}

		// 识别邮箱
		if emailIndex == -1 && (strings.Contains(fieldLower, "邮箱") || strings.Contains(fieldLower, "email")) {
			emailIndex = i
		}

		// 识别手机号
		if phoneIndex == -1 && (strings.Contains(fieldLower, "手机") || strings.Contains(fieldLower, "phone")) {
			phoneIndex = i
		}

		// 识别班级
		if classNameIndex == -1 && (strings.Contains(fieldLower, "班级") || strings.Contains(fieldLower, "class")) {
			classNameIndex = i
		}
	}

	// 验证必要字段是否存在
	if nameIndex == -1 {
		return nil, errors.New("无法识别姓名字段")
	}
	if studentIDIndex == -1 {
		return nil, errors.New("无法识别学号字段")
	}
	if emailIndex == -1 {
		return nil, errors.New("无法识别邮箱字段")
	}

	// 解析数据行
	var students []StudentExcelData
	for i, row := range rows[1:] {
		// 跳过空行
		if len(row) == 0 || (len(row) > nameIndex && strings.TrimSpace(row[nameIndex]) == "") {
			continue
		}

		// 获取字段值
		name := ""
		if len(row) > nameIndex {
			name = strings.TrimSpace(row[nameIndex])
		}

		studentID := ""
		if len(row) > studentIDIndex {
			studentID = strings.TrimSpace(row[studentIDIndex])
		}

		email := ""
		if len(row) > emailIndex {
			email = strings.TrimSpace(row[emailIndex])
		}

		phone := ""
		if len(row) > phoneIndex && phoneIndex != -1 {
			phone = strings.TrimSpace(row[phoneIndex])
		}

		className := ""
		if len(row) > classNameIndex && classNameIndex != -1 {
			className = strings.TrimSpace(row[classNameIndex])
		}

		// 验证必要字段
		if name == "" || studentID == "" || email == "" {
			continue
		}

		// 根据班级名称获取班级ID
		var classID uint = 0
		if className != "" {
			var class Class
			if err := s.db.Where("name = ?", className).First(&class).Error; err == nil {
				classID = class.ID
			}
		}

		// 创建学生数据
		student := StudentExcelData{
			Index:      i + 2, // 行号从2开始
			Name:       name,
			StudentID:  studentID,
			Email:      email,
			Phone:      phone,
			ClassName:  className,
			ClassID:    classID,
			TeacherIDs: []uint{}, // 暂时为空，后续处理
		}

		students = append(students, student)
	}

	return students, nil
}

// 内部方法：解析CSV文件
func (s *UserService) parseCSVFile(file multipart.File) ([]StudentExcelData, error) {
	// 创建csv reader
	reader := csv.NewReader(file)

	// 设置字段映射
	reader.FieldsPerRecord = -1 // 允许可变字段数

	// 读取所有行
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取CSV文件失败: %v", err)
	}

	if len(rows) < 2 {
		return nil, errors.New("CSV文件中至少需要包含表头和一行数据")
	}

	// 解析表头，自动识别关键字段
	head := rows[0]
	nameIndex := -1
	studentIDIndex := -1
	emailIndex := -1
	phoneIndex := -1
	classNameIndex := -1

	// 自动识别字段
	for i, field := range head {
		field = strings.TrimSpace(field)
		fieldLower := strings.ToLower(field)

		// 识别姓名
		if nameIndex == -1 && (strings.Contains(fieldLower, "姓名") || strings.Contains(fieldLower, "name")) {
			nameIndex = i
		}

		// 识别学号
		if studentIDIndex == -1 && (strings.Contains(fieldLower, "学号") || strings.Contains(fieldLower, "student") || strings.Contains(fieldLower, "id")) {
			studentIDIndex = i
		}

		// 识别邮箱
		if emailIndex == -1 && (strings.Contains(fieldLower, "邮箱") || strings.Contains(fieldLower, "email")) {
			emailIndex = i
		}

		// 识别手机号
		if phoneIndex == -1 && (strings.Contains(fieldLower, "手机") || strings.Contains(fieldLower, "phone")) {
			phoneIndex = i
		}

		// 识别班级
		if classNameIndex == -1 && (strings.Contains(fieldLower, "班级") || strings.Contains(fieldLower, "class")) {
			classNameIndex = i
		}
	}

	// 验证必要字段是否存在
	if nameIndex == -1 {
		return nil, errors.New("无法识别姓名字段")
	}
	if studentIDIndex == -1 {
		return nil, errors.New("无法识别学号字段")
	}
	if emailIndex == -1 {
		return nil, errors.New("无法识别邮箱字段")
	}

	// 解析数据行
	var students []StudentExcelData
	for i, row := range rows[1:] {
		// 跳过空行
		if len(row) == 0 || (len(row) > nameIndex && strings.TrimSpace(row[nameIndex]) == "") {
			continue
		}

		// 获取字段值
		name := ""
		if len(row) > nameIndex {
			name = strings.TrimSpace(row[nameIndex])
		}

		studentID := ""
		if len(row) > studentIDIndex {
			studentID = strings.TrimSpace(row[studentIDIndex])
		}

		email := ""
		if len(row) > emailIndex {
			email = strings.TrimSpace(row[emailIndex])
		}

		phone := ""
		if len(row) > phoneIndex && phoneIndex != -1 {
			phone = strings.TrimSpace(row[phoneIndex])
		}

		className := ""
		if len(row) > classNameIndex && classNameIndex != -1 {
			className = strings.TrimSpace(row[classNameIndex])
		}

		// 验证必要字段
		if name == "" || studentID == "" || email == "" {
			continue
		}

		// 根据班级名称获取班级ID
		var classID uint = 0
		if className != "" {
			var class Class
			if err := s.db.Where("name = ?", className).First(&class).Error; err == nil {
				classID = class.ID
			}
		}

		// 创建学生数据
		student := StudentExcelData{
			Index:      i + 2, // 行号从2开始
			Name:       name,
			StudentID:  studentID,
			Email:      email,
			Phone:      phone,
			ClassName:  className,
			ClassID:    classID,
			TeacherIDs: []uint{}, // 暂时为空，后续处理
		}

		students = append(students, student)
	}

	return students, nil
}

// 批量注册学生
func (s *UserService) BatchRegister(students []StudentExcelData) (*BatchRegisterResponse, error) {
	var successCount int
	var failedCount int
	var results []BatchRegisterResult
	var errors []map[string]interface{}

	// 遍历学生数据，逐个注册
	for _, student := range students {
		// 创建注册请求
		req := &RegisterRequest{
			Username:   student.StudentID, // 使用学号作为用户名
			Password:   "123456",          // 默认密码
			Email:      student.Email,
			Phone:      student.Phone,
			Nickname:   student.Name,
			RoleID:     3, // 假设学生角色ID为3
			StudentID:  student.StudentID,
			ClassID:    student.ClassID,
			TeacherIDs: student.TeacherIDs,
		}

		// 注册学生
		_, err := s.Register(req)
		var result BatchRegisterResult
		result.Index = student.Index
		result.Username = student.StudentID
		result.StudentID = student.StudentID

		if err != nil {
			failedCount++
			result.Status = "failed"
			result.Message = err.Error()
			errors = append(errors, map[string]interface{}{
				"index":    student.Index,
				"username": student.StudentID,
				"error":    err.Error(),
			})
		} else {
			successCount++
			result.Status = "success"
			result.Message = "注册成功"
		}

		results = append(results, result)
	}

	// 创建响应
	response := &BatchRegisterResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Results:      results,
		Errors:       errors,
	}

	return response, nil
}
