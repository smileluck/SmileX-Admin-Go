// 异步导出 handler 与传输层通用小工具。
package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	bizexport "github.com/smilex/smilex-admin-gin/internal/biz/export"
	bizfile "github.com/smilex/smilex-admin-gin/internal/biz/file"
	"github.com/smilex/smilex-admin-gin/internal/server/middleware"
	"github.com/smilex/smilex-admin-gin/pkg/response"
)

// ---- 异步导出 ----

// submitExport 提交导出任务：原始查询条件（query，剔除分页参数）快照进记录，
// worker 按同一套条件分批拉数，保证导出结果与列表页所见一致
func (s *HTTPServer) submitExport(c *gin.Context, biz string) {
	sub := middleware.Subject(c)
	params := c.Request.URL.Query()
	params.Del("page")
	params.Del("page_size")
	// 敏感明文导出 fail-closed：无对应权限码时剔除 reveal，导出仍按 export.mask 脱敏
	if params.Get("reveal") == "1" {
		if permCode, ok := bizexport.SensitivePermByBiz[biz]; !ok ||
			!s.auth.HasPermissionCode(c.Request.Context(), sub.UserID, permCode) {
			params.Del("reveal")
		}
	}
	vo, err := s.export.Submit(c.Request.Context(), biz, params, sub.UserID, sub.Username)
	if err != nil {
		switch {
		case errors.Is(err, bizexport.ErrQueueFull):
			response.FailI18n(c, http.StatusTooManyRequests, response.CodeErr, err)
		case errors.Is(err, bizexport.ErrUnsupportedBiz):
			response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
		default:
			response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
		}
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) listExports(c *gin.Context) {
	sub := middleware.Subject(c)
	if c.Query("recent") == "1" {
		vos, err := s.export.Recent(c.Request.Context(), sub.UserID, 5)
		if err != nil {
			response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
			return
		}
		response.OK(c, vos)
		return
	}
	page, size := s.pageParams(c)
	vos, pg, err := s.export.List(c.Request.Context(), sub.UserID, page, size)
	if err != nil {
		response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
		return
	}
	response.OK(c, listResult{List: vos, Page: pg})
}

func (s *HTTPServer) downloadExport(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	sub := middleware.Subject(c)
	d, err := s.export.ResolveDownload(c.Request.Context(), id, sub.UserID)
	if err != nil {
		s.exportErr(c, err)
		return
	}
	// 云存储：鉴权通过后 302 到短时效预签名 URL
	if d.URL != "" {
		c.Redirect(http.StatusFound, d.URL)
		return
	}
	// 本地存储：后端代理流式输出（强制 attachment + nosniff，CSV 不内联渲染）
	defer d.Body.Close()
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Content-Disposition", contentDisposition("attachment", d.Record.Name))
	if d.Record.Size > 0 {
		c.Header("Content-Length", strconv.FormatInt(d.Record.Size, 10))
	}
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, d.Body)
}

func (s *HTTPServer) deleteExport(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	sub := middleware.Subject(c)
	if err := s.export.Delete(c.Request.Context(), id, sub.UserID); err != nil {
		s.exportErr(c, err)
		return
	}
	response.OK(c, nil)
}

// exportErr 导出操作错误映射：不存在 404，越权 403，未完成 409，队列满 429，存储后端未配置 503，其余 500
func (s *HTTPServer) exportErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, bizexport.ErrNotFound):
		response.FailI18n(c, http.StatusNotFound, response.CodeErr, err)
	case errors.Is(err, bizexport.ErrNotOwner):
		response.FailI18n(c, http.StatusForbidden, response.CodeForbidden, err)
	case errors.Is(err, bizexport.ErrNotReady):
		response.FailI18n(c, http.StatusConflict, response.CodeErr, err)
	case errors.Is(err, bizexport.ErrQueueFull):
		response.FailI18n(c, http.StatusTooManyRequests, response.CodeErr, err)
	case errors.Is(err, bizfile.ErrDriverUnavailable):
		response.FailI18n(c, http.StatusServiceUnavailable, response.CodeErr, err)
	default:
		response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
	}
}

// contentDisposition 生成 Content-Disposition 头：ASCII 回退名 + RFC 5987 UTF-8 编码名
func contentDisposition(disposition, filename string) string {
	fallback := strings.Map(func(r rune) rune {
		if r < 32 || r > 126 || r == '"' || r == '\\' {
			return '_'
		}
		return r
	}, filename)
	return fmt.Sprintf(`%s; filename="%s"; filename*=UTF-8''%s`, disposition, fallback, url.PathEscape(filename))
}

// parseUnixParam 解析 unix 秒级时间戳查询参数（空/非法返回 false 表示不限）
func parseUnixParam(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, false
	}
	return time.Unix(n, 0), true
}

// revealParam 敏感数据明文查看开关：?reveal=1 且当前用户持有对应权限码时为 true；
// 无权限时静默忽略（仍脱敏返回，不暴露功能存在性）
func (s *HTTPServer) revealParam(c *gin.Context, permCode string) bool {
	if c.Query("reveal") != "1" {
		return false
	}
	sub := middleware.Subject(c)
	return s.auth.HasPermissionCode(c.Request.Context(), sub.UserID, permCode)
}

// pageParams 解析分页参数：单页上限取运行时参数 page.sizeMax（系统参数页可调，未配置回退 100）
func (s *HTTPServer) pageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	max := s.syscfg.IntDefault(c.Request.Context(), "page.sizeMax", 100)
	if size < 1 {
		size = 10
	}
	if size > max {
		size = max // 运行时上限夹取（page.sizeMax 可调）
	}
	return page, size
}

func idParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "invalid id")
		return 0, false
	}
	return uint(id), true
}

// truncate 按字节长度截断（UA 摘要存储用）
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
