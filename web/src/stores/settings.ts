import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import i18n, { type Locale } from '@/i18n'

const KEY = 'ra_settings'

interface SettingsState {
  theme: 'light' | 'dark'
  i18nEnabled: boolean
  defaultLocale: Locale
  locale: Locale
}

function load(): SettingsState {
  const def: SettingsState = {
    theme: 'light',
    i18nEnabled: true,
    defaultLocale: 'zh-CN',
    locale: 'zh-CN',
  }
  try {
    const raw = localStorage.getItem(KEY)
    if (raw) return { ...def, ...JSON.parse(raw) }
  } catch {
    // ignore
  }
  return def
}

export const useSettingsStore = defineStore('settings', () => {
  const state = load()

  const theme = ref<'light' | 'dark'>(state.theme)
  const i18nEnabled = ref<boolean>(state.i18nEnabled)
  const defaultLocale = ref<Locale>(state.defaultLocale)
  const locale = ref<Locale>(state.locale)

  function applyTheme() {
    document.documentElement.classList.toggle('dark', theme.value === 'dark')
  }

  function applyLocale() {
    const target = i18nEnabled.value ? locale.value : defaultLocale.value
    i18n.global.locale.value = target
    document.documentElement.lang = target
  }

  function persist() {
    try {
      localStorage.setItem(
        KEY,
        JSON.stringify({
          theme: theme.value,
          i18nEnabled: i18nEnabled.value,
          defaultLocale: defaultLocale.value,
          locale: locale.value,
        }),
      )
    } catch {
      // ignore
    }
  }

  // 初始化时立即生效（在 main.ts 挂载前调用）。
  applyTheme()
  applyLocale()

  watch([theme, i18nEnabled, defaultLocale, locale], () => {
    applyTheme()
    applyLocale()
    persist()
  })

  function setTheme(t: 'light' | 'dark') {
    theme.value = t
  }

  function setI18nEnabled(v: boolean) {
    i18nEnabled.value = v
  }

  function setDefaultLocale(l: Locale) {
    defaultLocale.value = l
  }

  function setLocale(l: Locale) {
    locale.value = l
  }

  return {
    theme,
    i18nEnabled,
    defaultLocale,
    locale,
    setTheme,
    setI18nEnabled,
    setDefaultLocale,
    setLocale,
  }
})
