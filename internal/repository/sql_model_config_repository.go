package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"personatrip/internal/models"
)

// SQLModelConfigRepository 是使用原生SQL实现的模型配置仓库
type SQLModelConfigRepository struct {
	db *sql.DB
}

// NewSQLModelConfigRepository 创建新的SQL模型配置仓库
func NewSQLModelConfigRepository(db *sql.DB) ModelConfigRepository {
	return &SQLModelConfigRepository{
		db: db,
	}
}

// Create 创建新的模型配置
func (r *SQLModelConfigRepository) Create(ctx context.Context, config *models.ModelConfig) error {
	query := `
	INSERT INTO model_configs (
		name, model_type, model_name, api_key, base_url, 
		is_active, temperature, max_tokens, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.ExecContext(
		ctx, query,
		config.Name, config.ModelType, config.ModelName, config.APIKey, config.BaseURL,
		config.IsActive, config.Temperature, config.MaxTokens, now, now,
	)
	if err != nil {
		return err
	}

	// 获取自增ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	config.ID = uint(id)

	// 如果设置为活跃，则将其他配置设置为非活跃
	if config.IsActive {
		return r.setActiveOnly(ctx, config.ID)
	}
	return nil
}

// Update 更新模型配置
func (r *SQLModelConfigRepository) Update(ctx context.Context, config *models.ModelConfig) error {
	query := `
	UPDATE model_configs
	SET name = ?, model_type = ?, model_name = ?, api_key = ?, base_url = ?,
		is_active = ?, temperature = ?, max_tokens = ?, updated_at = ?
	WHERE id = ?
	`
	_, err := r.db.ExecContext(
		ctx, query,
		config.Name, config.ModelType, config.ModelName, config.APIKey, config.BaseURL,
		config.IsActive, config.Temperature, config.MaxTokens, time.Now(), config.ID,
	)
	if err != nil {
		return err
	}

	// 如果设置为活跃，则将其他配置设置为非活跃
	if config.IsActive {
		return r.setActiveOnly(ctx, config.ID)
	}
	return nil
}

// Delete 删除模型配置
func (r *SQLModelConfigRepository) Delete(ctx context.Context, id uint) error {
	// 检查是否为活跃配置
	var isActive bool
	err := r.db.QueryRowContext(ctx, "SELECT is_active FROM model_configs WHERE id = ?", id).Scan(&isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("model config not found")
		}
		return err
	}

	// 删除配置
	_, err = r.db.ExecContext(ctx, "DELETE FROM model_configs WHERE id = ?", id)
	return err
}

// GetByID 根据ID获取模型配置
func (r *SQLModelConfigRepository) GetByID(ctx context.Context, id uint) (*models.ModelConfig, error) {
	query := `
	SELECT id, name, model_type, model_name, api_key, base_url, 
		is_active, temperature, max_tokens, created_at, updated_at
	FROM model_configs
	WHERE id = ?
	`
	var config models.ModelConfig
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID, &config.Name, &config.ModelType, &config.ModelName, &config.APIKey, &config.BaseURL,
		&config.IsActive, &config.Temperature, &config.MaxTokens, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("model config not found")
		}
		return nil, err
	}
	return &config, nil
}

// GetAll 获取所有模型配置
func (r *SQLModelConfigRepository) GetAll(ctx context.Context) ([]models.ModelConfig, error) {
	query := `
	SELECT id, name, model_type, model_name, api_key, base_url, 
		is_active, temperature, max_tokens, created_at, updated_at
	FROM model_configs
	ORDER BY id
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []models.ModelConfig
	for rows.Next() {
		var config models.ModelConfig
		err := rows.Scan(
			&config.ID, &config.Name, &config.ModelType, &config.ModelName, &config.APIKey, &config.BaseURL,
			&config.IsActive, &config.Temperature, &config.MaxTokens, &config.CreatedAt, &config.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return configs, nil
}

// GetActive 获取当前活跃的模型配置
func (r *SQLModelConfigRepository) GetActive(ctx context.Context) (*models.ModelConfig, error) {
	query := `
	SELECT id, name, model_type, model_name, api_key, base_url, 
		is_active, temperature, max_tokens, created_at, updated_at
	FROM model_configs
	WHERE is_active = TRUE
	LIMIT 1
	`
	var config models.ModelConfig
	err := r.db.QueryRowContext(ctx, query).Scan(
		&config.ID, &config.Name, &config.ModelType, &config.ModelName, &config.APIKey, &config.BaseURL,
		&config.IsActive, &config.Temperature, &config.MaxTokens, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no active model config found")
		}
		return nil, err
	}
	return &config, nil
}

// SetActive 设置指定ID的配置为活跃
func (r *SQLModelConfigRepository) SetActive(ctx context.Context, id uint) error {
	// 检查配置是否存在
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM model_configs WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("model config not found")
	}

	return r.setActiveOnly(ctx, id)
}

// setActiveOnly 将指定ID的配置设置为唯一活跃配置
func (r *SQLModelConfigRepository) setActiveOnly(ctx context.Context, id uint) error {
	// 开始事务
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 将所有配置设置为非活跃
	_, err = tx.ExecContext(ctx, "UPDATE model_configs SET is_active = FALSE")
	if err != nil {
		return err
	}

	// 将指定ID的配置设置为活跃
	_, err = tx.ExecContext(ctx, "UPDATE model_configs SET is_active = TRUE WHERE id = ?", id)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}
