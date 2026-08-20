<template>
  <n-grid :cols="4" :x-gap="16">
    <n-gi v-for="card in cards" :key="card.label">
      <n-card>
        <n-statistic :label="card.label" :value="card.value" />
      </n-card>
    </n-gi>
  </n-grid>
  <n-card title="欢迎" style="margin-top: 16px">
    <p>{{ userStore.user?.nickname || userStore.user?.username }}，欢迎使用 SmileX Admin 管理系统。</p>
    <p style="color: #999; font-size: 13px">
      DDD + Gin + GORM 后端 · Vue3 + Naive UI 前端 · 多数据库（MySQL/PostgreSQL/SQLite）
    </p>
  </n-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NGrid, NGi, NCard, NStatistic } from 'naive-ui'
import { useUserStore } from '../../stores/user'
import { listUsers, listRoles, listPermissions } from '../../api'

const userStore = useUserStore()
const cards = ref([
  { label: '用户数', value: 0 },
  { label: '角色数', value: 0 },
  { label: '权限数', value: 0 },
  { label: '我的权限码', value: userStore.permissions.length },
])

onMounted(async () => {
  const [u, r, p] = await Promise.all([
    listUsers({ page: 1, page_size: 1 }),
    listRoles({ page: 1, page_size: 1 }),
    listPermissions({ page: 1, page_size: 1 }),
  ])
  cards.value[0].value = u.data.data.page.total
  cards.value[1].value = r.data.data.page.total
  cards.value[2].value = p.data.data.page.total
})
</script>
