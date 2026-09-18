package export

import (
	"context"
	"net/url"
	"testing"
	"time"

	bizuser "github.com/smilex/smilex-admin-gin/internal/biz/user"
	"github.com/smilex/smilex-admin-gin/internal/conf"
)

// ---- 测试假件 ----

type fakeUserRepo struct{ users []*bizuser.User }

func (r *fakeUserRepo) Create(context.Context, *bizuser.User) error        { return nil }
func (r *fakeUserRepo) Update(context.Context, *bizuser.User) error        { return nil }
func (r *fakeUserRepo) UpdatePassword(context.Context, uint, string) error { return nil }
func (r *fakeUserRepo) Delete(context.Context, uint) error                 { return nil }
func (r *fakeUserRepo) FindByID(context.Context, uint) (*bizuser.User, error) {
	return nil, bizuser.ErrUserNotFound
}
func (r *fakeUserRepo) FindByIDWithPassword(context.Context, uint) (*bizuser.User, error) {
	return nil, bizuser.ErrUserNotFound
}
func (r *fakeUserRepo) FindByUsername(context.Context, string) (*bizuser.User, error) {
	return nil, bizuser.ErrUserNotFound
}
func (r *fakeUserRepo) List(_ context.Context, _ bizuser.Query, _ int, _ int) ([]*bizuser.User, int64, error) {
	out := make([]*bizuser.User, 0, len(r.users))
	for _, u := range r.users {
		cp := *u
		out = append(out, &cp)
	}
	return out, int64(len(out)), nil
}
func (r *fakeUserRepo) SetRoles(context.Context, uint, []uint) error { return nil }

// TestUserExporter_Reveal reveal=1 跳过 export.mask 脱敏（提交入口已按权限码剔除无权限请求的该参数）
func TestUserExporter_Reveal(t *testing.T) {
	cfg := &conf.Bootstrap{Export: conf.Export{Mask: map[string]string{"phone": "phone", "email": "email"}}}
	e := NewUserExporter(&fakeUserRepo{users: []*bizuser.User{{
		ID: 1, Username: "alice", Nickname: "Alice", Phone: "13812341234",
		Email: "alice@example.com", Status: bizuser.StatusEnabled, CreatedAt: time.Now(),
	}}}, cfg)
	ctx := context.Background()

	rows, total, err := e.Fetch(ctx, url.Values{}, 0, 10)
	if err != nil || total != 1 || len(rows) != 1 {
		t.Fatalf("masked fetch: total=%d rows=%v err=%v", total, rows, err)
	}
	if rows[0][3] != "138****1234" || rows[0][4] != "a***@example.com" {
		t.Errorf("默认应按 export.mask 脱敏, got %v", rows[0])
	}

	rows, _, err = e.Fetch(ctx, url.Values{"reveal": {"1"}}, 0, 10)
	if err != nil || len(rows) != 1 {
		t.Fatalf("reveal fetch: %v %v", rows, err)
	}
	if rows[0][3] != "13812341234" || rows[0][4] != "alice@example.com" {
		t.Errorf("reveal=1 应输出明文, got %v", rows[0])
	}
}
