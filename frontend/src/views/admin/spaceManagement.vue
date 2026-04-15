<template>
  <div class="space-management">
    <div class="header-actions">
      <div class="page-title">空间管理</div>
      <div class="header-right">
        <t-input
          v-model="searchKeyword"
          placeholder="搜索空间名称"
          clearable
          @enter="handleSearch"
          @clear="handleSearch"
        >
          <template #suffix-icon>
            <t-icon name="search" />
          </template>
        </t-input>
        <t-button theme="primary" @click="handleAddSpace">
          <template #icon><t-icon name="add" /></template>
          新建空间
        </t-button>
      </div>

    </div>

    <div class="space-list-container">
      <div class="space-list">
        <div
          v-for="space in spaceList"
          :key="space.id"
          class="space-list-item"
        >
          <div class="item-cover">
            <img
              v-if="space.cover_url"
              :src="space.cover_url"
              :alt="space.name"
            />
            <div v-else class="no-cover">
              <t-icon name="image" size="32" />
            </div>
            <div class="space-status" :class="{ active: space.status === 1 }">
              {{ space.status === 1 ? '启用' : '禁用' }}
            </div>
          </div>

          <div class="item-content">
            <div class="item-main">
              <div class="item-title">{{ space.name }}</div>
              <div class="item-description">{{ space.description || '暂无描述' }}</div>
            </div>

            <div class="item-info">
              <div class="info-item">
                <t-icon name="location" size="14" />
                <span>{{ space.province || '-' }} {{ space.city || '' }}</span>
              </div>
              <div class="info-item">
                <t-icon name="view-module" size="14" />
                <span>场景数: {{ space.scene_count || 0 }}</span>
              </div>
              <div class="info-item">
                <t-icon name="time" size="14" />
                <span>{{ formatDate(space.created_at) }}</span>
              </div>
            </div>
          </div>

          <div class="item-actions">
            <t-button
              theme="default"
              variant="text"
              size="small"
              @click="goToGraphEditor(space.id)"
            >
              <template #icon><t-icon name="map-location" /></template>
              图编辑
            </t-button>
            <t-button
              theme="primary"
              variant="text"
              size="small"
              @click="handleEditSpace(space)"
            >
              <template #icon><t-icon name="edit" /></template>
              编辑
            </t-button>
            <t-button
              theme="danger"
              variant="text"
              size="small"
              @click="openDeleteDialog(space.id)"
            >
              <template #icon><t-icon name="delete" /></template>
              删除
            </t-button>
          </div>
        </div>
      </div>
    </div>

    <div class="pagination-wrapper">
      <t-pagination
        v-model="paginationConfig.current"
        v-model:pageSize="paginationConfig.pageSize"
        :total="paginationConfig.total"
        :show-jumper="paginationConfig.showJumper"
        @change="handlePageChange"
      />
    </div>

    <t-dialog
      v-model:visible="deleteDialogVisible"
      header="确认删除"
      width="400px"
      @confirm="confirmDelete"
      @cancel="deleteDialogVisible = false"
    >
      <p>确定要删除该空间吗？此操作不可撤销。</p>
    </t-dialog>

    <FormDialog
      ref="formDialogRef"
      :form-fields="formFields"
      :default-data="defaultFormData"
      @submit="handleSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import FormDialog, { type FormField } from '@/components/admin/FormDialog.vue'
import SpaceApi from '@/services/api/space.api'
import type { SpaceListItem } from '@/models/space.model'

const router = useRouter()
const searchKeyword = ref('')
const spaceList = ref<SpaceListItem[]>([])

const paginationConfig = reactive({
  current: 1,
  pageSize: 12,
  total: 0,
  defaultCurrent: 1,
  defaultPageSize: 12,
  showJumper: true
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 12,
  total: 0
})

const formDialogRef = ref()
const submitLoading = ref(false)
const deleteDialogVisible = ref(false)
const deletingId = ref<number | null>(null)

const formFields: FormField[] = [
  {
    name: 'name',
    label: '空间名称',
    type: 'input',
    required: true,
    placeholder: '请输入空间名称'
  },
  {
    name: 'description',
    label: '描述',
    type: 'textarea',
    placeholder: '请输入空间描述'
  },
  {
    name: 'province',
    label: '省份',
    type: 'input',
    placeholder: '请输入省份'
  },
  {
    name: 'city',
    label: '城市',
    type: 'input',
    placeholder: '请输入城市'
  },
  {
    name: 'longitude',
    label: '经度',
    type: 'input',
    placeholder: '请输入经度'
  },
  {
    name: 'latitude',
    label: '纬度',
    type: 'input',
    placeholder: '请输入纬度'
  },
  {
    name: 'zoom_level',
    label: '缩放级别',
    type: 'input',
    placeholder: '请输入缩放级别'
  },
  {
    name: 'sort_order',
    label: '排序',
    type: 'input',
    placeholder: '请输入排序值'
  },
  {
    name: 'status',
    label: '状态',
    type: 'switch'
  }
]

