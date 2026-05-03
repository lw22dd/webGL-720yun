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
          <div v-if="getSceneProgress(row.scene_code)" class="scene-progress">
            <!-- 上传进度显示进度条 -->
            <template v-if="getSceneProgress(row.scene_code)!.type === 'upload'">
              <t-progress
                :percentage="getSceneProgress(row.scene_code)!.percentage"
                :status="getSceneProgress(row.scene_code)!.status"
                :label="getSceneProgress(row.scene_code)!.label"
                :theme="getSceneProgressTheme(row.scene_code)"
                :color="getSceneProgressColor(row.scene_code)"
                size="small"
              />
              <div class="progress-actions">
                <span class="progress-message">{{ getSceneProgress(row.scene_code)!.message }}</span>
                <t-space>
                  <t-button
                    v-if="getSceneProgress(row.scene_code)!.status === 'active'"
                    theme="primary"
                    variant="text"
                    size="small"
                    @click="handlePauseUpload(row.scene_code)"
                  >
                    暂停
                  </t-button>
                  <t-button
                    v-if="getSceneProgress(row.scene_code)!.status === 'warning'"
                    theme="primary"
                    variant="text"
                    size="small"
                    @click="handleResumeUpload(row.scene_code)"
                  >
                    继续
                  </t-button>
                  <t-button
                    theme="primary"
                    variant="text"
                    size="small"
                    @click="handleCancelUpload(row.scene_code)"
                  >
                    取消
                  </t-button>
                </t-space>
              </div>
            </template>
            <!-- 切片进度显示百分比+文字 -->
            <template v-else>
              <div class="slice-status-display">
                <span class="slice-percentage">{{ getSceneProgress(row.scene_code)!.percentage }}%</span>
                <span class="slice-message">{{ getSceneProgress(row.scene_code)!.message }}</span>
              </div>
            </template>
          </div>
          <t-tag v-else :theme="getSliceStatusTheme(row.slice_status)">
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
              theme="primary"
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
              theme="primary"
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
      <div class="empty-icon">
        <t-icon name="folder-open" size="64" />
      </div>
      <p class="empty-title">暂无场景数据</p>
      <p class="empty-desc">点击右上角"新增场景"按钮创建第一个场景</p>
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
      v-model:visible="addDialogVisible"
      :header="isEditing ? '编辑场景' : '新增场景'"
      width="600px"
      :confirm-btn="{ content: isEditing ? '保存' : '下一步', loading: submitLoading }"
      @confirm="handleSubmit"
      @close="resetForm"
    >
      <t-form :data="formData" :rules="formRules" ref="formRef">
        <t-form-item label="场景标题" name="title">
          <div style="display: flex; gap: 8px; width: 100%;">
            <t-input v-model="formData.title" placeholder="请输入场景标题" style="flex: 1;" />
            <t-button theme="primary" variant="text" size="small" @click="openLocationPicker">
              <template #icon><t-icon name="location" /></template>
              搜索地点
            </t-button>
          </div>
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

    <t-dialog
      v-model:visible="locationPickerVisible"
      header="选择地点"
      width="680px"
      :footer="false"
    >
      <LocationPicker
        :initial-keyword="locationPickerKeyword"
        @select="handleLocationSelect"
      />
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import SceneApi from '@/services/api/scene.api'
import uploadService from '@/services/uploadService'
import wsClient from '@/services/websocket.service'
import { useSceneStore } from '@/stores/scene/scene.store'
import type { SpaceListItem } from '@/models/space.model'
import type { SceneListItem } from '@/models/scene.model'
import type { CompleteData } from '@/models/upload.model'
import ScenePreviewDialog from '@/components/admin/ScenePreviewDialog.vue'
//import LocationPicker from '@/components/admin/LocationPicker.vue'

interface SceneProgress {
  percentage: number
  status: 'active' | 'success' | 'error' | 'warning'
  label: string
  message: string
  type: 'upload' | 'slice'
}

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

