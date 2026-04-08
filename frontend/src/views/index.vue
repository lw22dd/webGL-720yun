<template>
  <div class="page-wrapper">
    <Header />
    <main class="main-content">
      <div class="map-section">
        <ChinaMap ref="chinaMapRef" @space-click="handleSpaceClick" />
        <SpaceDetailPanel
          :visible="panelVisible"
          :space="selectedSpace"
          @close="handlePanelClose"
          @scene-click="handleSceneClick"
        />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Header from '@/components/Header.vue'
import ChinaMap from '@/components/ChinaMap.vue'
import SpaceDetailPanel from '@/components/SpaceDetailPanel.vue'
import type { MockSpace, MockScene } from '@/utils/mockData'

const router = useRouter()
const chinaMapRef = ref<InstanceType<typeof ChinaMap> | null>(null)
const panelVisible = ref(false)
const selectedSpace = ref<MockSpace | null>(null)

const handleSpaceClick = (space: MockSpace) => {
  selectedSpace.value = space
  panelVisible.value = true
}

const handlePanelClose = () => {
  panelVisible.value = false
}

const handleSceneClick = (scene: MockScene) => {
  router.push(`/panorama?scene=${scene.scene_code}`)
}
</script>

<style scoped>
.page-wrapper {
  min-height: 100vh;
  background: #0a0e27;
}

.main-content {
  width: 100%;
  height: calc(100vh - 64px);
  padding-top: 64px;
}

.map-section {
  position: relative;
  width: 100%;
  height: 100%;
}
</style>
