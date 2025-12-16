<template>
  <div class="max-w-4xl mx-auto p-5">
    <el-card shadow="hover" class="rounded-lg">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg font-semibold">用户详情</span>
        </div>
      </template>

      <div v-if="loading" class="flex justify-center py-10">
        <el-spinner type="primary" />
      </div>

      <div v-else-if="userInfo" class="space-y-6">
        <!-- 基本信息 -->
        <div class="bg-gray-50 p-5 rounded-lg">
          <h3 class="text-md font-medium mb-3">基本信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="用户ID">{{ userInfo.id }}</el-descriptions-item>
            <el-descriptions-item label="用户名">{{ userInfo.name }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ userInfo.email }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ userInfo.status === 1 ? '启用' : '禁用' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(userInfo.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatDate(userInfo.updated_at) }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 操作按钮 -->
        <div class="flex justify-center gap-3">
          <el-button type="primary" @click="handleEdit">
            <el-icon><EditPen /></el-icon>
            编辑
          </el-button>
          <el-button type="danger" @click="handleDelete">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
          <el-button @click="handleBack">
            <el-icon><ArrowLeft /></el-icon>
            返回列表
          </el-button>
        </div>
      </div>

      <div v-else class="text-center py-10 text-gray-500">
        <el-icon class="text-4xl mb-2"><WarningFilled /></el-icon>
        <p>用户不存在或已被删除</p>
        <el-button type="primary" @click="handleBack" class="mt-3">
          <el-icon><ArrowLeft /></el-icon>
          返回列表
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { EditPen, Delete, ArrowLeft, WarningFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import UserApi from '@/apis/userApi'

// 路由和导航
const route = useRoute()
const router = useRouter()

// 状态
const loading = ref(false)
const userInfo = ref<any>(null)

// 获取用户ID
const userId = ref<string>(route.params.id as string)

// 格式化日期
const formatDate = (dateString: string): string => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 获取用户详情
const getUserDetail = async () => {
  if (!userId.value) return
  
  loading.value = true
  try {
    const result = await UserApi.getUserDetail(userId.value)
    if (result.code === 200) {
      userInfo.value = result.data
    } else {
      ElMessage.error(result.msg || '获取用户详情失败')
    }
  } catch (error) {
    ElMessage.error('获取用户详情失败，请重试')
  } finally {
    loading.value = false
  }
}

// 编辑用户
const handleEdit = () => {
  router.push(`/user/edit/${userId.value}`)
}

// 删除用户
const handleDelete = async () => {
  if (!userId.value) return
  
  try {
    const result = await UserApi.deleteUser(userId.value)
    if (result.code === 200) {
      ElMessage.success('删除用户成功')
      router.push('/user/list')
    } else {
      ElMessage.error(result.msg || '删除用户失败')
    }
  } catch (error) {
    ElMessage.error('删除用户失败，请重试')
  }
}

// 返回列表
const handleBack = () => {
  router.push('/user/list')
}

// 初始化
onMounted(() => {
  getUserDetail()
})
</script>