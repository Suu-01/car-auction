import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  css: {
    postcss: './postcss.config.js',
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: path => path.replace(/^\/api/, '')
      },
      '/ws': {
        target: 'http://localhost:8080',
        ws: true,                 
        changeOrigin: true,
        rewrite: path => path,
      },
      '/static': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: path => path,
      },
    },
  },
})
