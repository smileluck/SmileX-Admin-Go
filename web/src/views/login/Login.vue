<template>
  <div class="login-page">
    <div class="bg-orb orb-a"></div>
    <div class="bg-orb orb-b"></div>
    <div class="bg-grid"></div>

    <div class="login-card">
      <div class="login-head">
        <div class="logo-mark">S</div>
        <h1 class="title">SmileX Admin</h1>
        <p class="subtitle">欢迎回来，请登录你的账号</p>
      </div>

      <n-form ref="formRef" :model="form" :rules="rules" :show-label="false">
        <n-form-item path="username">
          <n-input v-model:value="form.username" size="large" placeholder="用户名" @keyup.enter="onLogin" />
        </n-form-item>
        <n-form-item path="password">
          <n-input v-model:value="form.password" size="large" type="password" show-password-on="click" placeholder="密码" @keyup.enter="onLogin" />
        </n-form-item>
        <n-button class="login-btn" type="primary" size="large" block :loading="loading" @click="onLogin">
          登 录
        </n-button>
      </n-form>

      <p class="footer">© 2026 SmileX Admin · All rights reserved</p>
    </div>
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
  position: relative;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: linear-gradient(135deg, #0f172a 0%, #1e1b4b 50%, #312e81 100%);
}

/* 背景光斑 */
.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  opacity: 0.5;
  pointer-events: none;
}
.orb-a {
  width: 480px;
  height: 480px;
  background: #6366f1;
  top: -120px;
  left: -100px;
  animation: float 12s ease-in-out infinite;
}
.orb-b {
  width: 420px;
  height: 420px;
  background: #a855f7;
  bottom: -140px;
  right: -80px;
  animation: float 14s ease-in-out infinite reverse;
}
@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(30px, 40px); }
}

/* 细网格纹理 */
.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.04) 1px, transparent 1px);
  background-size: 44px 44px;
  mask-image: radial-gradient(ellipse at center, black 30%, transparent 75%);
  pointer-events: none;
}

/* 玻璃拟态卡片 */
.login-card {
  position: relative;
  width: 400px;
  padding: 40px 36px 28px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.4);
}

.login-head {
  text-align: center;
  margin-bottom: 32px;
}
.logo-mark {
  width: 56px;
  height: 56px;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  font-size: 26px;
  font-weight: 800;
  color: #fff;
  background: linear-gradient(135deg, #6366f1, #a855f7);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.45);
}
.title {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.5px;
}
.subtitle {
  margin: 8px 0 0;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.55);
}

.login-btn {
  margin-top: 8px;
  border-radius: 10px;
  font-weight: 600;
  letter-spacing: 4px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  transition: transform 0.2s ease, box-shadow 0.2s ease, filter 0.2s ease;
}
.login-btn:not(:disabled):hover {
  transform: translateY(-1px);
  filter: brightness(1.1);
  box-shadow: 0 10px 28px rgba(99, 102, 241, 0.5);
}

.footer {
  margin: 24px 0 0;
  text-align: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}

/* 暗色背景下的输入框样式 */
.login-card :deep(.n-input) {
  border-radius: 10px;
  --n-border: 1px solid rgba(255, 255, 255, 0.18) !important;
  --n-border-hover: 1px solid rgba(255, 255, 255, 0.4) !important;
  --n-border-focus: 1px solid #818cf8 !important;
  --n-color: rgba(255, 255, 255, 0.06) !important;
  --n-color-focus: rgba(255, 255, 255, 0.1) !important;
  --n-text-color: #fff !important;
  --n-placeholder-color: rgba(255, 255, 255, 0.4) !important;
  --n-caret-color: #fff !important;
}
.login-card :deep(.n-input .n-input__eye) {
  color: rgba(255, 255, 255, 0.6);
}
</style>
