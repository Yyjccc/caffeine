import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    host: "127.0.0.1", // 强制使用 IPv4
    port: 6000,        // 可修改为其他端口，例如 3000
  },
})
