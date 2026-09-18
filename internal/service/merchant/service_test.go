package merchant

import (
	"context"
	"testing"
	"time"

	bizmerchant "github.com/smilex/smilex-admin-gin/internal/biz/merchant"
)

// ---- 测试假件 ----

type fakeRepo struct{ m *bizmerchant.Merchant }

func (r *fakeRepo) Create(context.Context, *bizmerchant.Merchant) error { return nil }
func (r *fakeRepo) Update(context.Context, *bizmerchant.Merchant) error { return nil }
func (r *fakeRepo) Delete(context.Context, uint) error                  { return nil }
func (r *fakeRepo) Get(ctx context.Context, id uint) (*bizmerchant.Merchant, error) {
	if r.m != nil && r.m.ID == id {
		cp := *r.m
		return &cp, nil
	}
	return nil, bizmerchant.ErrMerchantNotFound
}
func (r *fakeRepo) FindByAppKey(context.Context, string) (*bizmerchant.Merchant, error) {
	return nil, bizmerchant.ErrMerchantNotFound
}
func (r *fakeRepo) List(context.Context, bizmerchant.Query, int, int) ([]*bizmerchant.Merchant, int64, error) {
	if r.m == nil {
		return nil, 0, nil
	}
	cp := *r.m
	return []*bizmerchant.Merchant{&cp}, 1, nil
}

type fakeLogRepo struct{ logs []*bizmerchant.APILog }

func (r *fakeLogRepo) Record(*bizmerchant.APILog) {}
func (r *fakeLogRepo) List(context.Context, bizmerchant.APILogQuery, int, int) ([]*bizmerchant.APILog, int64, error) {
	out := make([]*bizmerchant.APILog, 0, len(r.logs))
	for _, l := range r.logs {
		cp := *l
		out = append(out, &cp)
	}
	return out, int64(len(r.logs)), nil
}

func newTestMerchant() *bizmerchant.Merchant {
	return &bizmerchant.Merchant{
		ID: 1, Name: "示例商户", Code: "demo", AppKey: "mk_demo",
		ContactName: "张三", ContactPhone: "13812341234", ContactEmail: "alice@example.com",
		Status: bizmerchant.StatusEnabled,
	}
}

func newTestService(m *bizmerchant.Merchant, logs []*bizmerchant.APILog) *Service {
	return NewService(bizmerchant.NewUsecase(&fakeRepo{m: m}, &fakeLogRepo{logs: logs}))
}

// TestToVO_Masked 默认脱敏：导出的 ToVO（开放 API ping 复用）恒脱敏
func TestToVO_Masked(t *testing.T) {
	vo := ToVO(newTestMerchant())
	if vo.ContactName != "张*" || vo.ContactPhone != "138****1234" || vo.ContactEmail != "a***@example.com" {
		t.Errorf("默认应脱敏: %+v", vo)
	}
}

// TestGet_Reveal reveal 开关控制联系方式明/脱敏
func TestGet_Reveal(t *testing.T) {
	s := newTestService(newTestMerchant(), nil)
	ctx := context.Background()

	vo, err := s.Get(ctx, 1, true)
	if err != nil {
		t.Fatal(err)
	}
	if vo.ContactName != "张三" || vo.ContactPhone != "13812341234" || vo.ContactEmail != "alice@example.com" {
		t.Errorf("reveal=true 应输出明文: %+v", vo)
	}

	vo, err = s.Get(ctx, 1, false)
	if err != nil {
		t.Fatal(err)
	}
	if vo.ContactPhone != "138****1234" {
		t.Errorf("reveal=false 应脱敏: %+v", vo)
	}
}

// TestList_Reveal 列表 reveal 开关
func TestList_Reveal(t *testing.T) {
	s := newTestService(newTestMerchant(), nil)
	ctx := context.Background()

	list, _, err := s.List(ctx, bizmerchant.Query{}, 1, 10, true)
	if err != nil || len(list) != 1 {
		t.Fatalf("list reveal: %v %v", list, err)
	}
	if list[0].ContactPhone != "13812341234" {
		t.Errorf("reveal=true 应明文, got %q", list[0].ContactPhone)
	}

	list, _, err = s.List(ctx, bizmerchant.Query{}, 1, 10, false)
	if err != nil || len(list) != 1 {
		t.Fatalf("list masked: %v %v", list, err)
	}
	if list[0].ContactPhone != "138****1234" {
		t.Errorf("reveal=false 应脱敏, got %q", list[0].ContactPhone)
	}
}

// TestListAPILogs_Reveal 日志 IP 明/脱敏
func TestListAPILogs_Reveal(t *testing.T) {
	logs := []*bizmerchant.APILog{{
		ID: 1, MerchantID: 1, AppKey: "mk_demo", Method: "GET", Path: "/open-api/v1/ping",
		IP: "192.168.1.10", StatusCode: 200, LatencyMs: 5, CreatedAt: time.Now(),
	}}
	s := newTestService(newTestMerchant(), logs)
	ctx := context.Background()

	out, _, err := s.ListAPILogs(ctx, bizmerchant.APILogQuery{}, 1, 10, false)
	if err != nil || len(out) != 1 {
		t.Fatalf("masked: %v %v", out, err)
	}
	if out[0].IP != "192.168.*.*" {
		t.Errorf("reveal=false 应脱敏 IP, got %q", out[0].IP)
	}

	out, _, err = s.ListAPILogs(ctx, bizmerchant.APILogQuery{}, 1, 10, true)
	if err != nil || len(out) != 1 {
		t.Fatalf("reveal: %v %v", out, err)
	}
	if out[0].IP != "192.168.1.10" {
		t.Errorf("reveal=true 应输出明文 IP, got %q", out[0].IP)
	}
}
