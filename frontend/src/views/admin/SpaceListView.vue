<template>
  <div class="space-list-view">
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
      <div class="space-grid" v-loading="loading">
        <div
          v-for="space in spaceList"
          :key="space.id"
          class="space-card"
          @click="handleSpaceClick(space)"
        >
          <div class="card-cover">
            <img
              v-if="space.slug"
              :src="getCoverUrl(space.slug)"
              :alt="space.name"
              @error="handleCoverError($event)"
            />
            <div v-else class="no-cover">
              <t-icon name="image" size="32" />
            </div>
            <div class="card-status" :class="{ active: space.status === 1 }">
              {{ space.status === 1 ? '启用' : '禁用' }}
            </div>
          </div>

          <div class="card-body">
            <div class="card-main">
              <h3 class="card-title">{{ space.name }}</h3>
              <p class="card-description">{{ space.description || '暂无描述' }}</p>
            </div>

            <div class="card-info">
              <div class="info-item">
                <t-icon name="location" size="14" />
                <span>{{ space.province || '-' }} {{ space.city || '' }}</span>
              </div>
              <div class="info-item">
                <t-icon name="view-module" size="14" />
                <span>{{ space.scene_count || 0 }} 个场景</span>
              </div>
            </div>
          </div>

          <div class="card-actions" @click.stop>
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

      <div v-if="!loading && spaceList.length === 0" class="empty-state">
        <t-icon name="folder-open" size="64" />
        <p>暂无空间数据</p>
        <t-button theme="primary" @click="handleAddSpace">新建空间</t-button>
      </div>
    </div>

    <div class="pagination-wrapper" v-if="pagination.total > 0">
      <t-pagination
        v-model="pagination.current"
        v-model:pageSize="pagination.pageSize"
        :total="pagination.total"
        :show-jumper="true"
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
import { MessagePlugin } from 'tdesign-vue-next'
import FormDialog, { type FormField } from '@/components/admin/FormDialog.vue'
import SpaceApi from '@/apis/space.api'
import type { SpaceListItem } from '@/models/space.model'
import { geocode } from '@/utils/amap'

const emit = defineEmits<{
  (e: 'spaceClick', space: SpaceListItem): void
}>()

const loading = ref(false)
const searchKeyword = ref('')
const spaceList = ref<SpaceListItem[]>([])

const pagination = reactive({
  current: 1,
  pageSize: 12,
  total: 0
})

const formDialogRef = ref()
const deleteDialogVisible = ref(false)
const deletingId = ref<number | null>(null)

const formFields: FormField[] = [
  {
    name: 'name',
    label: '空间名称',
    type: 'input',
    required: true,
    placeholder: '请输入空间名称，输入时自动搜索地点'
  },
  {
    name: 'longitude',
    label: '经度',
    type: 'input',
    placeholder: '选择地点后自动填写'
  },
  {
    name: 'latitude',
    label: '纬度',
    type: 'input',
    placeholder: '选择地点后自动填写'
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
    placeholder: '选择地点后自动填写'
  },
  {
    name: 'city',
    label: '城市',
    type: 'input',
    placeholder: '选择地点后自动填写'
  }
]

const defaultFormData = {
  status: 1
}

const getCoverUrl = (slug: string) => {
  return `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'}/api/v1/res/covers/${slug}`
}

const handleCoverError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.style.display = 'none'
  const parent = img.parentElement
  if (parent) {
    const placeholder = parent.querySelector('.no-cover') as HTMLElement
    if (placeholder) {
      placeholder.style.display = 'flex'
    }
  }
}

const loadSpaceList = async () => {
  try {
    loading.value = true
    const result = await SpaceApi.getSpaceList({
      keyword: searchKeyword.value,
      page: pagination.current,
      page_size: pagination.pageSize
    })
    if (result.code === 200 && result.data) {
      spaceList.value = result.data.spaces
      pagination.total = result.data.page_info.total
    } else {
      MessagePlugin.error(result.msg || '获取空间列表失败')
    }
  } catch (error) {
    MessagePlugin.error('获取空间列表失败')
  } finally {
    loading.value = false
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
  try {
    let longitude = data.longitude ? Number(data.longitude) : null
    let latitude = data.latitude ? Number(data.latitude) : null

    if (!data.id && (longitude === null || latitude === null)) {
      const address = [data.province, data.city].filter(Boolean).join(' ')
      if (address) {
        const coords = await geocode(address)
        if (coords) {
          longitude = coords.lng
          latitude = coords.lat
        }
      }
    }

    const submitData = {
      ...data,
      longitude,
      latitude,
      zoom_level: data.zoom_level ? Number(data.zoom_level) : null,
      sort_order: data.sort_order ? Number(data.sort_order) : null
    }

    let result
    if (data.id) {
      result = await SpaceApi.updateSpace(Number(data.id), submitData as any)
    } else {
      result = await SpaceApi.createSpace(submitData as any)
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
  pagination.current = context.current
  pagination.pageSize = context.pageSize
  loadSpaceList()
}

const handleSearch = () => {
  pagination.current = 1
  loadSpaceList()
}

const handleSpaceClick = (space: SpaceListItem) => {
  emit('spaceClick', space)
}

onMounted(() => {
  loadSpaceList()
})

defineExpose({
  loadSpaceList
})
</script>

<style scoped>
.space-list-view {
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

.space-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.space-card {
  background: var(--td-bg-color-container);
  border-radius: 12px;
  border: 1px solid var(--td-component-border);
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.space-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transform: translateY(-4px);
  border-color: var(--td-brand-color);
}

.card-cover {
  position: relative;
  height: 160px;
  background: linear-gradient(135deg, #1E3A5F 0%, #0F1826 100%);
  overflow: hidden;
}

.card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card-cover .no-cover {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.3);
}

.card-status {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 2px 10px;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
  background: rgba(245, 63, 63, 0.9);
  color: #fff;
}

.card-status.active {
  background: rgba(0, 168, 112, 0.9);
}

.card-body {
  padding: 16px;
}

.card-main {
  margin-bottom: 12px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--td-text-color-primary);
  margin: 0 0 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-description {
  font-size: 13px;
  color: var(--td-text-color-secondary);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-info {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-top: 12px;
  border-top: 1px solid var(--td-component-border);
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--td-text-color-secondary);
}

.info-item :deep(.t-icon) {
  color: var(--td-brand-color);
}

.card-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 16px;
  background: var(--td-bg-color-container-hover);
  border-top: 1px solid var(--td-component-border);
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

.pagination-wrapper {
  margin-top: auto;
  padding-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
