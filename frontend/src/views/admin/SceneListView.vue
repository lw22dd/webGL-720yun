<template>
  <div class="scene-list-view">
    <div class="breadcrumb">
      <span class="breadcrumb-item" @click="emit('back')">空间列表</span>
      <span class="breadcrumb-separator">/</span>
      <span class="breadcrumb-current">{{ spaceName }}</span>
      <span class="breadcrumb-separator">/</span>
      <span class="breadcrumb-current">场景列表</span>
    </div>

    <div class="header-actions">
      <div class="header-left">
        <t-input
          v-model="searchKeyword"
          placeholder="搜索场景"
          clearable
          @enter="handleSearch"
          @clear="handleSearch"
        >
          <template #suffix-icon>
            <t-icon name="search" />
          </template>
        </t-input>
      </div>
      <div class="header-right">
        <t-button theme="primary" @click="openAddDialog">
          <template #icon><t-icon name="add" /></template>
          新增场景
        </t-button>
      </div>
    </div>

    <div v-if="uploadingTasks.length > 0" class="upload-progress-section">
      <t-card title="上传进度" :bordered="false">
        <div v-for="task in uploadingTasks" :key="task.uploadId" class="upload-task-item">
          <div class="task-info">
            <span class="file-name">{{ task.fileName }}</span>
            <span class="file-size">{{ uploadService.formatFileSize(task.fileSize) }}</span>
          </div>
          <t-progress
            :percentage="task.percentage"
            :theme="task.status === 'failed' ? 'error' : 'primary'"
            :label="`${task.percentage.toFixed(1)}%`"
          />
          <div class="task-status">
            <t-tag :theme="getStatusTheme(task.status)">
              {{ getStatusText(task.status) }}
            </t-tag>
            <span v-if="task.speed" class="speed">{{ task.speed }}</span>
          </div>
        </div>
      </t-card>
    </div>

    <div class="scene-table-container">
      <t-table
        :data="sceneList"
        :columns="columns"
        :pagination="paginationConfig"
        row-key="id"
        hover
        stripe
        @page-change="handlePageChange"
      >
        <template #thumbnail_url="{ row }">
          <div class="scene-thumbnail" v-if="row.scene_code">
            <img
              :src="getPreviewUrl(row.scene_code)"
              :alt="row.title"
              @error="handleThumbError($event)"
            />
          </div>
          <div v-else class="no-thumbnail">
            <t-icon name="image" />
          </div>
        </template>

        <template #source_info="{ row }">
          <span v-if="row.source_width && row.source_height">
            {{ row.source_width }} x {{ row.source_height }}
          </span>
          <span v-else>-</span>
        </template>

        <template #slice_status="{ row }">
          <t-tag :theme="getSliceStatusTheme(row.slice_status)">
            {{ getSliceStatusText(row.slice_status) }}
          </t-tag>
        </template>

        <template #status="{ row }">
          <span :class="['status-badge', row.status === 1 ? 'active' : 'disabled']">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </span>
        </template>

        <template #operations="{ row }">
          <t-space>
            <t-button
              v-if="!row.source_url"
              theme="primary"
              variant="text"
              size="small"
              @click="openUploadDialog(row)"
            >
              上传
            </t-button>
            <t-button
              v-if="row.source_url"
              theme="primary"
              variant="text"
              size="small"
              @click="openPreviewDialog(row)"
            >
              预览
            </t-button>
            <t-button
              v-if="row.slice_status === 'completed' || row.slice_status === 'ready'"
              theme="success"
              variant="text"
              size="small"
              @click="handleViewPanorama(row)"
            >
              查看
            </t-button>
            <t-button
              theme="primary"
              variant="text"
              size="small"
              @click="handleEditScene(row)"
            >
              编辑
            </t-button>
            <t-button
              theme="danger"
              variant="text"
              size="small"
              @click="openDeleteDialog(row.id)"
            >
              删除
            </t-button>
          </t-space>
        </template>
      </t-table>
    </div>

    <div v-if="!loading && sceneList.length === 0" class="empty-state">
      <t-icon name="folder-open" size="64" />
      <p>暂无场景数据</p>
      <t-button theme="primary" @click="openAddDialog">新增场景</t-button>
    </div>

    <t-dialog
      v-model:visible="deleteDialogVisible"
      header="确认删除"
      width="400px"
      @confirm="confirmDelete"
      @cancel="deleteDialogVisible = false"
    >
      <p>确定要删除该场景吗？此操作不可撤销。</p>
    </t-dialog>

    <t-dialog
      v-model:visible="uploadDialogVisible"
      header="上传全景图"
      width="500px"
      :confirm-btn="{ content: '开始上传', loading: uploadLoading }"
      @confirm="handleUpload"
      @close="resetUploadForm"
    >
      <div class="upload-dialog-content">
        <p class="upload-scene-info">
          场景：<strong>{{ uploadingScene?.title }}</strong>
        </p>
        <t-form>
          <t-form-item label="全景图文件">
            <t-upload
              v-model="panoramaFiles"
              :auto-upload="false"
              :multiple="false"
              accept="image/jpeg,image/png"
              :size-limit="{ size: 500, unit: 'MB' }"
              theme="file-input"
              placeholder="选择全景图文件（最大500MB）"
            />
          </t-form-item>
        </t-form>
      </div>
    </t-dialog>

    <t-dialog
      v-model:visible="addDialogVisible"
      :header="isEditing ? '编辑场景' : '新增场景'"
      width="600px"
      :confirm-btn="{ content: isEditing ? '保存' : '下一步', loading: submitLoading }"
      @confirm="handleSubmit"
      @close="resetForm"
    >
      <t-form :data="formData" :rules="formRules" ref="formRef">
        <t-form-item label="场景标题" name="title">
          <t-input v-model="formData.title" placeholder="请输入场景标题" />
        </t-form-item>
        <t-form-item label="场景编码" name="scene_code">
          <t-input
            v-model="formData.scene_code"
            placeholder="请输入场景编码（英文唯一）"
            :disabled="isEditing"
          />
        </t-form-item>
        <t-form-item label="初始FOV">
          <t-input-number v-model="formData.initial_fov" :min="30" :max="150" />
        </t-form-item>
        <t-form-item label="初始俯仰角">
          <t-input-number v-model="formData.initial_pitch" :min="-90" :max="90" />
        </t-form-item>
        <t-form-item label="初始偏航角">
          <t-input-number v-model="formData.initial_yaw" :min="-180" :max="180" />
        </t-form-item>
        <t-form-item label="经度">
          <t-input-number v-model="formData.longitude" :min="-180" :max="180" :decimal-places="7" />
        </t-form-item>
        <t-form-item label="纬度">
          <t-input-number v-model="formData.latitude" :min="-90" :max="90" :decimal-places="7" />
        </t-form-item>
        <t-form-item label="排序">
          <t-input-number v-model="formData.sort_order" />
        </t-form-item>
        <t-form-item v-if="isEditing" label="状态">
          <t-switch v-model="formData.status" :true-value="1" :false-value="0" />
        </t-form-item>
      </t-form>
    </t-dialog>

    <ScenePreviewDialog
      v-model:visible="previewDialogVisible"
      :scene="previewScene"
      @close="previewDialogVisible = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import SceneApi from '@/services/api/scene.api'
