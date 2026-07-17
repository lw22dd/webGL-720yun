<template>
  <t-dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    :header="mode === 'create' ? '新增热点' : '编辑热点'"
    width="560px"
    :footer="false"
    attach=".panorama-container"
    :close-on-overlay-click="false"
  >
    <t-form ref="formRef" :data="localForm" :rules="formRules" label-width="80px" label-align="top">
      <!-- 类型选择 -->
      <t-form-item label="类型" name="type">
        <t-radio-group v-model="localForm.type" @change="onTypeChange">
          <t-radio :value="1">场景跳转</t-radio>
          <t-radio :value="2">信息弹窗</t-radio>
        </t-radio-group>
      </t-form-item>

      <!-- 标题 -->
      <t-form-item label="标题" name="title">
        <t-input v-model="localForm.title" placeholder="请输入热点标题" maxlength="100" />
      </t-form-item>

      <!-- 坐标 -->
      <t-form-item label="坐标">
        <div class="coord-display">
          <span class="coord-item">Pitch: {{ localForm.pitch.toFixed(4) }}</span>
          <span class="coord-item">Yaw: {{ localForm.yaw.toFixed(4) }}</span>
        </div>
      </t-form-item>

      <!-- 图标选择 -->
      <t-form-item label="图标">
        <div class="icon-selector">
          <div class="icon-presets">
            <div
              v-for="preset in availablePresets"
              :key="preset.key"
              class="icon-preset-item"
              :class="{ active: localForm.icon_source === 'preset' && localForm.icon_preset_key === preset.key }"
              @click="selectPreset(preset.key)"
              v-html="preset.svg"
            />
          </div>
          <div class="icon-upload-area">
            <t-upload
              v-model="iconUploadFiles"
              :request-method="handleIconUpload"
              :max="1"
              :auto-upload="true"
              accept="image/svg+xml,image/png,image/jpeg,image/gif"
              :show-upload-progress="false"
              theme="image"
              :placeholder="localForm.icon_source === 'custom' ? '已上传自定义图标' : '上传自定义图标'"
            />
          </div>
        </div>
      </t-form-item>

      <!-- 跳转类：目标场景 + 转场效果 -->
      <template v-if="localForm.type === 1">
        <t-form-item label="目标场景" name="target_scene_id">
          <t-select v-model="localForm.target_scene_id" placeholder="请选择目标场景" filterable>
            <t-option
              v-for="scene in sceneList"
              :key="scene.id"
              :value="scene.id"
              :label="scene.title"
            />
          </t-select>
        </t-form-item>
       
      </template>

      <!-- 信息类：文字内容 + 配图 -->
      <template v-if="localForm.type === 2">
        <t-form-item label="文字内容" name="content">
          <t-textarea
            v-model="localForm.content"
            placeholder="请输入信息内容"
            :autosize="{ minRows: 3, maxRows: 6 }"
            maxlength="500"
          />
        </t-form-item>
        <t-form-item label="配图">
          <t-upload
            v-model="mediaUploadFiles"
            :request-method="handleMediaUpload"
            :max="1"
            :auto-upload="true"
            accept="image/png,image/jpeg,image/gif"
            :show-upload-progress="false"
            theme="image"
            :placeholder="localForm.media_url ? '已上传配图' : '上传配图（可选）'"
          />
        </t-form-item>
      </template>

      <!-- 操作按钮 -->
      <t-form-item>
        <div class="dialog-footer">
          <t-button v-if="mode === 'edit'" theme="danger" variant="outline" @click="handleDelete" :loading="deleting">
            删除热点
          </t-button>
          <div class="footer-right">
            <t-button theme="default" variant="outline" @click="$emit('update:visible', false)">取消</t-button>
            <t-button theme="primary" @click="handleSubmit" :loading="submitting">
              {{ mode === 'create' ? '创建' : '保存' }}
            </t-button>
          </div>
        </div>
      </t-form-item>
    </t-form>
  </t-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import HotspotApi from '@/apis/hotspot.api'
