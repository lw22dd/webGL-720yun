<script setup lang="ts">
import { ref } from 'vue'

defineProps<{ msg: string }>()

const count = ref(0)
</script>

<template>
  <h1 class="text-2xl font-bold mb-6">{{ msg }}</h1>

  <div class="bg-white p-6 rounded-lg shadow-md mb-6">
    <button 
      type="button" 
      @click="count++"
      class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
    >
      count is {{ count }}
    </button>
    <p class="mt-4 text-gray-600">
      Edit
      <code class="bg-gray-100 px-1 py-0.5 rounded">components/HelloWorld.vue</code> to test HMR
    </p>
  </div>

  <p class="mb-4 text-gray-700">
    Check out
    <a 
      href="https://vuejs.org/guide/quick-start.html#local" 
      target="_blank"
      class="text-blue-600 hover:underline"
    >
      create-vue
    </a>, the official Vue + Vite starter
  </p>
  <p class="mb-4 text-gray-700">
    Learn more about IDE Support for Vue in the
    <a
      href="https://vuejs.org/guide/scaling-up/tooling.html#ide-support"
      target="_blank"
      class="text-blue-600 hover:underline"
    >
      Vue Docs Scaling up Guide
    </a>.
  </p>
  <p class="text-gray-500">Click on the Vite and Vue logos to learn more</p>
</template>

<style scoped>
/* Tailwind CSS 已应用，不再需要自定义样式 */
</style>
