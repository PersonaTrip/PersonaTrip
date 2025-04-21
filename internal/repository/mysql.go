package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"personatrip/internal/models"
)

// MySQL 实现用户数据存储
type MySQL struct {
	DB *sql.DB // 公开的数据库连接，可以被其他仓库使用
}

// NewMySQL 创建新的MySQL存储实例
func NewMySQL(dsn string) (*MySQL, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// 检查连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 创建用户表（如果不存在）
	if err := createTables(db); err != nil {
		return nil, err
	}

	return &MySQL{DB: db}, nil
}

// 创建必要的表
func createTables(db *sql.DB) error {
	// 创建用户表
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id VARCHAR (100) NOT NULL UNIQUE,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(userTable)
	if err != nil {
		return err
	}

	// 创建管理员表
	adminTable := `
	CREATE TABLE IF NOT EXISTS admins (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		role VARCHAR(20) NOT NULL DEFAULT 'admin',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(adminTable)
	if err != nil {
		return err
	}

	// 创建模型配置表
	modelConfigTable := `
	CREATE TABLE IF NOT EXISTS model_configs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		model_type VARCHAR(50) NOT NULL,
		model_name VARCHAR(100) NOT NULL,
		api_key VARCHAR(255),
		base_url VARCHAR(255),
		is_active BOOLEAN DEFAULT FALSE,
		temperature FLOAT DEFAULT 0.7,
		max_tokens INT DEFAULT 2000,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(modelConfigTable)
	return err
}

// AutoMigrate 自动迁移数据库表
func (m *MySQL) AutoMigrate(models ...interface{}) error {
	// 我们已经在createTables函数中创建了表，这里只是一个占位符
	// 如果需要添加其他表或字段，可以在createTables函数中添加
	return nil
}

// Close 关闭数据库连接
func (m *MySQL) Close() error {
	return m.DB.Close()
}

// CreateUser 创建新用户
func (m *MySQL) CreateUser(ctx *gin.Context, user *models.UserMySQL) (*models.UserMySQL, error) {
	// 检查用户名是否已存在
	var count int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("username already exists")
	}

	// 检查邮箱是否已存在
	err = m.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 插入用户
	query := `
	INSERT INTO users (username, user_id, email, password, created_at, updated_at)
	VALUES (?, ?, ?, ?, NOW(), NOW())`

	result, err := m.DB.Exec(query, user.Username, user.UserID, user.Email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	user.ID = uint(id)
	user.Password = "" // 清除密码
	return user, nil
}

// GetUserByID 通过ID获取用户
func (m *MySQL) GetUserByID(ctx *gin.Context, userID string) (*models.UserMySQL, error) {
	var user models.UserMySQL

	row := m.DB.QueryRow(
		"SELECT id, user_id, username, email, created_at, updated_at FROM users WHERE user_id = ?",
		userID,
	)
	err := row.Scan(
		&user.ID,
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByUsername 通过用户名获取用户
func (m *MySQL) GetUserByUsername(ctx *gin.Context, username string) (*models.UserMySQL, error) {
	var user models.UserMySQL

	row := m.DB.QueryRow(
		"SELECT id, user_id, username, email, password, created_at, updated_at FROM users WHERE username = ?",
		username,
	)
	err := row.Scan(
		&user.ID,
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// CheckUserCredentials 检查用户凭据
func (m *MySQL) CheckUserCredentials(ctx *gin.Context, username, password string) (*models.UserMySQL, error) {
	// 获取用户（包括密码）
	user, err := m.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 清除密码
	user.Password = ""
	return user, nil
}
