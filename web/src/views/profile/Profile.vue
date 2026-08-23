<template>
  <div class="profile-page">
    <n-grid :cols="24" responsive="screen" :x-gap="12" :y-gap="12">
      <!-- 左：身份摘要 -->
      <n-gi span="24 s:8">
        <n-card class="sum-card">
          <div class="sum">
            <div class="sum-avatar">{{ avatarChar }}</div>
            <div class="sum-name">{{ user?.nickname || user?.username }}</div>
            <div class="sum-username mono">@{{ user?.username }}</div>
            <div class="sum-roles">
              <n-tag v-for="r in user?.role_names || []" :key="r" size="small" round>{{ r }}</n-tag>
              <span v-if="!user?.role_names?.length" class="sum-empty">{{ t('profile.noRole') }}</span>
            </div>
            <div class="sum-meta">
              <div class="sum-meta-row">
                <span class="muted">{{ t('profile.email') }}</span>
                <span class="mono">{{ user?.email || '—' }}</span>
              </div>
              <div class="sum-meta-row">
                <span class="muted">{{ t('common.status') }}</span>
                <n-tag :type="user?.status === 1 ? 'success' : 'error'" size="small" round>
                  {{ user?.status === 1 ? t('common.enabled') : t('common.disabled') }}
                </n-tag>
              </div>
              <div class="sum-meta-row">
                <span class="muted">{{ t('profile.joinTime') }}</span>
                <span class="mono">{{ user?.created_at || '—' }}</span>
              </div>
            </div>
          </div>
        </n-card>
      </n-gi>

      <!-- 右：资料编辑 + 修改密码 -->
      <n-gi span="24 s:16">
        <n-space vertical :size="12">
          <n-card :title="t('profile.basicInfo')">
            <n-form ref="infoFormRef" :model="infoForm" :rules="infoRules" label-placement="left" label-width="80" class="form">
              <n-form-item :label="t('profile.nickname')" path="nickname">
                <n-input v-model:value="infoForm.nickname" :maxlength="20" show-word-limit :placeholder="t('profile.nicknamePlaceholder')" />
              </n-form-item>
              <n-form-item :label="t('profile.email')" path="email">
                <n-input v-model:value="infoForm.email" :placeholder="t('profile.emailPlaceholder')" />
              </n-form-item>
              <div class="form-actions">
                <n-button type="primary" :loading="infoSaving" @click="saveInfo">{{ t('common.save') }}</n-button>
              </div>
            </n-form>
          </n-card>

          <n-card :title="t('profile.changePassword')">
            <n-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-placement="left" label-width="80" class="form">
              <n-form-item :label="t('profile.oldPassword')" path="oldPassword">
                <n-input v-model:value="pwdForm.oldPassword" type="password" show-password-on="click" :placeholder="t('profile.oldPasswordPlaceholder')" />
              </n-form-item>
              <n-form-item :label="t('profile.newPassword')" path="newPassword">
                <n-input v-model:value="pwdForm.newPassword" type="password" show-password-on="click" :maxlength="20" :placeholder="t('profile.newPasswordPlaceholder')" />
              </n-form-item>
              <n-form-item :label="t('profile.confirmPassword')" path="confirmPassword">
                <n-input v-model:value="pwdForm.confirmPassword" type="password" show-password-on="click" :maxlength="20" :placeholder="t('profile.confirmPasswordPlaceholder')" />
              </n-form-item>
              <div class="form-actions">
                <n-button type="primary" :loading="pwdSaving" @click="savePwd">{{ t('profile.changePassword') }}</n-button>
              </div>
            </n-form>
          </n-card>
        </n-space>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NGrid, NGi, NCard, NSpace, NForm, NFormItem, NInput, NButton, NTag, useMessage,
  type FormInst, type FormRules,
} from 'naive-ui'
import { changePassword, getProfile, updateProfile } from '../../api'
import { useUserStore } from '../../stores/user'

const router = useRouter()
const message = useMessage()
const userStore = useUserStore()
const { t } = useI18n()

const user = computed(() => userStore.user)
const avatarChar = computed(() => (user.value?.nickname || user.value?.username || 'U').charAt(0).toUpperCase())

