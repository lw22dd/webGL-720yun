<template>
  <div class="property-panel">
    <div class="panel-header">
      <div class="panel-title">
        <t-icon :name="node ? 'location' : 'link'" />
        {{ node ? '节点属性' : '连线属性' }}
      </div>
      <t-button variant="text" size="small" @click="$emit('close')">
        <template #icon><t-icon name="close" /></template>
      </t-button>
    </div>

    <div class="panel-content">
      <template v-if="node">
        <div class="property-section">
          <div class="section-title">基本信息</div>
          <div class="property-item">
            <label>ID</label>
            <span>{{ node.id }}</span>
          </div>
          <div class="property-item">
            <label>场景编码</label>
            <span>{{ node.scene_code }}</span>
          </div>
          <div class="property-item">
            <label>标题</label>
            <t-input v-model="editForm.title" size="small" />
          </div>
          <div class="property-item">
            <label>浏览次数</label>
            <span>{{ node.view_count }}</span>
          </div>
        </div>

        <div class="property-section">
          <div class="section-title">坐标信息</div>
          <div class="property-item">
            <label>经度</label>
            <t-input-number
              v-model="editForm.longitude"
              size="small"
              :decimalPlaces="6"
              :min="-180"
              :max="180"
            />
          </div>
          <div class="property-item">
            <label>纬度</label>
            <t-input-number
              v-model="editForm.latitude"
              size="small"
              :decimalPlaces="6"
              :min="-90"
              :max="90"
            />
          </div>
          <div class="property-item">
            <label>状态</label>
            <t-tag :theme="node.has_position ? 'success' : 'warning'" size="small">
              {{ node.has_position ? '已定位' : '未定位' }}
            </t-tag>
          </div>
        </div>

        <div class="panel-actions">
          <t-button theme="primary" size="small" block @click="handleSaveNode">
            保存修改
          </t-button>
        </div>
      </template>

      <template v-if="edge">
        <div class="property-section">
          <div class="section-title">连线信息</div>
          <div class="property-item">
            <label>ID</label>
            <span>{{ edge.id }}</span>
          </div>
          <div class="property-item">
            <label>热点标题</label>
            <t-input v-model="editForm.hotspot_title" size="small" />
          </div>
          <div class="property-item">
            <label>源场景</label>
            <span>{{ edge.source_id }}</span>
          </div>
          <div class="property-item">
            <label>目标场景</label>
            <span>{{ edge.target_id }}</span>
          </div>
        </div>

        <div class="panel-actions">
          <t-button theme="primary" size="small" block @click="handleSaveEdge">
            保存修改
          </t-button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { SceneNodeData, EdgeData } from '@/models/SceneModel'

interface Props {
  node: SceneNodeData | null
  edge: EdgeData | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'updateNode', updates: Partial<SceneNodeData>): void
  (e: 'updateEdge', updates: Partial<EdgeData>): void
  (e: 'close'): void
}>()

const editForm = reactive({
  title: '',
  longitude: 0,
  latitude: 0,
  hotspot_title: ''
})

watch(() => props.node, (newNode) => {
  if (newNode) {
    editForm.title = newNode.title
    editForm.longitude = newNode.longitude
    editForm.latitude = newNode.latitude
  }
}, { immediate: true })

watch(() => props.edge, (newEdge) => {
  if (newEdge) {
    editForm.hotspot_title = newEdge.hotspot_title
  }
}, { immediate: true })

const handleSaveNode = () => {
  emit('updateNode', {
    title: editForm.title,
    longitude: editForm.longitude,
    latitude: editForm.latitude
  })
}

const handleSaveEdge = () => {
  emit('updateEdge', {
    hotspot_title: editForm.hotspot_title
  })
}
</script>

<style scoped>
.property-panel {
  position: absolute;
  right: 16px;
  top: 80px;
  width: 280px;
  background-color: var(--td-bg-color-container);
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  z-index: 100;
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--td-component-border);
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--td-text-color-primary);
}

.panel-content {
  padding: 16px;
  max-height: calc(100vh - 200px);
  overflow-y: auto;
}

.property-section {
  margin-bottom: 16px;
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--td-text-color-secondary);
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.property-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.property-item label {
  font-size: 13px;
  color: var(--td-text-color-secondary);
  flex-shrink: 0;
  width: 70px;
}

.property-item span {
  font-size: 13px;
  color: var(--td-text-color-primary);
}

.property-item :deep(.t-input),
.property-item :deep(.t-input-number) {
  flex: 1;
  max-width: 160px;
}

.panel-actions {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--td-component-border);
}
</style>
