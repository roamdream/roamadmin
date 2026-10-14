<template>
  <div class="dashboard">
    <div class="banner">
      <div class="banner-text">
        <h2>{{ greeting }}，{{ auth.userInfo?.nickname || auth.userInfo?.username }}</h2>
        <p>{{ t('dashboard.welcomeBack', { date: today }) }}</p>
      </div>
      <div class="banner-art">
        <div class="ring ring-a"></div>
        <div class="ring ring-b"></div>
        <div class="ring ring-c"></div>
      </div>
    </div>

    <el-row :gutter="20" class="stats">
      <el-col :span="8">
        <div class="stat-card">
          <div class="stat-icon icon-user"><el-icon><User /></el-icon></div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.users }}</div>
            <div class="stat-label">{{ t('dashboard.users') }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-card">
          <div class="stat-icon icon-role"><el-icon><UserFilled /></el-icon></div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.roles }}</div>
            <div class="stat-label">{{ t('dashboard.roles') }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-card">
          <div class="stat-icon icon-tenant"><el-icon><OfficeBuilding /></el-icon></div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.tenants }}</div>
            <div class="stat-label">{{ t('dashboard.tenants') }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="panels">
      <el-col :span="16">
        <el-card>
          <template #header><span class="card-title">{{ t('dashboard.quickStart') }}</span></template>
          <div class="quick">
            <div class="quick-item" @click="router.push('/system/user')">
              <el-icon><User /></el-icon><span>{{ t('menu.user') }}</span>
            </div>
            <div class="quick-item" @click="router.push('/system/role')">
              <el-icon><UserFilled /></el-icon><span>{{ t('menu.role') }}</span>
            </div>
            <div class="quick-item" @click="router.push('/system/menu')">
              <el-icon><Menu /></el-icon><span>{{ t('menu.menu') }}</span>
            </div>
            <div class="quick-item" @click="router.push('/system/dict')">
              <el-icon><Notebook /></el-icon><span>{{ t('menu.dict') }}</span>
            </div>
            <div class="quick-item" @click="router.push('/system/tenant')">
              <el-icon><OfficeBuilding /></el-icon><span>{{ t('menu.tenant') }}</span>
            </div>
            <div class="quick-item" @click="router.push('/system/operlog')">
              <el-icon><Document /></el-icon><span>{{ t('menu.operlog') }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header><span class="card-title">{{ t('dashboard.about') }}</span></template>
          <ul class="about">
            <li><span>{{ t('dashboard.framework') }}</span><b>{{ t('dashboard.frameworkVal') }}</b></li>
            <li><span>{{ t('dashboard.backend') }}</span><b>{{ t('dashboard.backendVal') }}</b></li>
            <li><span>{{ t('dashboard.frontend') }}</span><b>{{ t('dashboard.frontendVal') }}</b></li>
            <li>
              <span>{{ t('dashboard.docs') }}</span>
              <a href="/swagger" target="_blank"><b>{{ t('dashboard.onlineDebug') }}</b></a>
            </li>
          </ul>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import * as api from '@/api'

const router = useRouter()
const { t, locale } = useI18n()
const auth = useAuthStore()
const stats = reactive({ users: 0, roles: 0, tenants: 0 })

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 6) return t('dashboard.night')
  if (h < 12) return t('dashboard.morning')
  if (h < 18) return t('dashboard.afternoon')
  return t('dashboard.evening')
})

const today = computed(() => {
  const loc = locale.value === 'zh-CN' ? 'zh-CN' : 'en-US'
  return new Intl.DateTimeFormat(loc, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long',
  }).format(new Date())
})

onMounted(async () => {
  try {
    const [u, r, t2] = await Promise.all([
      api.listUsers({ page: 1, size: 1 }),
      api.listRoles({ page: 1, size: 1 }),
      api.listTenants({ page: 1, size: 1 }),
    ])
    stats.users = u.total ?? 0
    stats.roles = r.total ?? 0
    stats.tenants = t2.total ?? 0
  } catch {
    // 统计失败不阻塞页面
  }
})
</script>

<style scoped>
.banner {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28px 32px;
  margin-bottom: 20px;
  border-radius: 16px;
  color: #fff;
  overflow: hidden;
  background: linear-gradient(120deg, #4f46e5 0%, #7c3aed 55%, #a855f7 100%);
  box-shadow: 0 16px 40px -16px rgba(99, 102, 241, 0.6);
}
.banner-text h2 {
  margin: 0 0 8px;
  font-size: 22px;
}
.banner-text p {
  margin: 0;
  opacity: 0.85;
  font-size: 14px;
}
.banner-art {
  position: relative;
  width: 160px;
  height: 80px;
}
.ring {
  position: absolute;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.35);
}
.ring-a {
  width: 80px;
  height: 80px;
  right: 30px;
  top: 0;
}
.ring-b {
  width: 130px;
  height: 130px;
  right: 5px;
  top: -25px;
}
.ring-c {
  width: 180px;
  height: 180px;
  right: -20px;
  top: -50px;
}

.stats {
  margin-bottom: 20px;
}
.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  border-radius: 14px;
  background: var(--ra-card);
  border: 1px solid var(--ra-card-border);
  box-shadow: var(--ra-shadow);
}
.stat-icon {
  width: 52px;
  height: 52px;
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: #fff;
  border-radius: 14px;
}
.icon-user {
  background: linear-gradient(135deg, #6366f1, #818cf8);
}
.icon-role {
  background: linear-gradient(135deg, #a855f7, #d946ef);
}
.icon-tenant {
  background: linear-gradient(135deg, #06b6d4, #22d3ee);
}
.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: var(--ra-text);
}
.stat-label {
  font-size: 13px;
  color: var(--ra-text-sub);
}

.card-title {
  font-weight: 600;
}
.quick {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}
.quick-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px;
  border-radius: 10px;
  border: 1px solid var(--ra-card-border);
  color: var(--ra-text-regular);
  cursor: pointer;
  transition: all 0.15s;
}
.quick-item:hover {
  border-color: var(--el-color-primary-light-7);
  background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}
.quick-item .el-icon {
  color: var(--el-color-primary);
}
.about {
  list-style: none;
  margin: 0;
  padding: 0;
}
.about li {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px dashed var(--ra-card-border);
  font-size: 13px;
}
.about li:last-child {
  border-bottom: none;
}
.about span {
  color: var(--ra-text-sub);
}
.about b {
  color: var(--ra-text-regular);
  font-weight: 500;
}
.about a b {
  color: var(--el-color-primary);
}
</style>
