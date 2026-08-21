<template>
  <n-config-provider :theme="null" :theme-overrides="themeOverrides" :locale="zhCN" :date-locale="dateZhCN">
    <n-message-provider>
      <n-dialog-provider>
        <router-view />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { NConfigProvider, NMessageProvider, NDialogProvider, zhCN, dateZhCN, type GlobalThemeOverrides } from 'naive-ui'

// 主题规范：色值与 src/styles/tokens.css 同源，改主题两处同步调整
const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#3F75AB',
    primaryColorHover: '#518CC8',
    primaryColorPressed: '#315E8C',
    primaryColorSuppl: '#518CC8',
    borderRadius: '8px',
    borderRadiusSmall: '5px',
    bodyColor: '#F5F7FA',
    cardColor: '#FFFFFF',
    textColorBase: '#151E2B',
    borderColor: '#E3E8EF',
    fontFamily: "-apple-system, BlinkMacSystemFont, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Segoe UI', sans-serif",
  },
}
</script>

<style>
@import './styles/tokens.css';
</style>
