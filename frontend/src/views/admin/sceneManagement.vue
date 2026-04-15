<template>
  <div class="scene-management">
    <div class="header-actions">
      <t-space>
        <t-select
          v-model="selectedSpaceId"
          placeholder="选择空间"
          clearable
          style="width: 200px"
          @change="handleSpaceChange"
        >
          <t-option
            v-for="space in spaceList"
            :key="space.id"
            :value="space.id"
            :label="space.name"
          />
        </t-select>
        <t-button theme="primary" @click="openAddDialog" :disabled="!selectedSpaceId">
          <template #icon><t-icon name="add" /></template>
          新增场景
        </t-button>
      </t-space>
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

    <AdminTable
      ref="tableRef"
      :data="sceneList"
      :columns="columns"
      :pagination="paginationConfig"
      delete-item-name="场景"
      @page-change="handlePageChange"
      @delete="handleDeleteScene"
    >
      <template #operations="{ row }">
        <t-space>
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
    </AdminTable>

    <t-dialog
      v-model:visible="addDialogVisible"
      :header="isEditing ? '编辑场景' : '新增场景'"
      width="600px"
      :confirm-btn="{ content: isEditing ? '保存' : '上传', loading: submitLoading }"
      @confirm="handleSubmit"
      @close="resetForm"
    >
      <t-form :data="formData" :rules="formRules" ref="formRef">
        <t-form-item label="场景标题" name="title">
          <t-input v-model="formData.title" placeholder="请输入场景标题" />
        </t-form-item>
        <t-form-item label="场景编码" name="scene_code">
          <t-input v-model="formData.scene_code" placeholder="请输入场景编码（英文唯一）" :disabled="isEditing" />
        </t-form-item>
        <t-form-item v-if="!isEditing" label="全景图" name="panorama">
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
        <t-form-item label="状态">
          <t-switch v-model="formData.status" :true-value="1" :false-value="0" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import AdminTable from '@/components/admin/AdminTable.vue'
import SpaceApi from '@/services/api/space.api'
import SceneApi from '@/services/api/scene.api'
import uploadService from '@/services/uploadService'
import { useSceneStore } from '@/stores/scene/scene.store'
import type { SpaceListItem } from '@/models/space.model'
import type { SceneListItem } from '@/models/scene.model'

const sceneStore = useSceneStore()
const uploadingTasks = computed(() => sceneStore.uploadingTasks)

const selectedSpaceId = ref<number | null>(null)
const spaceList = ref<SpaceListItem[]>([])
const sceneList = ref<SceneListItem[]>([])

const paginationConfig = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  defaultCurrent: 1,
  defaultPageSize: 10,
  showJumper: true
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const tableRef = ref()
const formRef = ref()
const addDialogVisible = ref(false)
const isEditing = ref(false)
const submitLoading = ref(false)
const editingSceneId = ref<number | null>(null)
const panoramaFiles = ref<any[]>([])

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

const columns = computed(() => [
  {
    colKey: 'id',
    title: 'ID',
    width: 80,
    align: 'center' as const,
    fixed: 'left' as const
  },
  {
    colKey: 'thumbnail_url',
    title: '缩略图',
    width: 100,
    cell: ({ row }: { row?: any }) => {
      if (!row?.thumbnail_url) return '-'
      return `<img src="${row.thumbnail_url}" style="width: 80px; height: 40px; object-fit: cover; border-radius: 4px;" />`
    }
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
    colKey: 'source_width',
    title: '分辨率',
    width: 120,
    cell: ({ row }: { row?: any }) => {
      if (!row?.source_width) return '-'
      return `${row.source_width} x ${row.source_height}`
    }
  },
  {
    colKey: 'view_count',
    title: '浏览次数',
    width: 100,
    align: 'center' as const
  },
  {
    colKey: 'status',
    title: '状态',
    width: 80,
    align: 'center' as const,
    cell: ({ row }: { row?: any }) => {
      if (!row) return '-'
      return row.status === 1 ? '启用' : '禁用'
    }
  },
  {
    colKey: 'created_at',
    title: '创建时间',
    width: 180,
    ellipsis: true
  },
  {
    colKey: 'operations',
    title: '操作',
    width: 140,
    align: 'center' as const,
    fixed: 'right' as const
  }
])