import uploadService from '@/services/uploadService'
import type { CreateHotspotRequest, UpdateHotspotRequest, HotspotEditorForm } from '@/models/hotspot.model'
import type { SceneListItem } from '@/models/scene.model'
import { getPresetsForType } from '@/utils/hotspotIcons'

const props = defineProps<{
  visible: boolean
  form: HotspotEditorForm
  mode: 'create' | 'edit'
  sceneList: SceneListItem[]
  spaceId?: number
  sceneCode?: string
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'submitted'): void
  (e: 'deleted'): void
}>()

const formRef = ref()
const submitting = ref(false)
const deleting = ref(false)
const iconUploadFiles = ref<any[]>([])
const mediaUploadFiles = ref<any[]>([])

const localForm = reactive<HotspotEditorForm>({
  scene_id: 0,
  type: 1,
  pitch: 0,
  yaw: 0,
  title: '',
  icon_source: 'preset',
  icon_preset_key: 'arrow-blue',
  transition_effect: 'fade',
})

// 同步 prop 到 localForm
watch(
  () => props.form,
  (newForm) => {
    if (newForm) {
      Object.assign(localForm, newForm)
      // 重置上传文件列表
      iconUploadFiles.value = []
      mediaUploadFiles.value = []
    }
  },
  { immediate: true, deep: true }
)

const availablePresets = computed(() => {
  return getPresetsForType(localForm.type)
})

const formRules = {
  title: [
    { required: true, message: '请输入标题', type: 'error' },
    { min: 1, max: 100, message: '标题长度 1-100 字符', type: 'error' },
  ],
  target_scene_id: [
    { required: true, message: '请选择目标场景', type: 'error' },
  ],
  content: [
    { required: true, message: '请输入信息内容', type: 'error' },
  ],
}

function onTypeChange() {
  // 切换类型时重置图标为该类型的默认预设
  const presets = getPresetsForType(localForm.type)
  if (presets.length > 0 && localForm.icon_source === 'preset') {
    localForm.icon_preset_key = presets[0].key
  }
}

function selectPreset(key: string) {
  localForm.icon_source = 'preset'
  localForm.icon_preset_key = key
  iconUploadFiles.value = []
}

async function handleIconUpload(file: any): Promise<any> {
  try {
    if (!props.spaceId || !props.sceneCode) {
      MessagePlugin.error('缺少场景信息，无法上传')
      return { status: 'fail', error: '缺少场景信息' }
    }

    const result = await uploadService.uploadFile(file.raw, {
      space_id: props.spaceId,
      scene_code: props.sceneCode,
      title: `hotspot-icon-${Date.now()}`,
    })

    if (result.source_url) {
      localForm.icon_url = result.source_url
      localForm.icon_source = 'custom'
      MessagePlugin.success('图标上传成功')
      return { status: 'success', response: { url: result.source_url } }
    }
    return { status: 'fail', error: '上传失败' }
  } catch (e: any) {
    MessagePlugin.error(e.message || '图标上传失败')
    return { status: 'fail', error: e.message }
  }
}

async function handleMediaUpload(file: any): Promise<any> {
  try {
    if (!props.spaceId || !props.sceneCode) {
      MessagePlugin.error('缺少场景信息，无法上传')
      return { status: 'fail', error: '缺少场景信息' }
    }

    const result = await uploadService.uploadFile(file.raw, {
      space_id: props.spaceId,
      scene_code: props.sceneCode,
      title: `hotspot-media-${Date.now()}`,
    })

    if (result.source_url) {
      localForm.media_url = result.source_url
      localForm.media_type = 'image'
      MessagePlugin.success('配图上传成功')
      return { status: 'success', response: { url: result.source_url } }
    }
    return { status: 'fail', error: '上传失败' }
  } catch (e: any) {
    MessagePlugin.error(e.message || '配图上传失败')
    return { status: 'fail', error: e.message }
  }
}

