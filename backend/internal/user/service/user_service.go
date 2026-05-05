package service

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"webGL-720yun/internal/model"
	"webGL-720yun/internal/user/dto"
	"webGL-720yun/internal/user/repository"
	"webGL-720yun/pkg/jwt"
	"webGL-720yun/pkg/redis"
	"webGL-720yun/pkg/utils"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type UserService struct {
	repo         *repository.UserRepository
	jwtService   *jwt.JWTService
	redisService *redis.RedisService
}

func NewUserService(db *gorm.DB, jwtService *jwt.JWTService, redisService *redis.RedisService) *UserService {
	return &UserService{
		repo:         repository.NewUserRepository(db),
		jwtService:   jwtService,
		redisService: redisService,
	}
}

func (s *UserService) Register(req *dto.RegisterRequest) (*model.User, error) {
	if err := s.checkUserExists(req.Username, req.Email, req.Phone, 0); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	var user *model.User

	err = s.repo.Transaction(func(tx *gorm.DB) error {
		user = &model.User{
			Username: req.Username,
			Password: hashedPassword,
			Nickname: req.Nickname,
			RoleID:   req.RoleID,
			ClassID:  req.ClassID,
			Status:   model.UserStatusActive,
		}

		if req.Email != "" {
			user.Email = req.Email
		}

		if req.Phone != "" {
			user.Phone = req.Phone
		}

		role, err := s.repo.FindRoleByID(req.RoleID)
		if err != nil {
			return fmt.Errorf("获取角色信息失败: %v", err)
		}

		// 学生角色时，使用传入的 ID 作为学号
		if role.Name == model.RoleStudent {
			if req.ID == 0 {
				return errors.New("学生角色必须提供学号(ID)")
			}
			user.ID = req.ID
		}

		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("用户创建失败: %v", err)
		}

		// 学生角色处理 - 关联教师
		if role.Name == model.RoleStudent {
			if len(req.TeacherIDs) > 0 {
				for _, teacherID := range req.TeacherIDs {
					studentTeacher := &model.StudentTeacher{
						StudentID: user.ID,
						TeacherID: teacherID,
					}
					if err := tx.Create(studentTeacher).Error; err != nil {
						return fmt.Errorf("关联教师失败: %v", err)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindUserByUsernameOrEmail(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != model.UserStatusActive {
		return nil, errors.New("账户已被禁用")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Username, user.Role.Name, user.RoleID, user.IsSuperAdmin)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %v", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %v", err)
	}

	expiresIn := time.Duration(24) * time.Hour
	if err := s.redisService.SaveUserSession(user.ID, accessToken, refreshToken, expiresIn); err != nil {
		return nil, fmt.Errorf("保存会话失败: %v", err)
	}

	onlineExpiresIn := time.Duration(30) * time.Minute
	if err := s.redisService.SetUserOnline(user.ID, onlineExpiresIn); err != nil {
		return nil, fmt.Errorf("设置在线状态失败: %v", err)
	}

	if err := s.redisService.CacheUserInfo(user.ID, user, time.Duration(1)*time.Hour); err != nil {
		return nil, fmt.Errorf("缓存用户信息失败: %v", err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

func (s *UserService) Logout(userID uint, accessToken string) error {
	if err := s.redisService.AddToBlacklist(accessToken, time.Duration(24)*time.Hour); err != nil {
		return fmt.Errorf("添加令牌到黑名单失败: %v", err)
	}

	if err := s.redisService.DeleteUserSession(userID); err != nil {
		return fmt.Errorf("删除用户会话失败: %v", err)
	}

	if err := s.redisService.SetUserOffline(userID); err != nil {
		return fmt.Errorf("设置用户离线失败: %v", err)
	}

	if err := s.redisService.DeleteCachedUserInfo(userID); err != nil {
		return fmt.Errorf("删除缓存用户信息失败: %v", err)
	}

	return nil
}

func (s *UserService) RefreshToken(refreshToken string) (*dto.RefreshTokenResponse, error) {
	claims, err := s.jwtService.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	userID := claims.UserID

	session, err := s.redisService.GetUserSession(userID)
	if err != nil {
		return nil, errors.New("会话不存在或已过期")
	}

	if session["refresh_token"] != refreshToken {
		return nil, errors.New("刷新令牌不匹配")
	}

	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status != model.UserStatusActive {
		return nil, errors.New("账户已被禁用")
	}

	newAccessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Username, user.Role.Name, user.RoleID, user.IsSuperAdmin)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %v", err)
	}

	if err := s.redisService.SaveUserSession(userID, newAccessToken, refreshToken, time.Duration(24)*time.Hour); err != nil {
		return nil, fmt.Errorf("更新会话失败: %v", err)
	}

	return &dto.RefreshTokenResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
	}, nil
}

func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
	cachedUser, err := s.redisService.GetCachedUserInfo(userID)
	if err == nil && cachedUser != nil {
		var user model.User
		if userData, err := json.Marshal(cachedUser); err == nil {
			if err := json.Unmarshal(userData, &user); err == nil {
				return &user, nil
			}
		}
	}

	if cachedUser == nil && err == nil {
		return nil, errors.New("用户不存在")
	}

	lockKey := fmt.Sprintf("user:%d", userID)
	expiresIn := time.Duration(5) * time.Second

	if s.redisService.AcquireLock(lockKey, expiresIn) {
		defer s.redisService.ReleaseLock(lockKey)

		if cachedUser, err := s.redisService.GetCachedUserInfo(userID); err == nil && cachedUser != nil {
			var user model.User
			if userData, err := json.Marshal(cachedUser); err == nil {
				if err := json.Unmarshal(userData, &user); err == nil {
					return &user, nil
				}
			}
		}

		user, err := s.repo.FindUserByID(userID)
		if err != nil {
			s.redisService.CacheEmptyUserInfo(userID, time.Duration(30)*time.Minute)
			return nil, err
		}

		s.redisService.CacheUserInfo(userID, *user, time.Duration(1)*time.Hour)
		return user, nil
	}

	time.Sleep(100 * time.Millisecond)
	for i := 0; i < 3; i++ {
		if cachedUser, err := s.redisService.GetCachedUserInfo(userID); err == nil && cachedUser != nil {
			var user model.User
			if userData, err := json.Marshal(cachedUser); err == nil {
				if err := json.Unmarshal(userData, &user); err == nil {
					return &user, nil
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
	}

	return s.repo.FindUserByID(userID)
}

func (s *UserService) UpdateUser(userID uint, req *dto.UpdateUserRequest) (*model.User, error) {
	if err := s.checkUserExists("", req.Email, req.Phone, userID); err != nil {
		return nil, err
	}

	var user *model.User

	err := s.repo.Transaction(func(tx *gorm.DB) error {
		var err error
		user, err = s.repo.FindUserByID(userID)
		if err != nil {
			return err
		}

		if req.Email != "" {
			user.Email = req.Email
		}

		if req.Phone != "" {
			user.Phone = req.Phone
		}

		if req.Nickname != "" {
			user.Nickname = req.Nickname
		}

		if req.Avatar != "" {
			user.Avatar = req.Avatar
		}

		if req.Status != nil {
			user.Status = *req.Status
		}

		if err := tx.Save(user).Error; err != nil {
			return fmt.Errorf("更新用户信息失败: %v", err)
		}

		role, err := s.repo.FindRoleByID(user.RoleID)
		if err != nil {
			return fmt.Errorf("获取角色信息失败: %v", err)
		}

		if role.Name == model.RoleStudent {
			student, err := s.repo.FindStudentByID(userID)
			if err != nil {
				return fmt.Errorf("获取学生信息失败: %v", err)
			}

			if req.ClassID > 0 && req.ClassID != student.ClassID {
				student.ClassID = req.ClassID
				if err := tx.Save(student).Error; err != nil {
					return fmt.Errorf("更新班级信息失败: %v", err)
				}
			}

			if len(req.TeacherIDs) > 0 {
				if err := s.repo.DeleteStudentTeachers(userID); err != nil {
					return fmt.Errorf("删除旧教师关联失败: %v", err)
				}

				for _, teacherID := range req.TeacherIDs {
					studentTeacher := &model.StudentTeacher{
						StudentID: userID,
						TeacherID: teacherID,
					}
					if err := tx.Create(studentTeacher).Error; err != nil {
						return fmt.Errorf("添加新教师关联失败: %v", err)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.redisService.CacheUserInfo(userID, *user, time.Duration(1)*time.Hour)
	user, _ = s.repo.FindUserByID(userID)

	return user, nil
}

func (s *UserService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	err := s.repo.Transaction(func(tx *gorm.DB) error {
		user, err := s.repo.FindUserByID(userID)
		if err != nil {
			return err
		}

		if err := user.ChangePassword(req.OldPassword, req.NewPassword); err != nil {
			return err
		}

		if err := tx.Save(user).Error; err != nil {
			return fmt.Errorf("更新密码失败: %v", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.redisService.DeleteCachedUserInfo(userID)
	return nil
}

func (s *UserService) GetUserList(req *dto.UserListRequest) (*dto.UserListResponse, error) {
	users, total, err := s.repo.GetUserList(req)
	if err != nil {
		return nil, err
	}

	return &dto.UserListResponse{
		PageInfo: utils.PageInfo{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
			Pages:    int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
		},
		Users: users,
	}, nil
}

func (s *UserService) ResetPassword(req *dto.ResetPasswordRequest) (string, error) {
	var tempPassword string

	err := s.repo.Transaction(func(tx *gorm.DB) error {
		user, err := s.repo.FindUserByUsername(req.Username)
		if err != nil {
			return err
		}

		tempPassword, err = user.ResetPassword()
		if err != nil {
			return fmt.Errorf("密码重置失败: %v", err)
		}

		if err := tx.Save(user).Error; err != nil {
			return fmt.Errorf("保存密码失败: %v", err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return tempPassword, nil
}

func (s *UserService) DeleteUser(userID uint) error {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return err
	}

	err = s.repo.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.DeleteStudentTeachers(userID); err != nil {
			return fmt.Errorf("删除学生-教师关系失败: %v", err)
		}

		if err := tx.Delete(user).Error; err != nil {
			return fmt.Errorf("删除用户失败: %v", err)
		}

		if err := s.repo.DeleteUserSessions(userID); err != nil {
			return fmt.Errorf("删除用户会话失败: %v", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.redisService.DeleteCachedUserInfo(userID)
	s.redisService.DeleteUserSession(userID)
	s.redisService.SetUserOffline(userID)

	return nil
}

func (s *UserService) ParseExcel(file multipart.File) ([]dto.StudentExcelData, error) {
	return s.parseExcelFile(file)
}

func (s *UserService) ParseCSV(file multipart.File) ([]dto.StudentExcelData, error) {
	return s.parseCSVFile(file)
}

func (s *UserService) ParseFile(file multipart.File, filename string) ([]dto.StudentExcelData, error) {
	if strings.HasSuffix(strings.ToLower(filename), ".xlsx") || strings.HasSuffix(strings.ToLower(filename), ".xls") {
		return s.parseExcelFile(file)
	} else if strings.HasSuffix(strings.ToLower(filename), ".csv") {
		return s.parseCSVFile(file)
	}
	return nil, errors.New("不支持的文件格式，仅支持.xlsx, .xls和.csv文件")
}

func (s *UserService) parseExcelFile(file multipart.File) ([]dto.StudentExcelData, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, errors.New("Excel文件中没有工作表")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取Excel行失败: %v", err)
	}

	return s.parseStudentRows(rows, "Excel")
}

func (s *UserService) parseCSVFile(file multipart.File) ([]dto.StudentExcelData, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取CSV文件失败: %v", err)
	}

	return s.parseStudentRows(rows, "CSV")
}

func (s *UserService) parseStudentRows(rows [][]string, sourceType string) ([]dto.StudentExcelData, error) {
	if len(rows) < 2 {
		return nil, fmt.Errorf("%s文件中至少需要包含表头和一行数据", sourceType)
	}

	head := rows[0]
	indices := s.resolveColumnIndices(head)

	if indices.Name == -1 {
		return nil, errors.New("无法识别姓名字段")
	}
	if indices.StudentID == -1 {
		return nil, errors.New("无法识别学号字段")
	}
	if indices.Email == -1 {
		return nil, errors.New("无法识别邮箱字段")
	}

	var students []dto.StudentExcelData
	for i, row := range rows[1:] {
		if len(row) == 0 || (len(row) > indices.Name && strings.TrimSpace(row[indices.Name]) == "") {
			continue
		}

		student := s.extractStudentData(row, indices, i+2)
		if student == nil {
			continue
		}

		students = append(students, *student)
	}

	return students, nil
}

type ColumnIndices struct {
	Name      int
	StudentID int
	Email     int
	Phone     int
	ClassName int
}

func (s *UserService) resolveColumnIndices(headers []string) ColumnIndices {
	var indices ColumnIndices
	indices.Name = -1
	indices.StudentID = -1
	indices.Email = -1
	indices.Phone = -1
	indices.ClassName = -1

	for i, field := range headers {
		field = strings.TrimSpace(field)
		fieldLower := strings.ToLower(field)

		if indices.Name == -1 && (strings.Contains(fieldLower, "姓名") || strings.Contains(fieldLower, "name")) {
			indices.Name = i
		}

		if indices.StudentID == -1 && (strings.Contains(fieldLower, "学号") || strings.Contains(fieldLower, "student") || strings.Contains(fieldLower, "id")) {
			indices.StudentID = i
		}

		if indices.Email == -1 && (strings.Contains(fieldLower, "邮箱") || strings.Contains(fieldLower, "email")) {
			indices.Email = i
		}

		if indices.Phone == -1 && (strings.Contains(fieldLower, "手机") || strings.Contains(fieldLower, "phone")) {
			indices.Phone = i
		}

		if indices.ClassName == -1 && (strings.Contains(fieldLower, "班级") || strings.Contains(fieldLower, "class")) {
			indices.ClassName = i
		}
	}

	return indices
}

func (s *UserService) extractStudentData(row []string, indices ColumnIndices, rowIndex int) *dto.StudentExcelData {
	name := ""
	if len(row) > indices.Name {
		name = strings.TrimSpace(row[indices.Name])
	}

	studentID := ""
	if len(row) > indices.StudentID {
		studentID = strings.TrimSpace(row[indices.StudentID])
	}

	email := ""
	if len(row) > indices.Email {
		email = strings.TrimSpace(row[indices.Email])
	}

	phone := ""
	if len(row) > indices.Phone && indices.Phone != -1 {
		phone = strings.TrimSpace(row[indices.Phone])
	}

	className := ""
	if len(row) > indices.ClassName && indices.ClassName != -1 {
		className = strings.TrimSpace(row[indices.ClassName])
	}

	if name == "" || studentID == "" || email == "" {
		return nil
	}

	var classID uint = 0
	if className != "" {
		class, _ := s.repo.FindClassByName(className)
		if class != nil {
			classID = class.ID
		}
	}

	return &dto.StudentExcelData{
		Index:      rowIndex,
		Name:       name,
		StudentID:  studentID,
		Email:      email,
		Phone:      phone,
		ClassName:  className,
		ClassID:    classID,
		TeacherIDs: []uint{},
	}
}

func (s *UserService) BatchRegister(students []dto.StudentExcelData) (*dto.BatchRegisterResponse, error) {
	var successCount int
	var failedCount int
	var results []dto.BatchRegisterResult
	var errs []map[string]interface{}

	for _, student := range students {
		// 批量注册学生 - 使用学号作为用户名
		req := &dto.RegisterRequest{
			Username:   student.StudentID,
			Password:   "123456",
			Email:      student.Email,
			Phone:      student.Phone,
			Nickname:   student.Name,
			RoleID:     3,
			ClassID:    student.ClassID,
			TeacherIDs: student.TeacherIDs,
		}

		_, err := s.Register(req)
		var result dto.BatchRegisterResult
		result.Index = student.Index
		result.Username = student.StudentID
		result.StudentID = student.StudentID

		if err != nil {
			failedCount++
			result.Status = "failed"
			result.Message = err.Error()
			errs = append(errs, map[string]interface{}{
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

	response := &dto.BatchRegisterResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Results:      results,
		Errors:       errs,
	}

	return response, nil
}

// checkUserExists 通用用户存在性检查
func (s *UserService) checkUserExists(username, email, phone string, excludeID uint) error {
	if username != "" {
		exists, err := s.repo.CheckFieldExists("username", username, excludeID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("用户名已存在")
		}
	}

	if email != "" {
		exists, err := s.repo.CheckFieldExists("email", email, excludeID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("邮箱已存在")
		}
	}

	if phone != "" {
		exists, err := s.repo.CheckFieldExists("phone", phone, excludeID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("手机号已存在")
		}
	}

	return nil
}
