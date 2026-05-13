<template>
  <t-dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    :header="hotspot?.title || '知识问答'"
    width="500px"
    :footer="false"
    attach=".panorama-container"
  >
    <div class="hotspot-quiz-content">
      <p class="quiz-question">{{ hotspot?.question }}</p>
      <div class="quiz-options">
        <div
          v-for="(option, index) in quizOptions"
          :key="index"
          class="quiz-option"
          :class="{ selected: modelValue === index }"
          @click="$emit('update:modelValue', Number(index))"
        >
          {{ String.fromCharCode(65 + Number(index)) }}. {{ String(option) }}
        </div>
      </div>
      <t-button
        theme="primary"
        block
        :disabled="modelValue === null"
        @click="$emit('submit')"
      >
        提交答案
      </t-button>
    </div>
  </t-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { HotspotForViewer } from '@/models/hotspot.model'

const props = defineProps<{
  visible: boolean
  hotspot: HotspotForViewer | null
  modelValue: number | null
}>()

defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'update:modelValue', value: number): void
  (e: 'submit'): void
}>()

const quizOptions = computed(() => {
  if (!props.hotspot?.options) return []
  try {
    return JSON.parse(props.hotspot.options)
  } catch {
    return []
  }
})
</script>

<style scoped>
.hotspot-quiz-content {
  padding: 8px 0;
}

.quiz-question {
  font-size: 16px;
  font-weight: 500;
  color: #1d2129;
  margin-bottom: 16px;
}

.quiz-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.quiz-option {
  padding: 12px 16px;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.quiz-option:hover {
  border-color: #0EA5E9;
  background: rgba(14, 165, 233, 0.05);
}

.quiz-option.selected {
  border-color: #0EA5E9;
  background: rgba(14, 165, 233, 0.1);
  color: #0EA5E9;
}
</style>