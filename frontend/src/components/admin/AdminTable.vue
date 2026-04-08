<template>
  <div class="admin-table">
    <div v-if="$slots.actionBar || showBatchActions" class="table-header">
      <div class="table-title">
        <slot name="title" />
      </div>
      <div class="table-actions">
        <slot name="actionBar" />
        <t-button
          v-if="showBatchActions && selectedRowKeys.length > 0"
          theme="danger"
          size="medium"
          @click="handleBatchDelete"
        >
          批量删除 ({{ selectedRowKeys.length }})
        </t-button>
      </div>
    </div>

    <div class="table-container">
      <t-table
        :data="data"
        :columns="columns"
        row-key="id"
        hover
        :pagination="paginationConfig"
        :selected-row-keys="selectedRowKeys"
        :select-on-row-click="selectOnRowClick"
        @page-change="handlePageChange"
        @select-change="handleSelectChange"
      >
        <template v-for="slot in Object.keys($slots)" #[slot]="slotProps">
          <slot :name="slot" v-bind="slotProps" />
        </template>
      </t-table>
    </div>

    <t-dialog
      v-model:visible="deleteDialogVisible"
      header="确认删除"
      width="400px"
      @confirm="confirmDelete"
      @cancel="deleteDialogVisible = false"
    >
      <p>确定要删除该{{ deleteItemName }}吗？此操作不可撤销。</p>
    </t-dialog>

    <t-dialog
      v-model:visible="batchDeleteDialogVisible"
      header="确认批量删除"
      width="400px"
      @confirm="confirmBatchDelete"
      @cancel="batchDeleteDialogVisible = false"
    >
      <p>确定要删除选中的 {{ selectedRowKeys.length }} 项{{ deleteItemName }}吗？此操作不可撤销。</p>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { PageInfo } from 'tdesign-vue-next'

export interface PaginationConfig {
  current: number
  pageSize: number
  total: number
  defaultCurrent?: number
  defaultPageSize?: number
  showJumper?: boolean
}

export interface TableColumn {
  colKey: string
  title?: string
  width?: number | string
  minWidth?: number | string
  align?: 'left' | 'center' | 'right'
  fixed?: 'left' | 'right'
  ellipsis?: boolean
  cell?: (params: { row: any; rowIndex: number }) => any
  type?: 'multiple' | 'single'
}

const props = withDefaults(defineProps<{
  data: any[]
  columns: TableColumn[]
  pagination?: PaginationConfig
  deleteItemName?: string
  showBatchActions?: boolean
  selectOnRowClick?: boolean
}>(), {
  deleteItemName: '该项',
  showBatchActions: true,
  selectOnRowClick: false
})

const emit = defineEmits<{
  (e: 'page-change', context: { current: number; pageSize: number }): void
  (e: 'delete', id: string | number): void
  (e: 'batch-delete', ids: (string | number)[]): void
  (e: 'edit', row: any): void
  (e: 'create'): void
  (e: 'select-change', context: { selectedRowKeys: (string | number)[]; selectedRowData: any[] }): void
}>()

const paginationConfig = computed<PaginationConfig>(() => ({
  current: 1,
  pageSize: 10,
  total: 0,
  defaultCurrent: 1,
  defaultPageSize: 10,
  showJumper: true,
  ...props.pagination
}))

const selectedRowKeys = ref<(string | number)[]>([])

const deleteDialogVisible = ref(false)
const batchDeleteDialogVisible = ref(false)
const deletingId = ref<string | number | null>(null)

const handlePageChange = (context: PageInfo) => {
  emit('page-change', { current: context.current, pageSize: context.pageSize })
}

const handleSelectChange = (value: (string | number)[], context: { selectedRowData: any[] }) => {
  selectedRowKeys.value = value
  emit('select-change', { selectedRowKeys: value, selectedRowData: context.selectedRowData })
}

const openDeleteDialog = (id: string | number) => {
  deletingId.value = id
  deleteDialogVisible.value = true
}

const openEditDialog = (row: any) => {
  emit('edit', row)
}

const handleCreate = () => {
  emit('create')
}

const handleBatchDelete = () => {
  if (selectedRowKeys.value.length > 0) {
    batchDeleteDialogVisible.value = true
  }
}

const clearSelection = () => {
  selectedRowKeys.value = []
}

defineExpose({
  openDeleteDialog,
  openEditDialog,
  handleCreate,
  clearSelection,
  selectedRowKeys
})

const confirmDelete = () => {
  if (deletingId.value !== null) {
    emit('delete', deletingId.value)
  }
  deleteDialogVisible.value = false
  deletingId.value = null
}

const confirmBatchDelete = () => {
  if (selectedRowKeys.value.length > 0) {
    emit('batch-delete', [...selectedRowKeys.value])
    clearSelection()
  }
  batchDeleteDialogVisible.value = false
}
</script>

<style scoped>
.admin-table {
  width: 100%;
  height: 100%;
  min-height: 100%;
  display: flex;
  flex-direction: column;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 0 8px;
}

.table-title {
  font-size: 18px;
  font-weight: 500;
  color: var(--td-text-color-primary);
}

.table-actions {
  display: flex;
  gap: 12px;
  align-items: center;
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
