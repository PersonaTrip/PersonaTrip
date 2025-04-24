package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"`
	Preferences  UserPreferences    `json:"preferences" bson:"preferences"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserPreferences 用户旅行偏好
type UserPreferences struct {
	TravelStyle     []string `json:"travel_style" bson:"travel_style"`         // 旅行风格: 文化、自然、美食、冒险等
	Budget          string   `json:"budget" bson:"budget"`                     // 预算等级: 经济、中等、豪华
	Accommodation   []string `json:"accommodation" bson:"accommodation"`       // 住宿偏好: 酒店、民宿、露营等
	Transportation  []string `json:"transportation" bson:"transportation"`     // 交通偏好: 公共交通、自驾、步行等
	Activities      []string `json:"activities" bson:"activities"`             // 活动偏好: 博物馆、徒步、购物等
	FoodPreferences []string `json:"food_preferences" bson:"food_preferences"` // 饮食偏好: 当地美食、素食、特定菜系等
}

// TripPlan 旅行计划模型
type TripPlan struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Destination string             `json:"destination" bson:"destination"`
	StartDate   string             `json:"start_date" bson:"start_date"`
	EndDate     string             `json:"end_date" bson:"end_date"`
	Days        []TripDay          `json:"days" bson:"days"`
	Budget      Budget             `json:"budget" bson:"budget"`
	Notes       string             `json:"notes" bson:"notes"`
	IsPublic    bool               `json:"is_public" bson:"is_public"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// TripDay 旅行日程
type TripDay struct {
	Day           int           `json:"day" bson:"day"`
	Date          string        `json:"date" bson:"date"`
	Activities    []Activity    `json:"activities" bson:"activities"`
	Meals         []Meal        `json:"meals" bson:"meals"`
	Accommodation Accommodation `json:"accommodation" bson:"accommodation"`
}

// Activity 活动项目
type Activity struct {
	Name        string   `json:"name" bson:"name"`
	Type        string   `json:"type" bson:"type"` // 景点、体验、交通等
	Location    Location `json:"location" bson:"location"`
	StartTime   string   `json:"start_time" bson:"start_time"`
	EndTime     string   `json:"end_time" bson:"end_time"`
	Description string   `json:"description" bson:"description"`
	Cost        float64  `json:"cost" bson:"cost"`
	ImageURL    string   `json:"image_url" bson:"image_url"`
}

// Meal 餐饮
type Meal struct {
	Type        string   `json:"type" bson:"type"` // 早餐、午餐、晚餐、小吃
	Venue       string   `json:"venue" bson:"venue"`
	Location    Location `json:"location" bson:"location"`
	Description string   `json:"description" bson:"description"`
	Cost        float64  `json:"cost" bson:"cost"`
}

// Accommodation 住宿
type Accommodation struct {
	Name        string   `json:"name" bson:"name"`
	Type        string   `json:"type" bson:"type"` // 酒店、民宿、青旅等
	Location    Location `json:"location" bson:"location"`
	Description string   `json:"description" bson:"description"`
	Cost        float64  `json:"cost" bson:"cost"`
	ImageURL    string   `json:"image_url" bson:"image_url"`
}

// Location 地理位置
type Location struct {
	Name      string  `json:"name" bson:"name"`
	Address   string  `json:"address" bson:"address"`
	City      string  `json:"city" bson:"city"`
	Country   string  `json:"country" bson:"country"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

// Budget 预算
type Budget struct {
	Currency       string  `json:"currency" bson:"currency"`
	TotalEstimate  float64 `json:"total_estimate" bson:"total_estimate"`
	Accommodation  float64 `json:"accommodation" bson:"accommodation"`
	Transportation float64 `json:"transportation" bson:"transportation"`
	Food           float64 `json:"food" bson:"food"`
	Activities     float64 `json:"activities" bson:"activities"`
	Other          float64 `json:"other" bson:"other"`
}

// PlanRequest 创建旅行计划的请求
type PlanRequest struct {
	Destination     string    `json:"destination" binding:"required"`
	StartDate       time.Time `json:"start_date" binding:"required"`
	EndDate         time.Time `json:"end_date" binding:"required"`
	Budget          string    `json:"budget"`           // 预算等级: 经济、中等、豪华
	TravelStyle     []string  `json:"travel_style"`     // 旅行风格
	Accommodation   []string  `json:"accommodation"`    // 住宿偏好
	Transportation  []string  `json:"transportation"`   // 交通偏好
	Activities      []string  `json:"activities"`       // 活动偏好
	FoodPreferences []string  `json:"food_preferences"` // 饮食偏好
	SpecialRequests string    `json:"special_requests"` // 特殊要求
}
