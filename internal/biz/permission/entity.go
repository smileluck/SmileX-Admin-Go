// Package permission 权限限界上下文 —— 领域层
package permission

import (
	"sort"
	"time"
)

// Type 权限类型：menu 前端菜单 / button 按钮权限点（可选绑定接口参与 RBAC 校验）
type Type string

const (
	TypeMenu   Type = "menu"
	TypeButton Type = "button"
	// TypeAPI 已废弃：接口绑定能力并入 button，存量 api 记录由启动迁移转为 button
	TypeAPI Type = "api"
)

// Permission 权限聚合根（type=menu 时为前端菜单，type=button 时为按钮权限点）
type Permission struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Type      Type      `json:"type"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	ParentID  uint      `json:"parent_id"`
	Icon      string    `json:"icon"`
	Sort      int       `json:"sort"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Match 判断权限是否能命中请求：button（及存量未迁移的 api）绑定了 method/path 即参与校验
func (p *Permission) Match(method, path string) bool {
	if p.Type == TypeMenu || p.Method == "" || p.Path == "" {
		return false
	}
	if p.Path == "*" {
		return true
	}
	if p.Method != "*" && p.Method != method {
		return false
	}
	if p.Path == path {
		return true
	}
	// 前缀通配：/api/v1/users/*
	if n := len(p.Path); n > 1 && p.Path[n-1] == '*' {
		prefix := p.Path[:n-1]
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// MenuNode 菜单树节点
type MenuNode struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	Code     string      `json:"code"`
	Path     string      `json:"path"`
	Icon     string      `json:"icon"`
	Sort     int         `json:"sort"`
	Children []*MenuNode `json:"children"`
}

// BuildMenuTree 将菜单权限列表组装为树（同级按 sort 升序、id 稳定排序）
func BuildMenuTree(items []*Permission, parentID uint) []*MenuNode {
	var nodes []*MenuNode
	for _, p := range items {
		if p.Type != TypeMenu || p.ParentID != parentID {
			continue
		}
		node := &MenuNode{ID: p.ID, Name: p.Name, Code: p.Code, Path: p.Path, Icon: p.Icon, Sort: p.Sort}
		node.Children = BuildMenuTree(items, p.ID)
		nodes = append(nodes, node)
	}
	sort.SliceStable(nodes, func(i, j int) bool {
		if nodes[i].Sort != nodes[j].Sort {
			return nodes[i].Sort < nodes[j].Sort
		}
		return nodes[i].ID < nodes[j].ID
	})
	return nodes
}
