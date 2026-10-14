<template>
  <el-card>
    <div class="toolbar">
      <el-button type="success" @click="openCreate()">新增根菜单</el-button>
    </div>

    <el-table :data="list" v-loading="loading" row-key="id" border default-expand-all>
      <el-table-column prop="name" label="菜单名称" min-width="180" />
      <el-table-column label="类型" width="90">
        <template #default="{ row }">
          <el-tag :type="typeTag(row.type)">{{ typeLabel(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="path" label="路由" min-width="140" />
      <el-table-column prop="component" label="组件" min-width="140" />
      <el-table-column prop="perm" label="权限码" min-width="150" />
      <el-table-column prop="icon" label="图标" width="100" />
      <el-table-column prop="sort" label="排序" width="70" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openCreate(row.id)">新增子级</el-button>
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑菜单' : '新增菜单'" width="520px" @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="80px">
        <el-form-item label="上级菜单">
          <el-tree-select
            v-model="form.parentId"
            :data="parentOptions"
            :props="{ label: 'name', children: 'children' }"
            node-key="id"
            check-strictly
            clearable
            placeholder="根菜单"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio-button value="dir">目录</el-radio-button>
            <el-radio-button value="menu">菜单</el-radio-button>
            <el-radio-button value="button">按钮</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item v-if="form.type !== 'button'" label="路由">
          <el-input v-model="form.path" placeholder="如 /system/user" />
        </el-form-item>
        <el-form-item v-if="form.type === 'menu'" label="组件">
          <el-input v-model="form.component" placeholder="如 system/user" />
        </el-form-item>
        <el-form-item label="权限码">
          <el-input v-model="form.perm" placeholder="如 system:user:list" />
        </el-form-item>
        <el-form-item v-if="form.type !== 'button'" label="图标">
          <el-input v-model="form.icon" placeholder="如 user" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as api from '@/api'
import type { MenuTreeItem } from '@/api/types'

const loading = ref(false)
const list = ref<MenuTreeItem[]>([])
const parentOptions = ref<MenuTreeItem[]>([])

const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref()

const emptyForm = () => ({
  parentId: 0,
  name: '',
  type: 'menu' as 'dir' | 'menu' | 'button',
  perm: '',
  path: '',
  component: '',
  icon: '',
  sort: 0,
  status: 1,
})
const form = reactive(emptyForm())

const formRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
}

const typeLabel = (t: string) => ({ dir: '目录', menu: '菜单', button: '按钮' }[t] ?? t)
const typeTag = (t: string) => ({ dir: 'info', menu: 'primary', button: 'warning' }[t] ?? 'info')

async function load() {
  loading.value = true
  try {
    const res = await api.listMenus()
    list.value = res.list ?? []
    parentOptions.value = [{ id: 0, name: '根菜单', children: list.value } as MenuTreeItem]
  } catch (e: any) {
    ElMessage.error(e?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

function resetForm() {
  Object.assign(form, emptyForm())
}

function openCreate(parentId = 0) {
  editingId.value = null
  resetForm()
  form.parentId = parentId
  dialogVisible.value = true
}

async function openEdit(row: MenuTreeItem) {
  editingId.value = row.id
  const res = await api.getMenu(row.id)
  Object.assign(form, {
    parentId: res.menu?.parentId ?? 0,
    name: res.menu?.name,
    type: (res.menu?.type ?? 'menu') as 'dir' | 'menu' | 'button',
    perm: res.menu?.perm,
    path: res.menu?.path,
    component: res.menu?.component,
    icon: res.menu?.icon,
    sort: res.menu?.sort,
    status: res.menu?.status,
  })
  dialogVisible.value = true
}

async function save() {
  await formRef.value.validate().catch(() => Promise.reject())
  saving.value = true
  try {
    const payload = {
      parentId: form.parentId,
      name: form.name,
      type: form.type,
      perm: form.perm,
      path: form.path,
      component: form.component,
      icon: form.icon,
      sort: form.sort,
      status: form.status,
    }
    if (editingId.value) {
      await api.updateMenu(editingId.value, payload)
      ElMessage.success('已更新')
    } else {
      await api.createMenu(payload)
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

async function remove(row: MenuTreeItem) {
  try {
    await ElMessageBox.confirm(`确定删除菜单「${row.name}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await api.deleteMenu(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e: any) {
    ElMessage.error(e?.message || '删除失败')
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar {
  margin-bottom: 12px;
}
</style>