function formToCreateReq(form: HotspotEditorForm): CreateHotspotRequest {
  return {
    scene_id: form.scene_id,
    type: form.type,
    pitch: form.pitch,
    yaw: form.yaw,
    title: form.title,
    icon_url: form.icon_source === 'custom' ? (form.icon_url || '') : '',
    style: form.icon_source === 'preset' ? (form.icon_preset_key || '') : 'custom',
    target_scene_id: form.type === 1 ? form.target_scene_id : undefined,
    content: form.type === 2 ? form.content : undefined,
    media_type: form.type === 2 ? (form.media_url ? 'image' : 'text') : undefined,
    media_url: form.type === 2 ? form.media_url : undefined,
    transition_effect: form.type === 1 ? form.transition_effect : undefined,
  }
}

function formToUpdateReq(form: HotspotEditorForm): UpdateHotspotRequest {
  return {
    target_scene_id: form.type === 1 ? form.target_scene_id : undefined,
    pitch: form.pitch,
    yaw: form.yaw,
    title: form.title,
    icon_url: form.icon_source === 'custom' ? (form.icon_url || '') : '',
    style: form.icon_source === 'preset' ? (form.icon_preset_key || '') : 'custom',
    content: form.type === 2 ? form.content : undefined,
    media_type: form.type === 2 ? (form.media_url ? 'image' : 'text') : undefined,
    media_url: form.type === 2 ? form.media_url : undefined,
    transition_effect: form.type === 1 ? form.transition_effect : undefined,
  }
}

async function handleSubmit() {
  try {
    const valid = await formRef.value?.validate()
    if (valid !== true) return

    submitting.value = true

    if (props.mode === 'create') {
      const result = await HotspotApi.createHotspot(formToCreateReq(localForm))
      if (result.code === 200) {
        MessagePlugin.success('热点创建成功')
        emit('submitted')
      } else {
        MessagePlugin.error(result.msg || '创建失败')
      }
    } else {
      if (!localForm.id) return
      const result = await HotspotApi.updateHotspot(localForm.id, formToUpdateReq(localForm))
      if (result.code === 200) {
        MessagePlugin.success('热点更新成功')
        emit('submitted')
      } else {
        MessagePlugin.error(result.msg || '更新失败')
      }
    }
  } catch (e: any) {
    MessagePlugin.error(e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function handleDelete() {
  const confirmDialog = DialogPlugin.confirm({
    header: '确认删除',
    body: '确定要删除此热点吗？此操作不可撤销。',
    theme: 'warning',
    onConfirm: async () => {
      try {
        deleting.value = true
        if (!localForm.id) return
        const result = await HotspotApi.deleteHotspot(localForm.id)
        if (result.code === 200) {
          MessagePlugin.success('热点已删除')
          emit('deleted')
        } else {
          MessagePlugin.error(result.msg || '删除失败')
        }
      } catch (e: any) {
        MessagePlugin.error(e.message || '删除失败')
      } finally {
        deleting.value = false
        confirmDialog.destroy()
      }
    },
  })
}
</script>

<style scoped>
.coord-display {
  display: flex;
  gap: 24px;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.04);
  border-radius: 6px;
}

.coord-item {
  font-size: 13px;
  color: #4e5969;
  font-family: monospace;
}

.icon-selector {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
}

.icon-presets {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.icon-preset-item {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid #e5e6eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
}

.icon-preset-item:hover {
  border-color: #0ea5e9;
}

.icon-preset-item.active {
  border-color: #0ea5e9;
  background: rgba(14, 165, 233, 0.08);
}

.icon-preset-item :deep(svg) {
  width: 36px;
  height: 36px;
}

.icon-upload-area {
  width: 100%;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.footer-right {
  display: flex;
  gap: 12px;
}
</style>
