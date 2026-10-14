<template>
  <el-container class="layout">
    <el-aside :width="collapsed ? '68px' : '224px'" class="aside">
      <Sidebar :collapsed="collapsed" />
    </el-aside>
    <el-container class="body">
      <el-header class="header">
        <Header :collapsed="collapsed" @toggle="collapsed = !collapsed" />
      </el-header>
      <TabsView />
      <el-main class="main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" :key="route.fullPath + ':' + viewKey" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, provide } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from './components/Sidebar.vue'
import Header from './components/Header.vue'
import TabsView from './components/TabsView.vue'

const route = useRoute()
const collapsed = ref(false)
const viewKey = ref(0)

// 供标签页「刷新」使用：改变 key 强制重渲染当前路由组件。
provide('reload', () => {
  viewKey.value++
})
</script>

<style scoped>
.layout {
  height: 100%;
}
.aside {
  background: var(--ra-card);
  border-right: 1px solid var(--ra-card-border);
  transition: width 0.2s ease;
  overflow-x: hidden;
}
.body {
  min-width: 0;
}
.header {
  height: 60px;
  background: var(--ra-card);
  border-bottom: 1px solid var(--ra-card-border);
  padding: 0 20px;
}
.main {
  padding: 16px 20px 20px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.18s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
