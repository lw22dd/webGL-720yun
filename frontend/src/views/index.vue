<template>
  <div class="page-wrapper">
    <Header />
    <main class="main-content">
      <div class="map-section">
        <ChinaMap ref="chinaMapRef" @space-click="handleSpaceClick" />
      </div>
      <SpaceDetailPanel
        :visible="panelVisible"
        :space="selectedSpace"
        @close="handlePanelClose"
        @scene-click="handleSceneClick"
      />
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
  height: 100vh;
  background: #f5f7fa;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.main-content {
  width: 100%;
  flex: 1;
  display: flex;
  overflow: hidden;
}

.map-section {
  position: relative;
  flex: 1;
  height: 100%;
  min-width: 0;
}
</style>
