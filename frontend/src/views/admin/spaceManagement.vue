<template>
  <div class="space-management-container">
    <SpaceListView
      v-if="!currentSpace"
      ref="spaceListRef"
      @spaceClick="handleSpaceClick"
    />
    <SceneListView
      v-else
      ref="sceneListRef"
      :space="currentSpace"
      @back="handleBack"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import SpaceListView from './SpaceListView.vue'
import SceneListView from './SceneListView.vue'
import type { SpaceListItem } from '@/models/space.model'

const route = useRoute()
const router = useRouter()

const spaceListRef = ref<InstanceType<typeof SpaceListView> | null>(null)
const sceneListRef = ref<InstanceType<typeof SceneListView> | null>(null)
const currentSpace = ref<SpaceListItem | null>(null)

const handleSpaceClick = (space: SpaceListItem) => {
  currentSpace.value = space
  router.replace({ query: { spaceId: space.id.toString() } })
}

const handleBack = () => {
  currentSpace.value = null
  router.replace({ query: {} })
}

const initFromRoute = () => {
  const spaceId = route.query.spaceId
  if (spaceId) {
    console.log('Initializing with spaceId:', spaceId)
  }
}

initFromRoute()
</script>

<style scoped>
.space-management-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>
