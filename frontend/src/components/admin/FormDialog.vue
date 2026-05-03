<template>
  <t-dialog
    v-model:visible="dialogVisible"
    :header="dialogTitle"
    width="500px"
    :footer="false"
  >
    <t-form
      ref="formRef"
      :data="formData"
      :rules="formRules"
      label-width="100px"
      @submit="handleSubmit"
    >
      <t-form-item
        v-for="field in formFields"
        :key="field.name"
        :label="field.label"
        :name="field.name"
      >
        <div class="field-with-action">
          <template v-if="field.type === 'input'">
            <div v-if="field.name === 'name'" class="name-field-wrapper">
              <t-input
                v-model="formData[field.name]"
                :placeholder="field.placeholder || `请输入${field.label}`"
                @input="handleNameInput"
                @focus="handleNameFocus"
                @blur="handleNameBlur"
              />
              <div
                v-if="showSuggestions && searchResults.length > 0"
                class="suggestion-dropdown"
              >
                <div
                  v-for="(item, index) in searchResults"
                  :key="index"
                  class="suggestion-item"
                  @mousedown.prevent="selectLocation(item)"
                >
                  <div class="suggestion-icon">
                    <t-icon name="location" size="14" />
                  </div>
                  <div class="suggestion-content">
                    <div class="suggestion-name" v-html="highlightKeyword(item.name)"></div>
                    <div class="suggestion-address">{{ item.district }} {{ item.address }}</div>
                  </div>
                </div>
              </div>
            </div>
            <t-input
              v-else
              v-model="formData[field.name]"
              :placeholder="field.placeholder || `请输入${field.label}`"
            />
          </template>
          <t-select
            v-else-if="field.type === 'select'"
            v-model="formData[field.name]"
            :options="field.options"
            :placeholder="field.placeholder || `请选择${field.label}`"
          />
          <t-switch
            v-else-if="field.type === 'switch'"
            v-model="formData[field.name]"
            :label="['启用', '禁用']"
          />
          <t-textarea
            v-else-if="field.type === 'textarea'"
            v-model="formData[field.name]"
            :placeholder="field.placeholder || `请输入${field.label}`"
            :rows="3"
          />
        </div>
      </t-form-item>

      <t-form-item style="margin-top: 20px">
        <t-space>
          <t-button type="submit" theme="primary" :loading="submitLoading">
            提交
          </t-button>
          <t-button @click="dialogVisible = false">取消</t-button>
        </t-space>
      </t-form-item>
    </t-form>
  </t-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { FormRule } from 'tdesign-vue-next'
import { autoComplete } from '@/utils/amap'

export interface FormField {
  name: string
  label: string
  type: 'input' | 'select' | 'switch' | 'textarea'
  placeholder?: string
  options?: { label: string; value: any }[]
  required?: boolean
  rules?: FormRule[]
}

interface LocationResult {
  name: string
  address: string
  district: string
  province: string
  city: string
  lng: number
  lat: number
}

const props = defineProps<{
  formFields: FormField[]
  defaultData?: Record<string, any>
}>()

const emit = defineEmits<{
  (e: 'submit', data: Record<string, any>): void
}>()

const dialogVisible = ref(false)
const dialogTitle = ref('新增')
const formRef = ref()
const submitLoading = ref(false)

const isEdit = ref(false)
const editingId = ref<string | number | null>(null)

const formData = reactive<Record<string, any>>({})

const showSuggestions = ref(false)
const searchResults = ref<LocationResult[]>([])
const searchTimer = ref<ReturnType<typeof setTimeout> | null>(null)

watch(
  () => props.formFields,
  () => {
    props.formFields.forEach(field => {
      if (formData[field.name] === undefined) {
        formData[field.name] = field.type === 'switch' ? 1 : ''
      }
    })
  },
  { immediate: true }
)

