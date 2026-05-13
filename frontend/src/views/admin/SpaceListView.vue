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

    <FormDialog ref="formDialogRef">
      <template #default="{ formData, submitLoading, closeDialog, editingId }">
        <t-form
          ref="formRef"
          :data="formData"
          :rules="formRules"
          label-width="100px"
          @submit="handleSubmit(editingId)"
        >
          <t-form-item label="空间名称" name="name">
            <div class="name-field-wrapper">
              <t-input
                v-model="formData.name"
                placeholder="请输入空间名称，输入时自动搜索地点"
                @input="handleNameInput"
                @focus="handleNameFocus"
                @blur="handleNameBlur"
              >
                <template #suffix-icon>
                  <t-icon name="search" style="cursor: pointer" @click="handleSearchClick" />
                </template>
              </t-input>

              <div
                v-if="showSuggestions"
                class="suggestion-dropdown"
              >
                <div v-if="searchLoading" class="suggestion-loading">
                  <t-loading size="small" text="正在搜索地点..." />
                </div>
                <template v-else-if="searchResults.length > 0">
                  <div
                    v-for="(item, index) in searchResults"
                    :key="index"
                    class="suggestion-item"
                    @mousedown.prevent="selectLocation(item, formData)"
                  >
                    <div class="suggestion-icon">
                      <t-icon name="location" size="14" />
                    </div>
                    <div class="suggestion-content">
                      <div class="suggestion-name" v-html="highlightKeyword(item.name)"></div>
                      <div class="suggestion-address">{{ item.district }} {{ item.address }}</div>
                    </div>
                  </div>
                </template>
                <div v-else-if="formData.name" class="suggestion-empty">
                  未找到相关地点
                </div>
              </div>
            </div>
          </t-form-item>

          <t-form-item label="经度" name="longitude">
            <t-input
              v-model="formData.longitude"
              placeholder="选择地点后自动填写"
              disabled
            >
              <template #suffix>
                <span class="coord-unit">°E</span>
              </template>
            </t-input>
          </t-form-item>

          <t-form-item label="纬度" name="latitude">
            <t-input
              v-model="formData.latitude"
              placeholder="选择地点后自动填写"
              disabled
            >
              <template #suffix>
                <span class="coord-unit">°N</span>
              </template>
            </t-input>
          </t-form-item>

          <t-form-item label="省份" name="province">
            <t-input v-model="formData.province" placeholder="选择地点后自动填写" disabled />
          </t-form-item>

          <t-form-item label="城市" name="city">
            <t-input v-model="formData.city" placeholder="选择地点后自动填写" disabled />
          </t-form-item>

          <t-form-item label="封面图" name="cover">
            <div class="cover-upload-wrapper">
              <t-tag
                v-if="coverStatus === 'idle'"
                theme="warning"
                variant="light"
                size="small"
              >未上传</t-tag>
              <t-tag
                v-else-if="coverStatus === 'done'"
                theme="success"
                variant="light"
                size="small"
              >已上传</t-tag>

              <t-button
                v-if="coverStatus !== 'uploading'"
                variant="outline"
                size="medium"
                @click="triggerCoverInput"
              >
                <template #icon>
                  <t-icon name="upload" />
                </template>
                {{ coverStatus === 'done' ? '更换封面' : '上传封面图' }}
              </t-button>
              <t-loading v-else size="small" text="封面上传中..." />
              <input
                ref="coverInputRef"
                type="file"
                accept="image/jpeg,image/png"
                class="cover-input-hidden"
                @change="handleCoverChange"
              />
              <t-button
                v-if="coverStatus === 'done'"
                variant="text"
                theme="danger"
                size="small"
                @click="removeCover"
              >
                移除
              </t-button>
              <span v-if="coverFileName" class="cover-name">{{ coverFileName }}</span>
            </div>
          </t-form-item>

          <t-form-item label="描述" name="description">
            <t-textarea v-model="formData.description" placeholder="请输入空间描述" :rows="3" />
          </t-form-item>

          <t-form-item style="margin-top: 32px">
            <t-space>
              <t-button type="submit" theme="primary" size="large" :loading="submitLoading">
                提交保存
              </t-button>
              <t-button variant="outline" size="large" @click="closeDialog">取消</t-button>
            </t-space>
          </t-form-item>
        </t-form>
      </template>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import type { FormRule } from 'tdesign-vue-next'
import FormDialog from '@/components/admin/FormDialog.vue'
import SpaceApi from '@/apis/space.api'
import type { SpaceListItem } from '@/models/space.model'
import { geocode, autoComplete } from '@/utils/amap'

interface LocationResult {
  name: string
  address: string
  district: string
  province: string
  city: string
  lng: number
  lat: number
}

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
const formRef = ref()
const deleteDialogVisible = ref(false)
const deletingId = ref<number | null>(null)

const showSuggestions = ref(false)
const searchLoading = ref(false)
const searchResults = ref<LocationResult[]>([])
const searchTimer = ref<ReturnType<typeof setTimeout> | null>(null)

const coverFile = ref<File | null>(null)
const coverFileName = ref('')
const coverStatus = ref<'idle' | 'uploading' | 'done' | null>(null)
const coverInputRef = ref<HTMLInputElement>()

