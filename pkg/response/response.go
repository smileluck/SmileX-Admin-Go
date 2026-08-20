// Package response 统一 HTTP 响应
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const (
	CodeOK      = 0
	CodeErr     = 1
	CodeUnauthorized = 401
	CodeForbidden    = 403
)

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: CodeOK, Msg: "ok", Data: data})
}

func Fail(c *gin.Context, httpStatus, code int, msg string) {
	c.JSON(httpStatus, Body{Code: code, Msg: msg})
}

func BadRequest(c *gin.Context, msg string) { Fail(c, http.StatusBadRequest, CodeErr, msg) }
func Unauthorized(c *gin.Context, msg string) { Fail(c, http.StatusUnauthorized, CodeUnauthorized, msg) }
func Forbidden(c *gin.Context, msg string)    { Fail(c, http.StatusForbidden, CodeForbidden, msg) }
func NotFound(c *gin.Context, msg string)     { Fail(c, http.StatusNotFound, CodeErr, msg) }
func ServerError(c *gin.Context, msg string)  { Fail(c, http.StatusInternalServerError, CodeErr, msg) }
