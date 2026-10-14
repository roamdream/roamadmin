<template>
  <div class="login">
    <div class="bg-grid"></div>
    <div class="glow glow-a"></div>
    <div class="glow glow-b"></div>

    <div class="panel">
      <div class="left">
        <div class="brand">
          <div class="brand-mark">R</div>
          <h1 class="brand-name">{{ t('login.title') }}</h1>
          <p class="brand-sub">{{ t('login.subtitle') }}</p>
          <p class="brand-tag">{{ t('login.tagline') }}</p>
        </div>

        <div class="radar-wrap">
          <div class="radar">
            <div class="ring r1"></div>
            <div class="ring r2"></div>
            <div class="ring r3"></div>
            <div class="sweep"></div>
            <div class="blip b1"></div>
            <div class="blip b2"></div>
            <div class="blip b3"></div>
            <div class="core"></div>
          </div>
        </div>

        <div class="stats">
          <div class="stat">
            <div class="stat-val">{{ devices }}</div>
            <div class="stat-label">{{ t('login.statDevices') }}</div>
          </div>
          <div class="stat">
            <div class="stat-val">{{ nodes }}</div>
            <div class="stat-label">{{ t('login.statNodes') }}</div>
          </div>
          <div class="stat">
            <div class="stat-val">99.99%</div>
            <div class="stat-label">{{ t('login.statUptime') }}</div>
          </div>
        </div>
      </div>

      <div class="right">
        <div class="form-card">
          <h2 class="form-title">{{ t('login.login') }}</h2>
          <el-form ref="formRef" :model="form" :rules="rules" size="large">
            <el-form-item prop="username">
              <el-input v-model="form.username" :placeholder="t('login.username')">
                <template #prefix><el-icon><User /></el-icon></template>
              </el-input>
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="form.password"
                type="password"
                :placeholder="t('login.password')"
                show-password
                @keyup.enter="submit"
              >
                <template #prefix><el-icon><Lock /></el-icon></template>
              </el-input>
            </el-form-item>
            <el-button type="primary" class="btn" :loading="loading" @click="submit">
              {{ t('login.login') }}
            </el-button>
          </el-form>
          <p class="hint">{{ t('login.hint') }}</p>
        </div>
      </div>
    </div>

    <p class="copyright">© 2026 RoamAdmin · {{ t('login.copyright') }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const auth = useAuthStore()

const formRef = ref()
const loading = ref(false)
const form = reactive({ username: 'admin', password: 'admin123' })
const rules = {
  username: [{ required: true, message: () => t('login.usernameRequired'), trigger: 'blur' }],
  password: [{ required: true, message: () => t('login.passwordRequired'), trigger: 'blur' }],
}

// 数字滚动效果
const devices = ref(0)
const nodes = ref(0)
function countUp(target: number, ref: { value: number }, duration = 1600) {
  const start = performance.now()
  const tick = (now: number) => {
    const p = Math.min((now - start) / duration, 1)
    ref.value = Math.floor(target * (1 - Math.pow(1 - p, 3)))
    if (p < 1) requestAnimationFrame(tick)
  }
  requestAnimationFrame(tick)
}
onMounted(() => {
  countUp(12480, devices)
  countUp(3210, nodes)
})

async function submit() {
  await formRef.value.validate().catch(() => Promise.reject())
  loading.value = true
  try {
    await auth.login(form.username, form.password)
    ElMessage.success(t('login.loginSuccess'))
    router.push((route.query.redirect as string) || '/dashboard')
  } catch (e: any) {
    ElMessage.error(e?.message || t('login.loginFailed'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* 主题变量：跟随后台浅色/暗色主题，登录页与后台配色保持一致。 */
.login {
  --lg-bg: radial-gradient(1200px 800px at 20% 0%, #e0e7ff 0%, transparent 55%),
    radial-gradient(1000px 700px at 100% 100%, #f3e8ff 0%, transparent 55%), #f5f7fa;
  --lg-grid: rgba(99, 102, 241, 0.12);
  --lg-glow-a: #818cf8;
  --lg-glow-b: #c084fc;
  --lg-glow-opacity: 0.4;
  --lg-panel-bg: rgba(255, 255, 255, 0.72);
  --lg-panel-border: rgba(99, 102, 241, 0.14);
  --lg-panel-shadow: 0 40px 120px -30px rgba(99, 102, 241, 0.35);
  --lg-text: #1f2329;
  --lg-brand-sub: #4f46e5;
  --lg-tag: #6b7280;
  --lg-ring: rgba(99, 102, 241, 0.3);
  --lg-sweep: rgba(99, 102, 241, 0.55);
  --lg-core: #6366f1;
  --lg-stat-bg: #ffffff;
  --lg-stat-border: #eef0f4;
  --lg-stat-val: #4f46e5;
  --lg-stat-label: #6b7280;
  --lg-right-bg: rgba(255, 255, 255, 0.55);
  --lg-right-border: #eef0f4;
  --lg-form-title: #1f2329;
  --lg-input-bg: #ffffff;
  --lg-input-border: rgba(148, 163, 184, 0.4);
  --lg-input-text: #1f2329;
  --lg-input-ph: #9ca3af;
  --lg-input-prefix: #6366f1;
  --lg-hint: #9ca3af;
  --lg-copyright: #6b7280;
}

html.dark .login {
  --lg-bg: radial-gradient(1200px 800px at 20% 0%, #1e1b4b 0%, transparent 55%),
    radial-gradient(1000px 700px at 100% 100%, #312e81 0%, transparent 55%), #0b1020;
  --lg-grid: rgba(99, 102, 241, 0.14);
  --lg-glow-a: #6366f1;
  --lg-glow-b: #a855f7;
  --lg-glow-opacity: 0.5;
  --lg-panel-bg: rgba(17, 24, 39, 0.6);
  --lg-panel-border: rgba(148, 163, 184, 0.16);
  --lg-panel-shadow: 0 40px 120px -30px rgba(0, 0, 0, 0.8);
  --lg-text: #ffffff;
  --lg-brand-sub: #c7d2fe;
  --lg-tag: rgba(255, 255, 255, 0.5);
  --lg-ring: rgba(99, 102, 241, 0.35);
  --lg-sweep: rgba(99, 102, 241, 0.7);
  --lg-core: #818cf8;
  --lg-stat-bg: rgba(255, 255, 255, 0.05);
  --lg-stat-border: rgba(148, 163, 184, 0.14);
  --lg-stat-val: #a5b4fc;
  --lg-stat-label: rgba(255, 255, 255, 0.55);
  --lg-right-bg: rgba(255, 255, 255, 0.04);
  --lg-right-border: rgba(148, 163, 184, 0.14);
  --lg-form-title: #ffffff;
  --lg-input-bg: rgba(255, 255, 255, 0.06);
  --lg-input-border: rgba(148, 163, 184, 0.2);
  --lg-input-text: #ffffff;
  --lg-input-ph: rgba(255, 255, 255, 0.4);
  --lg-input-prefix: #a5b4fc;
  --lg-hint: rgba(255, 255, 255, 0.35);
  --lg-copyright: rgba(255, 255, 255, 0.35);
}

.login {
  position: relative;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: var(--lg-bg);
  transition: background 0.25s ease;
}

/* 透视网格背景 */
.bg-grid {
  position: absolute;
  inset: 0;
  background-image: linear-gradient(var(--lg-grid) 1px, transparent 1px),
    linear-gradient(90deg, var(--lg-grid) 1px, transparent 1px);
  background-size: 44px 44px;
  transform: perspective(600px) rotateX(55deg) scale(2.2);
  transform-origin: center 20%;
  -webkit-mask-image: linear-gradient(180deg, #000 0%, transparent 90%);
  mask-image: linear-gradient(180deg, #000 0%, transparent 90%);
  animation: grid-move 20s linear infinite;
}
@keyframes grid-move {
  to {
    background-position: 0 44px, 44px 0;
  }
}

.glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  opacity: var(--lg-glow-opacity);
}
.glow-a {
  width: 480px;
  height: 480px;
  top: -160px;
  left: -120px;
  background: var(--lg-glow-a);
}
.glow-b {
  width: 420px;
  height: 420px;
  bottom: -140px;
  right: -100px;
  background: var(--lg-glow-b);
}

.panel {
  position: relative;
  z-index: 2;
  display: flex;
  width: 960px;
  max-width: 94vw;
  border-radius: 24px;
  overflow: hidden;
  background: var(--lg-panel-bg);
  border: 1px solid var(--lg-panel-border);
  backdrop-filter: blur(20px);
  box-shadow: var(--lg-panel-shadow);
}

/* 左侧品牌 + 雷达可视化 */
.left {
  flex: 1.2;
  padding: 48px 44px;
  color: var(--lg-text);
}
.brand-mark {
  width: 52px;
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  border-radius: 14px;
  background: linear-gradient(135deg, #6366f1, #a855f7);
  box-shadow: 0 12px 30px -8px rgba(99, 102, 241, 0.5);
}
.brand-name {
  margin: 18px 0 6px;
  font-size: 30px;
  letter-spacing: 0.5px;
}
.brand-sub {
  margin: 0;
  font-size: 15px;
  color: var(--lg-brand-sub);
}
.brand-tag {
  margin: 22px 0 0;
  font-size: 13px;
  letter-spacing: 1px;
  color: var(--lg-tag);
}

.radar-wrap {
  display: flex;
  justify-content: center;
  margin: 30px 0 10px;
}
.radar {
  position: relative;
  width: 180px;
  height: 180px;
}
.radar .ring {
  position: absolute;
  inset: 0;
  margin: auto;
  border: 1px solid var(--lg-ring);
  border-radius: 50%;
}
.radar .r1 {
  width: 60px;
  height: 60px;
}
.radar .r2 {
  width: 120px;
  height: 120px;
}
.radar .r3 {
  width: 180px;
  height: 180px;
}
.sweep {
  position: absolute;
  inset: 0;
  margin: auto;
  width: 90px;
  height: 90px;
  border-radius: 50%;
  background: conic-gradient(from 0deg, var(--lg-sweep), transparent 70deg);
  -webkit-mask-image: radial-gradient(circle, transparent 0, #000 100%);
  mask-image: radial-gradient(circle, transparent 0, #000 100%);
  animation: spin 3s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
.blip {
  position: absolute;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #22d3ee;
  box-shadow: 0 0 10px #22d3ee;
  animation: pulse 2.4s ease-in-out infinite;
}
.b1 {
  top: 30px;
  left: 60px;
}
.b2 {
  top: 90px;
  right: 24px;
  animation-delay: 0.8s;
}
.b3 {
  bottom: 40px;
  left: 70px;
  animation-delay: 1.6s;
}
@keyframes pulse {
  0%,
  100% {
    opacity: 0.2;
  }
  50% {
    opacity: 1;
  }
}
.core {
  position: absolute;
  inset: 0;
  margin: auto;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--lg-core);
  box-shadow: 0 0 16px var(--lg-core);
}

.stats {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}
.stat {
  flex: 1;
  text-align: center;
  padding: 14px 8px;
  border-radius: 12px;
  background: var(--lg-stat-bg);
  border: 1px solid var(--lg-stat-border);
}
.stat-val {
  font-size: 22px;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
  color: var(--lg-stat-val);
}
.stat-label {
  margin-top: 4px;
  font-size: 12px;
  color: var(--lg-stat-label);
}

/* 右侧登录表单 */
.right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 40px;
  background: var(--lg-right-bg);
  border-left: 1px solid var(--lg-right-border);
}
.form-card {
  width: 100%;
}
.form-title {
  margin: 0 0 24px;
  font-size: 20px;
  color: var(--lg-form-title);
  text-align: center;
}
.form-card :deep(.el-input__wrapper) {
  background: var(--lg-input-bg);
  box-shadow: 0 0 0 1px var(--lg-input-border) inset;
}
.form-card :deep(.el-input__inner) {
  color: var(--lg-input-text);
}
.form-card :deep(.el-input__inner::placeholder) {
  color: var(--lg-input-ph);
}
.form-card :deep(.el-input__prefix) {
  color: var(--lg-input-prefix);
}
.btn {
  width: 100%;
  height: 46px;
  font-size: 15px;
  border-radius: 10px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border: none;
  box-shadow: 0 12px 28px -10px rgba(99, 102, 241, 0.7);
}
.btn:hover {
  opacity: 0.92;
}
.hint {
  margin: 16px 0 0;
  text-align: center;
  font-size: 12px;
  color: var(--lg-hint);
}

.copyright {
  position: absolute;
  bottom: 20px;
  width: 100%;
  text-align: center;
  font-size: 12px;
  color: var(--lg-copyright);
  z-index: 2;
}

@media (max-width: 860px) {
  .left {
    display: none;
  }
  .right {
    border-left: none;
  }
  .panel {
    width: 400px;
  }
}
</style>
