<template>
  <div class="max-w-4xl mx-auto p-5">
    <t-card shadow="hover" class="rounded-lg">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg font-semibold">用户详情</span>
        </div>
      </template>

      <div v-if="loading" class="flex justify-center py-10">
        <t-loading />
      </div>

      <div v-else-if="userInfo" class="space-y-6">
        <div class="bg-gray-50 p-5 rounded-lg">
          <h3 class="text-md font-medium mb-3">基本信息</h3>
          <t-descriptions :column="2" border>
            <t-description-item label="用户ID">{{ userInfo.id }}</t-description-item>
            <t-description-item label="用户名">{{ userInfo.name }}</t-description-item>
            <t-description-item label="邮箱">{{ userInfo.email }}</t-description-item>
            <t-description-item label="状态">{{ userInfo.status === 1 ? '启用' : '禁用' }}</t-description-item>
            <t-description-item label="创建时间">{{ formatDate(userInfo.created_at) }}</t-description-item>
            <t-description-item label="更新时间">{{ formatDate(userInfo.updated_at) }}</t-description-item>
          </t-descriptions>
        </div>

        <div class="flex justify-center gap-3">
          <t-button theme="primary" @click="handleEdit">
            编辑
          </t-button>
          <t-button theme="danger" @click="handleDelete">
            删除
          </t-button>
          <t-button @click="handleBack">
            返回列表
          </t-button>
        </div>
      </div>

      <div v-else class="text-center py-10 text-gray-500">
        <p>用户不存在或已被删除</p>
        <t-button theme="primary" @click="handleBack" class="mt-3">
          返回列表
        </t-button>
      </div>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import UserApi from '@/services/api/user.api'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const userInfo = ref<any>(null)

const userId = ref<string>(route.params.id as string)

const formatDate = (dateString: string): string => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

const getUserDetail = async () => {
  if (!userId.value) return

  loading.value = true
  try {
    const result = await UserApi.getUserDetail(userId.value)
    if (result.code === 200) {
      userInfo.value = result.data
    } else {
      MessagePlugin.error(result.msg || '获取用户详情失败')
    }
  } catch (error) {
    MessagePlugin.error('获取用户详情失败，请重试')
  } finally {
    loading.value = false
  }
}

const handleEdit = () => {
  router.push(`/user/edit/${userId.value}`)
}

const handleDelete = async () => {
  if (!userId.value) return

  try {
    const result = await UserApi.deleteUser(userId.value)
    if (result.code === 200) {
      MessagePlugin.success('删除用户成功')
      router.push('/user/list')
    } else {
      MessagePlugin.error(result.msg || '删除用户失败')
    }
  } catch (error) {
    MessagePlugin.error('删除用户失败，请重试')
  }
}

const handleBack = () => {
  router.push('/user/list')
}

onMounted(() => {
  getUserDetail()
})
</script>
