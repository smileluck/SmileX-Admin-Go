<template>
  <div class="about-page">
    <!-- 固定左右布局：左信息卡固定宽、右更新记录自适应，任何窗口宽度都不堆叠 -->
    <div class="about-layout">
      <!-- 左：产品信息卡 -->
      <n-card class="info-card">
        <div class="brand">
          <div class="seal">S</div>
          <div class="brand-text">
            <span class="brand-name">SmileX Admin</span>
            <n-tag size="small" round type="primary">v{{ version }}</n-tag>
          </div>
        </div>

        <div class="info-rows">
          <div v-for="row in infoRows" :key="row.label" class="info-row">
            <span class="info-icon"><n-icon :component="row.icon" :size="15" /></span>
            <span class="info-label">{{ row.label }}</span>
            <span class="info-value">{{ row.value }}</span>
          </div>
        </div>

        <div class="copyright mono">© 2026 SmileX · internal use only</div>
      </n-card>

      <!-- 右：更新记录（git 提交日志构建时生成） -->
      <n-card class="log-card" :title="t('about.updateLog')">
        <template #header-extra>
          <span class="gen-meta mono">{{ t('about.commitMeta', { n: commits.length }) }}</span>
        </template>
        <div class="log-list">
          <div v-for="c in commits" :key="c.hash" class="log-item">
            <n-tag :type="typeMeta(c.type).tag" size="small" round>{{ typeMeta(c.type).label }}</n-tag>
            <span v-if="c.scope" class="log-scope mono">{{ c.scope }}</span>
            <span class="log-message">{{ c.message }}</span>
            <span class="log-date mono">{{ c.date }}</span>
          </div>
        </div>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NCard, NTag, NIcon } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import {
  AppsOutline, GitCommitOutline, DesktopOutline, ServerOutline,
  DiscOutline, LayersOutline, PersonOutline,
} from '@vicons/ionicons5'
import pkg from '../../../package.json'
import changelog from '../../generated/changelog.json'

const { t } = useI18n()
const version = pkg.version
const commits = changelog.commits

// 左侧信息行：icon + 灰标签 + 粗体值
const infoRows = computed(() => [
  { icon: AppsOutline, label: t('about.info.appName'), value: 'SmileX Admin' },
  { icon: GitCommitOutline, label: t('about.info.currentVersion'), value: `v${version}` },
  { icon: DesktopOutline, label: t('about.info.frontend'), value: 'Vue 3 · Naive UI' },
  { icon: ServerOutline, label: t('about.info.backend'), value: 'Go · Gin · GORM' },
  { icon: DiscOutline, label: t('about.info.database'), value: 'MySQL / PostgreSQL / SQLite' },
  { icon: LayersOutline, label: t('about.info.arch'), value: t('about.archValue') },
  { icon: PersonOutline, label: t('about.info.developer'), value: 'SmileX' },
])

// 提交类型 -> 标签色（conventional commits），文案见 about.type.*
type TagType = 'primary' | 'error' | 'info' | 'warning' | 'success' | 'default'
const TYPE_TAGS: Record<string, TagType> = {
  feat: 'primary',
  fix: 'error',
  style: 'info',
  refactor: 'warning',
  perf: 'success',
  docs: 'default',
  test: 'default',
  chore: 'default',
}
const typeMeta = (type: string) => ({
  tag: TYPE_TAGS[type] ?? ('default' as const),
  label: t(`about.type.${TYPE_TAGS[type] ? type : 'other'}`),
})
</script>

<style scoped>
/* 撑满一屏：100vh - 顶栏 64px - 内容区上下 padding 8px×2（height:100% 链在
   naive 嵌套滚动容器下不精确，会导致整页滚动；改用视口确定性高度） */
.about-layout {
  display: flex;
  gap: 12px;
  align-items: stretch;
  height: calc(100vh - 64px - 16px);
}
/* 左右 4:6 固定比例布局，任何窗口宽度都并排 */
.info-card {
  flex: 4 1 0;
  min-width: 0;
}
.log-card {
  flex: 6 1 0;
  min-width: 0;
  display: flex;
  flex-direction: column;
}
/* 卡片内容纵向撑满，记录列表吃满剩余高度并内部滚动
   （naive 卡片内容类为 n-card-content，非 n-card__content） */
.log-card :deep(.n-card-content) {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}
/* 左卡内容在极小高度下兜底滚动，避免溢出 */
.info-card :deep(.n-card-content) {
  overflow: auto;
}
/* 窄窗口下隐藏次要信息，保证并排布局可用 */
@media (max-width: 540px) {
  .log-date {
    display: none;
  }
}

/* 品牌区：印章 + 名称 + 版本 */
.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--sx-line);
}
.seal {
  width: 46px;
  height: 46px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  font-family: var(--sx-font-mono);
  font-weight: 700;
  font-size: 22px;
  color: #fff;
  background: var(--sx-accent);
}
.brand-text {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}
.brand-name {
  font-size: 16px;
  font-weight: 700;
  color: var(--sx-ink);
}

/* 信息行：icon + 灰标签 … 黑色粗体值 */
.info-rows {
  padding: 14px 0 2px;
}
.info-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
}
.info-icon {
  display: inline-flex;
  align-items: center;
  color: var(--sx-muted);
}
.info-label {
  font-size: 12px;
  color: var(--sx-muted);
}
.info-value {
  margin-left: auto;
  font-size: 13px;
  font-weight: 600;
  color: var(--sx-ink);
  text-align: right;
}

.copyright {
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px solid var(--sx-line);
  font-size: 10px;
  color: var(--sx-muted);
}

/* 更新记录列表 */
.gen-meta {
  font-size: 11px;
  color: var(--sx-muted);
}
.log-list {
  flex: 1;
  min-height: 0;
  overflow: auto;
  scrollbar-width: thin;
}
.log-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 2px;
  border-bottom: 1px solid var(--sx-line);
}
.log-item:last-child {
  border-bottom: none;
}
.log-scope {
  flex-shrink: 0;
  font-size: 11px;
  color: var(--sx-accent);
}
.log-message {
  flex: 1;
  min-width: 0;
  font-size: 13px;
  color: var(--sx-ink);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.log-date {
  flex-shrink: 0;
  font-size: 11px;
  color: var(--sx-muted);
}
</style>