import uploadService from '@/services/uploadService'
import { useSceneStore } from '@/stores/scene/scene.store'
import type { SpaceListItem } from '@/models/space.model'
import type { SceneListItem } from '@/models/scene.model'
import ScenePreviewDialog from '@/components/admin/ScenePreviewDialog.vue'

const router = useRouter()

interface Props {
  space: SpaceListItem | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'back'): void
}>()

const sceneStore = useSceneStore()
const uploadingTasks = computed(() => sceneStore.uploadingTasks)

const loading = ref(false)
const searchKeyword = ref('')
const sceneList = ref<SceneListItem[]>([])
const spaceName = computed(() => props.space?.name || '未知空间')

const paginationConfig = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showJumper: true
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const formRef = ref()
const addDialogVisible = ref(false)
const isEditing = ref(false)
const submitLoading = ref(false)
const editingSceneId = ref<number | null>(null)

const deleteDialogVisible = ref(false)
const deletingId = ref<number | null>(null)

const uploadDialogVisible = ref(false)
const uploadLoading = ref(false)
const uploadingScene = ref<SceneListItem | null>(null)
const panoramaFiles = ref<any[]>([])

const previewDialogVisible = ref(false)
const previewScene = ref<SceneListItem | null>(null)

const formData = reactive({
  title: '',
  scene_code: '',
  initial_fov: 100,
  initial_pitch: 0,
  initial_yaw: 0,
  longitude: 0,
  latitude: 0,
  sort_order: 0,
  status: 1
})

const formRules = {
  title: [{ required: true, message: '请输入场景标题' }],
  scene_code: [{ required: true, message: '请输入场景编码' }]
}

