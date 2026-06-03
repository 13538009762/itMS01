<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <!-- 个人信息卡片 -->
      <el-col :xs="24" :sm="12">
        <el-card class="profile-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon"><User /></el-icon>
              <span>个人基本信息</span>
            </div>
          </template>

          <el-form :model="profileForm" label-position="top" class="form-style">
            <el-form-item label="用户名 (登录账号)" required>
              <el-input v-model="profileForm.username" placeholder="请输入用户名" prefix-icon="User" />
            </el-form-item>

            <el-form-item label="真实姓名" required>
              <el-input v-model="profileForm.real_name" placeholder="请输入真实姓名" prefix-icon="Edit" />
            </el-form-item>

            <el-form-item label="当前系统角色">
              <el-input :value="roleName" disabled prefix-icon="Key" />
              <div class="role-tip">
                <el-icon><InfoFilled /></el-icon>
                <span>系统角色权限由管理员分配，用户无法自行修改。</span>
              </div>
            </el-form-item>

            <el-form-item style="margin-top: 30px;">
              <el-button type="primary" :loading="profileSaving" @click="saveProfile" class="submit-btn">
                保存基本信息
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 修改密码卡片 -->
      <el-col :xs="24" :sm="12">
        <el-card class="profile-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon"><Lock /></el-icon>
              <span>修改登录密码</span>
            </div>
          </template>

          <el-form :model="passwordForm" label-position="top" class="form-style">
            <el-form-item label="原密码" required>
              <el-input v-model="passwordForm.old_password" type="password" show-password placeholder="请输入当前旧密码" prefix-icon="Lock" />
            </el-form-item>

            <el-form-item label="新密码" required>
              <el-input v-model="passwordForm.new_password" type="password" show-password placeholder="请输入新的登录密码" prefix-icon="Unlock" />
            </el-form-item>

            <el-form-item label="确认新密码" required>
              <el-input v-model="passwordForm.confirm_password" type="password" show-password placeholder="请再次输入新密码" prefix-icon="CircleCheck" />
            </el-form-item>

            <el-form-item style="margin-top: 30px;">
              <el-button type="success" :loading="passwordSaving" @click="savePassword" class="submit-btn">
                更新登录密码
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import request from '../utils/request'
import { ElMessage } from 'element-plus'

const currentUser = ref(JSON.parse(localStorage.getItem('user') || '{}'))

const profileForm = ref({
  username: currentUser.value.username || '',
  real_name: currentUser.value.real_name || ''
})

const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const profileSaving = ref(false)
const passwordSaving = ref(false)

const roleName = computed(() => {
  const map = { 1: '普通员工', 2: 'IT管理员', 3: '系统管理员' }
  return map[currentUser.value.role_id] || '未知角色'
})

const saveProfile = () => {
  if (!profileForm.value.username || !profileForm.value.real_name) {
    ElMessage.warning('用户名和真实姓名不能为空')
    return
  }

  profileSaving.value = true
  request.put('/employee/profile', {
    username: profileForm.value.username,
    real_name: profileForm.value.real_name
  }).then(res => {
    ElMessage.success('个人基本信息保存成功')
    
    // 更新本地存储的缓存
    currentUser.value.username = res.username
    currentUser.value.real_name = res.real_name
    localStorage.setItem('user', JSON.stringify(currentUser.value))
    
    // 派发自定义事件，通知 Layout.vue 更新头部显示
    window.dispatchEvent(new Event('user-update'))
  }).catch(err => {
    console.error(err)
  }).finally(() => {
    profileSaving.value = false
  })
}

const savePassword = () => {
  if (!passwordForm.value.old_password || !passwordForm.value.new_password || !passwordForm.value.confirm_password) {
    ElMessage.warning('请填写完整的密码信息')
    return
  }

  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    ElMessage.warning('新密码与确认密码不一致')
    return
  }

  passwordSaving.value = true
  request.put('/employee/profile/password', {
    old_password: passwordForm.value.old_password,
    new_password: passwordForm.value.new_password
  }).then(() => {
    ElMessage.success('密码修改成功')
    
    // 重置表单
    passwordForm.value = {
      old_password: '',
      new_password: '',
      confirm_password: ''
    }
  }).catch(err => {
    console.error(err)
  }).finally(() => {
    passwordSaving.value = false
  })
}
</script>

<style scoped>
.profile-container {
  padding: 10px;
}

.profile-card {
  margin-bottom: 20px;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.profile-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 20px -8px rgba(0, 0, 0, 0.15) !important;
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: bold;
  color: #303133;
}

.header-icon {
  margin-right: 8px;
  font-size: 18px;
  color: #409eff;
}

.form-style {
  padding: 10px 5px;
}

.role-tip {
  display: flex;
  align-items: center;
  margin-top: 6px;
  font-size: 12px;
  color: #909399;
}

.role-tip .el-icon {
  margin-right: 4px;
}

.submit-btn {
  width: 100%;
  padding: 12px 20px;
  font-size: 14px;
  border-radius: 4px;
  font-weight: 500;
  letter-spacing: 0.5px;
  transition: all 0.2s;
}

.submit-btn:hover {
  filter: brightness(1.05);
}
</style>
