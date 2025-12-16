<template>
  <div class="max-w-4xl mx-auto p-5">
    <el-card shadow="hover" class="rounded-lg">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg font-semibold">学生批量注册</span>
          <el-button type="primary" size="small" @click="handleDownloadTemplate">
            <el-icon><Download /></el-icon>
            下载模板
          </el-button>
        </div>
      </template>

      <div 
        class="relative border-2 border-dashed border-gray-300 rounded-lg p-16 text-center transition-all duration-300 cursor-pointer overflow-hidden hover:border-blue-500"
        :class="{'border-blue-500 bg-blue-50': isDragover}"
        @drop="handleDrop" 
        @dragover.prevent 
        @dragenter="handleDragEnter" 
        @dragleave="handleDragLeave"
      >
        
        <input 
          type="file" 
          ref="fileInput" 
          accept=".xlsx, .csv" 
          class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
          @change="handleFileChange" 
        />
        
        <div class="flex flex-col items-center gap-4">
          <el-icon class="text-blue-500 text-8xl"><UploadFilled /></el-icon>
          <h3 class="text-xl font-medium text-gray-800">点击或拖拽文件到此处上传</h3>
          <p class="text-gray-600">支持 .xlsx 和 .csv 格式文件</p>
          <el-button type="primary" @click="triggerFileInput" :disabled="isUploading">
            <el-icon><DocumentAdd /></el-icon>
            选择文件
          </el-button>
        </div>
      </div>

      <!-- 文件信息 -->
      <div v-if="selectedFile" class="mt-5 p-3 bg-gray-50 rounded-md">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="文件名">
            {{ selectedFile.name }}
          </el-descriptions-item>
          <el-descriptions-item label="文件大小">
            {{ formatFileSize(selectedFile.size) }}
          </el-descriptions-item>
          <el-descriptions-item label="上传状态">
            <el-tag :type="uploadStatus === 'success' ? 'success' : uploadStatus === 'error' ? 'danger' : ''">
              {{ uploadStatusMap[uploadStatus] }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- 上传进度 -->
      <div v-if="isUploading" class="mt-5 flex items-center gap-3">
        <el-progress 
          :percentage="uploadProgress" 
          status="active" 
          :text-inside="true"
          class="flex-1"
        />
        <el-button type="danger" size="small" @click="handleCancelUpload">
          取消上传
        </el-button>
      </div>

      <!-- 错误信息 -->
      <div v-if="errorMessage" class="mt-5">
        <el-alert 
          title="上传失败" 
          :description="errorMessage" 
          type="error" 
          show-icon 
          :closable="true"
          @close="errorMessage = ''"
        />
      </div>

      <!-- 成功信息 -->
      <div v-if="successMessage" class="mt-5">
        <el-alert 
          title="上传成功" 
          :description="successMessage" 
          type="success" 
          show-icon 
          :closable="true"
          @close="successMessage = ''"
        />
      </div>

      <!-- 注册结果 -->
      <div v-if="registerResult" class="mt-5">
        <el-divider content-position="left">注册结果</el-divider>
        <div class="flex gap-10 mt-5 p-5 bg-gray-50 rounded-lg">
          <div class="flex items-center gap-3">
            <span class="font-medium text-gray-600">总记录数:</span>
            <span class="text-2xl font-bold text-blue-600">{{ registerResult.success_count + registerResult.failed_count }}</span>
          </div>
          <div class="flex items-center gap-3">
            <span class="font-medium text-gray-600">成功:</span>
            <span class="text-2xl font-bold text-green-600">{{ registerResult.success_count }}</span>
          </div>
          <div class="flex items-center gap-3">
            <span class="font-medium text-gray-600">失败:</span>
            <span class="text-2xl font-bold text-red-600">{{ registerResult.failed_count }}</span>
          </div>
        </div>

        <!-- 失败详情 -->
        <div v-if="registerResult.failed_count > 0" class="mt-5">
          <el-collapse v-model="activeNames">
            <el-collapse-item title="查看失败详情" name="1">
              <el-table :data="registerResult.errors" stripe border size="small">
                <el-table-column prop="index" label="行号" width="80" />
                <el-table-column prop="username" label="用户名" width="120" />
                <el-table-column prop="error" label="错误信息" show-overflow-tooltip />
              </el-table>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div v-if="selectedFile && !isUploading" class="mt-5 flex justify-center gap-3">
        <el-button type="primary" @click="handleUpload" :loading="isUploading">
          <el-icon><Upload /></el-icon>
          开始上传
        </el-button>
        <el-button @click="handleReset">
          <el-icon><RefreshRight /></el-icon>
          重新选择
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { 
  UploadFilled, 
  Upload, 
  DocumentAdd, 
  Download, 
  RefreshRight 
} from '@element-plus/icons-vue'
import UserApi from '@/apis/userApi'

// 定义上传状态
const UPLOAD_STATUS = {
  IDLE: 'idle',
  UPLOADING: 'uploading',
  SUCCESS: 'success',
  ERROR: 'error'
} as const

// 状态映射
const uploadStatusMap = {
  [UPLOAD_STATUS.IDLE]: '未开始',
  [UPLOAD_STATUS.UPLOADING]: '上传中',
  [UPLOAD_STATUS.SUCCESS]: '上传成功',
  [UPLOAD_STATUS.ERROR]: '上传失败'
}

// 响应式数据
const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const isUploading = ref(false)
const isDragover = ref(false)
const uploadProgress = ref(0)
const uploadStatus = ref<typeof UPLOAD_STATUS[keyof typeof UPLOAD_STATUS]>(UPLOAD_STATUS.IDLE)
const errorMessage = ref('')
const successMessage = ref('')
const registerResult = ref<any>(null)
const activeNames = ref(['1'])

// 文件拖拽事件
const handleDragEnter = (e: DragEvent) => {
  e.preventDefault()
  isDragover.value = true
}

const handleDragLeave = (e: DragEvent) => {
  e.preventDefault()
  isDragover.value = false
}

const handleDragOver = (e: DragEvent) => {
  e.preventDefault()
}

const handleDrop = (e: DragEvent) => {
  e.preventDefault()
  isDragover.value = false
  
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    handleFileSelect(files[0])
  }
}

