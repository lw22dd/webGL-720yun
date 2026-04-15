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
  const spaceId = Number(route.params.id)
  if (spaceId) {
    sceneStore.loadScene(spaceId)
  }
})
</script>

<style scoped>
</style>