const columns = [
  {
    colKey: 'id',
    title: 'ID',
    width: 80,
    align: 'center'
  },
  {
    colKey: 'thumbnail_url',
    title: '缩略图',
    width: 100,
    align: 'center'
  },
  {
    colKey: 'title',
    title: '场景标题',
    minWidth: 150,
    ellipsis: true
  },
  {
    colKey: 'scene_code',
    title: '场景编码',
    minWidth: 120,
    ellipsis: true
  },
  {
    colKey: 'source_info',
    title: '分辨率',
    width: 120,
    align: 'center'
  },
  {
    colKey: 'slice_status',
    title: '切片状态',
    width: 100,
    align: 'center'
  },
  {
    colKey: 'status',
    title: '状态',
    width: 80,
    align: 'center'
  },
  {
    colKey: 'operations',
    title: '操作',
    width: 200,
    align: 'center',
    fixed: 'right'
  }
]

const getPreviewUrl = (sceneCode: string) => {
  const base = import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'
  return `${base}/api/v1/res/previews/${sceneCode}`
}

const handleThumbError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.style.display = 'none'
}

const loadSceneList = async () => {
  if (!props.space) return

  try {
    loading.value = true
    const result = await SceneApi.getSceneList({
      space_id: props.space.id,
      keyword: searchKeyword.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    if (result.code === 200 && result.data) {
      console.log('Received scene list:', result.data.scenes)
      
      sceneList.value = result.data.scenes || []
      pagination.total = result.data.page_info.total
      paginationConfig.total = result.data.page_info.total
    }
  } catch (error) {
    MessagePlugin.error('获取场景列表失败')
  } finally {
    loading.value = false
  }
}

const openAddDialog = () => {
  isEditing.value = false
  editingSceneId.value = null
  resetForm()
  addDialogVisible.value = true
}

const handleEditScene = (scene: SceneListItem) => {
  isEditing.value = true
  editingSceneId.value = scene.id
  Object.assign(formData, {
    title: scene.title,
    scene_code: scene.scene_code,
    initial_fov: (scene as any).initial_fov || 100,
    initial_pitch: (scene as any).initial_pitch || 0,
    initial_yaw: (scene as any).initial_yaw || 0,
    longitude: (scene as any).longitude || 0,
    latitude: (scene as any).latitude || 0,
    sort_order: scene.sort_order || 0,
    status: scene.status
  })
  addDialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate()
  if (valid !== true) return

  submitLoading.value = true

  try {
    if (isEditing.value && editingSceneId.value) {
      const result = await SceneApi.updateScene(editingSceneId.value, {
        title: formData.title,
        initial_fov: formData.initial_fov,
        initial_pitch: formData.initial_pitch,
        initial_yaw: formData.initial_yaw,
        longitude: formData.longitude,
        latitude: formData.latitude,
        sort_order: formData.sort_order,
        status: formData.status
      })
      if (result.code === 200) {
        MessagePlugin.success('编辑场景成功')
        addDialogVisible.value = false
        loadSceneList()
      } else {
        MessagePlugin.error(result.msg || '编辑失败')
      }
    } else {
      const result = await SceneApi.createScene({
        space_id: props.space!.id,
        title: formData.title,
        scene_code: formData.scene_code,
        initial_fov: formData.initial_fov,
        initial_pitch: formData.initial_pitch,
        initial_yaw: formData.initial_yaw,
        longitude: formData.longitude,
        latitude: formData.latitude,
        sort_order: formData.sort_order
      })
      if (result.code === 200) {
        MessagePlugin.success('创建场景成功，现在可以上传全景图')
        addDialogVisible.value = false
        loadSceneList()
      } else {
        MessagePlugin.error(result.msg || '创建失败')
      }
    }
  } catch (error: any) {
    MessagePlugin.error(error.message || '操作失败')
  } finally {
    submitLoading.value = false
  }
}

const resetForm = () => {
  Object.assign(formData, {
    title: '',
    scene_code: '',
    initial_fov: 100,
    initial_pitch: 0,
    initial_yaw: 0,
    longitude: 0,
    latitude: 0,
    sort_order: 0,
    status: 1
  })
  panoramaFiles.value = []
}

const openUploadDialog = (scene: SceneListItem) => {
  uploadingScene.value = scene
  uploadDialogVisible.value = true
}

const handleUpload = async () => {
  if (!uploadingScene.value || panoramaFiles.value.length === 0) {
    MessagePlugin.warning('请选择全景图文件')
    return
  }

  uploadLoading.value = true

  try {
    const file = panoramaFiles.value[0].raw
    await uploadService.uploadFile(file, {
      space_id: props.space!.id,
      scene_code: uploadingScene.value.scene_code,
      title: uploadingScene.value.title,
      onComplete: () => {
        MessagePlugin.success('上传成功')
        uploadDialogVisible.value = false
        resetUploadForm()
        loadSceneList()
      },
      onError: (error) => {
        MessagePlugin.error(error.message || '上传失败')
      }
    })
  } catch (error: any) {
    MessagePlugin.error(error.message || '上传失败')
  } finally {
    uploadLoading.value = false
  }
}

const resetUploadForm = () => {
  panoramaFiles.value = []
  uploadingScene.value = null
}

const openDeleteDialog = (id: number) => {
  deletingId.value = id
  deleteDialogVisible.value = true
}

const confirmDelete = async () => {
  if (deletingId.value === null) return

  try {
    const result = await SceneApi.deleteScene(deletingId.value)
    if (result.code === 200) {
      MessagePlugin.success('删除成功')
      loadSceneList()
    } else {
      MessagePlugin.error(result.msg || '删除失败')
    }
  } catch (error) {
    MessagePlugin.error('删除失败')
  } finally {
    deleteDialogVisible.value = false
    deletingId.value = null
  }
}

const openPreviewDialog = (scene: SceneListItem) => {
  previewScene.value = scene
  previewDialogVisible.value = true
}

const handleViewPanorama = (scene: SceneListItem) => {
  router.push(`/panorama?scene=${scene.scene_code}`)
}

const handlePageChange = (context: { current: number; pageSize: number }) => {
  pagination.currentPage = context.current
  paginationConfig.current = context.current
  pagination.pageSize = context.pageSize
  paginationConfig.pageSize = context.pageSize
  loadSceneList()
}

const handleSearch = () => {
  pagination.currentPage = 1
  paginationConfig.current = 1
  loadSceneList()
}

const getStatusTheme = (status: string) => {
  const themes: Record<string, string> = {
    pending: 'default',
    uploading: 'primary',
    merging: 'warning',
    completed: 'success',
    failed: 'error',
    paused: 'default'
  }
  return themes[status] || 'default'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '等待中',
    uploading: '上传中',
    merging: '合并中',
    completed: '已完成',
    failed: '失败',
    paused: '已暂停'
  }
  return texts[status] || status
}

const getSliceStatusTheme = (status: string) => {
  const themes: Record<string, string> = {
    pending: 'default',
    slicing: 'warning',
    completed: 'success',
    failed: 'error'
  }
  return themes[status] || 'default'
}

const getSliceStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待处理',
    slicing: '切片中',
    completed: '已完成',
    failed: '失败'
  }
  return texts[status] || status
}

