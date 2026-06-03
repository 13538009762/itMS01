<template>
  <div class="repair-list">
    <el-card>
      <template #header>
        <div class="clearfix" style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: bold; font-size: 16px;">确认送修处理 (仅管理员可用)</span>
          <el-button type="primary" size="small" @click="fetchRepairLogs">刷新</el-button>
        </div>
      </template>
      <el-table :data="displayedRepairLogs" style="width: 100%" v-loading="repairLoading" border stripe>
        <el-table-column prop="ID" label="报修单 ID" width="100" align="center"></el-table-column>
        <el-table-column label="报修人" width="120">
          <template #default="scope">
            {{ scope.row.user?.real_name || scope.row.user?.username || '用户 ID: ' + scope.row.user_id }}
          </template>
        </el-table-column>
        <el-table-column label="资产编号" width="160">
          <template #default="scope">
            {{ scope.row.asset?.asset_no || '资产 ID: ' + scope.row.asset_id }}
          </template>
        </el-table-column>
        <el-table-column label="资产名称" width="180">
          <template #default="scope">
            {{ scope.row.asset?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="报修原因"></el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="getRepairStatusType(scope.row.status)">
              {{ getRepairStatusName(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="报修时间" width="180" align="center">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="130" align="center">
          <template #default="scope">
            <el-button 
              v-if="scope.row.status === 0" 
              size="small" 
              type="primary" 
              @click="handleConfirm(scope.row)"
            >确认送修</el-button>
            <el-button 
              v-else-if="scope.row.status === 1" 
              size="small" 
              type="success" 
              @click="handleComplete(scope.row)"
            >完成维修</el-button>
            <span v-else style="color: #67C23A; font-size: 13px; font-weight: bold;">已完成</span>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 15px; display: flex; justify-content: flex-end;">
        <el-pagination
          v-model:current-page="repairCurrentPage"
          v-model:page-size="repairPageSize"
          :total="repairLogs.length"
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
import { ElMessage } from 'element-plus'

const repairLogs = ref([])
const repairLoading = ref(false)

// 分页状态 (一个界面最多显示15条数据)
const repairCurrentPage = ref(1)
const repairPageSize = ref(15)

// 分页截取计算属性
const displayedRepairLogs = computed(() => {
  const start = (repairCurrentPage.value - 1) * repairPageSize.value
  const end = start + repairPageSize.value
  return repairLogs.value.slice(start, end)
})

const fetchRepairLogs = () => {
  repairLoading.value = true
  request.get('/admin/repair/list').then(res => {
    repairLogs.value = res || []
    repairCurrentPage.value = 1
  }).catch(() => {
    ElMessage.error('获取报修申请列表失败')
  }).finally(() => {
    repairLoading.value = false
  })
}

const getRepairStatusName = (status) => {
  const map = { 0: '待处理', 1: '已送修', 2: '已完成' }
  return map[status] || '未知'
}

const getRepairStatusType = (status) => {
  const map = { 0: 'warning', 1: 'primary', 2: 'success' }
  return map[status] || 'info'
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', { hour12: false })
}

const handleConfirm = (row) => {
  request.put('/admin/repair/confirm', {
    repair_id: row.ID,
    asset_id: row.asset_id
  }).then(() => {
    ElMessage.success('送修处理成功')
    fetchRepairLogs()
  }).catch(err => {
    ElMessage.error('确认送修失败: ' + (err.response?.data?.message || err.message))
  })
}

const handleComplete = (row) => {
  request.put('/admin/repair/complete', {
    repair_id: row.ID,
    asset_id: row.asset_id
  }).then(() => {
    ElMessage.success('完成维修处理成功')
    fetchRepairLogs()
  }).catch(err => {
    ElMessage.error('完成维修失败: ' + (err.response?.data?.message || err.message))
  })
}

onMounted(() => {
  fetchRepairLogs()
})
</script>

<style scoped>
.repair-list {
  padding: 0;
}
</style>
