<template>
  <div class="excel-upload-wrapper">
    <t-card class="upload-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <t-icon name="file-excel" class="header-icon" />
            <span class="header-title">学生批量注册</span>
          </div>
          <t-button theme="default" size="small" @click="handleDownloadTemplate" variant="outline">
            下载模板
          </t-button>
        </div>
      </template>

      <div
        class="upload-dropzone"
        :class="{'dragover': isDragover}"
        @drop="handleDrop"
        @dragover.prevent
        @dragenter="handleDragEnter"
        @dragleave="handleDragLeave"
      >
        <input
          type="file"
          ref="fileInput"
          accept=".xlsx, .csv"
          class="file-input"
          @change="handleFileChange"
        />

        <div class="dropzone-content">
          <div class="dropzone-icon">
            <t-icon name="cloud-upload" size="64px" />
          </div>
          <h3 class="dropzone-title">点击或拖拽文件到此处上传</h3>
          <p class="dropzone-desc">支持 .xlsx 和 .csv 格式文件，最大 10MB</p>
          <t-button theme="primary" size="large" @click="triggerFileInput" :disabled="isUploading" class="select-btn">
            选择文件
          </t-button>
        </div>
      </div>

      <div v-if="selectedFile" class="file-info">
        <div class="file-info-header">
          <t-icon name="file" class="file-icon" />
          <span class="file-name">{{ selectedFile.name }}</span>
          <span class="file-size">{{ formatFileSize(selectedFile.size) }}</span>
        </div>
        <div class="file-status">
          <t-tag :theme="statusTheme" size="small">
            {{ uploadStatusMap[uploadStatus] }}
          </t-tag>
        </div>
      </div>

      <div v-if="isUploading" class="upload-progress">
        <div class="progress-bar-wrapper">
          <t-progress
            :percentage="uploadProgress"
            :status="uploadProgress === 100 ? 'success' : 'active'"
            size="large"
            :show-overlay="true"
          />
        </div>
        <t-button theme="default" size="small" @click="handleCancelUpload" variant="text" class="cancel-btn">
          取消上传
        </t-button>
      </div>

      <t-alert
        v-if="errorMessage"
        class="result-alert"
        variant="error"
        :message="errorMessage"
        :closable="true"
        @close="errorMessage = ''"
      />

      <t-alert
        v-if="successMessage"
        class="result-alert"
        variant="success"
        :message="successMessage"
        :closable="true"
        @close="successMessage = ''"
      />

      <div v-if="registerResult" class="register-result">
        <t-divider>
          <span class="divider-text">注册结果</span>
        </t-divider>
        <div class="result-stats">
          <div class="stat-item">
            <span class="stat-label">总记录数</span>
            <span class="stat-value">{{ registerResult.success_count + registerResult.failed_count }}</span>
          </div>
          <div class="stat-item success">
            <span class="stat-label">成功</span>
            <span class="stat-value">{{ registerResult.success_count }}</span>
          </div>
          <div class="stat-item error">
            <span class="stat-label">失败</span>
            <span class="stat-value">{{ registerResult.failed_count }}</span>
          </div>
        </div>

        <t-collapsible v-if="registerResult.failed_count > 0" class="error-collapsible">
          <template #header>
            <div class="collapsible-header">
              <span>查看失败详情</span>
            </div>
          </template>
          <t-table :data="registerResult.errors" stripe size="small" class="error-table">
            <t-table-column prop="index" label="行号" width="80" />
            <t-table-column prop="username" label="用户名" width="140" />
            <t-table-column prop="error" label="错误信息" />
          </t-table>
        </t-collapsible>
      </div>

      <div v-if="selectedFile && !isUploading" class="action-buttons">
        <t-button theme="primary" size="large" @click="handleUpload" :loading="isUploading" class="upload-btn">
          开始上传
        </t-button>
        <t-button theme="default" size="large" @click="handleReset" variant="outline">
          重新选择
        </t-button>
      </div>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import UserApi from '@/services/api/user.api'

const UPLOAD_STATUS = {
  IDLE: 'idle',
  UPLOADING: 'uploading',
  SUCCESS: 'success',
  ERROR: 'error'
} as const

const uploadStatusMap = {
  [UPLOAD_STATUS.IDLE]: '未开始',
  [UPLOAD_STATUS.UPLOADING]: '上传中',
  [UPLOAD_STATUS.SUCCESS]: '上传成功',
  [UPLOAD_STATUS.ERROR]: '上传失败'
}

const statusTheme = computed(() => {
  switch (uploadStatus.value) {
    case UPLOAD_STATUS.SUCCESS: return 'success'
    case UPLOAD_STATUS.ERROR: return 'danger'
    default: return 'default'
  }
})

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const isUploading = ref(false)
const isDragover = ref(false)
const uploadProgress = ref(0)
const uploadStatus = ref<typeof UPLOAD_STATUS[keyof typeof UPLOAD_STATUS]>(UPLOAD_STATUS.IDLE)
const errorMessage = ref('')
const successMessage = ref('')
const registerResult = ref<any>(null)

const handleDragEnter = (e: DragEvent) => {
  e.preventDefault()
  isDragover.value = true
}

