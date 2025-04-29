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

// DestinationInfo 目的地详细信息
type DestinationInfo struct {
	Name            string `json:"name" bson:"name"`
	Country         string `json:"country" bson:"country"`
	Language        string `json:"language" bson:"language"`
	Currency        string `json:"currency" bson:"currency"`
	TimeZone        string `json:"time_zone" bson:"time_zone"`
	BestTimeToVisit string `json:"best_time_to_visit" bson:"best_time_to_visit"`
}

// TravelInfo 旅行信息
type TravelInfo struct {
	VisaRequired         bool             `json:"visa_required" bson:"visa_required"`
	VisaTips             string           `json:"visa_tips" bson:"visa_tips"`
	PassportValidity     string           `json:"passport_validity" bson:"passport_validity"`
	VaccinationRequired  []string         `json:"vaccination_required" bson:"vaccination_required"`
	LocalCustoms         string           `json:"local_customs" bson:"local_customs"`
	EtiquetteTips        string           `json:"etiquette_tips" bson:"etiquette_tips"`
	SafetyTips           string           `json:"safety_tips" bson:"safety_tips"`
	HealthTips           string           `json:"health_tips" bson:"health_tips"`
	ElectricalSocketType string           `json:"electrical_socket_type" bson:"electrical_socket_type"`
	InternetAvailability string           `json:"internet_availability" bson:"internet_availability"`
	LanguagePhrases      []LanguagePhrase `json:"language_phrases" bson:"language_phrases"`
}

// LanguagePhrase 当地语言短语
type LanguagePhrase struct {
	Phrase        string `json:"phrase" bson:"phrase"`
	Pronunciation string `json:"pronunciation" bson:"pronunciation"`
	Meaning       string `json:"meaning" bson:"meaning"`
}

// WeatherForecast 天气预报
type WeatherForecast struct {
	ClimateOverview string          `json:"climate_overview" bson:"climate_overview"`
	DailyForecast   []DailyForecast `json:"daily_forecast" bson:"daily_forecast"`
}

// DailyForecast 每日天气预报
type DailyForecast struct {
	Date                string      `json:"date" bson:"date"`
	Temperature         Temperature `json:"temperature" bson:"temperature"`
	Conditions          string      `json:"conditions" bson:"conditions"`
	PrecipitationChance float64     `json:"precipitation_chance" bson:"precipitation_chance"`
	ClothingSuggestions []string    `json:"clothing_suggestions" bson:"clothing_suggestions"`
}

// Temperature 温度信息
type Temperature struct {
	Min  float64 `json:"min" bson:"min"`
	Max  float64 `json:"max" bson:"max"`
	Unit string  `json:"unit" bson:"unit"`
}

// DailyTemperature 当日温度信息
type DailyTemperature struct {
	Morning float64 `json:"morning" bson:"morning"`
	Day     float64 `json:"day" bson:"day"`
	Evening float64 `json:"evening" bson:"evening"`
	Unit    string  `json:"unit" bson:"unit"`
}

// PackingList 行李清单
type PackingList struct {
	Essentials  []string `json:"essentials" bson:"essentials"`
	Clothing    []string `json:"clothing" bson:"clothing"`
	Toiletries  []string `json:"toiletries" bson:"toiletries"`
	Electronics []string `json:"electronics" bson:"electronics"`
	Documents   []string `json:"documents" bson:"documents"`
	Other       []string `json:"other" bson:"other"`
}

// EmergencyContacts 紧急联系人
type EmergencyContacts struct {
	LocalEmergency string     `json:"local_emergency" bson:"local_emergency"`
	Police         string     `json:"police" bson:"police"`
	Ambulance      string     `json:"ambulance" bson:"ambulance"`
	Fire           string     `json:"fire" bson:"fire"`
	Embassy        string     `json:"embassy" bson:"embassy"`
	Hospitals      []Hospital `json:"hospitals" bson:"hospitals"`
}

