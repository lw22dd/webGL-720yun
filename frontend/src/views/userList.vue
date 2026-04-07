<template>
  <div class="max-w-7xl mx-auto p-5">
    <t-card shadow="hover" class="rounded-lg">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg font-semibold">用户列表</span>
          <t-button theme="primary" @click="handleAddUser">
            新增用户
          </t-button>
        </div>
      </template>

      <div class="flex flex-wrap gap-3 mb-5">
        <t-input
          v-model="searchKeyword"
          placeholder="搜索用户名或邮箱"
          size="small"
          class="w-64"
          @enter="handleSearch"
        >
        </t-input>

        <t-select
          v-model="statusFilter"
          placeholder="用户状态"
          size="small"
          class="w-40"
        >
          <t-option label="全部" value="" />
          <t-option label="启用" value="1" />
          <t-option label="禁用" value="0" />
        </t-select>

        <div class="flex gap-2 ml-auto">
          <t-button @click="handleResetFilter" size="small">
            重置
          </t-button>
        </div>
      </div>

      <div class="overflow-x-auto">
        <t-table :data="userList" stripe border size="small" class="min-w-full" hover>
          <t-table-column prop="id" label="用户ID" width="80" />
          <t-table-column prop="name" label="用户名" width="120" />
          <t-table-column prop="email" label="邮箱" />
          <t-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <t-tag :theme="scope.row.status === 1 ? 'success' : 'danger'">
                {{ scope.row.status === 1 ? '启用' : '禁用' }}
              </t-tag>
            </template>
          </t-table-column>
          <t-table-column prop="created_at" label="创建时间" width="180" />
          <t-table-column prop="updated_at" label="更新时间" width="180" />
          <t-table-column label="操作" width="200" fixed="right">
            <template #default="scope">
              <t-button
                theme="primary"
                size="small"
                @click="handleViewUser(scope.row.id)"
                class="mr-2"
              >
                查看
              </t-button>
              <t-button
                theme="warning"
                size="small"
                @click="handleEditUser(scope.row.id)"
                class="mr-2"
              >
                编辑
              </t-button>
              <t-button
                theme="danger"
                size="small"
                @click="handleDeleteUser(scope.row)"
              >
                删除
              </t-button>
            </template>
          </t-table-column>
        </t-table>
      </div>

      <div class="flex justify-end mt-4">
        <t-pagination
          v-model:current="pagination.current"
          v-model:pageSize="pagination.size"
          :page-size-options="[10, 20, 50, 100]"
          :total="pagination.total"
          @change="handlePageChange"
        />
      </div>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message, DialogPlugin } from 'tdesign-vue-next'
import UserApi from '@/apis/userApi'

const router = useRouter()

const searchKeyword = ref('')
const statusFilter = ref('')

const pagination = ref({
  current: 1,
  size: 10,
  total: 0
})

const userList = ref<any[]>([])
const loading = ref(false)

const getUserList = async () => {
  loading.value = true
  try {
    const result = await UserApi.getUserList({
      keyword: searchKeyword.value,
      status: statusFilter.value,
      page: pagination.value.current,
      page_size: pagination.value.size
    })
    if (result.code === 200) {
      userList.value = result.data.list || []
      pagination.value.total = result.data.total || 0
    } else {
      Message.error(result.msg || '获取用户列表失败')
    }
  } catch (error) {
    Message.error('获取用户列表失败，请重试')
  } finally {
    loading.value = false
  }
}

const handleViewUser = (userId: string) => {
  router.push(`/user/${userId}`)
}

const handleEditUser = (userId: string) => {
  router.push(`/user/edit/${userId}`)
}

const handleDeleteUser = async (user: any) => {
  const dialog = DialogPlugin.confirm({
    header: '删除确认',
    body: `确定要删除用户"${user.name}"吗？`,
    confirmBtn: '确定',
    cancelBtn: '取消',
    onConfirm: async () => {
      dialog.destroy()
      const result = await UserApi.deleteUser(user.id)
      if (result.code === 200) {
        Message.success('删除用户成功')
        getUserList()
      } else {
        Message.error(result.msg || '删除用户失败')
      }
    },
    onCancel: () => {
      dialog.destroy()
    }
  })
}

const handleAddUser = () => {
  router.push('/user/add')
}

const handleSearch = () => {
  pagination.value.current = 1
  getUserList()
}

const handleResetFilter = () => {
  searchKeyword.value = ''
  statusFilter.value = ''
  pagination.value.current = 1
  getUserList()
}

const handlePageChange = (context: any) => {
  pagination.value.current = context.current
  pagination.value.size = context.pageSize
  getUserList()
}

onMounted(() => {
  getUserList()
})
</script>
