package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"personatrip/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// MySQL 实现用户数据存储
type MySQL struct {
	db *sql.DB
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

	return &MySQL{db: db}, nil
}

// 创建必要的表
func createTables(db *sql.DB) error {
	// 创建用户表
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(userTable)
	return err
}

// Close 关闭数据库连接
func (m *MySQL) Close() error {
	return m.db.Close()
}

// CreateUser 创建新用户
func (m *MySQL) CreateUser(ctx *gin.Context, user *models.UserMySQL) (*models.UserMySQL, error) {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 插入用户
	query := `
	INSERT INTO users (username, email, password, created_at, updated_at)
	VALUES (?, ?, ?, NOW(), NOW())`

	result, err := m.db.Exec(query, user.Username, user.Email, string(hashedPassword))
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
func (m *MySQL) GetUserByID(ctx *gin.Context, id uint) (*models.UserMySQL, error) {
	query := `
	SELECT id, username, email, created_at, updated_at
	FROM users
	WHERE id = ?`

	var user models.UserMySQL
	err := m.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByUsername 通过用户名获取用户
func (m *MySQL) GetUserByUsername(ctx *gin.Context, username string) (*models.UserMySQL, error) {
	query := `
	SELECT id, username, email, password, created_at, updated_at
	FROM users
	WHERE username = ?`

	var user models.UserMySQL
	err := m.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
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
