<template>
  <div class="pending-list">
    <div class="list-header">
      <template v-if="!searchVisible">
        <t-icon name="search" @click="showSearch" class="search-icon" />
        <span>待处理节点 ({{ sceneStore.pendingNodes.length }})</span>
      </template>
      <template v-else>
        <t-input
          v-model="searchQuery"
          ref="searchInputRef"
          placeholder="搜索节点..."
          @blur="hideSearch"
          @enter="hideSearch"
          autofocus
          class="search-input"
        />
      </template>
    </div>

    <div class="list-body" v-if="filteredNodes.length > 0">
      <PendingNodeItem
        v-for="node in filteredNodes"
        :key="node.id"
        :node="node"
      />
    </div>

    <div class="list-empty" v-else>
      <t-icon name="search" size="48px" />
      <p v-if="searchQuery">未找到匹配的节点</p>
      <p v-else>所有节点已放置完成</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import { useGraphSceneStore } from '@/stores/graphSceneStore'
import PendingNodeItem from './PendingNodeItem.vue'

const sceneStore = useGraphSceneStore()

const searchVisible = ref(false)
const searchQuery = ref('')
const searchInputRef = ref<any>(null)

const filteredNodes = computed(() => {
  if (!searchQuery.value.trim()) {
    return sceneStore.pendingNodes
  }
  const query = searchQuery.value.toLowerCase()
  return sceneStore.pendingNodes.filter(node =>
    node.title.toLowerCase().includes(query) ||
    node.scene_code.toLowerCase().includes(query)
  )
})

const showSearch = async () => {
  searchVisible.value = true
  await nextTick()
  searchInputRef.value?.focus()
}

const hideSearch = () => {
  searchVisible.value = false
  searchQuery.value = ''
}
</script>

<style scoped>
.pending-list {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.list-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px;
  font-size: var(--font-size-base);
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}

.search-icon {
  color: #6b7280;
  cursor: pointer;
  transition: color 0.2s;
}

.search-icon:hover {
  color: #374151;
}

.search-input {
  flex: 1;
}

.list-body {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.list-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px 20px;
  color: #9ca3af;
}

.list-empty p {
  font-size: var(--font-size-sm);
}
</style>
