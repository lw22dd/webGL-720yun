<template>
  <div class="user-management">
    <Header />
    <main class="container mx-auto p-16">
      <h1 class="text-2xl font-bold mb-8">用户管理</h1>
      
      <!-- 操作按钮区 -->
      <div class="flex justify-between items-center mb-6">
        <div>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索用户名或邮箱"
            style="width: 300px; margin-right: 10px"
          >
            <template #append>
              <el-button @click="handleSearch">搜索</el-button>
            </template>
          </el-input>
        </div>
        <el-button type="primary" @click="handleAddUser">新增用户</el-button>
      </div>
      
      <!-- 用户列表 -->
      <el-table :data="userList" style="width: 100%" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="role.name" label="角色" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="scope">
            <el-switch
              v-model="scope.row.status"
              @change="(newStatus: any) => handleStatusChange(scope.row, newStatus)"
              active-value="1"
              inactive-value="0"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEditUser(scope.row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDeleteUser(scope.row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="flex justify-center mt-6">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
      
      <!-- 新增/编辑用户弹窗 -->
      <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" center>
        <el-form :model="formData" :rules="formRules" ref="formRef">
          <el-form-item prop="username">
            <el-input v-model="formData.username" placeholder="请输入用户名" />
          </el-form-item>
          <el-form-item prop="email">
            <el-input v-model="formData.email" placeholder="请输入邮箱" type="email" />
          </el-form-item>
          <el-form-item prop="password" v-if="!formData.id">
            <el-input v-model="formData.password" placeholder="请输入密码" type="password" />
          </el-form-item>
          <el-form-item prop="phone">
            <el-input v-model="formData.phone" placeholder="请输入手机号" />
          </el-form-item>
          <el-form-item prop="nickname">
            <el-input v-model="formData.nickname" placeholder="请输入昵称" />
          </el-form-item>
          <el-form-item prop="role_id">
            <el-select v-model="formData.role_id" placeholder="请选择角色">
              <el-option label="管理员" value="1" />
              <el-option label="教师" value="2" />
              <el-option label="学生" value="3" />
            </el-select>
          </el-form-item>
          <el-form-item prop="status">
            <el-select v-model="formData.status" placeholder="请选择状态">
              <el-option label="启用" value="1" />
              <el-option label="禁用" value="0" />
            </el-select>
          </el-form-item>
        </el-form>
        <template #footer>
          <span>
            <el-button @click="dialogVisible = false">取消</el-button>
            <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
              确定
            </el-button>
          </span>
        </template>
      </el-dialog>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import Header from '@/components/Header.vue'
import UserApi from '@/apis/userApi'

// 路由
import { useRouter } from 'vue-router'
const router = useRouter()

// 搜索关键字
const searchKeyword = ref('')

// 用户列表
const userList = ref([])

// 分页信息
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

// 弹窗状态
const dialogVisible = ref(false)
const dialogTitle = ref('新增用户')
const formRef = ref()
const submitLoading = ref(false)

// 表单数据
const formData = reactive({
  id: '',
  username: '',
  email: '',
  password: '',
  phone: '',
  nickname: '',
  role_id: 3,
  status: 1
})

// 表单验证规则
const formRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应在3-20个字符之间', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
  ],
  role_id: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 加载用户列表
const loadUserList = async () => {
  try {
    const result = await UserApi.getUserList({
      keyword: searchKeyword.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    console.log(result)
    if (result.code === 200) {
      userList.value = result.data.users
      pagination.total = result.data.total
    } else {
      ElMessage.error(result.msg || '获取用户列表失败')
    }
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  }
}

// 搜索用户
const handleSearch = () => {
  pagination.currentPage = 1
  loadUserList()
}

// 新增用户
const handleAddUser = () => {
  dialogTitle.value = '新增用户'
  Object.assign(formData, {
    id: '',
    username: '',
    email: '',
    password: '',
    phone: '',
    nickname: '',
    role_id: 3,
    status: 1
  })
  dialogVisible.value = true
}

// 编辑用户
const handleEditUser = (user: any) => {
  dialogTitle.value = '编辑用户'
  Object.assign(formData, {
    id: user.id,
    username: user.username,
    email: user.email,
    phone: user.phone || '',
    nickname: user.nickname || '',
    role_id: user.role_id,
    status: user.status
  })
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      submitLoading.value = true
      try {
        let result
        if (formData.id) {
          // 编辑用户
          result = await UserApi.updateUser(formData.id, formData)
        } else {
          // 新增用户
          result = await UserApi.addUser(formData)
        }
        
        if (result.code === 200) {
          ElMessage.success(formData.id ? '编辑用户成功' : '新增用户成功')
          dialogVisible.value = false
          loadUserList()
        } else {
          ElMessage.error(result.msg || '操作失败')
        }
      } catch (error) {
        ElMessage.error('操作失败')
      } finally {
        submitLoading.value = false
      }
    }
  })
}

// 删除用户
const handleDeleteUser = (userId: string) => {
  ElMessageBox.confirm('确定要删除该用户吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const result = await UserApi.deleteUser(userId)
      if (result.code === 200) {
        ElMessage.success('删除用户成功')
        loadUserList()
      } else {
        ElMessage.error(result.msg || '删除用户失败')
      }
    } catch (error) {
      ElMessage.error('删除用户失败')
    }
  }).catch(() => {
    // 取消删除
  })
}

// 状态变更
const handleStatusChange = async (user: any, newStatus: any) => {
  // 确保用户对象包含有效的id，避免发送无效请求
  if (!user || !user.id) {
    return
  }
  
  // 只有当新状态与旧状态不同时才发送请求，避免初始渲染时的无效请求
  if (user.status === newStatus) {
    return
  }
  
  try {
    const result = await UserApi.updateUser(user.id, { status: newStatus })
    if (result.code !== 200) {
      ElMessage.error(result.msg || '更新状态失败')
      // 恢复原状态
      user.status = user.status === '1' ? '0' : '1'
    } else {
      // 更新成功，同步状态
      user.status = newStatus
    }
  } catch (error) {
    ElMessage.error('更新状态失败')
    // 恢复原状态
    user.status = user.status === '1' ? '0' : '1'
  }
}

// 分页大小变更
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  loadUserList()
}

// 当前页变更
const handleCurrentChange = (current: number) => {
  pagination.currentPage = current
  loadUserList()
}

// 初始化
onMounted(() => {
  loadUserList()
})
</script>

<style scoped>
.user-management {
  width: 100%;
  min-height: 100vh;
}
</style>