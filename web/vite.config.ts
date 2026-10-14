import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 管理后台作为内嵌静态站点：构建产物直接输出到 Go 的 embed 目录 kernel/webui/dist，
// base 为 /admin/，与后端挂载前缀一致（单二进制、同源、零 CORS）。
export default defineConfig({
  plugins: [vue()],
  base: '/admin/',
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    outDir: '../kernel/webui/dist',
    emptyOutDir: true,
  },
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://127.0.0.1:8600',
    },
  },
})
