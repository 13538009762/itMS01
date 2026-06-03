<template>
  <div class="login-container">
    <!-- Top Header -->
    <header class="login-header">
      <div class="logo-area">
        <el-icon class="logo-icon" :size="24"><Monitor /></el-icon>
        <span class="system-name">itMS</span>
        <span class="system-desc">| IT 资产管理系统</span>
      </div>
    </header>

    <!-- Main Content Grid -->
    <main class="login-main">
      <!-- Left Hero Column -->
      <div class="hero-section">
        <h1 class="hero-title">智能资产 · 高效流转</h1>
        <p class="hero-subtitle">
          基于全生命周期的企业级 IT 资产管理平台，协助团队高效完成资产领用审批、精确把握报修周期，提升硬件资源周转效能。
        </p>

        <!-- Feature Grid -->
        <div class="feature-grid">
          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <el-icon><DataAnalysis /></el-icon>
            </div>
            <div class="feature-info">
              <h3>智能资产统计</h3>
              <p>实时统计、报修流转，让资产状态流向一目了然</p>
            </div>
          </div>
          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <el-icon><Key /></el-icon>
            </div>
            <div class="feature-info">
              <h3>权限安全可信</h3>
              <p>角色校验、Casbin 精细控制，保障系统数据防越权安全</p>
            </div>
          </div>
          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <el-icon><Connection /></el-icon>
            </div>
            <div class="feature-info">
              <h3>高效申请审批</h3>
              <p>一键申请领用及自动归还报修审批流，加速办公效能</p>
            </div>
          </div>
          <div class="feature-card">
            <div class="feature-icon-wrapper">
              <el-icon><Monitor /></el-icon>
            </div>
            <div class="feature-info">
              <h3>多端响应支持</h3>
              <p>全站响应式布局与美学重塑，完美适配各类浏览终端</p>
            </div>
          </div>
        </div>

        <!-- Bottom Knowledge Card -->
        <div class="knowledge-card">
          <div class="knowledge-bar"></div>
          <h4>企业级核心资产库</h4>
          <p>
            本系统致力于为企业提供一站式、智能化的 IT 资产全生命周期管理方案。通过打通“入库 - 领用 - 归还 - 报修 - 审批”的闭环流转，显著提升团队协同效率。在这里，每一次流转都可追溯，每一件资产都物尽其用。
          </p>
        </div>
      </div>

      <!-- Right Login Card -->
      <div class="login-card-wrapper">
        <div class="login-card">
          <h2 class="card-title">欢迎回来</h2>
          <p class="card-subtitle">登录您的 itMS 账户</p>
          
          <div class="login-tab">账号登录</div>

          <el-form :model="form" :rules="rules" ref="loginForm" label-position="top">
            <el-form-item prop="username">
              <el-input 
                v-model="form.username" 
                placeholder="请输入用户名" 
                prefix-icon="User"
                size="large"
              ></el-input>
            </el-form-item>
            <el-form-item prop="password">
              <el-input 
                v-model="form.password" 
                type="password" 
                placeholder="请输入密码" 
                prefix-icon="Lock"
                show-password
                size="large"
                @keyup.enter="handleLogin"
              ></el-input>
            </el-form-item>

            <div class="login-options">
              <el-checkbox v-model="rememberMe">记住我</el-checkbox>
              <span class="forgot-pwd">忘记密码？</span>
            </div>

            <el-form-item>
              <el-button 
                type="primary" 
                class="login-btn" 
                @click="handleLogin" 
                :loading="loading"
                size="large"
              >
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="login-footer">
      <div class="copyright">
        © 2026 itMS. All rights reserved.
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import request from '../utils/request'
import { ElMessage } from 'element-plus'

const router = useRouter()
const loginForm = ref(null)
const loading = ref(false)
const rememberMe = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = () => {
  loginForm.value.validate(valid => {
    if (valid) {
      loading.value = true
      request.post('/login', form).then(res => {
        localStorage.setItem('token', res.token)
        localStorage.setItem('user', JSON.stringify(res.user))
        ElMessage.success('登录成功')
        router.push('/')
      }).finally(() => {
        loading.value = false
      })
    }
  })
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  width: 100vw;
  position: relative;
  overflow: hidden;
  background: url('../assets/background.png') no-repeat center center;
  background-size: cover;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Arial', sans-serif;
  color: #333;
}

