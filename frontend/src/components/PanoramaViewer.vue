<template>
  <div class="panorama-container">
    <div ref="containerRef" class="panorama-canvas"></div>
    <div class="panorama-controls">
      <t-button theme="primary" @click="toggleFullscreen">全屏</t-button>
      <t-button @click="resetCamera">重置视角</t-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { Viewer } from '@photo-sphere-viewer/core';
import '@photo-sphere-viewer/core/index.css';

const containerRef = ref<HTMLDivElement | null>(null);
let viewer: Viewer | null = null;

const panoramaPath = new URL('@/assets/宝瓶口1.jpg', import.meta.url).href;

const initViewer = () => {
  if (!containerRef.value) return;

  viewer = new Viewer({
    container: containerRef.value,
    panorama: panoramaPath,
    defaultZoomLvl: 0,
    minFov: 30,
    maxFov: 90,
    navbar: false,
  });
};

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    containerRef.value?.requestFullscreen();
  } else {
    document.exitFullscreen();
  }
};

const resetCamera = () => {
  viewer?.rotate({ yaw: 0, pitch: 0 });
  viewer?.zoom(0);
};

onMounted(() => {
  initViewer();
});

onUnmounted(() => {
  viewer?.destroy();
});
</script>

<style scoped>
.panorama-container {
  width: 100%;
  height: 100vh;
  position: relative;
  overflow: hidden;
  background-color: #000;
}

.panorama-canvas {
  width: 100%;
  height: 100%;
}

.panorama-controls {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 100;
  display: flex;
  gap: 10px;
}
</style>
