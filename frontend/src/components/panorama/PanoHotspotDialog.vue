<template>
  <t-dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    :header="hotspot?.title"
    width="480px"
    :footer="false"
    attach=".panorama-container"
  >
    <div class="hotspot-info-content">
      <img
        v-if="hotspot?.media_type === 'image' && hotspot?.media_url"
        :src="hotspot.media_url"
        class="info-media"
        alt="热点配图"
      />
      <p v-if="hotspot?.content" class="info-text">{{ hotspot.content }}</p>
      <p v-if="!hotspot?.content && !hotspot?.media_url" class="info-empty">暂无内容</p>
    </div>
  </t-dialog>
</template>

<script setup lang="ts">
import type { HotspotForViewer } from '@/models/hotspot.model'

defineProps<{
  visible: boolean
  hotspot: HotspotForViewer | null
}>()

defineEmits<{
  (e: 'update:visible', value: boolean): void
}>()
</script>

<style scoped>
.hotspot-info-content {
  padding: 8px 0;
}

.info-media {
  max-width: 100%;
  border-radius: 8px;
  margin-bottom: 12px;
  display: block;
}

.info-text {
  color: #4e5969;
  line-height: 1.6;
  margin: 0;
  white-space: pre-wrap;
}

.info-empty {
  color: #86909c;
  text-align: center;
  padding: 16px 0;
  margin: 0;
}
</style>