// 监听切片任务变化，用于显示进度
const sliceTasksMap = computed(() => {
  const map = new Map<string, SceneProgress>()
  sceneStore.sliceTasks.forEach((task) => {
    const statusMap: Record<string, 'active' | 'success' | 'error' | 'warning'> = {
      'pending': 'warning',
      'slicing': 'warning',
      'completed': 'success',
      'failed': 'error'
    }
    // 显示格式: 75% - 开始上传瓦片...
    const label = `${task.progress}%`
    const message = task.message || '处理中...'
    map.set(task.sceneCode, {
      percentage: task.progress,
      status: statusMap[task.status] || 'warning',
      label: label,
      message: message,
      type: 'slice'
    })
  })
  return map
})

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

const uploadingScene = ref<SceneListItem | null>(null)
const uploadIdMap = ref<Map<string, string>>(new Map())

const previewDialogVisible = ref(false)
const previewScene = ref<SceneListItem | null>(null)

const locationPickerVisible = ref(false)
const locationPickerKeyword = ref('')

const sceneProgressMap = ref<Map<string, SceneProgress>>(new Map())

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
  title: [{ required: true, message: '请输入场景标题' }]
}

const getSceneProgress = (sceneCode: string) => {
  // 首先检查本地进度（上传进度）
  const localProgress = sceneProgressMap.value.get(sceneCode)
  if (localProgress) {
    return localProgress
  }

  // 然后检查切片任务进度（使用计算属性确保响应式）
  return sliceTasksMap.value.get(sceneCode) || null
}

const setSceneProgress = (sceneCode: string, progress: SceneProgress) => {
  sceneProgressMap.value.set(sceneCode, progress)
}

const getSceneProgressTheme = (sceneCode: string) => {
  const progress = getSceneProgress(sceneCode)
  if (!progress) return 'default'
  switch (progress.status) {
    case 'success': return 'success'
    case 'error': return 'danger'
    case 'warning': return 'warning'
    default: return 'primary'
  }
}

const getSceneProgressColor = (sceneCode: string) => {
  const progress = getSceneProgress(sceneCode)
  if (!progress) return ''
  switch (progress.status) {
    case 'success': return '#00A870'
    case 'error': return '#F53F3F'
    case 'warning': return '#FF9900'
    default: return '#0052D9'
  }
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

const previewTimestamps = ref<Record<string, number>>({})

const getPreviewUrl = (sceneCode: string) => {
  const base = import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'
  const timestamp = previewTimestamps.value[sceneCode]
  if (timestamp) {
    return `${base}/api/v1/res/previews/${sceneCode}?t=${timestamp}`
  }
  return `${base}/api/v1/res/previews/${sceneCode}`
}

const handleThumbError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.style.display = 'none'
  const container = img.parentElement
  if (container && !container.querySelector('.thumb-error-placeholder')) {
    const placeholder = document.createElement('div')
    placeholder.className = 'thumb-error-placeholder'
    placeholder.innerHTML = '<t-icon name="image" /><span class="placeholder-text">无封面</span>'
    container.appendChild(placeholder)
  }
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
}

