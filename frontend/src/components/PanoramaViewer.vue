<template>
  <div class="panorama-container">
    <div ref="containerRef" class="panorama-canvas"></div>
    <div class="panorama-controls">
      <el-button type="primary" @click="toggleFullscreen">全屏</el-button>
      <el-button @click="resetCamera">重置视角</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';

// 容器引用
const containerRef = ref<HTMLDivElement | null>(null);

// Three.js核心对象
let scene: THREE.Scene;
let camera: THREE.PerspectiveCamera;
let renderer: THREE.WebGLRenderer;
let controls: OrbitControls;
let sphere: THREE.Mesh;
let animationId: number;

// 全景图路径
const panoramaPath = new URL('@/assets/宝瓶口1.jpg', import.meta.url).href;

// 初始化Three.js场景
const initScene = () => {
  if (!containerRef.value) return;

  // 获取容器尺寸
  const container = containerRef.value;
  const width = container.clientWidth;
  const height = container.clientHeight;

  // 创建场景
  scene = new THREE.Scene();

  // 创建相机
  camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000);
  camera.position.set(0, 0, 0.1); // 相机放在球体内部

  // 创建渲染器
  renderer = new THREE.WebGLRenderer({ antialias: true });
  renderer.setSize(width, height);
  renderer.setPixelRatio(window.devicePixelRatio);
  container.appendChild(renderer.domElement);

  // 创建全景图球体
  const geometry = new THREE.SphereGeometry(500, 60, 40);
  // 反转球体法线方向，使纹理朝向内部
  geometry.scale(-1, 1, 1);

  // 加载纹理
  const textureLoader = new THREE.TextureLoader();
  const texture = textureLoader.load(panoramaPath);

  // 创建材质
  const material = new THREE.MeshBasicMaterial({ map: texture });

  // 创建网格并添加到场景
  sphere = new THREE.Mesh(geometry, material);
  scene.add(sphere);

  // 创建轨道控制器
  controls = new OrbitControls(camera, renderer.domElement);
  controls.enableDamping = true; // 启用阻尼效果
  controls.dampingFactor = 0.05;
  controls.enableZoom = true;
  controls.enablePan = false;
  controls.minDistance = 0.1;
  controls.maxDistance = 1;

  // 开始动画循环
  animate();

  // 添加窗口大小变化监听
  window.addEventListener('resize', onWindowResize);
};

// 动画循环
const animate = () => {
  animationId = requestAnimationFrame(animate);
  controls.update();
  renderer.render(scene, camera);
};

// 窗口大小变化处理
const onWindowResize = () => {
  if (!containerRef.value) return;

  const container = containerRef.value;
  const width = container.clientWidth;
  const height = container.clientHeight;

  camera.aspect = width / height;
  camera.updateProjectionMatrix();
  renderer.setSize(width, height);
};

// 切换全屏
const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    containerRef.value?.requestFullscreen();
  } else {
    document.exitFullscreen();
  }
};

// 重置相机位置
const resetCamera = () => {
  camera.position.set(0, 0, 0.1);
  controls.reset();
};

// 组件挂载时初始化
onMounted(() => {
  initScene();
});

// 组件卸载时清理
onUnmounted(() => {
  window.removeEventListener('resize', onWindowResize);
  cancelAnimationFrame(animationId);
  if (containerRef.value && renderer) {
    containerRef.value.removeChild(renderer.domElement);
  }
  // 释放资源
  if (sphere) {
    sphere.geometry.dispose();
    if (Array.isArray(sphere.material)) {
      sphere.material.forEach(material => material.dispose());
    } else {
      sphere.material.dispose();
    }
  }
  if (renderer) {
    renderer.dispose();
  }
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