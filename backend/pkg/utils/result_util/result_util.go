package result_util

import (
	"github.com/gin-gonic/gin"
)

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	// 如果为nil，则不返回data字段
	Data T `json:"data,omitempty"`
}

// HandleError 处理错误
func HandleError(c *gin.Context, code int, msg string) {
	c.JSON(code, Result[any]{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
