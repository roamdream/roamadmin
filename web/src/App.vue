<template>
  <el-config-provider :locale="epLocale">
    <router-view />
  </el-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import i18n from './i18n'

// 随 i18n 语言切换 Element Plus 组件内置文案（分页、日期、弹窗按钮等）。
const epLocale = computed(() => (i18n.global.locale.value === 'zh-CN' ? zhCn : en))
</script>
