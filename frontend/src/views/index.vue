<template>
  <div class="page-wrapper">
    <Header />
    <main class="main-content">
      <div class="hero-section">
        <div class="hero-content">
          <h1 class="hero-title">720° 全景云平台</h1>
          <p class="hero-subtitle">沉浸式全景体验，让世界触手可及</p>
          <div class="hero-buttons">
            <t-button theme="primary" size="large" @click="navigateToPanorama" class="hero-btn primary">
              查看全景图
            </t-button>
            <t-button theme="default" size="large" @click="navigateToUserManagement" class="hero-btn secondary">
              用户管理
            </t-button>
          </div>
        </div>
      </div>
      
      <div class="excel-section">
        <ExcelUpload />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import Header from '@/components/Header.vue'
import ExcelUpload from '@/components/ExcelUpload.vue'

const router = useRouter()

const navigateToPanorama = () => {
  router.push('/panorama')
}

const navigateToUserManagement = () => {
  router.push('/user/list')
}
</script>

<style scoped>
.page-wrapper {
  min-height: 100vh;
  background: linear-gradient(180deg, #f7f8fa 0%, #ffffff 100%);
}

.main-content {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 120px 24px 80px;
}

.hero-section {
  text-align: center;
  padding: 60px 0 80px;
}

.hero-content {
  max-width: 640px;
  margin: 0 auto;
}

.hero-title {
  font-size: 44px;
  font-weight: 700;
  color: #1d2129;
  margin: 0 0 16px;
  letter-spacing: -0.02em;
  line-height: 1.2;
}

.hero-subtitle {
  font-size: 16px;
  color: #86909c;
  margin: 0 0 40px;
  line-height: 1.6;
}

.hero-buttons {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.hero-btn {
  border-radius: 8px;
  font-weight: 500;
  font-size: 15px;
  height: 44px;
  padding: 0 28px;
}

.hero-btn.primary {
  box-shadow: 0 4px 12px rgba(0, 82, 217, 0.15);
  transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
}

.hero-btn.primary:hover {
  box-shadow: 0 6px 16px rgba(0, 82, 217, 0.2);
  transform: translateY(-1px);
}

.hero-btn.secondary {
  color: #4e5969;
  border: 1px solid #e5e6eb;
  background: #ffffff;
  transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
}

.hero-btn.secondary:hover {
  color: #0052d9;
  border-color: #0052d9;
  background: rgba(0, 82, 217, 0.04);
}

.excel-section {
  margin-top: 40px;
}
</style>
