package httputil

import (
	"net/http"

	"personatrip/internal/models"

	"github.com/gin-gonic/gin"
)

// ReturnSuccess 返回成功响应，带消息
func ReturnSuccess(c *gin.Context, message string) {
	c.JSON(http.StatusOK, models.NewSuccessResponse(message))
}

// ReturnSuccessWithData 返回成功响应，带通用数据
func ReturnSuccessWithData(c *gin.Context, message string, data interface{}) {
	response := models.NewSuccessResponse(message)

	// 如果data是nil，则创建一个空对象
	if data == nil {
		response.Data = map[string]interface{}{}
	} else {
		response.Data = data
	}

	c.JSON(http.StatusOK, response)
}

// ReturnSuccessWithBean 返回成功响应，带对象数据
func ReturnSuccessWithBean(c *gin.Context, message string, bean interface{}) {
	response := models.NewSuccessResponse(message)

	// 如果bean是nil，则创建一个空对象
	if bean == nil {
		response.Bean = map[string]interface{}{}
	} else {
		response.Bean = bean
	}

	c.JSON(http.StatusOK, response)
}

// ReturnSuccessWithList 返回成功响应，带列表数据
func ReturnSuccessWithList(c *gin.Context, message string, list interface{}) {
	// 确保空列表序列化为空数组而不是null
	response := models.NewSuccessResponse(message)

	// 如果list是nil，则创建一个空数组
	if list == nil {
		response.List = []interface{}{}
	} else {
		response.List = list
	}

	c.JSON(http.StatusOK, response)
}

// ReturnCreated 返回创建成功响应，带对象数据
func ReturnCreated(c *gin.Context, message string, bean interface{}) {
	response := models.NewSuccessResponse(message)

	// 如果bean是nil，则创建一个空对象
	if bean == nil {
		response.Bean = map[string]interface{}{}
	} else {
		response.Bean = bean
	}

	c.JSON(http.StatusCreated, response)
}

// ReturnError 返回错误响应
func ReturnError(c *gin.Context, code int, message string) {
	c.JSON(code, models.NewErrorResponse(code, message))
}

// ReturnBadRequest 返回400错误响应
func ReturnBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, message))
}

// ReturnUnauthorized 返回401错误响应
func ReturnUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, message))
}

// ReturnForbidden 返回403错误响应
func ReturnForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, models.NewErrorResponse(http.StatusForbidden, message))
}

// ReturnNotFound 返回404错误响应
func ReturnNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, models.NewErrorResponse(http.StatusNotFound, message))
}

// ReturnInternalError 返回500错误响应
func ReturnInternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, message))
}
