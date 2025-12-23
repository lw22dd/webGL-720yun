<template>
    <Header />
    <main class="container mx-auto p-16 ">
        <div class="mb-8">
            <el-button type="primary" size="large" @click="navigateToPanorama">
                查看全景图
            </el-button>
        </div>
        <ExcelUpload />
    </main>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import Header from '@/components/Header.vue'
import ExcelUpload from '@/components/ExcelUpload.vue'

const router = useRouter()

// 导航到全景图页面
const navigateToPanorama = () => {
  router.push('/panorama')
}
</script>

<style scoped>
/* Tailwind CSS 已应用，不再需要自定义样式 */
</style>
