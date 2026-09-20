// 商户（开放 API 授权）与商户 API 日志 handler。
package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	bizmerchant "github.com/smilex/smilex-admin-gin/internal/biz/merchant"
	merchantsvc "github.com/smilex/smilex-admin-gin/internal/service/merchant"
	"github.com/smilex/smilex-admin-gin/pkg/i18n"
	"github.com/smilex/smilex-admin-gin/pkg/response"
)

// ---- 商户（开放 API 授权） ----

func (s *HTTPServer) listMerchants(c *gin.Context) {
	page, size := pageParams(c)
	q := bizmerchant.Query{
		Name:   strings.TrimSpace(c.Query("name")),
		Code:   strings.TrimSpace(c.Query("code")),
		AppKey: strings.TrimSpace(c.Query("app_key")),
	}
	if v := c.Query("status"); v != "" {
		if st, err := strconv.Atoi(v); err == nil {
			q.Status = &st
		}
	}
	list, pg, err := s.merchant.List(c.Request.Context(), q, page, size, s.revealParam(c, "merchant:viewSensitive"))
	if err != nil {
		response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
		return
	}
	response.OK(c, listResult{List: list, Page: pg})
}

func (s *HTTPServer) createMerchant(c *gin.Context) {
	var req merchantsvc.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	vo, err := s.merchant.Create(c.Request.Context(), req)
	if err != nil {
		s.merchantErr(c, err)
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) getMerchant(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	vo, err := s.merchant.Get(c.Request.Context(), id, s.revealParam(c, "merchant:viewSensitive"))
	if err != nil {
		s.merchantErr(c, err)
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) updateMerchant(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req merchantsvc.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	if err := s.merchant.Update(c.Request.Context(), id, req); err != nil {
		s.merchantErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) deleteMerchant(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := s.merchant.Delete(c.Request.Context(), id); err != nil {
		s.merchantErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) resetMerchantSecret(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	vo, err := s.merchant.ResetSecret(c.Request.Context(), id)
	if err != nil {
		s.merchantErr(c, err)
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) setMerchantStatus(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req merchantsvc.SetStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
		return
	}
	if err := s.merchant.SetStatus(c.Request.Context(), id, req); err != nil {
		s.merchantErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) listMerchantAPILogs(c *gin.Context) {
	page, size := pageParams(c)
	q := bizmerchant.APILogQuery{AppKey: c.Query("app_key"), Path: c.Query("path")}
	if v := c.Query("status_code"); v != "" {
		if sc, err := strconv.Atoi(v); err == nil {
			q.StatusCode = &sc
		}
	}
	if t, ok := parseUnixParam(c.Query("start")); ok {
		q.Start = t
	}
	if t, ok := parseUnixParam(c.Query("end")); ok {
		q.End = t
	}
	logs, pg, err := s.merchant.ListAPILogs(c.Request.Context(), q, page, size, s.revealParam(c, "merchantLog:viewSensitive"))
	if err != nil {
		response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
		return
	}
	response.OK(c, listResult{List: logs, Page: pg})
}

// merchantErr 商户操作错误映射：不存在 404，其余 400
func (s *HTTPServer) merchantErr(c *gin.Context, err error) {
	if errors.Is(err, bizmerchant.ErrMerchantNotFound) {
		response.FailI18n(c, http.StatusNotFound, response.CodeErr, err)
		return
	}
	response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
}
