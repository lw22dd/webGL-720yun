// vite.config.ts
import { defineConfig } from "file:///D:/lwdd/code/%E6%AF%95%E8%AE%BE/webGL-720yun/frontend/node_modules/vite/dist/node/index.js";
import vue from "file:///D:/lwdd/code/%E6%AF%95%E8%AE%BE/webGL-720yun/frontend/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import path from "path";
import tailwindcss from "file:///D:/lwdd/code/%E6%AF%95%E8%AE%BE/webGL-720yun/frontend/node_modules/@tailwindcss/vite/dist/index.mjs";
var __vite_injected_original_dirname = "D:\\lwdd\\code\\\u6BD5\u8BBE\\webGL-720yun\\frontend";
var vite_config_default = defineConfig({
  plugins: [
    vue(),
    tailwindcss()
  ],
  resolve: {
    alias: {
      "@": path.resolve(__vite_injected_original_dirname, "./src")
      // 将 @ 指向 src 目录
      //
    }
  },
  server: {
    host: true,
    // 允许外部访问
    port: 5173,
    // Vite默认端口
    proxy: {
      "/api": {
        target: process.env.VITE_API_BASE_URL || "http://localhost:7000",
        changeOrigin: true
      }
    }
  },
  build: {
    // 优化构建输出，减少不必要的JS文件
    rollupOptions: {
      output: {
        manualChunks: void 0
        // 不拆分代码块
      }
    },
    // 禁用源映射
    sourcemap: false
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCJEOlxcXFxsd2RkXFxcXGNvZGVcXFxcXHU2QkQ1XHU4QkJFXFxcXHdlYkdMLTcyMHl1blxcXFxmcm9udGVuZFwiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiRDpcXFxcbHdkZFxcXFxjb2RlXFxcXFx1NkJENVx1OEJCRVxcXFx3ZWJHTC03MjB5dW5cXFxcZnJvbnRlbmRcXFxcdml0ZS5jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0Q6L2x3ZGQvY29kZS8lRTYlQUYlOTUlRTglQUUlQkUvd2ViR0wtNzIweXVuL2Zyb250ZW5kL3ZpdGUuY29uZmlnLnRzXCI7aW1wb3J0IHsgZGVmaW5lQ29uZmlnIH0gZnJvbSAndml0ZSc7XG5pbXBvcnQgdnVlIGZyb20gJ0B2aXRlanMvcGx1Z2luLXZ1ZSc7XG5pbXBvcnQgcGF0aCBmcm9tICdwYXRoJzsvL1x1NkNFOFx1NjEwRlx1OEZEOVx1OTFDQ1x1NjJBNVx1OTUxOVx1NjVGNlx1OTAxQVx1OEZDN1x1NUZFQlx1OTAxRlx1NEZFRVx1NTkwRFx1NTNFRlx1NEVFNVx1NUI4OVx1ODhDNVx1NUJGOVx1NUU5NFx1NEY5RFx1OEQ1NlxuaW1wb3J0IHRhaWx3aW5kY3NzIGZyb20gJ0B0YWlsd2luZGNzcy92aXRlJztcblxuLy8gaHR0cHM6Ly92aXRlLmRldi9jb25maWcvXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoe1xuICBwbHVnaW5zOiBbdnVlKClcbiAgICAsdGFpbHdpbmRjc3MoKSAsXG4gIF0sXG4gIHJlc29sdmU6IHtcbiAgICBhbGlhczoge1xuICAgICAgJ0AnOiBwYXRoLnJlc29sdmUoX19kaXJuYW1lLCAnLi9zcmMnKSwgLy8gXHU1QzA2IEAgXHU2MzA3XHU1NDExIHNyYyBcdTc2RUVcdTVGNTVcbiAgICAgIC8vXG4gICAgfSxcbiAgfSxcbiAgc2VydmVyOiB7XG4gICAgaG9zdDogdHJ1ZSwgLy8gXHU1MTQxXHU4QkI4XHU1OTE2XHU5MEU4XHU4QkJGXHU5NUVFXG4gICAgcG9ydDogNTE3MywgLy8gVml0ZVx1OUVEOFx1OEJBNFx1N0FFRlx1NTNFM1xuICAgIHByb3h5OiB7XG4gICAgICAnL2FwaSc6IHtcbiAgICAgICAgdGFyZ2V0OiBwcm9jZXNzLmVudi5WSVRFX0FQSV9CQVNFX1VSTCB8fCAnaHR0cDovL2xvY2FsaG9zdDo3MDAwJyxcbiAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlLFxuICAgICAgfSxcbiAgICB9LFxuICB9LFxuICBidWlsZDoge1xuICAgIC8vIFx1NEYxOFx1NTMxNlx1Njc4NFx1NUVGQVx1OEY5M1x1NTFGQVx1RkYwQ1x1NTFDRlx1NUMxMVx1NEUwRFx1NUZDNVx1ODk4MVx1NzY4NEpTXHU2NTg3XHU0RUY2XG4gICAgcm9sbHVwT3B0aW9uczoge1xuICAgICAgb3V0cHV0OiB7XG4gICAgICAgIG1hbnVhbENodW5rczogdW5kZWZpbmVkLCAvLyBcdTRFMERcdTYyQzZcdTUyMDZcdTRFRTNcdTc4MDFcdTU3NTdcbiAgICAgIH0sXG4gICAgfSxcbiAgICAvLyBcdTc5ODFcdTc1MjhcdTZFOTBcdTY2MjBcdTVDMDRcbiAgICBzb3VyY2VtYXA6IGZhbHNlLFxuICB9LFxufSlcblxuXG4iXSwKICAibWFwcGluZ3MiOiAiO0FBQTZULFNBQVMsb0JBQW9CO0FBQzFWLE9BQU8sU0FBUztBQUNoQixPQUFPLFVBQVU7QUFDakIsT0FBTyxpQkFBaUI7QUFIeEIsSUFBTSxtQ0FBbUM7QUFNekMsSUFBTyxzQkFBUSxhQUFhO0FBQUEsRUFDMUIsU0FBUztBQUFBLElBQUMsSUFBSTtBQUFBLElBQ1gsWUFBWTtBQUFBLEVBQ2Y7QUFBQSxFQUNBLFNBQVM7QUFBQSxJQUNQLE9BQU87QUFBQSxNQUNMLEtBQUssS0FBSyxRQUFRLGtDQUFXLE9BQU87QUFBQTtBQUFBO0FBQUEsSUFFdEM7QUFBQSxFQUNGO0FBQUEsRUFDQSxRQUFRO0FBQUEsSUFDTixNQUFNO0FBQUE7QUFBQSxJQUNOLE1BQU07QUFBQTtBQUFBLElBQ04sT0FBTztBQUFBLE1BQ0wsUUFBUTtBQUFBLFFBQ04sUUFBUSxRQUFRLElBQUkscUJBQXFCO0FBQUEsUUFDekMsY0FBYztBQUFBLE1BQ2hCO0FBQUEsSUFDRjtBQUFBLEVBQ0Y7QUFBQSxFQUNBLE9BQU87QUFBQTtBQUFBLElBRUwsZUFBZTtBQUFBLE1BQ2IsUUFBUTtBQUFBLFFBQ04sY0FBYztBQUFBO0FBQUEsTUFDaEI7QUFBQSxJQUNGO0FBQUE7QUFBQSxJQUVBLFdBQVc7QUFBQSxFQUNiO0FBQ0YsQ0FBQzsiLAogICJuYW1lcyI6IFtdCn0K