const formRules: Record<string, FormRule[]> = {
  name: [
    { required: true, message: '请输入空间名称', trigger: 'blur' }
  ]
}

const triggerCoverInput = () => {
  coverInputRef.value?.click()
}

const handleCoverChange = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  if (!['image/jpeg', 'image/png'].includes(file.type)) {
    MessagePlugin.warning('仅支持 JPG / PNG 格式的图片')
    return
  }

  if (file.size > 10 * 1024 * 1024) {
    MessagePlugin.warning('封面图片不能超过 10MB')
    return
  }

  coverFile.value = file
  coverFileName.value = file.name
  coverStatus.value = 'idle'
}

const removeCover = () => {
  coverFile.value = null
  coverFileName.value = ''
  coverStatus.value = null
  if (coverInputRef.value) {
    coverInputRef.value.value = ''
  }
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
  coverFile.value = null
  coverFileName.value = ''
  coverStatus.value = null
  if (coverInputRef.value) {
    coverInputRef.value.value = ''
  }
  formDialogRef.value?.openAddDialog('新增空间')
}

const handleEditSpace = (space: SpaceListItem) => {
  coverFile.value = null
  coverFileName.value = ''
  if (coverInputRef.value) {
    coverInputRef.value.value = ''
  }
  if (space.cover_url) {
    coverStatus.value = 'done'
    coverFileName.value = 'cover.jpg'
  } else {
    coverStatus.value = 'idle'
  }
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
    space.id,
    '编辑空间'
  )
}

const handleSubmit = async (editingId: string | number | null) => {
  if (!formRef.value) return
  const valid = await formRef.value.validate()
  if (!valid) return

  const data = formDialogRef.value?.formData || {}
  formDialogRef.value?.setSubmitLoading(true)
  try {
    let longitude = data.longitude ? Number(data.longitude) : null
    let latitude = data.latitude ? Number(data.latitude) : null

    if (!editingId && (longitude === null || latitude === null)) {
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
    if (coverFile.value) {
      coverStatus.value = 'uploading'
    }
    if (editingId) {
      result = await SpaceApi.updateSpace(Number(editingId), submitData as any, coverFile.value || undefined)
    } else {
      result = await SpaceApi.createSpace(submitData as any, coverFile.value || undefined)
    }

    if (result.code === 200) {
      MessagePlugin.success(editingId ? '编辑空间成功' : '新增空间成功')
      if (coverFile.value) {
        coverStatus.value = 'done'
      } else {
        coverStatus.value = 'idle'
      }
      loadSpaceList()
      // 延迟关闭弹窗，让用户看到上传状态
      setTimeout(() => {
        coverFile.value = null
        coverFileName.value = ''
        coverStatus.value = null
        formDialogRef.value?.closeDialog()
      }, 1500)
    } else {
      MessagePlugin.error(result.msg || '操作失败')
      coverStatus.value = 'idle'
    }
  } catch (error: any) {
    console.error('提交失败:', error)
    const errorMessage = error?.message || error?.msg || '操作失败'
    MessagePlugin.error(`操作失败: ${errorMessage}`)
  } finally {
    formDialogRef.value?.setSubmitLoading(false)
  }
}

// 地点搜索逻辑
const handleSearchClick = () => {
  const name = formDialogRef.value?.formData?.name
  if (name) {
    showSuggestions.value = true
    doSearch(name)
  }
}

const handleNameInput = () => {
  if (searchTimer.value) clearTimeout(searchTimer.value)
  const data = formDialogRef.value?.formData || {}
  const keyword = data.name?.trim()
  if (!keyword) {
    searchResults.value = []
    showSuggestions.value = false
    searchLoading.value = false
    return
  }

  showSuggestions.value = true
  searchLoading.value = true

  searchTimer.value = setTimeout(() => { doSearch(keyword) }, 300)
}

const handleNameFocus = () => {
  const data = formDialogRef.value?.formData || {}
  if (data.name) {
    showSuggestions.value = true
    if (searchResults.value.length === 0 && !searchLoading.value) {
      doSearch(data.name)
    }
  }
}

const handleNameBlur = () => {
  setTimeout(() => { showSuggestions.value = false }, 250)
}

const doSearch = async (keyword: string) => {
  searchLoading.value = true
  try {
    const tips = await autoComplete(keyword)
    if (tips && tips.length > 0) {
      searchResults.value = tips
        .filter((tip: any) => tip.location && (tip.location.lng || tip.location.getLng))
        .map((tip: any) => {
          const loc = tip.location
          return {
            name: tip.name,
            address: tip.address || '',
            district: tip.district || '',
            province: tip.district ? extractProvince(tip.district) : '',
            city: tip.district ? extractCity(tip.district) : '',
            lng: typeof loc.getLng === 'function' ? loc.getLng() : parseFloat(loc.lng),
            lat: typeof loc.getLat === 'function' ? loc.getLat() : parseFloat(loc.lat)
          }
        })
        .slice(0, 8)
    } else {
      searchResults.value = []
    }
  } catch (error) {
    searchResults.value = []
  } finally {
    searchLoading.value = false
    showSuggestions.value = searchResults.value.length > 0 || (formDialogRef.value?.formData?.name || '') !== ''
  }
}

const selectLocation = (item: LocationResult, formData: Record<string, any>) => {
  if (item.lng) formData.longitude = String(item.lng.toFixed(6))
  if (item.lat) formData.latitude = String(item.lat.toFixed(6))
  if (item.name) formData.name = item.name
  if (item.province) formData.province = item.province
  if (item.city) formData.city = item.city
  searchResults.value = []
  showSuggestions.value = false
}

const extractProvince = (district: string): string => {
  if (!district) return ''
  const parts = district.split(/省|自治区|直辖市|特别行政区/)
  if (parts.length > 1) {
    const prefix = parts[0]
    if (district.includes('省')) return prefix + '省'
    if (district.includes('自治区')) return prefix + '自治区'
    if (district.includes('直辖市')) return prefix + '直辖市'
    if (district.includes('特别行政区')) return prefix + '特别行政区'
  }
  const match = district.match(/^(.+?)(省|市|自治区)/)
  return match ? match[0] : district.split(/市/)[0] + '省'
}

const extractCity = (district: string): string => {
  if (!district) return ''
  const parts = district.split(/省|自治区/)
  if (parts.length > 1) {
    const afterProvince = parts[1]
    const cityMatch = afterProvince.match(/^(.+?市)/)
    return cityMatch ? cityMatch[1] : afterProvince.split(/地区|州|盟/)[0]
  }
  const match = district.match(/市(.+?区|.+?县|.+?市)/)
  if (match) {
    const cityPart = district.substring(0, district.indexOf(match[1]) + match[1].length)
    const cityMatch2 = cityPart.match(/(.+?市)/)
    return cityMatch2 ? cityMatch2[1] : ''
  }
  return ''
}

const highlightKeyword = (text: string) => {
  const data = formDialogRef.value?.formData || {}
  const keyword = data.name?.trim()
  if (!keyword) return text
  const regex = new RegExp(`(${keyword})`, 'gi')
  return text.replace(regex, '<span class="highlight">$1</span>')
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
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0 16px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #1d2129;
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
  min-height: 0;
  overflow-y: auto;
}

.space-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  padding-bottom: 20px;
}

