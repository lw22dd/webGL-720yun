<template>
  <div class="space-management-container" v-loading="initializing">
    <SpaceListView
      v-if="!currentSpace && !initializing"
      ref="spaceListRef"
      @spaceClick="handleSpaceClick"
    />
    <SceneListView
      v-if="currentSpace"
      ref="sceneListRef"
      :space="currentSpace"
      @back="handleBack"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import SpaceListView from './SpaceListView.vue'
import SceneListView from './SceneListView.vue'
import SpaceApi from '@/services/api/space.api'
import type { SpaceListItem } from '@/models/space.model'

const route = useRoute()
const router = useRouter()

const spaceListRef = ref<InstanceType<typeof SpaceListView> | null>(null)
const sceneListRef = ref<InstanceType<typeof SceneListView> | null>(null)
const currentSpace = ref<SpaceListItem | null>(null)
const loading = ref(false)
const initializing = ref(false)

const handleSpaceClick = (space: SpaceListItem) => {
  currentSpace.value = space
  router.replace({ query: { spaceId: space.id.toString() } })
}

const handleBack = () => {
  currentSpace.value = null
  router.replace({ query: {} })
}

const loadSpaceDetail = async (spaceId: number) => {
  try {
    initializing.value = true
    const result = await SpaceApi.getSpaceDetail(spaceId)
    if (result.code === 200 && result.data) {
      currentSpace.value = result.data as SpaceListItem
    } else {
      MessagePlugin.error(result.msg || '获取空间详情失败')
    }
  } catch (error) {
    MessagePlugin.error('获取空间详情失败')
  } finally {
    initializing.value = false
  }
}

const initFromRoute = async () => {
  const spaceId = route.query.spaceId
  if (spaceId) {
    const id = parseInt(spaceId as string, 10)
    if (!isNaN(id)) {
      await loadSpaceDetail(id)
    }
  }
}

onMounted(() => {
  initFromRoute()
})
</script>

<style scoped>
.space-management-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>
