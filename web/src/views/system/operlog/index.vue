<template>
  <el-card>
    <el-form inline :model="query" @submit.prevent>
      <el-form-item>
        <el-input v-model="query.username" placeholder="用户名" clearable style="width: 150px" @clear="load" />
      </el-form-item>
      <el-form-item>
        <el-input v-model="query.module" placeholder="模块" clearable style="width: 150px" @clear="load" />
      </el-form-item>
      <el-form-item>
        <el-select v-model="query.status" placeholder="状态" clearable style="width: 120px">
          <el-option label="成功" :value="1" />
          <el-option label="失败" :value="0" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="load">查询</el-button>
        <el-button type="danger" @click="clear">清空日志</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="list" v-loading="loading" border>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="module" label="模块" width="120" />
      <el-table-column prop="action" label="操作" min-width="140" />
      <el-table-column label="请求" min-width="180">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ row.method }}</el-tag>
          <span class="path">{{ row.path }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="ip" label="IP" width="140" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="errMsg" label="错误信息" min-width="160" show-overflow-tooltip />
      <el-table-column label="时间" width="180">
        <template #default="{ row }">{{ formatDateTime(row.createdAt) }}</template>
      </el-table-column>
    </el-table>

    <el-pagination
      class="pager"
      background
      layout="total, sizes, prev, pager, next"
      :total="total"
      v-model:current-page="query.page"
      v-model:page-size="query.size"
      :page-sizes="[10, 20, 50, 100]"
      @current-change="load"
      @size-change="load"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as api from '@/api'
import { formatDateTime } from '@/utils/format'
import type { OperLogItem } from '@/api/types'

const loading = ref(false)
const list = ref<OperLogItem[]>([])
const total = ref(0)
const query = reactive({
  page: 1,
  size: 20,
  username: '',
  module: '',
  status: undefined as number | undefined,
})

async function load() {
  loading.value = true
  try {
    const res = await api.listOperLogs({ ...query })
    list.value = res.list ?? []
    total.value = res.total ?? 0
  } catch (e: any) {
    ElMessage.error(e?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function clear() {
  try {
    await ElMessageBox.confirm('确定清空全部操作日志吗？此操作不可恢复。', '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await api.clearOperLogs()
    ElMessage.success('已清空')
    load()
  } catch (e: any) {
    ElMessage.error(e?.message || '清空失败')
  }
}

onMounted(load)
</script>

<style scoped>
.path {
  margin-left: 6px;
  color: #606266;
}
</style>