const formRules = reactive<Record<string, FormRule[]>>({})
watch(
  () => props.formFields,
  (fields) => {
    Object.keys(formRules).forEach(key => delete formRules[key])
    fields.forEach(field => {
      if (field.required) {
        formRules[field.name] = field.rules || [
          { required: true, message: `请输入${field.label}`, trigger: 'blur' }
        ]
      } else if (field.rules) {
        formRules[field.name] = field.rules
      }
    })
  },
  { immediate: true, deep: true }
)

const openAddDialog = () => {
  isEdit.value = false
  editingId.value = null
  dialogTitle.value = '新增'
  Object.keys(formData).forEach(key => {
    const field = props.formFields.find(f => f.name === key)
    formData[key] = field?.type === 'switch' ? 1 : ''
  })
  if (props.defaultData) {
    Object.assign(formData, props.defaultData)
  }
  dialogVisible.value = true
}

const openEditDialog = (data: Record<string, any>, id: string | number) => {
  isEdit.value = true
  editingId.value = id
  dialogTitle.value = '编辑'
  Object.assign(formData, data)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  const valid = await (formRef.value as any).validate()
  if (valid) {
    submitLoading.value = true
    try {
      emit('submit', { ...formData, id: editingId.value })
    } finally {
      submitLoading.value = false
    }
  }
}

const closeDialog = () => {
  dialogVisible.value = false
}

const getFormData = () => {
  return { ...formData }
}

const setFieldValue = (field: string, value: any) => {
  formData[field] = value
}

const handleNameInput = () => {
  if (searchTimer.value) {
    clearTimeout(searchTimer.value)
  }

  const keyword = formData.name?.trim()
  if (!keyword) {
    searchResults.value = []
    showSuggestions.value = false
    return
  }

  searchTimer.value = setTimeout(() => {
    doSearch(keyword)
  }, 300)
}

const handleNameFocus = () => {
  if (searchResults.value.length > 0) {
    showSuggestions.value = true
  }
}

const handleNameBlur = () => {
  setTimeout(() => {
    showSuggestions.value = false
  }, 200)
}

const doSearch = async (keyword: string) => {
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
    showSuggestions.value = searchResults.value.length > 0
  } catch (error) {
    console.error('搜索失败:', error)
    searchResults.value = []
    showSuggestions.value = false
  }
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
  const keyword = formData.name?.trim()
  if (!keyword) return text
  const regex = new RegExp(`(${keyword})`, 'gi')
  return text.replace(regex, '<span class="highlight">$1</span>')
}

const selectLocation = (item: LocationResult) => {
  formData.name = item.name
  formData.longitude = String(item.lng)
  formData.latitude = String(item.lat)
  formData.province = item.province
  formData.city = item.city
  searchResults.value = []
  showSuggestions.value = false
}

defineExpose({
  openAddDialog,
  openEditDialog,
  closeDialog,
  getFormData,
  setFieldValue
})
</script>

<style scoped>
.field-with-action {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.field-with-action :deep(.t-input),
.field-with-action :deep(.t-select),
.field-with-action :deep(.t-textarea) {
  flex: 1;
}

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
  max-height: 280px;
  overflow-y: auto;
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-border);
  border-radius: 6px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  z-index: 1000;
}

.suggestion-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
  border-bottom: 1px solid var(--td-component-border);
  transition: background 0.15s;
}

.suggestion-item:last-child {
  border-bottom: none;
}

.suggestion-item:hover {
  background: var(--td-bg-color-container-hover);
}

.suggestion-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--td-brand-color-light);
  color: var(--td-brand-color);
  flex-shrink: 0;
  margin-top: 2px;
}

.suggestion-content {
  flex: 1;
  min-width: 0;
}

.suggestion-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--td-text-color-primary);
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.suggestion-name :deep(.highlight) {
  color: var(--td-brand-color);
  font-weight: 600;
}

.suggestion-address {
  font-size: 12px;
  color: var(--td-text-color-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