// 进入页面刷新一次本人信息（角色名/资料保持最新），并回填编辑表单
onMounted(async () => {
  try {
    const { data } = await getProfile()
    userStore.user = data.data.user
  } catch { /* 拉取失败沿用 store 现有数据 */ }
  const u = userStore.user
  if (u) {
    infoForm.nickname = u.nickname || ''
    infoForm.email = u.email || ''
  }
})

// ---- 基本信息 ----
const infoFormRef = ref<FormInst | null>(null)
const infoForm = reactive({ nickname: '', email: '' })
const infoSaving = ref(false)
// 与后端 binding 保持一致：昵称≤20、邮箱格式（选填）
const infoRules = computed<FormRules>(() => ({
  nickname: [{ max: 20, message: t('profile.rules.nicknameTooLong'), trigger: ['blur', 'input'] }],
  email: [
    {
      trigger: ['blur', 'input'],
      validator: (_rule, value: string) => !value || /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value),
      message: t('profile.rules.emailInvalid'),
    },
  ],
}))

async function saveInfo() {
  try {
    await infoFormRef.value?.validate()
  } catch {
    return // 校验失败，错误已在表单项上展示
  }
  infoSaving.value = true
  try {
    const { data } = await updateProfile({ nickname: infoForm.nickname.trim(), email: infoForm.email.trim() })
    userStore.user = data.data.user
    message.success(t('common.saveSuccess'))
  } catch (e: any) {
    message.error(e?.response?.data?.msg || t('profile.saveFailed'))
  } finally {
    infoSaving.value = false
  }
}

// ---- 修改密码（成功后强制重新登录） ----
const pwdFormRef = ref<FormInst | null>(null)
const pwdForm = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' })
const pwdSaving = ref(false)
const pwdRules = computed<FormRules>(() => ({
  oldPassword: [{ required: true, message: t('profile.rules.oldPasswordRequired'), trigger: ['blur', 'input'] }],
  newPassword: [
    { required: true, message: t('profile.rules.newPasswordRequired'), trigger: ['blur', 'input'] },
    { min: 6, max: 20, message: t('profile.rules.passwordLength'), trigger: ['blur', 'input'] },
  ],
  confirmPassword: [
    { required: true, message: t('profile.rules.confirmPasswordRequired'), trigger: ['blur', 'input'] },
    { validator: (_rule, v: string) => v === pwdForm.newPassword, message: t('profile.rules.passwordMismatch'), trigger: ['blur', 'input'] },
  ],
}))

async function savePwd() {
  try {
    await pwdFormRef.value?.validate()
  } catch {
    return
  }
  pwdSaving.value = true
  try {
    await changePassword({ old_password: pwdForm.oldPassword, new_password: pwdForm.newPassword })
    message.success(t('profile.passwordChanged'))
    await userStore.logout()
    router.push('/login')
  } catch (e: any) {
    message.error(e?.response?.data?.msg || t('profile.changeFailed'))
  } finally {
    pwdSaving.value = false
  }
}
</script>

<style scoped>
.profile-page {
  min-height: 100%;
}
.sum {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 8px 8px;
  text-align: center;
}
.sum-avatar {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #fff;
  font-weight: 700;
  font-size: 26px;
  background: var(--sx-accent);
  margin-bottom: 12px;
}
.sum-name {
  font-size: 17px;
  font-weight: 600;
  color: var(--sx-ink);
}
.sum-username {
  font-size: 12px;
  color: var(--sx-muted);
  margin: 2px 0 12px;
}
.sum-roles {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 6px;
  margin-bottom: 16px;
}
.sum-empty {
  font-size: 12px;
  color: var(--sx-muted);
}
.sum-meta {
  width: 100%;
  border-top: 1px solid var(--sx-line);
  padding-top: 12px;
}
.sum-meta-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 4px;
  font-size: 13px;
  color: var(--sx-ink);
}
.muted {
  color: var(--sx-muted);
}
.form {
  max-width: 420px;
}
.form-actions {
  padding-left: 80px;
}
</style>
