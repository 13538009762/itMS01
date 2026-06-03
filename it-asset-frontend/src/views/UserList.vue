<template>
  <div class="user-list">
    <el-card>
      <template #header>
        <div class="clearfix" style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: bold; font-size: 16px;">用户管理</span>
          <div>
            <el-button type="primary" size="small" @click="fetchUsers">刷新</el-button>
            <el-button type="success" size="small" @click="showAddDialog = true">添加用户</el-button>
          </div>
        </div>
      </template>

      <el-table :data="users" style="width: 100%" v-loading="loading" border stripe>
        <el-table-column prop="ID" label="用户 ID" width="90" align="center"></el-table-column>
        <el-table-column prop="username" label="用户名" min-width="120"></el-table-column>
        <el-table-column prop="real_name" label="真实姓名" min-width="120"></el-table-column>
        <el-table-column prop="role_id" label="系统角色" width="160" align="center">
          <template #default="scope">
            <el-tag :type="getRoleType(scope.row.role_id)">
              {{ getRoleName(scope.row.role_id) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center">
          <template #default="scope">
            <el-button size="small" type="primary" @click="openEditDialog(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加用户弹窗 -->
    <el-dialog v-model="showAddDialog" title="添加用户" width="500px">
      <el-form :model="newUser" label-width="100px">
        <el-form-item label="用户名" required>
          <el-input v-model="newUser.username" placeholder="请输入登录用户名"></el-input>
        </el-form-item>
        <el-form-item label="真实姓名" required>
          <el-input v-model="newUser.real_name" placeholder="请输入员工真实姓名"></el-input>
        </el-form-item>
        <el-form-item label="登录密码" required>
          <el-input v-model="newUser.password" type="password" show-password placeholder="请输入登录密码"></el-input>
        </el-form-item>
        <el-form-item label="系统权限" required>
          <el-select v-model="newUser.role_id" placeholder="请选择系统角色权限" style="width: 100%;">
            <el-option label="普通员工" :value="1"></el-option>
            <el-option label="IT管理员" :value="2"></el-option>
            <el-option label="系统管理员" :value="3"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddDialog = false">取消</el-button>
          <el-button type="primary" @click="addUser">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑用户弹窗 -->
    <el-dialog v-model="showEditDialog" title="编辑用户" width="500px">
      <el-form :model="editingUser" label-width="100px">
        <el-form-item label="用户名" required>
          <el-input v-model="editingUser.username" placeholder="请输入登录用户名"></el-input>
        </el-form-item>
        <el-form-item label="真实姓名" required>
          <el-input v-model="editingUser.real_name" placeholder="请输入员工真实姓名"></el-input>
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="editingUser.password" type="password" show-password placeholder="不修改请留空"></el-input>
        </el-form-item>
        <el-form-item label="系统权限" required>
          <el-select v-model="editingUser.role_id" placeholder="请选择系统角色权限" style="width: 100%;">
            <el-option label="普通员工" :value="1"></el-option>
            <el-option label="IT管理员" :value="2"></el-option>
            <el-option label="系统管理员" :value="3"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="saveEditUser">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../utils/request'
import { ElMessage } from 'element-plus'

const users = ref([])
const loading = ref(false)

const showAddDialog = ref(false)
const newUser = ref({
  username: '',
  real_name: '',
  password: '',
  role_id: 1
})

const showEditDialog = ref(false)
const editingUser = ref({
  id: null,
  username: '',
  real_name: '',
  password: '',
  role_id: 1
})

const fetchUsers = () => {
  loading.value = true
  request.get('/admin/users').then(res => {
    users.value = res || []
  }).catch(err => {
    ElMessage.error('获取用户列表失败: ' + (err.response?.data?.message || err.message))
  }).finally(() => {
    loading.value = false
  })
}

const getRoleName = (roleId) => {
  const map = { 1: '普通员工', 2: 'IT管理员', 3: '系统管理员' }
  return map[roleId] || '未知'
}

const getRoleType = (roleId) => {
  const map = { 1: 'info', 2: 'primary', 3: 'success' }
  return map[roleId] || 'warning'
}

const addUser = () => {
  if (!newUser.value.username || !newUser.value.real_name || !newUser.value.password) {
    ElMessage.warning('请填写所有必填字段')
    return
  }
  request.post('/admin/users', newUser.value).then(() => {
    ElMessage.success('添加用户成功')
    showAddDialog.value = false
    newUser.value = {
      username: '',
      real_name: '',
      password: '',
      role_id: 1
    }
    fetchUsers()
  }).catch(err => {
    ElMessage.error('添加失败: ' + (err.response?.data?.message || err.message))
  })
}

const openEditDialog = (row) => {
  editingUser.value = {
    id: row.ID,
    username: row.username,
    real_name: row.real_name,
    password: '', // 留空
    role_id: row.role_id
  }
  showEditDialog.value = true
}

const saveEditUser = () => {
  if (!editingUser.value.username || !editingUser.value.real_name) {
    ElMessage.warning('请填写用户名与真实姓名')
    return
  }
  request.put('/admin/users', editingUser.value).then(() => {
    ElMessage.success('修改用户成功')
    showEditDialog.value = false
    fetchUsers()
  }).catch(err => {
    ElMessage.error('修改失败: ' + (err.response?.data?.message || err.message))
  })
}

const handleDelete = (row) => {
  import('element-plus').then(({ ElMessageBox }) => {
    ElMessageBox.confirm(`确定要删除用户 "${row.username}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      request.delete('/admin/users', { params: { id: row.ID } }).then(() => {
        ElMessage.success('删除成功')
        fetchUsers()
      }).catch(err => {
        ElMessage.error('删除失败: ' + (err.response?.data?.message || err.message))
      })
    }).catch(() => {})
  })
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.clearfix {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
