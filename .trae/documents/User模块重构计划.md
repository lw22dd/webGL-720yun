# 教育场景用户系统类结构设计与实现

## 1. 设计思路

### 1.1 抽象基类设计
- 将现有`User`对象重构为抽象基类`User`
- 包含所有用户共有的基础属性
- 定义密码管理的抽象方法
- 确保提供必要的抽象方法，使子类必须实现特定于角色的功能

### 1.2 子类设计
- `Teacher`类：继承自`User`抽象类，管理班级和学生
- `Student`类：继承自`User`抽象类，关联到班级和教师

### 1.3 密码安全设计
- 保持现有的bcrypt加密机制
- 实现密码重置功能（生成临时密码）
- 实现密码修改功能（验证原密码）

### 1.4 兼容性设计
- 提供适配层，确保现有代码可以平滑过渡
- 保持数据库表结构的兼容性

## 2. 类结构设计

### 2.1 抽象基类 User
```go
// User 用户抽象基类
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
    Password  string         `gorm:"type:varchar(255);not null" json:"-"`
    Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
    Phone     string         `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
    Nickname  string         `gorm:"type:varchar(50)" json:"nickname"`
    Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
    RoleID    uint           `gorm:"not null" json:"role_id"`
    Role      Role           `gorm:"foreignKey:RoleID" json:"role"`
    Status    int            `gorm:"default:1" json:"status"` // 1:正常 0:禁用
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ResetPassword 重置密码抽象方法
func (u *User) ResetPassword() (string, error) {
    // 生成临时密码
    tempPassword := utils.GenerateRandomPassword(8)
    // 加密临时密码
    hashedPassword, err := utils.HashPassword(tempPassword)
    if err != nil {
        return "", err
    }
    // 更新密码
    u.Password = hashedPassword
    return tempPassword, nil
}

// ChangePassword 修改密码方法
func (u *User) ChangePassword(oldPassword, newPassword string) error {
    // 验证原密码
    if !utils.CheckPassword(oldPassword, u.Password) {
        return errors.New("原密码错误")
    }
    // 加密新密码
    hashedPassword, err := utils.HashPassword(newPassword)
    if err != nil {
        return err
    }
    // 更新密码
    u.Password = hashedPassword
    return nil
}
```

### 2.2 Teacher 类
```go
// Teacher 教师类
type Teacher struct {
    User
    Classes []Class `gorm:"foreignKey:TeacherID" json:"classes"`
}
```

### 2.3 Student 类
```go
// Student 学生类
type Student struct {
    User
    StudentID string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"student_id"` // 学号
    ClassID   uint      `gorm:"not null" json:"class_id"`
    Class     Class     `gorm:"foreignKey:ClassID" json:"class"`
    Teachers  []Teacher `gorm:"many2many:student_teachers;" json:"teachers"` // 关联多个教师
}
```

### 2.4 Class 类
```go
// Class 班级类
type Class struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
    Description string    `gorm:"type:varchar(255)" json:"description"`
    TeacherID   uint      `gorm:"not null" json:"teacher_id"`
    Teacher     User      `gorm:"foreignKey:TeacherID" json:"teacher"`
    Students    []Student `gorm:"foreignKey:ClassID" json:"students"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 2.5 关联表
```go
// StudentTeacher 学生-教师关联表
type StudentTeacher struct {
    StudentID uint `gorm:"primaryKey" json:"student_id"`
    TeacherID uint `gorm:"primaryKey" json:"teacher_id"`
}
```

## 3. 实现步骤

### 3.1 更新 user_config.go
- 重构 User 结构体为抽象基类
- 添加 Teacher、Student、Class 结构体
- 添加 StudentTeacher 关联表
- 添加表名设置方法

### 3.2 更新数据库模型关系
- 确保 GORM 能够正确处理继承关系
- 确保关联表能够正确创建

### 3.3 更新服务层
- 更新 UserService，支持新的类结构
- 确保密码重置和修改功能正常工作
- 提供必要的适配方法

### 3.4 添加密码工具方法
- 添加生成随机密码的方法
- 确保密码复杂度验证

### 3.5 提供适配层
- 确保现有代码可以平滑过渡到新的类结构
- 提供转换方法，将旧的 User 对象转换为新的子类对象

## 4. 密码安全实现

### 4.1 生成随机密码
```go
// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}
```

### 4.2 密码复杂度验证
```go
// ValidatePasswordComplexity 验证密码复杂度
func ValidatePasswordComplexity(password string) bool {
    // 至少8个字符，包含大小写字母、数字和特殊字符
    if len(password) < 8 {
        return false
    }
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+]`).MatchString(password)
    return hasLower && hasUpper && hasDigit && hasSpecial
}
```

## 5. 兼容性处理

### 5.1 类型转换方法
```go
// ToTeacher 将User转换为Teacher
func (u *User) ToTeacher() *Teacher {
    return &Teacher{User: *u}
}

// ToStudent 将User转换为Student
func (u *User) ToStudent() *Student {
    return &Student{User: *u}
}
```

### 5.2 服务层适配
- 在 UserService 中添加方法，根据 RoleID 返回对应的子类对象
- 确保现有API接口可以正常工作

## 6. 测试和验证

### 6.1 单元测试
- 测试密码重置功能
- 测试密码修改功能
- 测试类型转换功能

### 6.2 集成测试
- 测试用户注册功能，确保学生和教师对象正确创建
- 测试用户登录功能，确保返回正确的用户类型
- 测试班级关联功能，确保学生和教师正确关联到班级

## 7. 预期效果

- 实现了面向教育场景的用户系统类结构
- 密码处理符合安全最佳实践
- 支持学生-教师-班级的关联关系
- 提供了必要的适配层，确保系统平滑过渡
- 代码结构清晰，扩展性强，便于后续功能扩展

## 8. 注意事项

- 确保数据库迁移脚本正确，避免数据丢失
- 确保密码处理符合安全规范，避免安全漏洞
- 确保适配层正确实现，避免系统崩溃
- 确保测试覆盖所有关键功能，避免功能缺失

## 9. 后续优化

- 添加更多面向教育场景的功能，如学生成绩管理、课程管理等
- 优化密码策略，支持更复杂的密码规则
- 添加多因素认证功能，提高系统安全性
- 优化数据库查询，提高系统性能