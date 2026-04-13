<template>
  <div class="scene-strip">
    <div class="strip-container">
      <button class="strip-nav-btn prev" @click="scrollPrev" :disabled="!canScrollPrev">
        ‹
      </button>
      <div ref="scrollContainer" class="strip-scroll" @scroll="updateScrollState">
        <div class="strip-content">
          <div
            v-for="scene in scenes"
            :key="scene.id"
            class="strip-item"
            :class="{ active: scene.id === currentSceneId }"
            @click="handleSelect(scene)"
          >
            <div class="strip-thumb">
              <img 
                v-if="scene.thumbnail_url" 
                :src="scene.thumbnail_url" 
                :alt="scene.title"
              />
              <div v-else class="strip-thumb-placeholder">
                <t-icon name="image" />
              </div>
            </div>
            <div class="strip-info">
              <span class="strip-name">{{ scene.title }}</span>
            </div>
          </div>
        </div>
      </div>
      <button class="strip-nav-btn next" @click="scrollNext" :disabled="!canScrollNext">
        ›
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

interface SceneItem {
  id: number
  title: string
  thumbnail_url?: string
  scene_code?: string
}

defineProps<{
  scenes: SceneItem[]
  currentSceneId: number | null
}>()

const emit = defineEmits<{
  (e: 'sceneSelect', sceneId: number): void
}>()

const scrollContainer = ref<HTMLElement | null>(null)
const scrollPosition = ref(0)

const canScrollPrev = computed(() => scrollPosition.value > 0)
const canScrollNext = computed(() => {
  if (!scrollContainer.value) return false
  const { scrollWidth, clientWidth } = scrollContainer.value
  return scrollPosition.value < scrollWidth - clientWidth - 10
})

const updateScrollState = () => {
  if (scrollContainer.value) {
    scrollPosition.value = scrollContainer.value.scrollLeft
  }
}

const scrollPrev = () => {
  if (scrollContainer.value) {
    scrollContainer.value.scrollBy({ left: -200, behavior: 'smooth' })
  }
}

const scrollNext = () => {
  if (scrollContainer.value) {
    scrollContainer.value.scrollBy({ left: 200, behavior: 'smooth' })
  }
}

const handleSelect = (scene: SceneItem) => {
  emit('sceneSelect', scene.id)
}

onMounted(() => {
  updateScrollState()
})
</script>

<style scoped>
.scene-strip {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(180deg, transparent 0%, rgba(0,0,0,0.8) 100%);
  padding: 20px 0 16px;
  z-index: 100;
}

.strip-container {
  display: flex;
  align-items: center;
  max-width: 100%;
  padding: 0 16px;
}

.strip-nav-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  font-size: 20px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.strip-nav-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.2);
}

.strip-nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.strip-scroll {
  flex: 1;
  overflow-x: auto;
  overflow-y: hidden;
  scroll-behavior: smooth;
  scrollbar-width: none;
  -ms-overflow-style: none;
  margin: 0 8px;
}

.strip-scroll::-webkit-scrollbar {
  display: none;
}

.strip-content {
  display: flex;
  gap: 12px;
  padding: 4px 0;
}

.strip-item {
  flex-shrink: 0;
  width: 120px;
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid transparent;
  background: rgba(255, 255, 255, 0.1);
}

.strip-item:hover {
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.strip-item.active {
  border-color: #0EA5E9;
  box-shadow: 0 0 12px rgba(14, 165, 233, 0.4);
}

.strip-thumb {
  width: 100%;
  height: 68px;
  overflow: hidden;
}

.strip-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.strip-thumb-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.3);
}

.strip-info {
  padding: 8px;
  text-align: center;
}

.strip-name {
  display: block;
  color: #fff;
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
