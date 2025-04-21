package repository

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"personatrip/internal/models"
)

// MemoryStore 提供基于内存的存储实现
type MemoryStore struct {
	mu        sync.RWMutex
	users     map[primitive.ObjectID]*models.User
	tripPlans map[primitive.ObjectID]*models.TripPlan
	userPlans map[primitive.ObjectID][]primitive.ObjectID // 用户ID到行程ID的映射
}

// NewMemoryStore 创建新的内存存储
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:     make(map[primitive.ObjectID]*models.User),
		tripPlans: make(map[primitive.ObjectID]*models.TripPlan),
		userPlans: make(map[primitive.ObjectID][]primitive.ObjectID),
	}
}

// CreateUser 创建新用户
func (m *MemoryStore) CreateUser(ctx *gin.Context, user *models.User) (*models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	m.users[user.ID] = user
	return user, nil
}

// GetUserByID 通过ID获取用户
func (m *MemoryStore) GetUserByID(ctx *gin.Context, id primitive.ObjectID) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

// GetUserByEmail 通过邮箱获取用户
func (m *MemoryStore) GetUserByEmail(ctx *gin.Context, email string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrNotFound
}

// UpdateUser 更新用户信息
func (m *MemoryStore) UpdateUser(ctx *gin.Context, user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[user.ID]; !ok {
		return ErrNotFound
	}

	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

// CreateTripPlan 创建旅行计划
func (m *MemoryStore) CreateTripPlan(ctx *gin.Context, plan *models.TripPlan) (*models.TripPlan, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	plan.ID = primitive.NewObjectID()
	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()

	m.tripPlans[plan.ID] = plan

	// 添加到用户的计划列表
	m.userPlans[plan.UserID] = append(m.userPlans[plan.UserID], plan.ID)

	return plan, nil
}

// GetTripPlanByID 通过ID获取旅行计划
func (m *MemoryStore) GetTripPlanByID(ctx *gin.Context, id primitive.ObjectID) (*models.TripPlan, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plan, ok := m.tripPlans[id]
	if !ok {
		return nil, ErrNotFound
	}
	return plan, nil
}

// GetTripPlansByUserID 获取用户的所有旅行计划
func (m *MemoryStore) GetTripPlansByUserID(ctx *gin.Context, userID primitive.ObjectID) ([]*models.TripPlan, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	planIDs, ok := m.userPlans[userID]
	if !ok {
		return []*models.TripPlan{}, nil
	}

	plans := make([]*models.TripPlan, 0, len(planIDs))
	for _, id := range planIDs {
		if plan, ok := m.tripPlans[id]; ok {
			plans = append(plans, plan)
		}
	}

	return plans, nil
}

// UpdateTripPlan 更新旅行计划
func (m *MemoryStore) UpdateTripPlan(ctx *gin.Context, plan *models.TripPlan) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tripPlans[plan.ID]; !ok {
		return ErrNotFound
	}

	plan.UpdatedAt = time.Now()
	m.tripPlans[plan.ID] = plan
	return nil
}

// DeleteTripPlan 删除旅行计划
func (m *MemoryStore) DeleteTripPlan(ctx *gin.Context, id primitive.ObjectID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plan, ok := m.tripPlans[id]
	if !ok {
		return ErrNotFound
	}

	// 从用户的计划列表中移除
	userID := plan.UserID
	planIDs := m.userPlans[userID]
	for i, planID := range planIDs {
		if planID == id {
			m.userPlans[userID] = append(planIDs[:i], planIDs[i+1:]...)
			break
		}
	}

	// 删除计划
	delete(m.tripPlans, id)
	return nil
}