// Hospital 医院信息
type Hospital struct {
	Name                    string `json:"name" bson:"name"`
	Address                 string `json:"address" bson:"address"`
	Phone                   string `json:"phone" bson:"phone"`
	HasEnglishSpeakingStaff bool   `json:"has_english_speaking_staff" bson:"has_english_speaking_staff"`
}

// TripPlan 旅行计划模型
type TripPlan struct {
	ID                     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID                 primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title                  string             `json:"title" bson:"title"`
	Destination            string             `json:"destination" bson:"destination"`
	DestinationInfo        DestinationInfo    `json:"destination_info" bson:"destination_info"`
	StartDate              string             `json:"start_date" bson:"start_date"`
	EndDate                string             `json:"end_date" bson:"end_date"`
	TravelInfo             TravelInfo         `json:"travel_info" bson:"travel_info"`
	WeatherForecast        WeatherForecast    `json:"weather_forecast" bson:"weather_forecast"`
	PackingList            PackingList        `json:"packing_list" bson:"packing_list"`
	EmergencyContacts      EmergencyContacts  `json:"emergency_contacts" bson:"emergency_contacts"`
	Days                   []TripDay          `json:"days" bson:"days"`
	Budget                 Budget             `json:"budget" bson:"budget"`
	LocalAttractions       []LocalAttraction  `json:"local_attractions" bson:"local_attractions"`
	LocalCuisine           []LocalCuisine     `json:"local_cuisine" bson:"local_cuisine"`
	Shopping               Shopping           `json:"shopping" bson:"shopping"`
	CulturalEvents         []CulturalEvent    `json:"cultural_events" bson:"cultural_events"`
	PracticalInfo          PracticalInfo      `json:"practical_information" bson:"practical_information"`
	Notes                  string             `json:"notes" bson:"notes"`
	SuggestedModifications string             `json:"suggested_modifications" bson:"suggested_modifications"`
	IsPublic               bool               `json:"is_public" bson:"is_public"`
	CreatedAt              time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at" bson:"updated_at"`
}

// TripDay 旅行日程
type TripDay struct {
	Day            int              `json:"day" bson:"day"`
	Date           string           `json:"date" bson:"date"`
	Weather        DayWeather       `json:"weather" bson:"weather"`
	Activities     []Activity       `json:"activities" bson:"activities"`
	Meals          []Meal           `json:"meals" bson:"meals"`
	Accommodation  Accommodation    `json:"accommodation" bson:"accommodation"`
	Transportation []Transportation `json:"transportation" bson:"transportation"`
	Tips           []string         `json:"tips" bson:"tips"`
}

// DayWeather 当天天气
type DayWeather struct {
	Temperature        DailyTemperature `json:"temperature" bson:"temperature"`
	Conditions         string           `json:"conditions" bson:"conditions"`
	ClothingSuggestion string           `json:"clothing_suggestion" bson:"clothing_suggestion"`
}

// Transportation 交通信息
type Transportation struct {
	Type             string  `json:"type" bson:"type"`
	From             string  `json:"from" bson:"from"`
	To               string  `json:"to" bson:"to"`
	DepartureTime    string  `json:"departure_time" bson:"departure_time"`
	ArrivalTime      string  `json:"arrival_time" bson:"arrival_time"`
	Cost             float64 `json:"cost" bson:"cost"`
	BookingReference string  `json:"booking_reference" bson:"booking_reference"`
	Notes            string  `json:"notes" bson:"notes"`
}

