// Package permission 权限限界上下文 —— 领域层
package permission

import "time"

// Type 权限类型：api 接口权限 / menu 前端菜单标识
type Type string

const (
	TypeAPI  Type = "api"
	TypeMenu Type = "menu"
)

// Permission 权限聚合根（type=api 时为接口权限，type=menu 时为前端菜单）
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

// Match 判断权限是否能命中请求
func (p *Permission) Match(method, path string) bool {
	if p.Type != TypeAPI {
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

// BuildMenuTree 将菜单权限列表组装为树
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
	return nodes
}
