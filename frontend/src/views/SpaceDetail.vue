<template>
  <div class="space-detail-page">
    <SceneInSpace @scene-click="handleSceneClick" />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import SceneInSpace from '@/components/space-detail/SceneInSpace.vue'
import type { SceneMarkerData } from '@/composables/useSceneMarkers'

const router = useRouter()

function handleSceneClick(scene: SceneMarkerData) {
  sessionStorage.setItem('panoramaFrom', router.currentRoute.value.fullPath)
  router.push(`/panorama?scene=${scene.scene_code}`)
}
</script>

<style scoped>
.space-detail-page {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
</style>
