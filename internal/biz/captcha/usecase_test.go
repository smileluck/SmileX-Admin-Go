package captcha

import (
	"sync"
	"testing"

	base64Captcha "github.com/mojocn/base64Captcha"
	"github.com/smilex/smilex-admin-gin/internal/conf"
)

// memStore 内存版 base64Captcha.Store，模拟 Redis 行为：不存在的 key 取回空串
type memStore struct {
	mu   sync.Mutex
	data map[string]string
}

func newMemStore() *memStore { return &memStore{data: map[string]string{}} }

func (s *memStore) Set(id, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = value
	return nil
}

func (s *memStore) Get(id string, clear bool) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data[id]
	if clear {
		delete(s.data, id)
	}
	if !ok {
		return ""
	}
	return v
}

func (s *memStore) Verify(id, answer string, clear bool) bool { return false }

func newTestUsecase(enabled bool) *Usecase {
	return NewUsecase(&conf.Bootstrap{Auth: conf.Auth{CaptchaEnabled: enabled}}, newMemStore())
}

func TestVerifyEmptyIdOrAnswerRejected(t *testing.T) {
	// 开启状态下，省略验证码字段（空 id + 空答案）必须被拒绝：
	// base64Captcha.Verify 对不存在的 id 取回空串，与空答案 EqualFold 相等会误放行
	uc := newTestUsecase(true)
	if uc.Verify("", "") {
		t.Fatal("empty id + empty answer must be rejected when captcha enabled")
	}
	if uc.Verify("", "abcd") {
		t.Fatal("empty id must be rejected")
	}
	if uc.Verify("some-id", "") {
		t.Fatal("empty answer must be rejected")
	}
}

func TestVerifyDisabledAlwaysPass(t *testing.T) {
	uc := newTestUsecase(false)
	if !uc.Verify("", "") {
		t.Fatal("disabled captcha must pass through")
	}
}

func TestVerifyCorrectAnswer(t *testing.T) {
	uc := newTestUsecase(true)
	uc.captcha = base64Captcha.NewCaptcha(nil, newMemStore())
	if err := uc.captcha.Store.Set("id1", "abcd"); err != nil {
		t.Fatal(err)
	}
	if !uc.Verify("id1", "ABCD") {
		t.Fatal("correct answer (case-insensitive) must pass")
	}
	if uc.Verify("id1", "abcd") {
		t.Fatal("captcha must be one-time: second verify must fail")
	}
}
