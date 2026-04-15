<template>
  <EditorLayout />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import EditorLayout from '@/components/editor/EditorLayout.vue'
import { useGraphSceneStore } from '@/stores/graphSceneStore'

const route = useRoute()
const sceneStore = useGraphSceneStore()

onMounted(() => {
  const sceneId = route.params.id as string
  if (sceneId) {
    sceneStore.loadScene(sceneId)
  }
})
</script>

<style scoped>
</style>
