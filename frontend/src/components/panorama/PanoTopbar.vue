<template>
  <div class="pano-topbar">
    <div class="pano-topbar-left">
      <t-button theme="default" variant="outline" @click="$emit('back')" class="back-btn">
        ← 返回
      </t-button>
      <span class="pano-title">{{ title }}</span>
    </div>
    <div class="pano-topbar-right">
      <t-button
        v-if="isAdmin"
        :theme="editMode ? 'primary' : 'default'"
        variant="outline"
        @click="$emit('update:editMode', !editMode)"
        class="action-btn"
      >
        {{ editMode ? '退出编辑' : '编辑热点' }}
      </t-button>
      <t-button theme="default" variant="outline" @click="$emit('minimap')" class="action-btn">
        小地图
      </t-button>
      <t-button theme="default" variant="outline" @click="$emit('fullscreen')" class="action-btn">
        ⛶ 全屏
      </t-button>
    </div>
    <div v-if="editMode" class="edit-hint-bar">
      编辑模式：点击空白处添加热点，点击热点进行编辑
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  isAdmin: boolean
  editMode: boolean
}>()

defineEmits<{
  (e: 'back'): void
  (e: 'minimap'): void
  (e: 'fullscreen'): void
  (e: 'update:editMode', value: boolean): void
}>()
</script>

<style scoped>
.pano-topbar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: linear-gradient(180deg, rgba(0,0,0,0.6) 0%, transparent 100%);
  z-index: 100;
}

.pano-topbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.pano-topbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.pano-title {
  color: #fff;
  font-size: 18px;
  font-weight: 500;
}

.back-btn,
.action-btn {
  background: rgba(255, 255, 255, 0.1) !important;
  border: 1px solid rgba(255, 255, 255, 0.2) !important;
  color: #fff !important;
  backdrop-filter: blur(10px);
}

.back-btn:hover,
.action-btn:hover {
  background: rgba(255, 255, 255, 0.2) !important;
}

.edit-hint-bar {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  padding: 6px 24px;
  background: rgba(245, 158, 11, 0.9);
  color: #fff;
  font-size: 13px;
  text-align: center;
  backdrop-filter: blur(10px);
}
</style>