// 触发文件选择
const triggerFileInput = () => {
  fileInput.value?.click()
}

// 文件选择事件
const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    handleFileSelect(target.files[0])
  }
}

// 处理文件选择
const handleFileSelect = (file: File) => {
  // 验证文件类型
  const fileExtension = file.name.split('.').pop()?.toLowerCase()
  if (fileExtension !== 'xlsx' && fileExtension !== 'csv') {
    errorMessage.value = '请选择 .xlsx 或 .csv 格式的文件'
    return
  }

  // 验证文件大小（最大10MB）
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

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 重置状态
const handleReset = () => {
  selectedFile.value = null
  isUploading.value = false
  uploadProgress.value = 0
  uploadStatus.value = UPLOAD_STATUS.IDLE
  errorMessage.value = ''
  successMessage.value = ''
  registerResult.value = null
  // 清空文件输入
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

// 取消上传
const handleCancelUpload = () => {
  isUploading.value = false
  uploadProgress.value = 0
  uploadStatus.value = UPLOAD_STATUS.IDLE
  errorMessage.value = '上传已取消'
}

// 下载模板
const handleDownloadTemplate = () => {
  // 这里可以实现模板下载逻辑
  // 目前只是模拟下载
  successMessage.value = '模板下载功能开发中...'
}

// 上传文件
const handleUpload = async () => {
  if (!selectedFile.value) return
  
  try {
    isUploading.value = true
    uploadStatus.value = UPLOAD_STATUS.UPLOADING
    uploadProgress.value = 0
    errorMessage.value = ''
    successMessage.value = ''
    registerResult.value = null

    // 调用实际的API，使用真实的进度跟踪
    const response = await UserApi.batchRegister(selectedFile.value, (progress) => {
      uploadProgress.value = progress
    })
    
    // 确保进度达到100%
    uploadProgress.value = 100

    // 处理响应
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

