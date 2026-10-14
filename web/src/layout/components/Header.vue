<template>
  <div class="header">
    <div class="left">
      <div class="toggle" @click="$emit('toggle')">
        <el-icon><Expand v-if="collapsed" /><Fold v-else /></el-icon>
      </div>
      <div class="crumb">
        <el-icon class="crumb-home"><HomeFilled /></el-icon>
        <span class="sep">/</span>
        <span class="current">{{ route.meta.title ? t(route.meta.title as string) : '' }}</span>
      </div>
    </div>

    <div class="right">
      <el-tooltip :content="t('settings.themeTip')" placement="bottom">
        <div class="icon-btn" @click="toggleTheme">
          <el-icon><Moon v-if="settings.theme === 'light'" /><Sunny v-else /></el-icon>
        </div>
      </el-tooltip>

      <el-dropdown
        v-if="settings.i18nEnabled"
        trigger="click"
        @command="(v) => settings.setLocale(v)"
      >
        <div class="icon-btn lang">
          <el-icon><Switch /></el-icon>
          <span>{{ currentLangLabel }}</span>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item
              v-for="l in LANGS"
              :key="l.value"
              :command="l.value"
              :disabled="l.value === settings.locale"
            >
              {{ l.label }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <el-tooltip :content="t('common.apiDoc')" placement="bottom">
        <a class="icon-btn" href="/swagger" target="_blank">
          <el-icon><Document /></el-icon>
        </a>
      </el-tooltip>

      <el-tooltip :content="t('settings.title')" placement="bottom">
        <div class="icon-btn" @click="settingsRef?.open()">
          <el-icon><Setting /></el-icon>
        </div>
      </el-tooltip>

      <el-dropdown trigger="click" @command="onCommand">
        <div class="user">
          <div class="avatar">{{ avatarText }}</div>
          <span class="name">{{ displayName }}</span>
          <el-icon class="caret"><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="logout">{{ t('common.logout') }}</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <SettingsDrawer ref="settingsRef" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { LANGS } from '@/i18n'
import SettingsDrawer from './SettingsDrawer.vue'

defineProps<{ collapsed: boolean }>()
defineEmits<{ (e: 'toggle'): void }>()

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const auth = useAuthStore()
const settings = useSettingsStore()
const settingsRef = ref()

const displayName = computed(
  () => auth.userInfo?.nickname || auth.userInfo?.username || '—',
)
const avatarText = computed(() => displayName.value.charAt(0).toUpperCase())

const currentLangLabel = computed(
  () => LANGS.find((l) => l.value === settings.locale)?.label ?? '',
)

function toggleTheme() {
  settings.setTheme(settings.theme === 'dark' ? 'light' : 'dark')
}

async function onCommand(cmd: string) {
  if (cmd !== 'logout') return
  try {
    await ElMessageBox.confirm(t('common.confirmLogout'), t('common.tip'), { type: 'warning' })
  } catch {
    return
  }
  await auth.logout()
  ElMessage.success(t('common.loggedOut'))
  router.push('/login')
}
</script>

<style scoped>
.header {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.left {
  display: flex;
  align-items: center;
  gap: 16px;
}
.toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 8px;
  cursor: pointer;
  color: var(--ra-text-regular);
  transition: all 0.15s;
}
.toggle:hover {
  background: var(--ra-hover);
  color: var(--ra-text);
}
.crumb {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--ra-text-sub);
}
.crumb-home {
  color: var(--el-color-primary);
}
.current {
  color: var(--ra-text);
  font-weight: 600;
}
.right {
  display: flex;
  align-items: center;
  gap: 8px;
}
.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 34px;
  min-width: 34px;
  padding: 0 8px;
  border-radius: 8px;
  color: var(--ra-text-regular);
  transition: all 0.15s;
  cursor: pointer;
  text-decoration: none;
}
.icon-btn:hover {
  background: var(--ra-hover);
  color: var(--el-color-primary);
}
.icon-btn.lang {
  font-size: 13px;
}
.user {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 4px 8px;
  margin-left: 8px;
  border-radius: 10px;
  transition: background 0.15s;
}
.user:hover {
  background: var(--ra-hover);
}
.avatar {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  border-radius: 50%;
  background: linear-gradient(135deg, #6366f1, #a855f7);
}
.name {
  font-size: 14px;
  color: var(--ra-text);
}
.caret {
  font-size: 12px;
  color: var(--ra-text-sub);
}
</style>
