package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"webGL-720yun/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init(config *config.DatabaseConfig) error {
	var err error
	var steps []string
	var successSteps []string
	var failedSteps []string

	// 设置GORM日志级别为Error，只输出错误信息
	logConfig := logger.Config{
		SlowThreshold:             time.Second,  // 慢SQL阈值
		LogLevel:                  logger.Error, // 日志级别
		IgnoreRecordNotFoundError: true,         // 忽略记录未找到错误
		Colorful:                  false,        // 禁用彩色输出
	}

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)

	// 步骤1: 连接数据库
	steps = append(steps, "连接数据库")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // 使用复数表名，与SQL脚本保持一致
		},
		Logger: logger.New(
			log.New(log.Writer(), "", log.LstdFlags), // io writer
			logConfig,
		),
	})

	if err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("连接数据库: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return fmt.Errorf("数据库连接失败: %v", err)
	}
	successSteps = append(successSteps, "连接数据库: 成功")

	// 步骤2: 设置连接池
	steps = append(steps, "设置连接池")
	sqlDB, err := DB.DB()
	if err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("设置连接池: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	successSteps = append(successSteps, "设置连接池: 成功")

	// 步骤3: 执行初始化SQL脚本
	steps = append(steps, "执行初始化SQL脚本")
	if err := executeInitSQL(); err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("执行初始化SQL脚本: 失败 - %v", err))
		printInitResult(steps, successSteps, failedSteps)
		return fmt.Errorf("执行初始化SQL脚本失败: %v", err)
	}
	successSteps = append(successSteps, "执行初始化SQL脚本: 成功")

	// 打印最终结果
	printInitResult(steps, successSteps, failedSteps)

	return nil
}

// executeInitSQL 执行初始化SQL脚本
func executeInitSQL() error {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %v", err)
	}

	// 尝试多种可能的路径，直到找到init.sql文件
	var sqlFile string
	var found bool

	// 可能的路径列表
	possiblePaths := []string{
		// 直接在当前目录下查找
		currentDir + "/pkg/database/init.sql",
		// 向上一级目录查找
		currentDir + "/../pkg/database/init.sql",
		// 向上两级目录查找
		currentDir + "/../../pkg/database/init.sql",
		// 从backend目录下查找
		currentDir + "/backend/pkg/database/init.sql",
	}

	// 遍历所有可能的路径，直到找到文件
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			sqlFile = path
			found = true
			break
		}
	}

	// 如果没有找到文件，返回错误
	if !found {
		return fmt.Errorf("init.sql文件不存在，尝试的路径: %v", possiblePaths)
	}

	// 检查文件是否存在
	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		return fmt.Errorf("init.sql文件不存在: %v, 尝试的路径: %s", err, sqlFile)
	}

	// 打印SQL文件路径，用于调试
	fmt.Printf("正在执行SQL脚本: %s\n", sqlFile)

	// 读取SQL文件内容
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("读取SQL文件失败: %v, 路径: %s", err, sqlFile)
	}

	// 使用database/sql包的Exec方法执行SQL文件
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 使用原生的sql.DB对象执行SQL脚本
	// 注意：mysql驱动默认不支持多语句执行，所以我们需要手动分割并执行
	// 但是对于CREATE TABLE这样的DDL语句，我们可以使用一个简单的方法：
	// 将SQL脚本按分号分割，然后逐个执行

	// 简单的分割方法，假设分号只出现在语句末尾
	sqlScript := string(content)
	lines := strings.Split(sqlScript, "\n")
	var currentStmt strings.Builder
	var statements []string

	for _, line := range lines {
		// 跳过注释行
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		currentStmt.WriteString(line)
		currentStmt.WriteString("\n")

		// 如果行以分号结尾，说明是一个完整的语句
		if strings.HasSuffix(line, ";") {
			statements = append(statements, currentStmt.String())
			currentStmt.Reset()
		}
	}

	// 添加最后一个语句（如果有的话）
	if currentStmt.Len() > 0 {
		statements = append(statements, currentStmt.String())
	}

	// 打印SQL语句数量
	fmt.Printf("SQL语句数量: %d\n", len(statements))

	// 逐个执行SQL语句
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// 打印正在执行的SQL语句索引和前50个字符
		shortStmt := stmt
		if len(shortStmt) > 50 {
			shortStmt = shortStmt[:50] + "..."
		}
		fmt.Printf("正在执行SQL语句 %d: %s\n", i+1, shortStmt)

		_, err = sqlDB.Exec(stmt)
		if err != nil {
			return fmt.Errorf("执行SQL语句失败: %v, 语句: %s", err, stmt)
		}
	}

	// 打印SQL脚本执行成功
	fmt.Printf("SQL脚本执行成功，共执行 %d 条语句\n", len(statements))

	return nil
}

// splitSQL 将SQL内容按分号分割成多个语句
func splitSQL(sql string) []string {
	var statements []string
	var currentStmt strings.Builder
	inString := false
	stringChar := byte(0)
	inComment := false
	commentType := 0 // 0: 不在注释中, 1: 单行注释(//), 2: 多行注释(/*)

	for i := 0; i < len(sql); {
		char := sql[i]

		// 处理注释
		if !inString {
			// 检查单行注释 //
			if i+1 < len(sql) && char == '/' && sql[i+1] == '/' && commentType == 0 {
				inComment = true
				commentType = 1
				i += 2
				continue
			}
			// 检查多行注释开始 /*
			if i+1 < len(sql) && char == '/' && sql[i+1] == '*' && commentType == 0 {
				inComment = true
				commentType = 2
				i += 2
				continue
			}
			// 检查多行注释结束 */
			if commentType == 2 && i+1 < len(sql) && char == '*' && sql[i+1] == '/' {
				inComment = false
				commentType = 0
				i += 2
				continue
			}
			// 检查单行注释结束
			if commentType == 1 && char == '\n' {
				inComment = false
				commentType = 0
			}
		}

		// 如果在注释中，跳过当前字符
		if inComment {
			i++
			continue
		}

		// 处理字符串开始/结束
		if (char == '\'' || char == '"' || char == '`') && !inString {
			inString = true
			stringChar = char
			currentStmt.WriteByte(char)
			i++
		} else if inString && char == stringChar && (i == 0 || sql[i-1] != '\\') {
			inString = false
			stringChar = 0
			currentStmt.WriteByte(char)
			i++
		} else if inString {
			// 处理转义字符
			if char == '\\' && i+1 < len(sql) {
				currentStmt.WriteByte(char)
				i++
				currentStmt.WriteByte(sql[i])
				i++
			} else {
				currentStmt.WriteByte(char)
				i++
			}
		} else {
			// 处理分号分隔符
			if char == ';' {
				// 添加当前语句
				stmt := strings.TrimSpace(currentStmt.String())
				if stmt != "" {
					statements = append(statements, stmt)
				}
				currentStmt.Reset()
				i++
			} else {
				currentStmt.WriteByte(char)
				i++
			}
		}
	}

	// 添加最后一个语句
	stmt := strings.TrimSpace(currentStmt.String())
	if stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}

// printInitResult 打印初始化结果
func printInitResult(steps []string, successSteps []string, failedSteps []string) {
	fmt.Println("\n=== 数据库初始化结果 ===")
	for _, step := range successSteps {
		fmt.Printf("✅ %s\n", step)
	}
	for _, step := range failedSteps {
		fmt.Printf("❌ %s\n", step)
	}
	fmt.Printf("\n总步骤: %d, 成功: %d, 失败: %d\n", len(steps), len(successSteps), len(failedSteps))
	fmt.Println("========================\n")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