const openUploadDialog = (scene: SceneListItem) => {
  uploadingScene.value = scene
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/jpeg,image/png'
  input.onchange = async (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return

    setSceneProgress(scene.scene_code, {
      percentage: 0,
      status: 'active',
      label: '0%',
      message: '准备上传...',
      type: 'upload'
    })

    try {
      const uploadResult = await uploadService.uploadFile(file, {
        space_id: props.space!.id,
        scene_code: scene.scene_code,
        title: scene.title,
        onInit: (uploadId) => {
          uploadIdMap.value.set(scene.scene_code, uploadId)
        },
        onProgress: (progress) => {
          handleUploadProgress({
            percentage: progress.percentage,
            uploaded_chunks: progress.uploaded_chunks,
            total_chunks: progress.total_chunks
          })
        }
      })

      setSceneProgress(scene.scene_code, {
        percentage: 100,
        status: 'success',
        label: '完成',
        message: '上传完成，正在关联场景...',
        type: 'upload'
      })

      const updateResult = await SceneApi.updateScene(scene.id, {
        file_id: uploadResult.file_id
      })

      if (updateResult.code === 200 && updateResult.data && updateResult.data.task_id) {
        sceneStore.addSliceTask({
          taskId: updateResult.data.task_id,
          sceneId: scene.id,
          sceneCode: scene.scene_code,
          status: 'pending',
          progress: 0,
          stage: 'pending',
          message: '等待切片处理...'
        })
      }

      // 清除上传进度，让 getSceneProgress 从 sliceTasks 读取
      sceneProgressMap.value.delete(scene.scene_code)

      loadSceneList()
      uploadingScene.value = null
      uploadIdMap.value.delete(scene.scene_code)
    } catch (error: any) {
      MessagePlugin.error(error.message || '操作失败')
      setSceneProgress(scene.scene_code, {
        percentage: 0,
        status: 'error',
        label: '失败',
        message: error.message || '上传失败',
        type: 'upload'
      })
      uploadingScene.value = null
      uploadIdMap.value.delete(scene.scene_code)
    }
  }
  input.click()
}

function handleUploadProgress(data: any) {
  if (!uploadingScene.value) return
  const sceneCode = uploadingScene.value.scene_code
  setSceneProgress(sceneCode, {
    percentage: data.percentage || 0,
    status: 'active',
    label: `${(data.percentage || 0).toFixed(1)}%`,
    message: `上传中 ${data.uploaded_chunks}/${data.total_chunks}`,
    type: 'upload'
  })
}

function handleMergeProgress(data: any) {
  if (!uploadingScene.value) return
  const sceneCode = uploadingScene.value.scene_code
  setSceneProgress(sceneCode, {
    percentage: data.percentage || 0,
    status: 'warning',
    label: `${(data.percentage || 0).toFixed(0)}%`,
    message: data.message || '合并中...',
    type: 'upload'
  })
}

function handleUploadComplete(data: any) {
  if (!uploadingScene.value) return
  const sceneCode = uploadingScene.value.scene_code
  setSceneProgress(sceneCode, {
    percentage: 100,
    status: 'success',
    label: '完成',
    message: '上传完成',
    type: 'upload'
  })
}

function handlePauseUpload(sceneCode: string) {
  const uploadId = uploadIdMap.value.get(sceneCode)
  if (!uploadId) return
  
  uploadService.pauseUpload(uploadId)
  
  const progress = sceneProgressMap.value.get(sceneCode)
  if (progress) {
    setSceneProgress(sceneCode, {
      ...progress,
      status: 'warning',
      label: '已暂停',
      message: '上传已暂停'
    })
  }
  MessagePlugin.info('上传已暂停')
}

function handleResumeUpload(sceneCode: string) {
  const uploadId = uploadIdMap.value.get(sceneCode)
  if (!uploadId) return
  
  uploadService.resumeUpload(uploadId)
  
  const progress = sceneProgressMap.value.get(sceneCode)
  if (progress) {
    setSceneProgress(sceneCode, {
      ...progress,
      status: 'active',
      label: `${progress.percentage.toFixed(1)}%`,
      message: '继续上传...'
    })
  }
  MessagePlugin.info('继续上传')
}

async function handleCancelUpload(sceneCode: string) {
  const uploadId = uploadIdMap.value.get(sceneCode)
  if (!uploadId) return
  
  try {
    await uploadService.cancelUpload(uploadId)
    uploadIdMap.value.delete(sceneCode)
    
    setSceneProgress(sceneCode, {
      percentage: 0,
      status: 'error',
      label: '已取消',
      message: '上传已取消',
      type: 'upload'
    })
    uploadingScene.value = null
    MessagePlugin.info('上传已取消')
  } catch (error) {
    MessagePlugin.error('取消上传失败')
  }
}

