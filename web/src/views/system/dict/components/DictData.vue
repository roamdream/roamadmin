<template>
  <div>
    <el-form inline :model="query" @submit.prevent>
      <el-form-item>
        <el-select v-model="query.dictType" placeholder="选择字典类型" clearable filterable style="width: 200px">
          <el-option v-for="t in typeOptions" :key="t.type" :label="t.name" :value="t.type" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-input v-model="query.label" placeholder="数据标签" clearable style="width: 160px" @clear="load" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="load">查询</el-button>
        <el-button type="success" @click="openCreate">新增</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="list" v-loading="loading" border>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="dictType" label="字典类型" min-width="130" />
      <el-table-column prop="label" label="标签" min-width="130" />
      <el-table-column prop="value" label="键值" min-width="120" />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="140" />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
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

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑字典数据' : '新增字典数据'" width="460px" @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="80px">
        <el-form-item label="字典类型" prop="dictType">
          <el-select v-model="form.dictType" filterable style="width: 100%">
            <el-option v-for="t in typeOptions" :key="t.type" :label="`${t.name}（${t.type}）`" :value="t.type" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签" prop="label">
          <el-input v-model="form.label" />
        </el-form-item>
        <el-form-item label="键值" prop="value">
          <el-input v-model="form.value" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as api from '@/api'
import type { DictDataItem, DictTypeItem } from '@/api/types'

const loading = ref(false)
const list = ref<DictDataItem[]>([])
const total = ref(0)
const query = reactive({ page: 1, size: 20, dictType: '', label: '' })
const typeOptions = ref<DictTypeItem[]>([])

const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref()

const emptyForm = () => ({ dictType: '', label: '', value: '', sort: 0, status: 1, remark: '' })
const form = reactive(emptyForm())

const formRules = {
  dictType: [{ required: true, message: '请选择字典类型', trigger: 'change' }],
  label: [{ required: true, message: '请输入标签', trigger: 'blur' }],
  value: [{ required: true, message: '请输入键值', trigger: 'blur' }],
}

async function load() {
  loading.value = true
  try {
    const res = await api.listDictData({ ...query })
    list.value = res.list ?? []
    total.value = res.total ?? 0
  } catch (e: any) {
    ElMessage.error(e?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function loadTypes() {
  try {
    typeOptions.value = (await api.listDictTypes({ page: 1, size: 200 })).list ?? []
  } catch {
    // 忽略
  }
}

function resetForm() {
  Object.assign(form, emptyForm())
}

function openCreate() {
  editingId.value = null
  resetForm()
  if (query.dictType) form.dictType = query.dictType
  dialogVisible.value = true
}

async function openEdit(row: DictDataItem) {
  editingId.value = row.id
  Object.assign(form, {
    dictType: row.dictType,
    label: row.label,
    value: row.value,
    sort: row.sort,
    status: row.status,
    remark: row.remark,
  })
  dialogVisible.value = true
}

async function save() {
  await formRef.value.validate().catch(() => Promise.reject())
  saving.value = true
  try {
    const payload = {
      dictType: form.dictType,
      label: form.label,
      value: form.value,
      sort: form.sort,
      status: form.status,
      remark: form.remark,
    }
    if (editingId.value) {
      await api.updateDictData(editingId.value, payload)
      ElMessage.success('已更新')
    } else {
      await api.createDictData(payload)
      ElMessage.success('已创建')
    }
    dialogVisible.value = false
    load()
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

async function remove(row: DictDataItem) {
  try {
    await ElMessageBox.confirm(`确定删除字典数据「${row.label}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await api.deleteDictData(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e: any) {
    ElMessage.error(e?.message || '删除失败')
  }
}

onMounted(() => {
  load()
  loadTypes()
})
</script>
