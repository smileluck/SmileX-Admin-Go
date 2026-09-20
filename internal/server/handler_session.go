// 在线会话管理 handler。
package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	bizsession "github.com/smilex/smilex-admin-gin/internal/biz/session"
	"github.com/smilex/smilex-admin-gin/internal/biz/user"
	"github.com/smilex/smilex-admin-gin/internal/server/middleware"
	"github.com/smilex/smilex-admin-gin/pkg/i18n"
	"github.com/smilex/smilex-admin-gin/pkg/response"
)

// ---- helpers ----

// ---- 在线用户 ----

func (s *HTTPServer) listOnlineUsers(c *gin.Context) {
	page, size := pageParams(c)
	q := bizsession.Query{Username: c.Query("username"), Device: c.Query("device")}
	currentSid := ""
	if sub := middleware.Subject(c); sub != nil {
		currentSid = sub.SessionID
	}
	vos, pg, err := s.session.List(c.Request.Context(), q, page, size, currentSid)
	if err != nil {
		response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
		return
	}
	response.OK(c, listResult{List: vos, Page: pg})
}

func (s *HTTPServer) kickOnlineSession(c *gin.Context) {
	sid := c.Param("sid")
	if sid == "" {
		response.BadRequest(c, "invalid sid")
		return
	}
	operatorID := uint(0)
	if sub := middleware.Subject(c); sub != nil {
		operatorID = sub.UserID
	}
	if err := s.session.Kick(c.Request.Context(), sid, operatorID); err != nil {
		s.sessionErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) kickUserSessions(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	operatorID := uint(0)
	if sub := middleware.Subject(c); sub != nil {
		operatorID = sub.UserID
	}
	if _, err := s.session.KickUser(c.Request.Context(), id, operatorID); err != nil {
		s.sessionErr(c, err)
		return
	}
	response.OK(c, nil)
}

// sessionErr 会话操作错误映射：超管保护返回 403，会话不存在返回 404，其余返回 400
func (s *HTTPServer) sessionErr(c *gin.Context, err error) {
	if errors.Is(err, user.ErrSuperAdminProtected) {
		response.FailI18n(c, http.StatusForbidden, response.CodeForbidden, err)
		return
	}
	if errors.Is(err, bizsession.ErrSessionNotFound) {
		response.NotFound(c, i18n.T(c.Request.Context(), "session.not_found"))
		return
	}
	response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
}