watch(() => props.space, () => {
  if (props.space) {
    loadSceneList()
  }
}, { immediate: true })

import { watch } from 'vue'

onMounted(() => {
  sceneStore.initWebSocket()
})

defineExpose({
  loadSceneList
})
</script>

<style scoped>
.scene-list-view {
  width: 100%;
  height: 100%;
  min-height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
  font-size: 14px;
}

.breadcrumb-item {
  color: var(--td-text-color-secondary);
  cursor: pointer;
  transition: color 0.2s;
}

.breadcrumb-item:hover {
  color: var(--td-brand-color);
}

.breadcrumb-separator {
  color: var(--td-text-color-disabled);
}

.breadcrumb-current {
  color: var(--td-text-color-primary);
  font-weight: 500;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left :deep(.t-input) {
  width: 240px;
}

.upload-progress-section {
  margin-bottom: 16px;
}

.upload-task-item {
  padding: 12px 0;
  border-bottom: 1px solid #e5e6eb;
}

.upload-task-item:last-child {
  border-bottom: none;
}

.task-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.file-name {
  font-weight: 500;
  color: #1d2129;
}

.file-size {
  color: #86909c;
  font-size: 12px;
}

.task-status {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
}

.speed {
  color: #86909c;
  font-size: 12px;
}

.scene-table-container {
  flex: 1;
  overflow: auto;
}

.scene-thumbnail {
  width: 80px;
  height: 40px;
  border-radius: 4px;
  overflow: hidden;
  background: #f5f5f5;
  margin: 0 auto;
}

.scene-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-thumbnail {
  width: 80px;
  height: 40px;
  border-radius: 4px;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #86909c;
  margin: 0 auto;
}

.status-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-badge.active {
  background: rgba(0, 168, 112, 0.1);
  color: #00a870;
}

.status-badge.disabled {
  background: rgba(245, 63, 63, 0.1);
  color: #f53f3f;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--td-text-color-secondary);
  gap: 16px;
}

.empty-state p {
  margin: 0;
  font-size: 16px;
}

.upload-dialog-content {
  padding: 8px 0;
}

.upload-scene-info {
  margin-bottom: 16px;
  color: var(--td-text-color-secondary);
}

.upload-scene-info strong {
  color: var(--td-text-color-primary);
}
</style>
