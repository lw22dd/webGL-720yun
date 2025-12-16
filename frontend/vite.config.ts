import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';//注意这里报错时通过快速修复可以安装对应依赖
import tailwindcss from '@tailwindcss/vite';

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()
    ,tailwindcss() ,
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'), // 将 @ 指向 src 目录
      //
    },
  },
  server: {
    host: true, // 允许外部访问
    port: 5173, // Vite默认端口
    proxy: {
      '/api': {
        target: process.env.VITE_API_BASE_URL || 'http://localhost:7000',
        changeOrigin: true,
      },
    },
  },
  build: {
    // 优化构建输出，减少不必要的JS文件
    rollupOptions: {
      output: {
        manualChunks: undefined, // 不拆分代码块
      },
    },
    // 禁用源映射
    sourcemap: false,
  },
})