function handleSliceProgress(data: any) {
  console.log('[View] Slice progress:', data)
}

function handleSliceComplete(data: any) {
  console.log('[View] Slice complete:', data)
  if (data.scene_code) {
    // 先更新进度为100%，让用户看到完成状态
    const task = Array.from(sceneStore.sliceTasks.values()).find(t => t.sceneCode === data.scene_code)
    if (task) {
      sceneStore.updateSliceTask(task.taskId, {
        progress: 100,
        status: 'completed',
        message: '切片完成'
      })
    }

    // 延迟1秒后清除进度显示并刷新列表
    setTimeout(() => {
      // 清除该场景的进度显示
      sceneProgressMap.value.delete(data.scene_code)
      // 从 store 中移除已完成的任务
      const taskToRemove = Array.from(sceneStore.sliceTasks.values()).find(t => t.sceneCode === data.scene_code)
      if (taskToRemove) {
        sceneStore.removeSliceTask(taskToRemove.taskId)
      }
      // 刷新列表显示 ready 状态
      loadSceneList()
      // 延迟2秒后更新预览图时间戳并刷新，确保封面图已生成
      setTimeout(() => {
        previewTimestamps.value[data.scene_code] = Date.now()
        loadSceneList()
      }, 2000)
    }, 1000)
  }
}

function handleSliceError(data: any) {
  console.log('[View] Slice error:', data)
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

const openLocationPicker = () => {
  locationPickerKeyword.value = formData.title || props.space?.name || ''
  locationPickerVisible.value = true
}

const handleLocationSelect = (location: { name: string; lng: number; lat: number }) => {
  formData.longitude = location.lng
  formData.latitude = location.lat
  if (!formData.title) {
    formData.title = location.name
  }
  locationPickerVisible.value = false
  MessagePlugin.success(`已选择地点：${location.name}`)
}

const handleViewPanorama = (scene: SceneListItem) => {
  sessionStorage.setItem('panoramaFrom', router.currentRoute.value.fullPath)
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
  console.log('[SceneListView] Component mounted, initializing WebSocket...')
  sceneStore.initWebSocket()

  wsClient.on('progress', handleUploadProgress)
  wsClient.on('merge_progress', handleMergeProgress)
  wsClient.on('complete', handleUploadComplete)
  wsClient.on('slice_progress', handleSliceProgress)
  wsClient.on('slice_complete', handleSliceComplete)
  wsClient.on('slice_error', handleSliceError)

  console.log('[SceneListView] WebSocket event handlers registered')
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

.thumb-error-placeholder {
  width: 80px;
  height: 40px;
  border-radius: 4px;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #86909c;
  margin: 0 auto;
  gap: 2px;
}

.placeholder-text {
  font-size: 10px;
}

.thumb-error-placeholder :deep(.t-icon) {
  font-size: 16px;
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
  padding: 100px 20px;
  color: var(--td-text-color-secondary);
  gap: 12px;
  background: var(--td-bg-color-container);
  border-radius: 8px;
  margin-top: 16px;
}

.empty-icon {
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--td-bg-color-secondarycontainer);
  border-radius: 50%;
  margin-bottom: 8px;
}

.empty-icon :deep(.t-icon) {
  color: var(--td-text-color-placeholder);
}

.empty-title {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: var(--td-text-color-primary);
}

.empty-desc {
  margin: 0;
  font-size: 14px;
  color: var(--td-text-color-secondary);
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

.scene-progress {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 120px;
}

.scene-progress :deep(.t-progress) {
  margin-bottom: 0;
}

.slice-status-display {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.slice-percentage {
  font-weight: 600;
  color: var(--td-brand-color);
  min-width: 36px;
}

.slice-message {
  color: var(--td-text-color-secondary);
  font-size: 12px;
}

.progress-actions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.progress-message {
  font-size: 11px;
  color: var(--td-text-color-secondary);
  text-align: center;
}
</style>
