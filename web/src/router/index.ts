import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// meta.perm 用于按权限过滤侧边栏（admin 的 perms 含 "*"，可见全部）。
const routes = [
  {
    path: '/login',
    component: () => import('@/views/login/index.vue'),
    meta: { public: true, title: '登录' },
  },
  {
    path: '/',
    component: () => import('@/layout/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: 'menu.dashboard', icon: 'HomeFilled' },
      },
      {
        path: 'system/user',
        component: () => import('@/views/system/user/index.vue'),
        meta: { title: 'menu.user', icon: 'User', perm: 'system:user:list', group: 'menu.system' },
      },
      {
        path: 'system/role',
        component: () => import('@/views/system/role/index.vue'),
        meta: { title: 'menu.role', icon: 'UserFilled', perm: 'system:role:list', group: 'menu.system' },
      },
      {
        path: 'system/menu',
        component: () => import('@/views/system/menu/index.vue'),
        meta: { title: 'menu.menu', icon: 'Menu', perm: 'system:menu:list', group: 'menu.system' },
      },
      {
        path: 'system/dict',
        component: () => import('@/views/system/dict/index.vue'),
        meta: { title: 'menu.dict', icon: 'Notebook', perm: 'system:dict:list', group: 'menu.system' },
      },
      {
        path: 'system/tenant',
        component: () => import('@/views/system/tenant/index.vue'),
        meta: { title: 'menu.tenant', icon: 'OfficeBuilding', perm: 'system:tenant:list', group: 'menu.system' },
      },
      {
        path: 'system/operlog',
        component: () => import('@/views/system/operlog/index.vue'),
        meta: { title: 'menu.operlog', icon: 'Document', perm: 'system:operlog:list', group: 'menu.system' },
      },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/dashboard' },
]

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes,
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.token) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.path === '/login' && auth.token) {
    return { path: '/dashboard' }
  }
  // 刷新后 token 仍在但用户信息/权限尚未加载，先补齐，避免侧边栏菜单因无权限被过滤。
  if (auth.token && !auth.userInfo) {
    try {
      await auth.fetchInfo()
    } catch {
      auth.reset()
      return { path: '/login', query: { redirect: to.fullPath } }
    }
  }
  return true
})

export default router
