// 用户管理 handler。
package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smilex/smilex-admin-gin/internal/biz/user"
	"github.com/smilex/smilex-admin-gin/internal/server/middleware"
	usersvc "github.com/smilex/smilex-admin-gin/internal/service/user"
	"github.com/smilex/smilex-admin-gin/pkg/i18n"
	"github.com/smilex/smilex-admin-gin/pkg/response"
)

// ---- 用户 ----

type listResult struct {
	List interface{} `json:"list"`
	Page interface{} `json:"page"`
}

func (s *HTTPServer) listUsers(c *gin.Context) {
	page, size := s.pageParams(c)
	q := user.Query{Username: c.Query("username")}
	if v := c.Query("status"); v != "" {
		if st, err := strconv.Atoi(v); err == nil {
			q.Status = &st
		}
	}
	users, pg, err := s.user.List(c.Request.Context(), q, page, size, s.revealParam(c, "user:viewSensitive"))
	if err != nil {
		response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
		return
	}
	response.OK(c, listResult{List: users, Page: pg})
}

func (s *HTTPServer) createUser(c *gin.Context) {
	var req usersvc.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	vo, err := s.user.Create(user.WithOperator(c.Request.Context(), middleware.Subject(c).UserID), req)
	if err != nil {
		s.userErr(c, err)
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) getUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	vo, err := s.user.Get(c.Request.Context(), id, s.revealParam(c, "user:viewSensitive"))
	if err != nil {
		response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) updateUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req usersvc.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	if err := s.user.Update(user.WithOperator(c.Request.Context(), middleware.Subject(c).UserID), id, req); err != nil {
		s.userErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) deleteUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := s.user.Delete(c.Request.Context(), id); err != nil {
		s.userErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) setUserRoles(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req usersvc.SetRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	if err := s.user.SetRoles(user.WithOperator(c.Request.Context(), middleware.Subject(c).UserID), id, req); err != nil {
		s.userErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) resetUserPassword(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req usersvc.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	if err := s.user.ResetPassword(user.WithOperator(c.Request.Context(), middleware.Subject(c).UserID), id, req); err != nil {
		s.userErr(c, err)
		return
	}
	response.OK(c, nil)
}
func (s *HTTPServer) userErr(c *gin.Context, err error) {
	if errors.Is(err, user.ErrSuperAdminProtected) || errors.Is(err, user.ErrDeleteSuperAdmin) ||
		errors.Is(err, user.ErrAssignSuperRole) {
		response.FailI18n(c, http.StatusForbidden, response.CodeForbidden, err)
		return
	}
	response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
}
