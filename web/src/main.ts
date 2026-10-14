import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'

import App from './App.vue'
import router from './router'
import i18n from './i18n'
import { onUnauthorized } from './api/client'
import { useAuthStore } from './stores/auth'
import { useSettingsStore } from './stores/settings'
import './styles/index.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(i18n)
app.use(ElementPlus)

// 初始化时应用持久化的主题/语言设置（需在挂载前生效）。
useSettingsStore()

// 全局注册图标组件，侧边栏等处通过 <component :is="IconName" /> 动态渲染。
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 任意接口返回 401 时，清理登录态并跳转登录页。
onUnauthorized(() => {
  useAuthStore().reset()
  router.push('/login')
})

app.mount('#app')
