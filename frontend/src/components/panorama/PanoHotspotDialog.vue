<template>
  <t-dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    :header="hotspot?.title"
    width="400px"
    :footer="false"
    attach=".panorama-container"
  >
    <div class="hotspot-info-content">
      <p>{{ hotspot?.content }}</p>
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

.hotspot-info-content p {
  color: #4e5969;
  line-height: 1.6;
  margin: 0;
}
</style>