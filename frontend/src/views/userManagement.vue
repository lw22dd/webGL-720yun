<template>
  <div class="user-management">
    <Header />
    <main class="container mx-auto p-16">
      <h1 class="text-2xl font-bold mb-8">用户管理</h1>

      <div class="flex justify-between items-center mb-6">
        <div>
          <t-input
            v-model="searchKeyword"
            placeholder="搜索用户名或邮箱"
            style="width: 300px; margin-right: 10px"
            @enter="handleSearch"
          >
          </t-input>
        </div>
        <t-button theme="primary" @click="handleAddUser">新增用户</t-button>
      </div>

      

      <div class="flex justify-center mt-6">
        <t-pagination
          v-model:current="pagination.currentPage"
          v-model:pageSize="pagination.pageSize"
          :page-size-options="[10, 20, 50, 100]"
          :total="pagination.total"
          @change="handlePageChange"
        />
      </div>

      <t-dialog v-model:visible="dialogVisible" :header="dialogTitle" width="500px" center>
        <t-form :data="formData" :rules="formRules" ref="formRef">
          <t-form-item label="用户名" name="username">
            <t-input v-model="formData.username" placeholder="请输入用户名" />
          </t-form-item>
          <t-form-item label="邮箱" name="email">
            <t-input v-model="formData.email" placeholder="请输入邮箱" />
          </t-form-item>
          <t-form-item v-if="!formData.id" label="密码" name="password">
            <t-input v-model="formData.password" placeholder="请输入密码" type="password" />
          </t-form-item>
          <t-form-item label="手机号" name="phone">
            <t-input v-model="formData.phone" placeholder="请输入手机号" />
          </t-form-item>
          <t-form-item label="昵称" name="nickname">
            <t-input v-model="formData.nickname" placeholder="请输入昵称" />
          </t-form-item>
          <t-form-item label="角色" name="role_id">
            <t-select v-model="formData.role_id" placeholder="请选择角色">
              <t-option label="管理员" value="1" />
              <t-option label="教师" value="2" />
              <t-option label="学生" value="3" />
            </t-select>
          </t-form-item>
          <t-form-item label="状态" name="status">
            <t-select v-model="formData.status" placeholder="请选择状态">
              <t-option label="启用" value="1" />
              <t-option label="禁用" value="0" />
            </t-select>
          </t-form-item>
        </t-form>
        <template #footer>
          <div class="flex justify-center gap-3">
            <t-button @click="dialogVisible = false">取消</t-button>
            <t-button theme="primary" @click="handleSubmit" :loading="submitLoading">
              确定
            </t-button>
          </div>
        </template>
      </t-dialog>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Message } from 'tdesign-vue-next'
import Header from '@/components/Header.vue'
import UserApi from '@/apis/userApi'

const searchKeyword = ref('')

const userList = ref([])

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('新增用户')
const formRef = ref()
const submitLoading = ref(false)

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

const loadUserList = async () => {
  try {
    const result = await UserApi.getUserList({
      keyword: searchKeyword.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    console.log(result)
    if (result.code === 200) {
      console.log(result.data.users)
      userList.value = result.data.users
      pagination.total = result.data.total
    } else {
      Message.error(result.msg || '获取用户列表失败')
    }
  } catch (error) {
    Message.error('获取用户列表失败')
  }
}

const handleSearch = () => {
  pagination.currentPage = 1
  loadUserList()
}

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

const handleSubmit = async () => {
  if (!formRef.value) return

  const valid = await (formRef.value as any).validate()
  if (valid) {
    submitLoading.value = true
    try {
      let result
      if (formData.id) {
        result = await UserApi.updateUser(formData.id, formData)
      } else {
        result = await UserApi.addUser(formData)
      }

      if (result.code === 200) {
        Message.success(formData.id ? '编辑用户成功' : '新增用户成功')
        dialogVisible.value = false
        loadUserList()
      } else {
        Message.error(result.msg || '操作失败')
      }
    } catch (error) {
      Message.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  }
}

const handleDeleteUser = (_userId: string) => {
  Message.warning('删除功能已禁用，请联系管理员')
}

const handleStatusChange = async (user: any, newStatus: any) => {
  if (!user || !user.id) {
    return
  }

  if (user.status === newStatus) {
    return
  }

  try {
    const result = await UserApi.updateUser(user.id, { status: newStatus })
    if (result.code !== 200) {
      Message.error(result.msg || '更新状态失败')
      user.status = user.status === 1 ? 0 : 1
    } else {
      user.status = newStatus
    }
  } catch (error) {
    Message.error('更新状态失败')
    user.status = user.status === 1 ? 0 : 1
  }
}

const handlePageChange = (context: any) => {
  pagination.currentPage = context.current
  pagination.pageSize = context.pageSize
  loadUserList()
}

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
