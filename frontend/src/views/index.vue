<template>
    <Header />
    <main class="container mx-auto p-16 ">
        <ExcelUpload />
    </main>
</template>

<script setup lang="ts">
import Header from '@/components/Header.vue'
import ExcelUpload from '@/components/ExcelUpload.vue'
</script>

<style scoped>
/* Tailwind CSS 已应用，不再需要自定义样式 */
</style>
