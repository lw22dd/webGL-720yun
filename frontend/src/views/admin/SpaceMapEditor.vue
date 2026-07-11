<template>
  <div class="space-map-editor-page">
    <SceneInSpace :is-admin="true" @scene-click="handleSceneClick" />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import SceneInSpace from '@/components/space-detail/SceneInSpace.vue'
import type { SceneMarkerData } from '@/composables/useSceneMarkers'

const router = useRouter()

function handleSceneClick(scene: SceneMarkerData) {
  sessionStorage.setItem('panoramaFrom', router.currentRoute.value.fullPath)
  router.push(`/panorama?id=${scene.id}&scene=${scene.scene_code}`)
}
</script>

<style scoped>
.space-map-editor-page {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
</style>
