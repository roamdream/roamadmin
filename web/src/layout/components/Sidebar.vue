<template>
  <div class="sidebar">
    <div class="logo">
      <div class="logo-mark">R</div>
      <transition name="fade">
        <span v-show="!collapsed" class="logo-text">RoamAdmin</span>
      </transition>
    </div>

    <el-scrollbar class="scroll">
      <el-menu
        :default-active="route.path"
        :collapse="collapsed"
        :collapse-transition="false"
        router
        class="menu"
      >
        <el-menu-item v-for="m in menus.tops" :key="m.path" :index="m.path">
          <el-icon><component :is="m.icon" /></el-icon>
          <template #title>{{ t(m.title) }}</template>
        </el-menu-item>

        <el-sub-menu v-for="g in menus.groups" :key="g.name" :index="'group-' + g.name">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>{{ t(g.name) }}</span>
          </template>
          <el-menu-item v-for="m in g.items" :key="m.path" :index="m.path">
            <el-icon><component :is="m.icon" /></el-icon>
            <template #title>{{ t(m.title) }}</template>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

defineProps<{ collapsed: boolean }>()

interface MenuItem {
  path: string
  title: string
  icon: string
  group: string
}

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const { t } = useI18n()

const menus = computed(() => {
  const root = router.options.routes.find((r) => r.path === '/')
  const children = (root?.children ?? [])
    .filter((c) => {
      const perm = c.meta?.perm as string | undefined
      return perm ? auth.hasPerm(perm) : true
    })
    .map<MenuItem>((c) => ({
      path: `/${c.path}`,
      title: (c.meta?.title as string) ?? '',
      icon: (c.meta?.icon as string) ?? 'Menu',
      group: (c.meta?.group as string) ?? '',
    }))

  const tops = children.filter((i) => !i.group)
  const map = new Map<string, MenuItem[]>()
  for (const i of children) {
    if (!i.group) continue
    if (!map.has(i.group)) map.set(i.group, [])
    map.get(i.group)!.push(i)
  }
  const groups = [...map.entries()].map(([name, items]) => ({ name, items }))
  return { tops, groups }
})
</script>

<style scoped>
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--ra-card);
}
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  border-bottom: 1px solid var(--ra-card-border);
  overflow: hidden;
}
.logo-mark {
  width: 32px;
  height: 32px;
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
  border-radius: 9px;
  background: linear-gradient(135deg, #6366f1, #a855f7);
}
.logo-text {
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 0.3px;
  white-space: nowrap;
  background: linear-gradient(90deg, #6366f1, #a855f7);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}
.scroll {
  flex: 1;
}

.menu {
  border-right: none;
  padding: 10px 10px 0;
}
.menu:not(.el-menu--collapse) {
  width: 100%;
}
.menu.el-menu--collapse {
  width: 100%;
  padding: 10px 4px 0;
}
.menu.el-menu--collapse :deep(.el-menu-item),
.menu.el-menu--collapse :deep(.el-sub-menu__title) {
  justify-content: center;
  padding: 0 !important;
}
.menu :deep(.el-menu-item),
.menu :deep(.el-sub-menu__title) {
  height: 44px;
  line-height: 44px;
  margin: 4px 0;
  border-radius: 10px;
  color: var(--ra-text-regular);
}
.menu :deep(.el-menu-item:hover),
.menu :deep(.el-sub-menu__title:hover) {
  background: var(--ra-hover);
  color: var(--ra-text);
}
.menu :deep(.el-menu-item.is-active) {
  background: var(--ra-active-bg);
  color: var(--el-color-primary);
  font-weight: 600;
}
.menu :deep(.el-menu-item .el-icon),
.menu :deep(.el-sub-menu__title .el-icon) {
  font-size: 18px;
}
.menu :deep(.el-sub-menu .el-menu-item) {
  min-width: 0;
  padding-left: 46px !important;
}
.menu :deep(.el-sub-menu.is-active > .el-sub-menu__title) {
  color: var(--el-color-primary);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
