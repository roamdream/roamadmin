<template>
  <el-card>
    <el-form inline :model="query" @submit.prevent>
      <el-form-item>
        <el-input v-model="query.username" placeholder="用户名" clearable style="width: 160px" @clear="load" />
      </el-form-item>
      <el-form-item>
        <el-input v-model="query.nickname" placeholder="昵称" clearable style="width: 160px" @clear="load" />
      </el-form-item>
      <el-form-item>
        <el-select v-model="query.status" placeholder="状态" clearable style="width: 120px">
          <el-option label="启用" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="load">查询</el-button>
        <el-button type="success" @click="openCreate">新增</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="list" v-loading="loading" border>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="username" label="用户名" min-width="120" />
      <el-table-column prop="nickname" label="昵称" min-width="120" />
      <el-table-column prop="email" label="邮箱" min-width="160" />
      <el-table-column prop="phone" label="手机" min-width="130" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">{{ formatDateTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="300" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="warning" @click="openResetPwd(row)">重置密码</el-button>
          <el-button link :type="row.status === 1 ? 'danger' : 'success'" @click="toggleStatus(row)">
            {{ row.status === 1 ? '禁用' : '启用' }}
          </el-button>
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

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑用户' : '新增用户'" width="520px" @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="!!editingId" />
        </el-form-item>
        <el-form-item v-if="!editingId" label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="手机">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.roleIds" multiple placeholder="选择角色" style="width: 100%">
            <el-option v-for="r in roleOptions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
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
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as api from '@/api'
import { formatDateTime } from '@/utils/format'
import type { UserItem, RoleItem } from '@/api/types'

const loading = ref(false)
const list = ref<UserItem[]>([])
const total = ref(0)
const query = reactive({ page: 1, size: 20, username: '', nickname: '', status: undefined as number | undefined })

const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref()
const roleOptions = ref<RoleItem[]>([])

const emptyForm = () => ({
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  avatar: '',
  status: 1,
  remark: '',
  roleIds: [] as number[],
})
const form = reactive(emptyForm())

const formRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, min: 6, message: '密码至少 6 位', trigger: 'blur' }],
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
}

async function load() {
  loading.value = true
  try {
    const res = await api.listUsers({ ...query })
    list.value = res.list ?? []
    total.value = res.total ?? 0
  } catch (e: any) {
    ElMessage.error(e?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function loadRoles() {
  try {
    roleOptions.value = (await api.listRoles({ page: 1, size: 200 })).list ?? []
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
  dialogVisible.value = true
}

async function openEdit(row: UserItem) {
  editingId.value = row.id
  const res = await api.getUser(row.id)
  Object.assign(form, {
    username: res.user?.username,
    nickname: res.user?.nickname,
    email: res.user?.email,
    phone: res.user?.phone,
    avatar: res.user?.avatar,
    status: res.user?.status,
    remark: res.user?.remark,
    roleIds: res.roleIds ?? [],
    password: '',
  })
  dialogVisible.value = true
}

async function save() {
  await formRef.value.validate().catch(() => Promise.reject())
  saving.value = true
  try {
    if (editingId.value) {
      await api.updateUser(editingId.value, {
        username: form.username,
        nickname: form.nickname,
        email: form.email,
        phone: form.phone,
        avatar: form.avatar,
        status: form.status,
        remark: form.remark,
        roleIds: form.roleIds,
      })
      ElMessage.success('已更新')
    } else {
      await api.createUser({ ...form })
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

async function remove(row: UserItem) {
  try {
    await ElMessageBox.confirm(`确定删除用户「${row.username}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await api.deleteUser(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e: any) {
    ElMessage.error(e?.message || '删除失败')
  }
}

async function toggleStatus(row: UserItem) {
  const next = row.status === 1 ? 0 : 1
  try {
    await api.setUserStatus(row.id, next)
    ElMessage.success('已更新')
    load()
  } catch (e: any) {
    ElMessage.error(e?.message || '操作失败')
  }
}

async function openResetPwd(row: UserItem) {
  let value: string
  try {
    const r = await ElMessageBox.prompt(`为用户「${row.username}」设置新密码`, '重置密码', {
      inputType: 'password',
      inputValidator: (v) => (v && v.length >= 6 ? true : '密码至少 6 位'),
    })
    value = r.value
  } catch {
    return
  }
  try {
    await api.resetUserPassword(row.id, value)
    ElMessage.success('已重置')
  } catch (e: any) {
    ElMessage.error(e?.message || '重置失败')
  }
}

onMounted(() => {
  load()
  loadRoles()
})
</script>
