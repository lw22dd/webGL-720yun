<template>
  <div class="max-w-7xl mx-auto p-5">
    <el-card shadow="hover" class="rounded-lg">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg font-semibold">用户列表</span>
          <el-button type="primary" @click="handleAddUser">
            <el-icon><Plus /></el-icon>
            新增用户
          </el-button>
        </div>
      </template>

      <!-- 搜索和筛选 -->
      <div class="flex flex-wrap gap-3 mb-5">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索用户名或邮箱"
          size="small"
          prefix-icon="Search"
          class="w-64"
        >
          <template #append>
            <el-button @click="handleSearch" size="small" type="primary" icon="Search"></el-button>
          </template>
        </el-input>

        <el-select
          v-model="statusFilter"
          placeholder="用户状态"
          size="small"
          class="w-40"
        >
          <el-option label="全部" value="" />
          <el-option label="启用" value="1" />
          <el-option label="禁用" value="0" />
        </el-select>

        <div class="flex gap-2 ml-auto">
          <el-button @click="handleResetFilter" size="small">
            <el-icon><RefreshRight /></el-icon>
            重置
          </el-button>
        </div>
      </div>

      <!-- 用户表格 -->
      <div class="overflow-x-auto">
        <el-table :data="userList" stripe border size="small" class="min-w-full">
          <el-table-column prop="id" label="用户ID" width="80" />
          <el-table-column prop="name" label="用户名" width="120" />
          <el-table-column prop="email" label="邮箱" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                {{ scope.row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180" />
          <el-table-column prop="updated_at" label="更新时间" width="180" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="scope">
              <el-button
                type="primary"
                size="small"
                @click="handleViewUser(scope.row.id)"
                class="mr-2"
              >
                <el-icon><View /></el-icon>
                查看
              </el-button>
              <el-button
                type="warning"
                size="small"
                @click="handleEditUser(scope.row.id)"
                class="mr-2"
              >
                <el-icon><EditPen /></el-icon>
                编辑
              </el-button>
              <el-button
                type="danger"
                size="small"
                @click="handleDeleteUser(scope.row)"
              >
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="flex justify-end mt-4">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, RefreshRight, View, EditPen, Delete, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import UserApi from '@/apis/userApi'

// 路由
const router = useRouter()

// 搜索和筛选
const searchKeyword = ref('')
const statusFilter = ref('')

// 分页信息
const pagination = ref({
  current: 1,
  size: 10,
  total: 0
})

// 用户列表
const userList = ref<any[]>([])
const loading = ref(false)

// 获取用户列表
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
      ElMessage.error(result.msg || '获取用户列表失败')
    }
  } catch (error) {
    ElMessage.error('获取用户列表失败，请重试')
  } finally {
    loading.value = false
  }
}

// 查看用户详情
const handleViewUser = (userId: string) => {
  router.push(`/user/${userId}`)
}

// 编辑用户
const handleEditUser = (userId: string) => {
  router.push(`/user/edit/${userId}`)
}

// 删除用户
const handleDeleteUser = async (user: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户"${user.name}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const result = await UserApi.deleteUser(user.id)
    if (result.code === 200) {
      ElMessage.success('删除用户成功')
      getUserList() // 重新获取列表
    } else {
      ElMessage.error(result.msg || '删除用户失败')
    }
  } catch (error) {
    // 取消删除
  }
}

// 新增用户
const handleAddUser = () => {
  router.push('/user/add')
}

// 搜索
const handleSearch = () => {
  pagination.value.current = 1
  getUserList()
}

// 重置筛选
const handleResetFilter = () => {
  searchKeyword.value = ''
  statusFilter.value = ''
  pagination.value.current = 1
  getUserList()
}

// 分页大小变化
const handleSizeChange = (size: number) => {
  pagination.value.size = size
  pagination.value.current = 1
  getUserList()
}

// 当前页变化
const handleCurrentChange = (current: number) => {
  pagination.value.current = current
  getUserList()
}

// 初始化
onMounted(() => {
  getUserList()
})
</script>