const handleDragLeave = (e: DragEvent) => {
  e.preventDefault()
  isDragover.value = false
}

const handleDrop = (e: DragEvent) => {
  e.preventDefault()
  isDragover.value = false

  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    handleFileSelect(files[0])
  }
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    handleFileSelect(target.files[0])
  }
}

const handleFileSelect = (file: File) => {
  const fileExtension = file.name.split('.').pop()?.toLowerCase()
  if (fileExtension !== 'xlsx' && fileExtension !== 'csv') {
    errorMessage.value = '请选择 .xlsx 或 .csv 格式的文件'
    return
  }

  if (file.size > 10 * 1024 * 1024) {
    errorMessage.value = '文件大小不能超过 10MB'
    return
  }

  selectedFile.value = file
  uploadStatus.value = UPLOAD_STATUS.IDLE
  errorMessage.value = ''
  successMessage.value = ''
  registerResult.value = null
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const handleReset = () => {
  selectedFile.value = null
  isUploading.value = false
  uploadProgress.value = 0
  uploadStatus.value = UPLOAD_STATUS.IDLE
  errorMessage.value = ''
  successMessage.value = ''
  registerResult.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const handleCancelUpload = () => {
  isUploading.value = false
  uploadProgress.value = 0
  uploadStatus.value = UPLOAD_STATUS.IDLE
  errorMessage.value = '上传已取消'
}

const handleDownloadTemplate = () => {
  successMessage.value = '模板下载功能开发中...'
}

const handleUpload = async () => {
  if (!selectedFile.value) return

  try {
    isUploading.value = true
    uploadStatus.value = UPLOAD_STATUS.UPLOADING
    uploadProgress.value = 0
    errorMessage.value = ''
    successMessage.value = ''
    registerResult.value = null

    const response = await UserApi.batchRegister(selectedFile.value, (progress) => {
      uploadProgress.value = progress
    })

    uploadProgress.value = 100

    if (response.code === 200) {
      uploadStatus.value = UPLOAD_STATUS.SUCCESS
      successMessage.value = response.msg || '批量注册成功'
      registerResult.value = response.data
    } else {
      throw new Error(response.msg || '上传失败')
    }
  } catch (error) {
    uploadStatus.value = UPLOAD_STATUS.ERROR
    errorMessage.value = error instanceof Error ? error.message : '上传失败，请重试'
  } finally {
    isUploading.value = false
  }
}
</script>

<style scoped>
.excel-upload-wrapper {
  width: 100%;
  max-width: 880px;
  margin: 0 auto;
}

.upload-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-icon {
  font-size: 20px;
  color: #0052d9;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

.upload-dropzone {
  position: relative;
  border: 2px dashed #e5e6eb;
  border-radius: 12px;
  padding: 60px 24px;
  text-align: center;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
  background: #f7f8fa;
}

.upload-dropzone:hover,
.upload-dropzone.dragover {
  border-color: #0052d9;
  background: rgba(0, 82, 217, 0.04);
}

.file-input {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.dropzone-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.dropzone-icon {
  color: #c9cdd4;
  transition: color 0.2s;
}

.upload-dropzone:hover .dropzone-icon,
.upload-dropzone.dragover .dropzone-icon {
  color: #0052d9;
}

.dropzone-title {
  font-size: 18px;
  font-weight: 500;
  color: #1d2129;
  margin: 0;
}

.dropzone-desc {
  font-size: 14px;
  color: #86909c;
  margin: 0;
}

.select-btn {
  border-radius: 8px;
  font-weight: 500;
  margin-top: 8px;
}

.file-info {
  margin-top: 24px;
  padding: 16px 20px;
  background: #f7f8fa;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.file-info-header {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
}

.file-icon {
  font-size: 20px;
  color: #0052d9;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  color: #1d2129;
}

.file-size {
  font-size: 13px;
  color: #86909c;
}

.upload-progress {
  margin-top: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.progress-bar-wrapper {
  flex: 1;
}

.cancel-btn {
  color: #e34d59;
  flex-shrink: 0;
}

.result-alert {
  margin-top: 20px;
  border-radius: 8px;
}

.register-result {
  margin-top: 28px;
}

.divider-text {
  font-size: 14px;
  font-weight: 600;
  color: #4e5969;
}

.result-stats {
  display: flex;
  gap: 40px;
  margin-top: 20px;
  padding: 24px;
  background: linear-gradient(135deg, #f7f8fa 0%, #ffffff 100%);
  border-radius: 12px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-item.success .stat-value {
  color: #00a870;
}

.stat-item.error .stat-value {
  color: #e34d59;
}

.stat-label {
  font-size: 13px;
  color: #86909c;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #1d2129;
}

.error-collapsible {
  margin-top: 20px;
}

.collapsible-header {
  font-size: 14px;
  font-weight: 500;
  color: #4e5969;
}

.error-table {
  margin-top: 12px;
  border-radius: 8px;
}

.action-buttons {
  margin-top: 28px;
  display: flex;
  justify-content: center;
  gap: 16px;
}

.upload-btn {
  border-radius: 8px;
  font-weight: 500;
  min-width: 140px;
}
</style>
