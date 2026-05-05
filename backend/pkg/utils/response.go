package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

// PageInfo 分页信息
type PageInfo struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	Pages    int   `json:"pages"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  "success",
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code: code,
		Msg:  message,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// InternalServerError 500错误
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
