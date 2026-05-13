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
                theme="line"
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

    <FormDialog ref="formDialogRef">
      <template #default="{ formData, submitLoading, closeDialog, isEdit }">
        <t-form
          ref="sceneFormRef"
          :data="formData"
          :rules="sceneFormRules"
          label-width="100px"
          @submit="handleFormSubmit"
        >
          <t-form-item label="场景标题" name="title">
            <t-input v-model="formData.title" placeholder="请输入场景名称，或在右侧地图选点自动填充" />
          </t-form-item>

          <t-form-item label="经度" name="longitude">
            <t-input v-model="formData.longitude" placeholder="从地图中拾取坐标" disabled>
              <template #suffix>
                <span class="coord-unit">°E</span>
              </template>
            </t-input>
          </t-form-item>

          <t-form-item label="纬度" name="latitude">
            <t-input v-model="formData.latitude" placeholder="从地图中拾取坐标" disabled>
              <template #suffix>
                <span class="coord-unit">°N</span>
              </template>
            </t-input>
          </t-form-item>

          <t-form-item label="初始FOV" name="initial_fov">
            <t-input-number
              v-model="formData.initial_fov"
              placeholder="视角范围 (30-150)"
              :min="30"
              :max="150"
              style="width: 100%"
            />
          </t-form-item>

          <t-form-item label="初始俯仰" name="initial_pitch">
            <t-input-number
              v-model="formData.initial_pitch"
              placeholder="俯仰角度 (-90-90)"
              :min="-90"
              :max="90"
              style="width: 100%"
            />
          </t-form-item>

          <t-form-item label="初始偏航" name="initial_yaw">
            <t-input-number
              v-model="formData.initial_yaw"
              placeholder="偏航角度 (-180-180)"
              :min="-180"
              :max="180"
              style="width: 100%"
            />
          </t-form-item>

          <t-form-item label="排序权重" name="sort_order">
            <t-input-number v-model="formData.sort_order" placeholder="数字越大越靠前" style="width: 100%" />
          </t-form-item>

          <t-form-item label="启用状态" name="status">
            <t-switch v-model="formData.status" :label="['启用', '禁用']" />
          </t-form-item>

          <t-form-item style="margin-top: 32px">
            <t-space>
              <t-button type="submit" theme="primary" size="large" :loading="submitLoading">
                {{ isEdit ? '保存修改' : '创建场景' }}
              </t-button>
              <t-button variant="outline" size="large" @click="closeDialog">取消</t-button>
            </t-space>
          </t-form-item>
        </t-form>
      </template>
    </FormDialog>

    <ScenePreviewDialog
      v-model:visible="previewDialogVisible"
      :scene="previewScene"
      @close="previewDialogVisible = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import type { FormRule } from 'tdesign-vue-next'
import SceneApi from '@/apis/scene.api'
import uploadService from '@/services/uploadService'
import wsClient from '@/services/websocket.service'
import { useSceneStore } from '@/stores/scene/scene.store'
import type { SpaceListItem } from '@/models/space.model'
import type { SceneListItem } from '@/models/scene.model'
import FormDialog from '@/components/admin/FormDialog.vue'
import ScenePreviewDialog from '@/components/admin/ScenePreviewDialog.vue'

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
// ... (uploadingTasks and sliceTasksMap remain the same)
const uploadingTasks = computed(() => sceneStore.uploadingTasks)

