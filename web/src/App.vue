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
    primaryColor: '#0F7B5F',
    primaryColorHover: '#149373',
    primaryColorPressed: '#0B5F4A',
    primaryColorSuppl: '#149373',
    borderRadius: '8px',
    borderRadiusSmall: '5px',
    bodyColor: '#F6F7F5',
    cardColor: '#FFFFFF',
    textColorBase: '#16211D',
    borderColor: '#E4E7E4',
    fontFamily: "-apple-system, BlinkMacSystemFont, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Segoe UI', sans-serif",
  },
}
</script>

<style>
@import './styles/tokens.css';
</style>
