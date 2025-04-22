package repository

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"personatrip/internal/models"
)

// MongoDB 实现数据存储
type MongoDB struct {
	client    *mongo.Client
	database  *mongo.Database
	users     *mongo.Collection
	tripPlans *mongo.Collection
}

// NewMongoDB 创建新的MongoDB存储实例
func NewMongoDB(uri string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 获取数据库和集合
	database := client.Database("personatrip")
	users := database.Collection("users")
	tripPlans := database.Collection("trip_plans")

	return &MongoDB{
		client:    client,
		database:  database,
		users:     users,
		tripPlans: tripPlans,
	}, nil
}

// Close 关闭数据库连接
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// CreateUser 创建新用户
func (m *MongoDB) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := m.users.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID 通过ID获取用户
func (m *MongoDB) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := m.users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail 通过邮箱获取用户
func (m *MongoDB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := m.users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (m *MongoDB) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	_, err := m.users.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	return err
}

// CreateTripPlan 创建旅行计划
func (m *MongoDB) CreateTripPlan(ctx *gin.Context, plan *models.TripPlan) (*models.TripPlan, error) {
	plan.ID = primitive.NewObjectID()
	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()

	_, err := m.tripPlans.InsertOne(ctx, plan)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

// GetTripPlanByID 通过ID获取旅行计划
func (m *MongoDB) GetTripPlanByID(ctx *gin.Context, id primitive.ObjectID) (*models.TripPlan, error) {
	var plan models.TripPlan
	err := m.tripPlans.FindOne(ctx, bson.M{"_id": id}).Decode(&plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

// GetTripPlansByUserID 获取用户的所有旅行计划
func (m *MongoDB) GetTripPlansByUserID(ctx *gin.Context, userID primitive.ObjectID) ([]*models.TripPlan, error) {
	cursor, err := m.tripPlans.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plans []*models.TripPlan
	if err = cursor.All(ctx, &plans); err != nil {
		return nil, err
	}

	return plans, nil
}

// UpdateTripPlan 更新旅行计划
func (m *MongoDB) UpdateTripPlan(ctx *gin.Context, plan *models.TripPlan) error {
	plan.UpdatedAt = time.Now()

	_, err := m.tripPlans.ReplaceOne(ctx, bson.M{"_id": plan.ID}, plan)
	return err
}

// DeleteTripPlan 删除旅行计划
func (m *MongoDB) DeleteTripPlan(ctx *gin.Context, id primitive.ObjectID) error {
	_, err := m.tripPlans.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