const sliceTasksMap = computed(() => {
  const map = new Map<string, SceneProgress>()
  sceneStore.sliceTasks.forEach((task) => {
    const statusMap: Record<string, 'active' | 'success' | 'error' | 'warning'> = {
      'pending': 'warning',
      'slicing': 'warning',
      'completed': 'success',
      'failed': 'error'
    }
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

const formDialogRef = ref()
const sceneFormRef = ref()
const isEditing = ref(false)
const editingSceneId = ref<number | null>(null)

const deleteDialogVisible = ref(false)
const deletingId = ref<number | null>(null)

const previewDialogVisible = ref(false)
const previewScene = ref<SceneListItem | null>(null)

const sceneProgressMap = ref<Map<string, SceneProgress>>(new Map())

const sceneFormRules: Record<string, FormRule[]> = {
  title: [
    { required: true, message: '请输入场景标题', trigger: 'blur' }
  ]
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
  formDialogRef.value?.openAddDialog('新增场景')
}

const handleEditScene = (scene: SceneListItem) => {
  isEditing.value = true
  editingSceneId.value = scene.id
  formDialogRef.value?.openEditDialog({
    title: scene.title,
    initial_fov: (scene as any).initial_fov || 100,
    initial_pitch: (scene as any).initial_pitch || 0,
    initial_yaw: (scene as any).initial_yaw || 0,
    longitude: (scene as any).longitude || 0,
    latitude: (scene as any).latitude || 0,
    sort_order: scene.sort_order || 0,
    status: scene.status
  }, scene.id, '编辑场景')
}

const handleFormSubmit = async () => {
  if (!sceneFormRef.value) return
  const valid = await sceneFormRef.value.validate()
  if (!valid) return

  const data = formDialogRef.value?.formData || {}
  formDialogRef.value?.setSubmitLoading(true)
  try {
    if (isEditing.value && editingSceneId.value) {
      const result = await SceneApi.updateScene(editingSceneId.value, {
        title: data.title,
        initial_fov: Number(data.initial_fov),
        initial_pitch: Number(data.initial_pitch),
        initial_yaw: Number(data.initial_yaw),
        longitude: Number(data.longitude),
        latitude: Number(data.latitude),
        sort_order: Number(data.sort_order),
        status: data.status
      })
      if (result.code === 200) {
        MessagePlugin.success('编辑场景成功')
        formDialogRef.value?.closeDialog()
        loadSceneList()
      } else {
        MessagePlugin.error(result.msg || '编辑失败')
      }
    } else {
      const result = await SceneApi.createScene({
        space_id: props.space!.id,
        title: data.title,
        initial_fov: Number(data.initial_fov),
        initial_pitch: Number(data.initial_pitch),
        initial_yaw: Number(data.initial_yaw),
        longitude: Number(data.longitude),
        latitude: Number(data.latitude),
        sort_order: Number(data.sort_order)
      })
      if (result.code === 200) {
        MessagePlugin.success('创建场景成功，请点击"上传"按钮添加全景图')
        formDialogRef.value?.closeDialog()
        loadSceneList()
      } else {
        MessagePlugin.error(result.msg || '创建失败')
      }
    }
  } catch (error: any) {
    MessagePlugin.error(error.message || '操作失败')
  } finally {
    formDialogRef.value?.setSubmitLoading(false)
  }
}

// ... (rest of the upload and websocket logic remains unchanged)
const getSceneProgress = (sceneCode: string) => {
  const localProgress = sceneProgressMap.value.get(sceneCode)
  if (localProgress) return localProgress
  return sliceTasksMap.value.get(sceneCode) || null
}

const setSceneProgress = (sceneCode: string, progress: SceneProgress) => {
  sceneProgressMap.value.set(sceneCode, progress)
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
    colKey: 'created_at',
    title: '创建时间',
    width: 140,
    align: 'center',
    cell: (h: any, { row }: { row: SceneListItem }) => {
      const date = new Date(row.created_at)
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hours = String(date.getHours()).padStart(2, '0')
      const minutes = String(date.getMinutes()).padStart(2, '0')
      return h('div', {
        style: 'font-size: 12px; line-height: 1.5;'
      }, [
        h('div', `${year}-${month}-${day}`),
        h('div', { style: 'color: #86909c;' }, `${hours}:${minutes}`)
      ])
    }
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
    width: 150,
    ellipsis: true
  },
  {
    colKey: 'scene_code',
    title: '场景编码',
    width: 120,
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

const openUploadDialog = (scene: SceneListItem) => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/jpeg,image/png'
  input.onchange = async (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    setSceneProgress(scene.scene_code, { percentage: 0, status: 'active', label: '0%', message: '准备上传...', type: 'upload' })
    try {
      const uploadResult = await uploadService.uploadFile(file, {
        space_id: props.space!.id,
        scene_code: scene.scene_code,
        title: scene.title,
        onInit: (id) => uploadIdMap.value.set(scene.scene_code, id),
        onProgress: (p) => handleUploadProgress({ percentage: p.percentage, uploaded_chunks: p.uploaded_chunks, total_chunks: p.total_chunks })
      })
      setSceneProgress(scene.scene_code, { percentage: 100, status: 'success', label: '完成', message: '关联场景...', type: 'upload' })
      const updateResult = await SceneApi.updateScene(scene.id, { file_id: uploadResult.file_id })
      if (updateResult.code === 200 && updateResult.data?.task_id) {
        sceneStore.addSliceTask({ taskId: updateResult.data.task_id, sceneId: scene.id, sceneCode: scene.scene_code, status: 'pending', progress: 0, stage: 'pending', message: '等待处理...' })
      }
      sceneProgressMap.value.delete(scene.scene_code)
      loadSceneList()
    } catch (error: any) {
      MessagePlugin.error(error.message || '上传失败')
      setSceneProgress(scene.scene_code, { percentage: 0, status: 'error', label: '失败', message: error.message, type: 'upload' })
    }
  }
  input.click()
}

const uploadIdMap = ref<Map<string, string>>(new Map())

function handleUploadProgress(data: any) {
  const scenes = sceneList.value.filter(s => sceneProgressMap.value.has(s.scene_code))
  if (scenes.length === 0) return
  const sceneCode = scenes[0].scene_code
  setSceneProgress(sceneCode, {
    percentage: data.percentage || 0,
    status: 'active',
    label: `${(data.percentage || 0).toFixed(1)}%`,
    message: `上传中 ${data.uploaded_chunks}/${data.total_chunks}`,
    type: 'upload'
  })
}

function handleMergeProgress(data: any) {
  const scenes = sceneList.value.filter(s => sceneProgressMap.value.has(s.scene_code))
  if (scenes.length === 0) return
  setSceneProgress(scenes[0].scene_code, { percentage: data.percentage || 0, status: 'warning', label: `${(data.percentage || 0).toFixed(0)}%`, message: data.message || '合并中...', type: 'upload' })
}

function handleUploadComplete(data: any) {
  const scenes = sceneList.value.filter(s => sceneProgressMap.value.has(s.scene_code))
  if (scenes.length === 0) return
  setSceneProgress(scenes[0].scene_code, { percentage: 100, status: 'success', label: '完成', message: '上传完成', type: 'upload' })
}

function handlePauseUpload(sceneCode: string) {
  const id = uploadIdMap.value.get(sceneCode)
  if (id) { uploadService.pauseUpload(id); MessagePlugin.info('上传已暂停') }
}

function handleResumeUpload(sceneCode: string) {
  const id = uploadIdMap.value.get(sceneCode)
  if (id) { uploadService.resumeUpload(id); MessagePlugin.info('继续上传') }
}

async function handleCancelUpload(sceneCode: string) {
  const id = uploadIdMap.value.get(sceneCode)
  if (id) {
    try { await uploadService.cancelUpload(id); uploadIdMap.value.delete(sceneCode); MessagePlugin.info('上传已取消'); loadSceneList() }
    catch (e) { MessagePlugin.error('取消失败') }
  }
}

function handleSliceProgress(data: any) {}
function handleSliceComplete(data: any) {
  if (data.scene_code) {
    const task = Array.from(sceneStore.sliceTasks.values()).find(t => t.sceneCode === data.scene_code)
    if (task) sceneStore.updateSliceTask(task.taskId, { progress: 100, status: 'completed', message: '完成' })
    setTimeout(() => {
      sceneProgressMap.value.delete(data.scene_code)
      const t = Array.from(sceneStore.sliceTasks.values()).find(t => t.sceneCode === data.scene_code)
      if (t) sceneStore.removeSliceTask(t.taskId)
      loadSceneList()
      setTimeout(() => { previewTimestamps.value[data.scene_code] = Date.now(); loadSceneList() }, 2000)
    }, 1000)
  }
}
function handleSliceError(data: any) {}

const openDeleteDialog = (id: number) => {
  deletingId.value = id
  deleteDialogVisible.value = true
}

const confirmDelete = async () => {
  if (deletingId.value === null) return
  try {
    const result = await SceneApi.deleteScene(deletingId.value)
    if (result.code === 200) { MessagePlugin.success('删除成功'); loadSceneList() }
    else MessagePlugin.error(result.msg || '删除失败')
  } catch (error) { MessagePlugin.error('删除失败') }
  finally { deleteDialogVisible.value = false; deletingId.value = null }
}

const openPreviewDialog = (scene: SceneListItem) => {
  previewScene.value = scene
  previewDialogVisible.value = true
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

const getSliceStatusTheme = (status: string) => {
  const themes: Record<string, string> = { pending: 'default', slicing: 'warning', completed: 'success', failed: 'error' }
  return themes[status] || 'default'
}

const getSliceStatusText = (status: string) => {
  const texts: Record<string, string> = { pending: '待处理', slicing: '切片中', completed: '已完成', failed: '失败' }
  return texts[status] || status
}

watch(() => props.space, () => { if (props.space) loadSceneList() }, { immediate: true })

onMounted(() => {
  sceneStore.initWebSocket()
  wsClient.on('progress', handleUploadProgress)
  wsClient.on('merge_progress', handleMergeProgress)
  wsClient.on('complete', handleUploadComplete)
  wsClient.on('slice_progress', handleSliceProgress)
  wsClient.on('slice_complete', handleSliceComplete)
  wsClient.on('slice_error', handleSliceError)
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

.scene-table-container :deep(.t-table__content table) {
  table-layout: fixed;
}

.scene-table-container :deep(.t-table__content th),
.scene-table-container :deep(.t-table__content td) {
  overflow: hidden;
  text-overflow: ellipsis;
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

.coord-unit {
  color: #86909c;
  font-size: 12px;
  font-weight: 500;
}

.progress-message {
  font-size: 11px;
  color: var(--td-text-color-secondary);
  text-align: center;
}
</style>