const loadSpaceList = async () => {
  try {
    const result = await SpaceApi.getSpaceList({ page: 1, page_size: 100 })
    if (result.code === 200 && result.data) {
      spaceList.value = result.data.spaces
    }
  } catch (error) {
    MessagePlugin.error('获取空间列表失败')
  }
}

const loadSceneList = async () => {
  if (!selectedSpaceId.value) {
    sceneList.value = []
    return
  }

  try {
    const result = await SceneApi.getSceneList({
      space_id: selectedSpaceId.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    if (result.code === 200 && result.data) {
      sceneList.value = result.data.scenes
      pagination.total = result.data.page_info.total
      paginationConfig.total = result.data.page_info.total
    }
  } catch (error) {
    MessagePlugin.error('获取场景列表失败')
  }
}

const handleSpaceChange = () => {
  pagination.currentPage = 1
  paginationConfig.current = 1
  loadSceneList()
}

const openAddDialog = () => {
  if (!selectedSpaceId.value) {
    MessagePlugin.warning('请先选择空间')
    return
  }
  isEditing.value = false
  editingSceneId.value = null
  resetForm()
  addDialogVisible.value = true
}

const handleEditScene = (scene: any) => {
  isEditing.value = true
  editingSceneId.value = scene.id
  Object.assign(formData, {
    title: scene.title,
    scene_code: scene.scene_code,
    initial_fov: scene.initial_fov || 100,
    initial_pitch: scene.initial_pitch || 0,
    initial_yaw: scene.initial_yaw || 0,
    longitude: scene.longitude || 0,
    latitude: scene.latitude || 0,
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
      const formDataObj = new FormData()
      Object.keys(formData).forEach(key => {
        const value = formData[key as keyof typeof formData]
        if (value !== undefined && value !== '') {
          formDataObj.append(key, String(value))
        }
      })

      const result = await SceneApi.updateScene(editingSceneId.value, formData as any)
      if (result.code === 200) {
        MessagePlugin.success('编辑场景成功')
        addDialogVisible.value = false
        loadSceneList()
      } else {
        MessagePlugin.error(result.msg || '编辑失败')
      }
    } else {
      if (panoramaFiles.value.length === 0) {
        MessagePlugin.warning('请选择全景图文件')
        submitLoading.value = false
        return
      }

      const file = panoramaFiles.value[0].raw
      await uploadService.uploadFile(file, {
        space_id: selectedSpaceId.value!,
        scene_code: formData.scene_code,
        title: formData.title,
        onComplete: () => {
          MessagePlugin.success('上传成功')
          addDialogVisible.value = false
          loadSceneList()
        },
        onError: (error) => {
          MessagePlugin.error(error.message || '上传失败')
        }
      })
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

const openDeleteDialog = (id: number) => {
  tableRef.value?.openDeleteDialog(id)
}

const handleDeleteScene = async (id: string | number) => {
  try {
    const result = await SceneApi.deleteScene(Number(id))
    if (result.code === 200) {
      MessagePlugin.success('删除成功')
      loadSceneList()
    } else {
      MessagePlugin.error(result.msg || '删除失败')
    }
  } catch (error) {
    MessagePlugin.error('删除失败')
  }
}

const handlePageChange = (context: { current: number; pageSize: number }) => {
  pagination.currentPage = context.current
  paginationConfig.current = context.current
  pagination.pageSize = context.pageSize
  paginationConfig.pageSize = context.pageSize
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

onMounted(() => {
  loadSpaceList()
  sceneStore.initWebSocket()
})
</script>

<style scoped>
.scene-management {
  width: 100%;
  height: 100%;
  min-height: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
</style>
