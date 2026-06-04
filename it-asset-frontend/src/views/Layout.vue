<template>
  <el-container class="layout-container">
    <el-header class="layout-header">
      <div class="logo">IT 资产管理系统</div>
      <div class="user-info">
        <span>{{ user.real_name }} ({{ roleName }})</span>
        <el-button type="text" @click="handleLogout">退出登录</el-button>
      </div>
    </el-header>
    <el-container>
      <el-aside width="200px" class="layout-aside">
        <el-menu :default-active="$route.path" router>
          <el-menu-item index="/assets">
            <el-icon><Monitor /></el-icon>
            <span>资产列表</span>
          </el-menu-item>
          <el-menu-item index="/my-assets">
            <el-icon><List /></el-icon>
            <span>我的申请与资产</span>
          </el-menu-item>
          <el-menu-item v-if="user.role_id === 2 || user.role_id === 3" index="/applies">
            <el-icon><Stamp /></el-icon>
            <span>资产审批</span>
          </el-menu-item>
          <el-menu-item v-if="user.role_id === 2 || user.role_id === 3" index="/repairs">
            <el-icon><Tool /></el-icon>
            <span>资产报修</span>
          </el-menu-item>
          <el-menu-item v-if="user.role_id === 3" index="/users">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/profile">
            <el-icon><Avatar /></el-icon>
            <span>我的</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main class="layout-main">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))

const roleName = computed(() => {
  if (user.value.role_id === 1) return '普通员工'
  if (user.value.role_id === 2) return 'IT管理员'
  if (user.value.role_id === 3) return '系统管理员'
  return '未知角色'
})

const handleUserUpdate = () => {
  user.value = JSON.parse(localStorage.getItem('user') || '{}')
}

onMounted(() => {
  window.addEventListener('user-update', handleUserUpdate)
})

onUnmounted(() => {
  window.removeEventListener('user-update', handleUserUpdate)
})

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}
.layout-header {
  background-color: #409EFF;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}
.logo {
  font-size: 20px;
  font-weight: bold;
}
.user-info span {
  margin-right: 15px;
}
.user-info .el-button {
  color: white;
}
.layout-aside {
  border-right: solid 1px #e6e6e6;
}
.layout-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
