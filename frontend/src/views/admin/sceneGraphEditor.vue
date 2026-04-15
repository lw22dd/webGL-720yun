<template>
  <div class="graph-editor-wrapper">
    <EditorLayout />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import EditorLayout from '@/components/editor/EditorLayout.vue'
import { useGraphSceneStore } from '@/stores/graph/scene.store'

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
.graph-editor-wrapper {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>
