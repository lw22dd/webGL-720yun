<template>
  <div class="editor-layout">
    <EditorToolbar />
    <div class="editor-body">
      <EditorCanvas />
      <EditorSidebar />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, provide, readonly } from 'vue'
import EditorToolbar from './EditorToolbar.vue'
import EditorCanvas from './EditorCanvas.vue'
import EditorSidebar from './EditorSidebar.vue'
import { ZOOM_LIMITS } from '@/utils/graph-editor/constants'

const isLayoutComputing = ref(false)
const canUndo = ref(false)
const canRedo = ref(false)

const graphRef = ref<any>(null)

provide('layoutComputingState', {
  isLayoutComputing,
  setLayoutComputing: (val: boolean) => { isLayoutComputing.value = val }
})

provide('historyState', {
  canUndo: readonly(canUndo),
  canRedo: readonly(canRedo),
  setCanUndo: (val: boolean) => { canUndo.value = val },
  setCanRedo: (val: boolean) => { canRedo.value = val }
})

provide('graphActions', {
  undo: () => {
    console.log('[EditorLayout] undo 被调用, graphRef:', graphRef.value)
    if (graphRef.value) {
      return graphRef.value.undo()
    }
    console.log('[EditorLayout] undo 失败: graphRef 为 null')
  },
  redo: () => {
    console.log('[EditorLayout] redo 被调用, graphRef:', graphRef.value)
    if (graphRef.value) {
      return graphRef.value.redo()
    }
    console.log('[EditorLayout] redo 失败: graphRef 为 null')
  },
  fitCanvas: () => graphRef.value?.zoomToFit({ maxScale: 1 }),
  setGraph: (g: any) => { graphRef.value = g },
  zoomIn: () => graphRef.value?.zoom(ZOOM_LIMITS.step),
  zoomOut: () => graphRef.value?.zoom(-ZOOM_LIMITS.step)
})
</script>

<style scoped>
.editor-layout {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.editor-body {
  flex: 1;
  display: flex;
  overflow: hidden;
  position: relative;
}
</style>