// Activity 活动项目
type Activity struct {
	Name            string   `json:"name" bson:"name"`
	Type            string   `json:"type" bson:"type"` // 景点、体验、交通等
	Location        Location `json:"location" bson:"location"`
	StartTime       string   `json:"start_time" bson:"start_time"`
	EndTime         string   `json:"end_time" bson:"end_time"`
	Description     string   `json:"description" bson:"description"`
	Cost            float64  `json:"cost" bson:"cost"`
	BookingRequired bool     `json:"booking_required" bson:"booking_required"`
	BookingTips     string   `json:"booking_tips" bson:"booking_tips"`
	CrowdLevel      string   `json:"crowd_level" bson:"crowd_level"`
	SuitableWeather string   `json:"suitable_weather" bson:"suitable_weather"`
	IndoorOutdoor   string   `json:"indoor_outdoor" bson:"indoor_outdoor"`
	Accessibility   string   `json:"accessibility" bson:"accessibility"`
	Rating          float64  `json:"rating" bson:"rating"`
	Photos          []string `json:"photos" bson:"photos"`
	Tips            []string `json:"tips" bson:"tips"`
	ImageURL        string   `json:"image_url" bson:"image_url"`
}

// Meal 餐饮
type Meal struct {
	Type            string   `json:"type" bson:"type"` // 早餐、午餐、晚餐、小吃
	Venue           string   `json:"venue" bson:"venue"`
	Cuisine         string   `json:"cuisine" bson:"cuisine"`
	Description     string   `json:"description" bson:"description"`
	Specialties     []string `json:"specialties" bson:"specialties"`
	DietaryOptions  []string `json:"dietary_options" bson:"dietary_options"`
	Address         string   `json:"address" bson:"address"`
	BookingRequired bool     `json:"booking_required" bson:"booking_required"`
	Cost            float64  `json:"cost" bson:"cost"`
	Tips            string   `json:"tips" bson:"tips"`
	Location        Location `json:"location" bson:"location"`
}

// Accommodation 住宿
type Accommodation struct {
	Name                  string   `json:"name" bson:"name"`
	Type                  string   `json:"type" bson:"type"` // 酒店、民宿、青旅等
	Address               string   `json:"address" bson:"address"`
	Description           string   `json:"description" bson:"description"`
	Amenities             []string `json:"amenities" bson:"amenities"`
	CheckIn               string   `json:"check_in" bson:"check_in"`
	CheckOut              string   `json:"check_out" bson:"check_out"`
	Cost                  float64  `json:"cost" bson:"cost"`
	BookingReference      string   `json:"booking_reference" bson:"booking_reference"`
	Contact               string   `json:"contact" bson:"contact"`
	NearestLandmarks      []string `json:"nearest_landmarks" bson:"nearest_landmarks"`
	TransportationOptions []string `json:"transportation_options" bson:"transportation_options"`
	Location              Location `json:"location" bson:"location"`
	ImageURL              string   `json:"image_url" bson:"image_url"`
}

// Location 地理位置
type Location struct {
	Name        string      `json:"name" bson:"name"`
	Address     string      `json:"address" bson:"address"`
	City        string      `json:"city" bson:"city"`
	Country     string      `json:"country" bson:"country"`
	Coordinates Coordinates `json:"coordinates" bson:"coordinates"`
}

// Coordinates 地理坐标
type Coordinates struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

// Budget 预算
type Budget struct {
	Currency       string        `json:"currency" bson:"currency"`
	ExchangeRate   string        `json:"exchange_rate" bson:"exchange_rate"`
	TotalEstimate  float64       `json:"total_estimate" bson:"total_estimate"`
	Accommodation  float64       `json:"accommodation" bson:"accommodation"`
	Transportation float64       `json:"transportation" bson:"transportation"`
	Food           float64       `json:"food" bson:"food"`
	Activities     float64       `json:"activities" bson:"activities"`
	Shopping       float64       `json:"shopping" bson:"shopping"`
	Other          float64       `json:"other" bson:"other"`
	DailyBreakdown []DailyBudget `json:"daily_breakdown" bson:"daily_breakdown"`
	PaymentTips    PaymentTips   `json:"payment_tips" bson:"payment_tips"`
}

// DailyBudget 每日预算
type DailyBudget struct {
	Day     int                 `json:"day" bson:"day"`
	Date    string              `json:"date" bson:"date"`
	Total   float64             `json:"total" bson:"total"`
	Details DailyExpenseDetails `json:"details" bson:"details"`
}

