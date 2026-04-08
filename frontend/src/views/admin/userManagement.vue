<template>
  <div class="user-management">
    <div class="table-container">
      <t-table
        :data="userList"
        :columns="columns"
        row-key="id"
        hover
        :pagination="paginationConfig"
        @page-change="handlePageChange"
      >
        <template #operations="{ row }">
          <t-space>
            <t-button
              theme="primary"
              variant="text"
              size="small"
              @click="handleEditUser(row)"
            >
              编辑
            </t-button>
            <t-button
              theme="danger"
              variant="text"
              size="small"
              @click="handleDeleteUser(row.id.toString())"
            >
              删除
            </t-button>
          </t-space>
        </template>
      </t-table>
    </div>

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
        <t-form-item label="用户名" name="username">
          <t-input v-model="formData.username" placeholder="请输入用户名" />
        </t-form-item>

        <t-form-item label="邮箱" name="email">
          <t-input v-model="formData.email" placeholder="请输入邮箱" />
        </t-form-item>

        <t-form-item label="手机号" name="phone">
          <t-input v-model="formData.phone" placeholder="请输入手机号" />
        </t-form-item>

        <t-form-item label="昵称" name="nickname">
          <t-input v-model="formData.nickname" placeholder="请输入昵称" />
        </t-form-item>

        <t-form-item label="密码" name="password">
          <t-input v-model="formData.password" type="password" placeholder="请输入密码" />
        </t-form-item>

        <t-form-item label="角色" name="role_id">
          <t-select v-model="formData.role_id" :options="roleOptions" placeholder="请选择角色" />
        </t-form-item>

        <t-form-item label="状态" name="status">
          <t-switch v-model="formData.status" :label="['启用', '禁用']" />
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

    <t-dialog
      v-model:visible="deleteDialogVisible"
      header="确认删除"
      width="400px"
      @confirm="confirmDelete"
      @cancel="deleteDialogVisible = false"
    >
      <p>确定要删除该用户吗？此操作不可撤销。</p>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import UserApi from '@/apis/userApi'

interface UserItem {
  id: number
  username: string
  email: string
  phone: string
  nickname: string
  role_id: number
  status: number
  [key: string]: any
}

const searchKeyword = ref('')
const userList = ref<UserItem[]>([])

const paginationConfig = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  defaultCurrent: 1,
  defaultPageSize: 10,
  showJumper: true
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('新增用户')
const formRef = ref()
const submitLoading = ref(false)
const deleteDialogVisible = ref(false)
const deletingUserId = ref('')

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

const roleOptions = [
  { label: '超级管理员', value: 1 },
  { label: '教师', value: 2 },
  { label: '学生', value: 3 }
]

const columns = computed(() => [
  {
    colKey: 'id',
    title: 'ID',
    width: 80,
    align: 'center',
    fixed: 'left'
  },
  {
    colKey: 'username',
    title: '用户名',
    minWidth: 120,
    ellipsis: true
  },
  {
    colKey: 'nickname',
    title: '昵称',
    minWidth: 100,
    ellipsis: true
  },
  {
    colKey: 'email',
    title: '邮箱',
    minWidth: 180,
    ellipsis: true
  },
  {
    colKey: 'role_id',
    title: '角色',
    width: 100,
    align: 'center',
    cell: ({ row }: { row?: any }) => {
      if (!row) return '-'
      const role = roleOptions.find(r => r.value === row.role_id)
      return role ? role.label : '-'
    }
  },
  {
    colKey: 'status',
    title: '状态',
    width: 80,
    align: 'center',
    cell: ({ row }: { row?: any }) => {
      if (!row) return '-'
      return row.status === 1 ? '启用' : '禁用'
    }
  },
  {
    colKey: 'operations',
    title: '操作',
    width: 140,
    align: 'center',
    fixed: 'right'
  }
])

const loadUserList = async () => {
  try {
    const result = await UserApi.getUserList({
      keyword: searchKeyword.value,
      page: pagination.currentPage,
      page_size: pagination.pageSize
    })
    if (result.code === 200) {
      console.log(result.data.users)
      userList.value = result.data.users
      pagination.total = result.data.total
      paginationConfig.total = result.data.total
    } else {
      MessagePlugin.error(result.msg || '获取用户列表失败')
    }
  } catch (error) {
    MessagePlugin.error('获取用户列表失败')
  }
}

const handleFieldUpdate = async (row: UserItem, field: string, value: any) => {
  const originalValue = row[field]
  if (originalValue === value) {
    return
  }

  try {
    const result = await UserApi.updateUser(row.id.toString(), { [field]: value })
    if (result.code === 200) {
      MessagePlugin.success('更新成功')
    } else {
      MessagePlugin.error(result.msg || '更新失败')
      row[field] = originalValue
    }
  } catch (error) {
    MessagePlugin.error('更新失败')
    row[field] = originalValue
  }
}

const handleSearch = () => {
  pagination.currentPage = 1
  paginationConfig.current = 1
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

const handleEditUser = (user: UserItem) => {
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
        result = await UserApi.addUser(formData as any)
      }

      if (result.code === 200) {
        MessagePlugin.success(formData.id ? '编辑用户成功' : '新增用户成功')
        dialogVisible.value = false
        loadUserList()
      } else {
        MessagePlugin.error(result.msg || '操作失败')
      }
    } catch (error) {
      MessagePlugin.error('操作失败')
    } finally {
      submitLoading.value = false
    }
  }
}

const handleDeleteUser = (userId: string) => {
  deletingUserId.value = userId
  deleteDialogVisible.value = true
}

const confirmDelete = async () => {
  try {
    const result = await UserApi.deleteUser(deletingUserId.value)
    if (result.code === 200) {
      MessagePlugin.success('删除成功')
      deleteDialogVisible.value = false
      loadUserList()
    } else {
      MessagePlugin.error(result.msg || '删除失败')
    }
  } catch (error) {
    MessagePlugin.error('删除失败')
  }
}

const handlePageChange = (context: any) => {
  pagination.currentPage = context.current
  paginationConfig.current = context.current
  pagination.pageSize = context.pageSize
  paginationConfig.pageSize = context.pageSize
  loadUserList()
}

onMounted(() => {
  loadUserList()
})
</script>

<style scoped>
.user-management {
  width: 100%;
  height: 100%;
  min-height: 100%;
  display: flex;
  flex-direction: column;
}

.table-container {
  background-color: #fff;
  border-radius: 8px;
  flex: 1;
  min-height: 400px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
}

:deep(.t-table) {
  min-height: 100%;
  display: flex;
  flex-direction: column;
}

:deep(.t-table__content) {
  flex: 1;
  display: flex;
  flex-direction: column;
}

:deep(.t-table__body-wrapper) {
  flex: 1;
}

:deep(.t-pagination) {
  margin-top: auto;
}
</style>
