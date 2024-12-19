import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  base:'/plugin/prometheus',
  plugins: [vue()],
  resolve: {
    extensions:['.vue','.ts'],
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server:{
    host:'10.41.107.198',
  //   proxy:{
  //     "/": {
  //       target: 'http://10.41.107.32:8090',
  //       changeOrigin:true,
  //     }
  // },
}
})