/* Header */
.login-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 60px;
  z-index: 10;
}
.logo-area {
  display: flex;
  align-items: center;
  gap: 8px;
}
.logo-icon {
  color: #1890ff;
  font-size: 26px;
}
.system-name {
  font-size: 24px;
  font-weight: 800;
  color: #0050b3;
  letter-spacing: 0.5px;
}
.system-desc {
  font-size: 14px;
  color: #666;
}

/* Footer */
.login-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 60px;
  font-size: 12px;
  color: #888;
  z-index: 10;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(5px);
}

/* Main */
.login-main {
  display: flex;
  flex: 1;
  padding: 0 80px;
  align-items: center;
  justify-content: space-between;
  z-index: 10;
}

/* Hero section (left) */
.hero-section {
  flex: 1;
  max-width: 650px;
  padding-right: 40px;
}
.hero-title {
  font-size: 40px;
  font-weight: 800;
  line-height: 1.3;
  margin-bottom: 16px;
  background: linear-gradient(135deg, #001529 30%, #1890ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
.hero-subtitle {
  font-size: 15px;
  color: #555;
  line-height: 1.6;
  margin-bottom: 40px;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  margin-bottom: 40px;
}
.feature-card {
  display: flex;
  gap: 16px;
  background: rgba(255, 255, 255, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.02);
  transition: all 0.3s;
}
.feature-card:hover {
  transform: translateY(-2px);
  background: rgba(255, 255, 255, 0.8);
  box-shadow: 0 8px 24px rgba(24, 144, 255, 0.1);
}
.feature-icon-wrapper {
  width: 44px;
  height: 44px;
  background: rgba(24, 144, 255, 0.1);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #1890ff;
  font-size: 20px;
}
.feature-info h3 {
  font-size: 15px;
  font-weight: bold;
  margin-bottom: 4px;
  color: #1a1a1a;
}
.feature-info p {
  font-size: 12px;
  color: #666;
  line-height: 1.4;
}

.knowledge-card {
  position: relative;
  background: rgba(255, 255, 255, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.4);
  border-radius: 16px;
  padding: 20px 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.03);
}
.knowledge-bar {
  position: absolute;
  left: 0;
  top: 24px;
  width: 4px;
  height: 16px;
  background: #1890ff;
  border-radius: 0 2px 2px 0;
}
.knowledge-card h4 {
  font-size: 15px;
  font-weight: bold;
  margin-bottom: 8px;
  color: #001529;
}
.knowledge-card p {
  font-size: 13px;
  color: #555;
  line-height: 1.6;
}

/* Login card (right) */
.login-card-wrapper {
  width: 420px;
}
.login-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(15px);
  -webkit-backdrop-filter: blur(15px);
  border: 1px solid rgba(255, 255, 255, 0.6);
  border-radius: 24px;
  padding: 40px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.08);
}
.card-title {
  font-size: 28px;
  font-weight: bold;
  color: #1a1a1a;
  margin-bottom: 6px;
}
.card-subtitle {
  font-size: 14px;
  color: #8c8c8c;
  margin-bottom: 24px;
}
.login-tab {
  font-size: 15px;
  font-weight: bold;
  color: #1890ff;
  border-bottom: 2px solid #1890ff;
  width: max-content;
  padding-bottom: 6px;
  margin-bottom: 24px;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.forgot-pwd {
  font-size: 14px;
  color: #1890ff;
  cursor: pointer;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
  background: linear-gradient(135deg, #1890ff 0%, #0050b3 100%);
  border: none;
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
  transition: all 0.3s;
}
.login-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(24, 144, 255, 0.4);
}

/* Custom inputs styling */
:deep(.el-input__wrapper) {
  background-color: #f1f5fa !important;
  border: none !important;
  box-shadow: none !important;
  border-radius: 8px;
  padding: 8px 12px;
}
:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #1890ff !important;
}
:deep(.el-form-item) {
  margin-bottom: 20px;
}

@media (max-width: 992px) {
  .login-main {
    flex-direction: column;
    justify-content: center;
    padding: 0 20px;
  }
  .hero-section {
    display: none;
  }
  .login-header {
    padding: 16px 24px;
  }
  .login-footer {
    padding: 16px 24px;
    flex-direction: column;
    gap: 10px;
  }
}
</style>

