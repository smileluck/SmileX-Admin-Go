<template>
  <div class="login-page">
    <!-- 品牌面板：签名元素「系统脉搏」 -->
    <aside class="brand-panel">
      <div class="brand-top">
        <div class="seal">S</div>
        <span class="brand-name">SmileX Admin</span>
      </div>

      <div class="brand-copy">
        <p class="mono-label">system console</p>
        <h1 class="headline">让系统治理<br />安静地运转</h1>
        <p class="sub">用户、角色、权限、菜单——一处掌控。</p>
      </div>

      <!-- 系统脉搏波形 -->
      <div class="pulse">
        <div class="pulse-head mono-label">system pulse</div>
        <svg class="pulse-line" viewBox="0 0 560 80" preserveAspectRatio="none" aria-hidden="true">
          <line x1="0" y1="40" x2="560" y2="40" class="pulse-base" />
          <polyline
            points="0,40 60,40 80,40 95,12 110,66 125,40 190,40 210,40 225,12 240,66 255,40 330,40 350,40 365,12 380,66 395,40 470,40 490,40 505,12 520,66 535,40 560,40"
            class="pulse-trace" />
        </svg>
        <div class="pulse-meta mono">
          <span>uptime 99.98%</span>
          <span>region cn-east</span>
          <span class="live"><i></i>operational</span>
        </div>
      </div>

      <div class="brand-foot mono">© 2026 SmileX · internal use only</div>
    </aside>

    <!-- 表单面板 -->
    <main class="form-panel">
      <div class="form-box">
        <div class="form-head">
          <div class="seal seal--sm">S</div>
          <div>
            <h2 class="form-title">登录</h2>
            <p class="form-sub">使用你的管理员账号继续</p>
          </div>
        </div>

        <n-form ref="formRef" :model="form" :rules="rules" :show-label="false">
          <n-form-item path="username">
            <n-input v-model:value="form.username" size="large" placeholder="用户名" @keyup.enter="onLogin" />
          </n-form-item>
          <n-form-item path="password">
            <n-input v-model:value="form.password" size="large" type="password" show-password-on="click" placeholder="密码"
              @keyup.enter="onLogin" />
          </n-form-item>
          <n-button class="login-btn" type="primary" size="large" block :loading="loading" @click="onLogin">
            登 录
          </n-button>
        </n-form>

        <p class="form-foot mono">authorized personnel only</p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, useMessage, type FormInst } from 'naive-ui'
import { useUserStore } from '../../stores/user'

const router = useRouter()
const userStore = useUserStore()
const message = useMessage()
const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const form = reactive({ username: 'admin', password: '' })
const rules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
}

async function onLogin() {
  await formRef.value?.validate()
  loading.value = true
  try {
    await userStore.login(form.username, form.password)
    message.success('登录成功')
    // 不在这里 loadUserContext——交由路由守卫统一加载并注册动态路由
    router.push('/')
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  min-height: 100vh;
  background: var(--sx-bg);
}

/* ---- 品牌面板（签名区域，唯一的大胆处）---- */
.brand-panel {
  position: relative;
  flex: 7;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 40px 48px;
  background: var(--sx-ink);
  color: #EDF2EF;
  overflow: hidden;
}
/* 面板内的工程网格 */
.brand-panel::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.045) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.045) 1px, transparent 1px);
  background-size: 48px 48px;
  mask-image: radial-gradient(ellipse at 30% 40%, black 20%, transparent 80%);
  pointer-events: none;
}

.brand-top {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
}
.seal {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  font-family: var(--sx-font-mono);
  font-weight: 700;
  font-size: 20px;
  color: var(--sx-ink);
  background: #7FD4B8;
}
.brand-name {
  font-family: var(--sx-font-mono);
  font-size: 14px;
  letter-spacing: 0.06em;
  color: rgba(237, 242, 239, 0.85);
}

.brand-copy {
  position: relative;
}
.mono-label {
  font-family: var(--sx-font-mono);
  font-size: 11px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #7FD4B8;
  margin: 0 0 16px;
}
.headline {
  font-size: clamp(28px, 3.2vw, 40px);
  font-weight: 700;
  line-height: 1.3;
  margin: 0 0 14px;
  letter-spacing: 0.01em;
}
.sub {
  margin: 0;
  font-size: 15px;
  color: rgba(237, 242, 239, 0.6);
}

/* 系统脉搏 */
.pulse {
  position: relative;
  margin-top: 48px;
}
.pulse-head {
  color: rgba(237, 242, 239, 0.5);
  margin-bottom: 10px;
}
.pulse-line {
  width: 100%;
  height: 72px;
  display: block;
}
.pulse-base {
  stroke: rgba(255, 255, 255, 0.12);
  stroke-width: 1;
}
.pulse-trace {
  fill: none;
  stroke: #7FD4B8;
  stroke-width: 2;
  stroke-linejoin: round;
  stroke-linecap: round;
  stroke-dasharray: 1200;
  stroke-dashoffset: 1200;
  animation: trace 4.5s linear infinite;
}
@keyframes trace {
  0% { stroke-dashoffset: 1200; opacity: 0.35; }
  45% { opacity: 1; }
  100% { stroke-dashoffset: 0; opacity: 0.35; }
}
.pulse-meta {
  display: flex;
  gap: 20px;
  margin-top: 12px;
  font-family: var(--sx-font-mono);
  font-size: 11px;
  color: rgba(237, 242, 239, 0.45);
}
.live {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #7FD4B8;
}
.live i {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #7FD4B8;
  animation: blink 2.4s ease-in-out infinite;
}
@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.25; }
}

.brand-foot {
  position: relative;
  font-size: 11px;
  color: rgba(237, 242, 239, 0.3);
}

/* ---- 表单面板（克制、安静）---- */
.form-panel {
  flex: 3;
  min-width: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 24px;
}
.form-box {
  width: 100%;
  max-width: 360px;
}
.form-head {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 32px;
}
.seal--sm {
  width: 44px;
  height: 44px;
  border-radius: 11px;
  background: var(--sx-accent);
  color: #fff;
  font-size: 21px;
}
.form-title {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  color: var(--sx-ink);
}
.form-sub {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--sx-muted);
}

.login-btn {
  margin-top: 8px;
  font-weight: 600;
  letter-spacing: 4px;
}

.form-foot {
  margin: 28px 0 0;
  text-align: center;
}

/* 响应式：窄屏时品牌面板退化为顶部条 */
@media (max-width: 860px) {
  .login-page {
    flex-direction: column;
  }
  .brand-panel {
    flex: none;
    padding: 28px 32px;
  }
  .brand-copy {
    display: none;
  }
  .pulse {
    display: none;
  }
  .brand-foot {
    display: none;
  }
}
</style>
