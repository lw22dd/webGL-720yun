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
        <t-input
          v-if="field.type === 'input'"
          v-model="formData[field.name]"
          :placeholder="field.placeholder || `请输入${field.label}`"
        />
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

export interface FormField {
  name: string
  label: string
  type: 'input' | 'select' | 'switch' | 'textarea'
  placeholder?: string
  options?: { label: string; value: any }[]
  required?: boolean
  rules?: FormRule[]
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

defineExpose({
  openAddDialog,
  openEditDialog,
  closeDialog
})
</script>
