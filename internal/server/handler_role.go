// 角色管理 handler。
package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smilex/smilex-admin-gin/internal/biz/role"
	"github.com/smilex/smilex-admin-gin/internal/server/middleware"
	rolesvc "github.com/smilex/smilex-admin-gin/internal/service/role"
	"github.com/smilex/smilex-admin-gin/pkg/i18n"
	"github.com/smilex/smilex-admin-gin/pkg/response"
)

// ---- 角色 ----

func (s *HTTPServer) listRoles(c *gin.Context) {
	page, size := pageParams(c)
	roles, pg, err := s.role.List(c.Request.Context(), role.Query{Name: c.Query("name")}, page, size)
	if err != nil {
		response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
		return
	}
	response.OK(c, listResult{List: roles, Page: pg})
}

func (s *HTTPServer) createRole(c *gin.Context) {
	var req rolesvc.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	r, err := s.role.Create(c.Request.Context(), req)
	if err != nil {
		response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
		return
	}
	response.OK(c, r)
}

func (s *HTTPServer) getRole(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	r, err := s.role.Get(c.Request.Context(), id)
	if err != nil {
		response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
		return
	}
	response.OK(c, r)
}

func (s *HTTPServer) updateRole(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req rolesvc.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	if err := s.role.Update(c.Request.Context(), id, req); err != nil {
		s.roleErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) deleteRole(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := s.role.Delete(c.Request.Context(), id); err != nil {
		s.roleErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) setRolePermissions(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req rolesvc.SetPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	sub := middleware.Subject(c)
	if err := s.role.SetPermissions(c.Request.Context(), sub.UserID, id, req); err != nil {
		s.roleErr(c, err)
		return
	}
	response.OK(c, nil)
}
func (s *HTTPServer) roleErr(c *gin.Context, err error) {
	if errors.Is(err, role.ErrSuperRoleLocked) || errors.Is(err, role.ErrPermExceedsOperator) {
		response.FailI18n(c, http.StatusForbidden, response.CodeForbidden, err)
		return
	}
	response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
}
