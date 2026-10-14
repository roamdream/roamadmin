<template>
  <div class="tabs">
    <el-scrollbar class="tabs-scroll">
      <div class="tab-list">
        <div
          v-for="tab in visited"
          :key="tab.path"
          class="tab"
          :class="{ active: tab.path === route.path }"
          @click="go(tab)"
        >
          <el-icon class="tab-icon"><component :is="tab.icon" /></el-icon>
          <span>{{ t(tab.titleKey) }}</span>
          <el-icon v-if="tab.closable" class="tab-close" @click.stop="close(tab)">
            <Close />
          </el-icon>
        </div>
      </div>
    </el-scrollbar>

    <el-dropdown trigger="click" @command="onCommand">
      <div class="tabs-actions">
        <el-icon><ArrowDown /></el-icon>
      </div>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="refresh">{{ t('tabs.refresh') }}</el-dropdown-item>
          <el-dropdown-item command="closeOthers">{{ t('tabs.closeOthers') }}</el-dropdown-item>
          <el-dropdown-item command="closeAll">{{ t('tabs.closeAll') }}</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, inject } from 'vue'
import { useRoute, useRouter, type RouteRecordRaw } from 'vue-router'
import { useI18n } from 'vue-i18n'

interface Tab {
  path: string
  titleKey: string
  icon: string
  closable: boolean
}

const AFFIX = '/dashboard'
const STORE_KEY = 'ra_tabs'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const reload = inject<() => void>('reload', () => {})

// 所有可入标签页的布局子路由。
const candidates = (() => {
  const root = router.options.routes.find((r) => r.path === '/')
  return (root?.children ?? []) as RouteRecordRaw[]
})()

function buildTab(path: string): Tab | null {
  const rec = candidates.find((c) => `/${c.path}` === path)
  if (!rec?.meta?.title) return null
  return {
    path,
    titleKey: rec.meta.title as string,
    icon: (rec.meta.icon as string) ?? 'Menu',
    closable: path !== AFFIX,
  }
}

function loadSaved(): Tab[] {
  let paths: string[] = []
  try {
    const raw = sessionStorage.getItem(STORE_KEY)
    if (raw) paths = JSON.parse(raw)
  } catch {
    // ignore
  }
  const tabs = paths.map(buildTab).filter(Boolean) as Tab[]
  // 确保固定的首页始终存在。
  if (!tabs.some((x) => x.path === AFFIX)) {
    const home = buildTab(AFFIX)
    if (home) tabs.unshift(home)
  }
  return tabs
}

function persist() {
  try {
    sessionStorage.setItem(
      STORE_KEY,
      JSON.stringify(visited.value.map((x) => x.path)),
    )
  } catch {
    // ignore
  }
}

const visited = ref<Tab[]>(loadSaved())

watch(
  () => route.path,
  () => {
    const meta = route.meta
    if (!meta.title) return
    if (!visited.value.some((x) => x.path === route.path)) {
      const tab = buildTab(route.path)
      if (tab) {
        visited.value.push(tab)
        persist()
      }
    }
  },
  { immediate: true },
)

function go(tab: Tab) {
  if (tab.path !== route.path) router.push(tab.path)
}

function close(tab: Tab) {
  const idx = visited.value.findIndex((x) => x.path === tab.path)
  if (idx === -1) return
  visited.value.splice(idx, 1)
  persist()
  if (tab.path === route.path) {
    const next = visited.value[idx] ?? visited.value[idx - 1]
    router.push(next ? next.path : AFFIX)
  }
}

function onCommand(cmd: string) {
  if (cmd === 'refresh') {
    reload()
    return
  }
  if (cmd === 'closeOthers') {
    visited.value = visited.value.filter((x) => !x.closable || x.path === route.path)
    persist()
  }
  if (cmd === 'closeAll') {
    visited.value = visited.value.filter((x) => !x.closable)
    persist()
    if (route.path !== AFFIX) router.push(AFFIX)
  }
}
</script>

<style scoped>
.tabs {
  display: flex;
  align-items: center;
  height: 42px;
  background: var(--ra-card);
  border-bottom: 1px solid var(--ra-card-border);
  padding: 0 8px;
}
.tabs-scroll {
  flex: 1;
}
.tab-list {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 42px;
  white-space: nowrap;
}
.tab {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 30px;
  padding: 0 10px;
  font-size: 13px;
  color: var(--ra-text-regular);
  border: 1px solid var(--ra-card-border);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
  user-select: none;
}
.tab:hover {
  color: var(--el-color-primary);
  border-color: var(--el-color-primary-light-7);
}
.tab.active {
  color: #fff;
  background: var(--el-color-primary);
  border-color: var(--el-color-primary);
}
.tab-icon {
  font-size: 14px;
}
.tab-close {
  font-size: 12px;
  border-radius: 50%;
  padding: 2px;
}
.tab-close:hover {
  background: rgba(255, 255, 255, 0.25);
}
.tabs-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  margin-left: 6px;
  border-radius: 6px;
  color: var(--ra-text-regular);
  cursor: pointer;
}
.tabs-actions:hover {
  background: var(--ra-hover);
  color: var(--el-color-primary);
}
</style>
