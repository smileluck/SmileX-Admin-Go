// Package captcha 图形验证码限界上下文 —— 领域层。
// 单实例部署，答案存进程内存（库自带 TTL 回收），多实例时需换共享存储。
package captcha

import (
	"time"

	base64Captcha "github.com/mojocn/base64Captcha"
	"github.com/smilex/smilex-admin-gin/internal/conf"
)

const (
	// captchaTTL 验证码有效期
	captchaTTL = 3 * time.Minute
	// collectNum 存量达到该阈值时触发一次过期回收
	collectNum = 200
	// 图片尺寸与登录页 large 输入框（40px 高）视觉对齐
	imgHeight = 40
	imgWidth  = 132
	// textLength 验证码字符数
	textLength = 4
)

// sourceChars 字母数字表，去掉易混淆的 0/O、1/I/L
const sourceChars = "ABCDEFGHJKMNPQRSTUVXYZ23456789"

// Usecase 图形验证码领域用例
type Usecase struct {
	captcha *base64Captcha.Captcha
	enabled bool
}

// NewUsecase 构造图形验证码用例；conf.auth.captchaEnabled=false 时整体停用（本地调试用）
func NewUsecase(c *conf.Bootstrap) *Usecase {
	driver := base64Captcha.NewDriverString(
		imgHeight, imgWidth, 20,
		base64Captcha.OptionShowHollowLine|base64Captcha.OptionShowSlimeLine,
		textLength, sourceChars, nil, nil, nil)
	return &Usecase{
		captcha: base64Captcha.NewCaptcha(driver, base64Captcha.NewMemoryStore(collectNum, captchaTTL)),
		enabled: c.Auth.CaptchaEnabled,
	}
}

// Enabled 验证码功能是否开启（应用层据此决定是否下发图片）
func (uc *Usecase) Enabled() bool {
	return uc.enabled
}

// Generate 生成一对验证码：id 与 base64 图片返回给调用方，答案只留在服务端
func (uc *Usecase) Generate() (id, b64Image string, err error) {
	if !uc.enabled {
		return "", "", nil
	}
	id, b64Image, _, err = uc.captcha.Generate()
	return id, b64Image, err
}

// Verify 校验验证码（一次性：无论对错即失效；忽略大小写）；停用期间恒通过
func (uc *Usecase) Verify(id, answer string) bool {
	if !uc.enabled {
		return true
	}
	return uc.captcha.Verify(id, answer, true)
}
