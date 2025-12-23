package user

import (
	"testing"
	"webGL-720yun/config"
	"webGL-720yun/pkg/database"
	"webGL-720yun/pkg/services/jwt"
	"webGL-720yun/pkg/services/redis"
)

// 测试前的初始化
func TestMain(m *testing.M) {
	// 加载配置
	if err := config.Init(); err != nil {
		panic("配置加载失败: " + err.Error())
	}

	// 初始化数据库
	if err := database.Init(&config.Conf.Database); err != nil {
		panic("数据库初始化失败: " + err.Error())
	}

	// 运行测试
	m.Run()
}

// 测试用户注册功能
func TestUserService_Register(t *testing.T) {
	// 创建依赖服务
	jwtService := jwt.NewJWTService(&config.Conf.JWT)
	redisService := redis.NewRedisService(&config.Conf.Redis)

	// 创建用户服务
	userService := NewUserService(jwtService, redisService)

	// 测试用例1: 正常注册
	t.Run("正常注册", func(t *testing.T) {
		req := &RegisterRequest{
			Username:  "testuser",
			Password:  "123456",
			Email:     "test@example.com",
			Nickname:  "测试用户",
			RoleID:    2, // 学生角色
			StudentID: "20240001",
		}

		user, err := userService.Register(req)
		if err != nil {
			t.Errorf("注册失败: %v", err)
			return
		}

		if user == nil {
			t.Error("注册返回的用户为nil")
			return
		}

		if user.Username != req.Username {
			t.Errorf("用户名不匹配，期望: %s, 实际: %s", req.Username, user.Username)
		}

		if user.Email != req.Email {
			t.Errorf("邮箱不匹配，期望: %s, 实际: %s", req.Email, user.Email)
		}

		t.Logf("注册成功，用户ID: %d, 用户名: %s", user.ID, user.Username)
	})

	// 测试用例2: 用户名已存在
	t.Run("用户名已存在", func(t *testing.T) {
		req := &RegisterRequest{
			Username:  "testuser", // 与上一个测试用例相同的用户名
			Password:  "123456",
			Email:     "test2@example.com",
			Nickname:  "测试用户2",
			RoleID:    2,
			StudentID: "20240002",
		}

		user, err := userService.Register(req)
		if err == nil {
			t.Error("用户名已存在，注册应该失败")
			return
		}

		if user != nil {
			t.Error("注册失败时，返回的用户不应为nil")
			return
		}

		t.Logf("用户名已存在，注册失败: %v", err)
	})

	// 测试用例3: 邮箱已存在
	t.Run("邮箱已存在", func(t *testing.T) {
		req := &RegisterRequest{
			Username:  "testuser2",
			Password:  "123456",
			Email:     "test@example.com", // 与第一个测试用例相同的邮箱
			Nickname:  "测试用户3",
			RoleID:    2,
			StudentID: "20240003",
		}

		user, err := userService.Register(req)
		if err == nil {
			t.Error("邮箱已存在，注册应该失败")
			return
		}

		if user != nil {
			t.Error("注册失败时，返回的用户不应为nil")
			return
		}

		t.Logf("邮箱已存在，注册失败: %v", err)
	})
}

// 测试用户登录功能
func TestUserService_Login(t *testing.T) {
	// 创建依赖服务
	jwtService := jwt.NewJWTService(&config.Conf.JWT)
	redisService := redis.NewRedisService(&config.Conf.Redis)

	// 创建用户服务
	userService := NewUserService(jwtService, redisService)

	// 测试用例1: 正常登录
	t.Run("正常登录", func(t *testing.T) {
		req := &LoginRequest{
			Username: "testuser",
			Password: "123456",
		}

		response, err := userService.Login(req)
		if err != nil {
			t.Errorf("登录失败: %v", err)
			return
		}

		if response == nil {
			t.Error("登录返回的响应为nil")
			return
		}

		if response.AccessToken == "" {
			t.Error("登录返回的访问令牌为空")
			return
		}

		if response.RefreshToken == "" {
			t.Error("登录返回的刷新令牌为空")
			return
		}

		if response.User == nil {
			t.Error("登录返回的用户信息为nil")
			return
		}

		t.Logf("登录成功，用户ID: %d, 访问令牌: %s", response.User.ID, response.AccessToken)
	})

	// 测试用例2: 密码错误
	t.Run("密码错误", func(t *testing.T) {
		req := &LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}

		response, err := userService.Login(req)
		if err == nil {
			t.Error("密码错误，登录应该失败")
			return
		}

		if response != nil {
			t.Error("登录失败时，返回的响应不应为nil")
			return
		}

		t.Logf("密码错误，登录失败: %v", err)
	})

	// 测试用例3: 用户名不存在
	t.Run("用户名不存在", func(t *testing.T) {
		req := &LoginRequest{
			Username: "nonexistentuser",
			Password: "123456",
		}

		response, err := userService.Login(req)
		if err == nil {
			t.Error("用户名不存在，登录应该失败")
			return
		}

		if response != nil {
			t.Error("登录失败时，返回的响应不应为nil")
			return
		}

		t.Logf("用户名不存在，登录失败: %v", err)
	})
}

// 测试根据ID获取用户功能
func TestUserService_GetUserByID(t *testing.T) {
	// 创建依赖服务
	jwtService := jwt.NewJWTService(&config.Conf.JWT)
	redisService := redis.NewRedisService(&config.Conf.Redis)

	// 创建用户服务
	userService := NewUserService(jwtService, redisService)

	// 先注册一个用户，用于测试
	req := &RegisterRequest{
		Username:  "testuser4",
		Password:  "123456",
		Email:     "test4@example.com",
		Phone:     "13800138004",
		Nickname:  "测试用户4",
		RoleID:    2,
		StudentID: "20240004",
	}

	user, err := userService.Register(req)
	if err != nil {
		t.Fatalf("注册用户失败: %v", err)
	}

	// 测试用例1: 正常获取用户
	t.Run("正常获取用户", func(t *testing.T) {
		user, err := userService.GetUserByID(user.ID)
		if err != nil {
			t.Errorf("获取用户失败: %v", err)
			return
		}

		if user == nil {
			t.Error("获取的用户为nil")
			return
		}

		if user.Username != "testuser4" {
			t.Errorf("用户名不匹配，期望: testuser4, 实际: %s", user.Username)
			return
		}

		t.Logf("获取用户成功，用户ID: %d, 用户名: %s", user.ID, user.Username)
	})

	// 测试用例2: 获取不存在的用户
	t.Run("获取不存在的用户", func(t *testing.T) {
		user, err := userService.GetUserByID(9999)
		if err == nil {
			t.Error("获取不存在的用户，应该返回错误")
			return
		}

		if user != nil {
			t.Error("获取不存在的用户，返回的用户不应为nil")
			return
		}

		t.Logf("获取不存在的用户，返回错误: %v", err)
	})
}
