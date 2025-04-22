package models

// ApiResponse 定义统一的API响应格式
type ApiResponse struct {
	Code    int         `json:"code"`           // 状态码
	Message string      `json:"message"`        // 消息
	Data    interface{} `json:"data,omitempty"` // 数据，当没有数据时省略
	List    interface{} `json:"list,omitempty"` // 列表数据，当没有列表数据时省略
	Bean    interface{} `json:"bean,omitempty"` // 对象数据，当没有对象数据时省略
}

// NewSuccessResponse 创建一个成功的响应
func NewSuccessResponse(message string) ApiResponse {
	return ApiResponse{
		Code:    200,
		Message: message,
	}
}

// NewErrorResponse 创建一个错误响应
func NewErrorResponse(code int, message string) ApiResponse {
	return ApiResponse{
		Code:    code,
		Message: message,
	}
}

// WithData 添加通用数据
func (r ApiResponse) WithData(data interface{}) ApiResponse {
	r.Data = data
	return r
}

// WithList 添加列表数据
func (r ApiResponse) WithList(list interface{}) ApiResponse {
	r.List = list
	return r
}

// WithBean 添加对象数据
func (r ApiResponse) WithBean(bean interface{}) ApiResponse {
	r.Bean = bean
	return r
}
