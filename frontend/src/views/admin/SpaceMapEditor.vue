<template>
  <div class="space-map-editor-page">
    <SpaceDetailMap :is-admin="true" @scene-click="handleSceneClick" />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import SpaceDetailMap from '@/components/space-detail/SpaceDetailMap.vue'
import type { SceneMarkerData } from '@/composables/useSceneMarkers'

const router = useRouter()

function handleSceneClick(scene: SceneMarkerData) {
  router.push(`/panorama?scene=${scene.scene_code}`)
}
</script>

<style scoped>
.space-map-editor-page {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
</style>
