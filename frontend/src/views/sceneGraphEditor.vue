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
import EditorToolbar from './EditorToolbar.vue'
import EditorCanvas from './EditorCanvas.vue'
import EditorSidebar from './EditorSidebar.vue'
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