const defaultFormData = {
  status: 1
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

const loadSpaceList = async () => {
  try {
    const result = await SpaceApi.getSpaceList({
      keyword: searchKeyword.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    if (result.code === 200 && result.data) {
      spaceList.value = result.data.spaces
      pagination.total = result.data.page_info.total
      paginationConfig.total = result.data.page_info.total
    } else {
      MessagePlugin.error(result.msg || '获取空间列表失败')
    }
  } catch (error) {
    MessagePlugin.error('获取空间列表失败')
  }
}

const handleAddSpace = () => {
  formDialogRef.value?.openAddDialog()
}

const handleEditSpace = (space: SpaceListItem) => {
  formDialogRef.value?.openEditDialog(
    {
      name: space.name,
      description: space.description,
      province: space.province,
      city: space.city,
      longitude: space.longitude,
      latitude: space.latitude,
      zoom_level: space.zoom_level,
      sort_order: space.sort_order,
      status: space.status
    },
    space.id
  )
}

const handleSubmit = async (data: Record<string, any>) => {
  submitLoading.value = true
  try {
    let result
    const formDataObj = new FormData()
    Object.keys(data).forEach(key => {
      if (key !== 'id' && data[key] !== undefined && data[key] !== '') {
        formDataObj.append(key, data[key])
      }
    })

    if (data.id) {
      result = await SpaceApi.updateSpace(Number(data.id), data as any)
    } else {
      result = await SpaceApi.createSpace(data as any)
    }

    if (result.code === 200) {
      MessagePlugin.success(data.id ? '编辑空间成功' : '新增空间成功')
      formDialogRef.value?.closeDialog()
      loadSpaceList()
    } else {
      MessagePlugin.error(result.msg || '操作失败')
    }
  } catch (error: any) {
    console.error('提交失败:', error)
    const errorMessage = error?.message || error?.msg || '操作失败'
    MessagePlugin.error(`操作失败: ${errorMessage}`)
  } finally {
    submitLoading.value = false
  }
}

const openDeleteDialog = (id: number) => {
  deletingId.value = id
  deleteDialogVisible.value = true
}

const confirmDelete = async () => {
  if (deletingId.value === null) return
  
  try {
    const result = await SpaceApi.deleteSpace(deletingId.value)
    if (result.code === 200) {
      MessagePlugin.success('删除成功')
      loadSpaceList()
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

const handlePageChange = (context: { current: number; pageSize: number }) => {
  pagination.currentPage = context.current
  paginationConfig.current = context.current
  pagination.pageSize = context.pageSize
  paginationConfig.pageSize = context.pageSize
  loadSpaceList()
}

const handleSearch = () => {
  pagination.currentPage = 1
  paginationConfig.current = 1
  loadSpaceList()
}

const goToGraphEditor = (spaceId: number) => {
  router.push(`/admin/spaces/${spaceId}/graph`)
}

onMounted(() => {
  loadSpaceList()
})

defineExpose({
  handleAddSpace
})
</script>

<style scoped>
.space-management {
  width: 100%;
  height: 100%;
  min-height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--td-text-color-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-right :deep(.t-input) {
  width: 240px;
}

.space-list-container {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 20px;
}

.space-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.space-list-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background-color: var(--td-bg-color-container);
  border-radius: 8px;
  border: 1px solid var(--td-component-border);
  transition: all 0.3s ease;
}

.space-list-item:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.item-cover {
  position: relative;
  width: 120px;
  height: 80px;
  flex-shrink: 0;
  border-radius: 6px;
  overflow: hidden;
  background-color: #f5f5f5;
}

.item-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-cover .no-cover {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #86909c;
}

.space-status {
  position: absolute;
  top: 4px;
  right: 4px;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
  background-color: rgba(245, 63, 63, 0.9);
  color: #fff;
}

.space-status.active {
  background-color: rgba(0, 168, 112, 0.9);
}

.item-content {
  flex: 1;
  margin-left: 16px;
  display: flex;
  align-items: center;
  gap: 24px;
  min-width: 0;
}

.item-main {
  flex: 1;
  min-width: 200px;
  max-width: 400px;
}

.item-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--td-text-color-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 6px;
}

.item-description {
  font-size: 13px;
  color: var(--td-text-color-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-info {
  display: flex;
  align-items: center;
  gap: 24px;
  flex-shrink: 0;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--td-text-color-secondary);
}

.info-item :deep(.t-icon) {
  color: var(--td-brand-color);
}

.item-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: 16px;
  flex-shrink: 0;
}

.pagination-wrapper {
  margin-top: auto;
  padding-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
