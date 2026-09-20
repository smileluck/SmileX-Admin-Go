// 应用用户（管理端 CRUD 与独立认证 app-auth）handler。
package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	bizappuser "github.com/smilex/smilex-admin-gin/internal/biz/appuser"
	"github.com/smilex/smilex-admin-gin/internal/server/middleware"
	appusersvc "github.com/smilex/smilex-admin-gin/internal/service/appuser"
	"github.com/smilex/smilex-admin-gin/pkg/i18n"
	"github.com/smilex/smilex-admin-gin/pkg/response"
)

// ---- 应用用户（管理端） ----

func (s *HTTPServer) listAppUsers(c *gin.Context) {
	page, size := pageParams(c)
	q := bizappuser.Query{
		Keyword: strings.TrimSpace(c.Query("kw")),
		Phone:   strings.TrimSpace(c.Query("phone")),
	}
	if v := c.Query("status"); v != "" {
		if st, err := strconv.Atoi(v); err == nil {
			q.Status = &st
		}
	}
	if v := c.Query("tenant_id"); v != "" {
		if tid, err := strconv.ParseUint(v, 10, 64); err == nil && tid > 0 {
			t := uint(tid)
			q.TenantID = &t
		}
	}
	list, pg, err := s.appuser.List(c.Request.Context(), q, page, size)
	if err != nil {
		response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
		return
	}
	response.OK(c, listResult{List: list, Page: pg})
}

func (s *HTTPServer) createAppUser(c *gin.Context) {
	var req appusersvc.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	vo, err := s.appuser.Create(c.Request.Context(), req)
	if err != nil {
		s.appuserErr(c, err)
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) getAppUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	vo, err := s.appuser.Get(c.Request.Context(), id)
	if err != nil {
		s.appuserErr(c, err)
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) updateAppUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req appusersvc.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	if err := s.appuser.Update(c.Request.Context(), id, req); err != nil {
		s.appuserErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) deleteAppUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := s.appuser.Delete(c.Request.Context(), id); err != nil {
		s.appuserErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) resetAppUserPassword(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req appusersvc.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	if err := s.appuser.ResetPassword(c.Request.Context(), id, req); err != nil {
		s.appuserErr(c, err)
		return
	}
	response.OK(c, nil)
}

// appuserErr 应用用户操作错误映射：不存在 404，其余 400
func (s *HTTPServer) appuserErr(c *gin.Context, err error) {
	if errors.Is(err, bizappuser.ErrAppUserNotFound) {
		response.FailI18n(c, http.StatusNotFound, response.CodeErr, err)
		return
	}
	response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
}

// ---- 应用用户独立认证（app-auth） ----

// appLogin 应用用户登录：成功返回令牌对 + 用户信息；失败统一 401（防爆破计数由 LoginIPGuard 负责）
func (s *HTTPServer) appLogin(c *gin.Context) {
	var req appusersvc.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	vo, err := s.appuser.Login(c.Request.Context(), req)
	if err != nil {
		response.FailI18n(c, http.StatusUnauthorized, response.CodeUnauthorized, err)
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) appRefresh(c *gin.Context) {
	var req appusersvc.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	tp, err := s.appuser.Refresh(c.Request.Context(), req)
	if err != nil {
		response.FailI18n(c, http.StatusUnauthorized, response.CodeUnauthorized, err)
		return
	}
	response.OK(c, tp)
}

func (s *HTTPServer) appProfile(c *gin.Context) {
	sub := middleware.AppSubject(c)
	vo, err := s.appuser.Profile(c.Request.Context(), sub.UserID)
	if err != nil {
		response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
		return
	}
	response.OK(c, vo)
}

// appChangePassword 本人修改密码（校验旧密码）
func (s *HTTPServer) appChangePassword(c *gin.Context) {
	sub := middleware.AppSubject(c)
	var req appusersvc.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	if err := s.appuser.ChangePassword(c.Request.Context(), sub.Username, req); err != nil {
		if errors.Is(err, bizappuser.ErrBadCredentials) {
			response.BadRequest(c, i18n.T(c.Request.Context(), "profile.wrong_old_password"))
			return
		}
		response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
		return
	}
	response.OK(c, nil)
}
