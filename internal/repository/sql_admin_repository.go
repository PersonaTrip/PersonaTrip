package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"personatrip/internal/models"
)

// SQLAdminRepository 是使用原生SQL实现的管理员仓库
type SQLAdminRepository struct {
	db *sql.DB
}

// NewSQLAdminRepository 创建新的SQL管理员仓库
func NewSQLAdminRepository(db *sql.DB) AdminRepository {
	return &SQLAdminRepository{
		db: db,
	}
}

// Create 创建新的管理员
func (r *SQLAdminRepository) Create(ctx context.Context, admin *models.Admin) error {
	// 检查用户名是否已存在
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM admins WHERE username = ?", admin.Username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}

	// 插入管理员记录
	query := `
	INSERT INTO admins (username, email, password, role, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, admin.Username, admin.Email, admin.Password, admin.Role, now, now)
	if err != nil {
		return err
	}

	// 获取自增ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	admin.ID = uint(id)
	return nil
}

// Update 更新管理员信息
func (r *SQLAdminRepository) Update(ctx context.Context, admin *models.Admin) error {
	query := `
	UPDATE admins
	SET username = ?, email = ?, password = ?, role = ?, updated_at = ?
	WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, admin.Username, admin.Email, admin.Password, admin.Role, time.Now(), admin.ID)
	return err
}

// Delete 删除管理员
func (r *SQLAdminRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM admins WHERE id = ?", id)
	return err
}

// GetByID 根据ID获取管理员
func (r *SQLAdminRepository) GetByID(ctx context.Context, id uint) (*models.Admin, error) {
	query := `
	SELECT id, username, email, password, role, created_at, updated_at
	FROM admins
	WHERE id = ?
	`
	var admin models.Admin
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.Password, &admin.Role, &admin.CreatedAt, &admin.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}
	return &admin, nil
}

// GetByUsername 根据用户名获取管理员
func (r *SQLAdminRepository) GetByUsername(ctx context.Context, username string) (*models.Admin, error) {
	query := `
	SELECT id, username, email, password, role, created_at, updated_at
	FROM admins
	WHERE username = ?
	`
	var admin models.Admin
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.Password, &admin.Role, &admin.CreatedAt, &admin.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}
	return &admin, nil
}

// GetAll 获取所有管理员
func (r *SQLAdminRepository) GetAll(ctx context.Context) ([]models.Admin, error) {
	query := `
	SELECT id, username, email, password, role, created_at, updated_at
	FROM admins
	ORDER BY id
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []models.Admin
	for rows.Next() {
		var admin models.Admin
		err := rows.Scan(
			&admin.ID, &admin.Username, &admin.Email, &admin.Password, &admin.Role, &admin.CreatedAt, &admin.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return admins, nil
}