// DailyExpenseDetails 每日支出详情
type DailyExpenseDetails struct {
	Accommodation  float64 `json:"accommodation" bson:"accommodation"`
	Transportation float64 `json:"transportation" bson:"transportation"`
	Food           float64 `json:"food" bson:"food"`
	Activities     float64 `json:"activities" bson:"activities"`
	Other          float64 `json:"other" bson:"other"`
}

// PaymentTips 支付相关提示
type PaymentTips struct {
	CreditCardsAccepted       bool     `json:"credit_cards_accepted" bson:"credit_cards_accepted"`
	AtmAvailability           string   `json:"atm_availability" bson:"atm_availability"`
	TippingCulture            string   `json:"tipping_culture" bson:"tipping_culture"`
	RecommendedPaymentMethods []string `json:"recommended_payment_methods" bson:"recommended_payment_methods"`
}

// LocalAttraction 当地景点
type LocalAttraction struct {
	Name            string   `json:"name" bson:"name"`
	Category        string   `json:"category" bson:"category"`
	Description     string   `json:"description" bson:"description"`
	MustSee         bool     `json:"must_see" bson:"must_see"`
	Address         string   `json:"address" bson:"address"`
	OpeningHours    string   `json:"opening_hours" bson:"opening_hours"`
	Cost            float64  `json:"cost" bson:"cost"`
	TimeRequired    string   `json:"time_required" bson:"time_required"`
	BestTimeToVisit string   `json:"best_time_to_visit" bson:"best_time_to_visit"`
	Tips            []string `json:"tips" bson:"tips"`
}

// LocalCuisine 当地美食
type LocalCuisine struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	MustTry     bool     `json:"must_try" bson:"must_try"`
	WhereToFind []string `json:"where_to_find" bson:"where_to_find"`
	PriceRange  string   `json:"price_range" bson:"price_range"`
	Photos      []string `json:"photos" bson:"photos"`
}

// Shopping 购物信息
type Shopping struct {
	RecommendedItems []string       `json:"recommended_items" bson:"recommended_items"`
	MarketsAndMalls  []MarketOrMall `json:"markets_and_malls" bson:"markets_and_malls"`
	Souvenirs        []string       `json:"souvenirs" bson:"souvenirs"`
}

// MarketOrMall 市场或商场
type MarketOrMall struct {
	Name         string `json:"name" bson:"name"`
	Type         string `json:"type" bson:"type"`
	Address      string `json:"address" bson:"address"`
	Specialty    string `json:"specialty" bson:"specialty"`
	OpeningHours string `json:"opening_hours" bson:"opening_hours"`
}

// CulturalEvent 文化活动
type CulturalEvent struct {
	Name        string  `json:"name" bson:"name"`
	Date        string  `json:"date" bson:"date"`
	Description string  `json:"description" bson:"description"`
	Location    string  `json:"location" bson:"location"`
	Cost        float64 `json:"cost" bson:"cost"`
	Tips        string  `json:"tips" bson:"tips"`
}

// PracticalInfo 实用信息
type PracticalInfo struct {
	LocalTransportation LocalTransportation `json:"local_transportation" bson:"local_transportation"`
	Communication       Communication       `json:"communication" bson:"communication"`
}

// LocalTransportation 当地交通
type LocalTransportation struct {
	Options     []string `json:"options" bson:"options"`
	Recommended string   `json:"recommended" bson:"recommended"`
	Cost        string   `json:"cost" bson:"cost"`
	Passes      string   `json:"passes" bson:"passes"`
	Apps        []string `json:"apps" bson:"apps"`
}

// Communication 通信
type Communication struct {
	LocalSim         string   `json:"local_sim" bson:"local_sim"`
	WifiAvailability string   `json:"wifi_availability" bson:"wifi_availability"`
	UsefulApps       []string `json:"useful_apps" bson:"useful_apps"`
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
