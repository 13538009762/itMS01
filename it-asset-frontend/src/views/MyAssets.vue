<template>
  <div class="my-assets">
    <!-- 我持有的资产 -->
    <el-card style="margin-bottom: 20px;">
      <template #header>
        <div class="clearfix" style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: bold; font-size: 16px;">我持有的资产</span>
          <el-button type="primary" size="small" @click="fetchHeldAssets">刷新</el-button>
        </div>
      </template>

      <el-table :data="heldAssets" style="width: 100%" v-loading="heldLoading" border stripe>
        <el-table-column prop="asset_no" label="资产编号" width="180"></el-table-column>
        <el-table-column prop="name" label="资产名称" min-width="150"></el-table-column>
        <el-table-column prop="holding_status" label="当前状态" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getHoldingStatusType(scope.row.holding_status)">
              {{ scope.row.holding_status || '使用中' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" align="center">
          <template #default="scope">
            <el-button 
              size="small" 
              type="warning" 
              :disabled="scope.row.holding_status !== '使用中'"
              @click="handleRepair(scope.row)"
            >一键报修</el-button>
            <el-button 
              size="small" 
              type="success" 
              :disabled="scope.row.holding_status !== '使用中'"
              @click="handleReturn(scope.row)"
            >申请归还</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 我的申请记录 -->
    <el-card>
      <template #header>
        <div class="clearfix" style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: bold; font-size: 16px;">我的申请历史记录</span>
          <el-button type="primary" size="small" @click="fetchApplyLogs">刷新</el-button>
        </div>
      </template>

      <el-table :data="displayedApplyLogs" style="width: 100%" v-loading="logsLoading" border stripe>
        <el-table-column prop="ID" label="申请 ID" width="95" align="center"></el-table-column>
        <el-table-column label="资产编号" width="160">
          <template #default="scope">
            {{ scope.row.asset?.asset_no || '资产 ID: ' + scope.row.asset_id }}
          </template>
        </el-table-column>
        <el-table-column label="资产名称" min-width="150">
          <template #default="scope">
            {{ scope.row.asset?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="申请类型" width="110" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.type === 1 ? 'success' : 'info'">
              {{ scope.row.type === 1 ? '领用' : '归还' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="申请原因" min-width="120">
          <template #default="scope">
            {{ scope.row.reason || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="reject_reason" label="驳回原因/审批意见" min-width="150">
          <template #default="scope">
            {{ scope.row.reject_reason || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getApplyStatusType(scope.row.status)">
              {{ getApplyStatusName(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="申请时间" width="180" align="center">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 15px; display: flex; justify-content: flex-end;">
        <el-pagination
          v-model:current-page="applyCurrentPage"
          v-model:page-size="applyPageSize"
          :total="applyLogs.length"
          layout="total, prev, pager, next"
          size="small"
          background
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import request from '../utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))
const heldAssets = ref([])
const heldLoading = ref(false)

const applyLogs = ref([])
const logsLoading = ref(false)

// 申请历史记录分页设置 (最多显示15行否则翻页)
const applyCurrentPage = ref(1)
const applyPageSize = ref(15)

const displayedApplyLogs = computed(() => {
  const start = (applyCurrentPage.value - 1) * applyPageSize.value
  const end = start + applyPageSize.value
  return applyLogs.value.slice(start, end)
})

const fetchHeldAssets = () => {
  heldLoading.value = true
  request.get('/employee/assets/held').then(res => {
    heldAssets.value = res || []
  }).catch(err => {
    ElMessage.error('获取我持有的资产失败')
  }).finally(() => {
    heldLoading.value = false
  })
}

const fetchApplyLogs = () => {
  logsLoading.value = true
  request.get('/employee/apply/my').then(res => {
    applyLogs.value = res || []
    applyCurrentPage.value = 1
  }).catch(err => {
    ElMessage.error('获取申请记录失败')
  }).finally(() => {
    logsLoading.value = false
  })
}

const getApplyStatusName = (status) => {
  const map = { 0: '待审批', 1: '已通过', 2: '已驳回' }
  return map[status] || '未知'
}

const getApplyStatusType = (status) => {
  const map = { 0: 'warning', 1: 'success', 2: 'danger' }
  return map[status] || 'info'
}

const getHoldingStatusType = (status) => {
  const map = {
    '使用中': 'primary',
    '待归还': 'warning',
    '待报修': 'danger',
    '维修中': 'info'
  }
  return map[status] || 'info'
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', { hour12: false })
}

const handleReturn = (row) => {
  request.post('/employee/apply', {
    asset_id: row.ID,
    type: 2
  }).then(() => {
    ElMessage.success('归还申请已提交，等待管理员审批')
    fetchHeldAssets()
    fetchApplyLogs()
  }).catch(err => {
    ElMessage.error('提交归还申请失败: ' + (err.response?.data?.message || err.message))
  })
}

const handleRepair = (row) => {
  ElMessageBox.prompt('请输入报修原因', '资产报修', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /\S+/,
    inputErrorMessage: '报修原因不能为空'
  }).then(({ value }) => {
    request.post('/employee/repair', {
      asset_id: row.ID,
      reason: value
    }).then(() => {
      ElMessage.success('报修申请已提交，请将设备送至维修点')
      fetchHeldAssets()
      fetchApplyLogs()
    }).catch(err => {
      ElMessage.error('提交报修申请失败: ' + (err.response?.data?.message || err.message))
    })
  }).catch(() => {
    // 取消输入
  })
}

onMounted(() => {
  fetchHeldAssets()
  fetchApplyLogs()
})
</script>

<style scoped>
.clearfix {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