.space-card {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #f0f1f2;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
}

.space-card:hover {
  border-color: #0052d9;
  box-shadow: 0 8px 24px rgba(0, 82, 217, 0.08);
  transform: translateY(-2px);
}

.card-cover {
  position: relative;
  width: 100%;
  height: 160px;
  background: linear-gradient(135deg, #f2f3f5, #e8eaf0);
  overflow: hidden;
}

.card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-cover {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #c9cdd4;
}

.card-status {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
  background: rgba(255, 255, 255, 0.9);
  color: #86909c;
  border: 1px solid #e5e6eb;
}

.card-status.active {
  background: rgba(0, 210, 110, 0.1);
  color: #00a870;
  border-color: rgba(0, 210, 110, 0.3);
}

.card-body {
  padding: 14px 16px 8px;
  flex: 1;
}

.card-main {
  margin-bottom: 10px;
}

.card-title {
  margin: 0 0 4px;
  font-size: 16px;
  font-weight: 700;
  color: #1d2129;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-description {
  margin: 0;
  font-size: 12px;
  color: #86909c;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-info {
  display: flex;
  gap: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #86909c;
}

.card-actions {
  padding: 8px 16px 12px;
  display: flex;
  gap: 8px;
  border-top: 1px solid #f2f3f5;
  margin-top: 4px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: #86909c;
}

.empty-state .t-icon {
  margin-bottom: 16px;
  color: #c9cdd4;
}

.empty-state p {
  font-size: 14px;
  margin-bottom: 20px;
}

.pagination-wrapper {
  padding: 16px 0;
  display: flex;
  justify-content: flex-end;
}

/* 表单搜索下拉 */
.name-field-wrapper {
  position: relative;
  width: 100%;
}

.suggestion-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  max-height: 240px;
  overflow-y: auto;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.suggestion-loading {
  padding: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #86909c;
}

.suggestion-empty {
  padding: 16px;
  text-align: center;
  color: #86909c;
  font-size: 13px;
}

.suggestion-item {
  display: flex;
  padding: 10px 12px;
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 1px solid #f2f3f5;
}

.suggestion-item:hover { background: #f2f3f5; }
.suggestion-item:last-child { border-bottom: none; }

.suggestion-icon {
  margin-right: 10px;
  color: #0052d9;
  margin-top: 2px;
  flex-shrink: 0;
}

.suggestion-content {
  min-width: 0;
  flex: 1;
}

.suggestion-name {
  font-weight: 500;
  font-size: 14px;
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.suggestion-address {
  font-size: 12px;
  color: #86909c;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.coord-unit {
  color: #86909c;
  font-size: 12px;
  font-weight: 500;
}

/* 封面上传 */
.cover-upload-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.cover-name {
  font-size: 13px;
  color: #4e5969;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cover-input-hidden {
  display: none;
}
</style>