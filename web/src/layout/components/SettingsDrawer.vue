<template>
  <el-drawer v-model="visible" :title="t('settings.title')" size="320px" :append-to-body="true">
    <div class="settings">
      <div class="item">
        <div class="item-main">
          <div class="item-label">{{ t('settings.theme') }}</div>
          <div class="item-desc">{{ t('settings.themeTip') }}</div>
        </div>
        <el-switch
          :model-value="settings.theme === 'dark'"
          @change="(v) => settings.setTheme(v ? 'dark' : 'light')"
        />
      </div>

      <div class="item">
        <div class="item-main">
          <div class="item-label">{{ t('settings.i18n') }}</div>
          <div class="item-desc">{{ t('settings.i18nTip') }}</div>
        </div>
        <el-switch :model-value="settings.i18nEnabled" @change="settings.setI18nEnabled" />
      </div>

      <div class="item column">
        <div class="item-main">
          <div class="item-label">{{ t('settings.defaultLang') }}</div>
          <div class="item-desc">{{ t('settings.i18nTip') }}</div>
        </div>
        <el-select
          :model-value="settings.defaultLocale"
          style="width: 100%"
          @change="settings.setDefaultLocale"
        >
          <el-option v-for="l in LANGS" :key="l.value" :label="l.label" :value="l.value" />
        </el-select>
      </div>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useSettingsStore } from '@/stores/settings'
import { LANGS } from '@/i18n'

const { t } = useI18n()
const settings = useSettingsStore()
const visible = ref(false)

defineExpose({ open: () => (visible.value = true) })
</script>

<style scoped>
.settings {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 0;
  border-bottom: 1px solid var(--ra-card-border);
}
.item.column {
  flex-direction: column;
  align-items: stretch;
}
.item-label {
  font-size: 14px;
  font-weight: 600;
  color: var(--ra-text);
}
.item-desc {
  margin-top: 2px;
  font-size: 12px;
  color: var(--ra-text-sub);
}
</style>
