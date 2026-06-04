<template>
  <div class="apply-list">
    <el-card>
      <template #header>
        <div class="clearfix" style="display: flex; justify-content: space-between; align-items: center;">
          <span style="font-weight: bold; font-size: 16px;">领用/归还审批 (仅管理员可用)</span>
          <el-button type="primary" size="small" @click="fetchApplyLogs">刷新</el-button>
        </div>
      </template>
      <el-table :data="displayedApplyLogs" style="width: 100%" v-loading="applyLoading" border stripe>
        <el-table-column prop="ID" label="申请单 ID" width="100" align="center"></el-table-column>
        <el-table-column label="申请人" width="120">
          <template #default="scope">
            {{ scope.row.user?.real_name || scope.row.user?.username || '用户 ID: ' + scope.row.user_id }}
          </template>
        </el-table-column>
        <el-table-column label="资产编号" width="160">
          <template #default="scope">
            {{ scope.row.asset?.asset_no || '资产 ID: ' + scope.row.asset_id }}
          </template>
        </el-table-column>
        <el-table-column label="资产名称">
          <template #default="scope">
            {{ scope.row.asset?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="申请类型" width="100" align="center">
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
        <el-table-column prop="reject_reason" label="驳回原因" min-width="120">
          <template #default="scope">
            {{ scope.row.reject_reason || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
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
        <el-table-column label="操作" width="160" align="center">
          <template #default="scope">
            <div v-if="scope.row.status === 0">
              <el-button size="small" type="success" @click="handleAudit(scope.row, 1)">同意</el-button>
              <el-button size="small" type="danger" @click="handleAudit(scope.row, 2)">驳回</el-button>
            </div>
            <span v-else style="color: #909399; font-size: 13px;">已处理</span>
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

const applyLogs = ref([])
const applyLoading = ref(false)

// 分页状态 (一个界面最多显示12条数据)
const applyCurrentPage = ref(1)
const applyPageSize = ref(12)

// 分页截取计算属性
const displayedApplyLogs = computed(() => {
  const start = (applyCurrentPage.value - 1) * applyPageSize.value
  const end = start + applyPageSize.value
  return applyLogs.value.slice(start, end)
})

const fetchApplyLogs = () => {
  applyLoading.value = true
  request.get('/admin/apply/list').then(res => {
    applyLogs.value = res || []
    applyCurrentPage.value = 1
  }).catch(() => {
    ElMessage.error('获取领用申请列表失败')
  }).finally(() => {
    applyLoading.value = false
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

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', { hour12: false })
}

const handleAudit = (row, status) => {
  if (status === 2) {
    ElMessageBox.prompt('请输入驳回原因', '拒绝申请', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /\S+/,
      inputErrorMessage: '驳回原因不能为空'
    }).then(({ value }) => {
      request.put('/admin/apply/audit', {
        id: row.ID,
        status: status,
        reject_reason: value
      }).then(() => {
        ElMessage.success('已驳回该申请')
        fetchApplyLogs()
      }).catch(err => {
        ElMessage.error('驳回操作失败: ' + (err.response?.data?.message || err.message))
      })
    }).catch(() => {
      // 取消操作
    })
  } else {
    request.put('/admin/apply/audit', {
      id: row.ID,
      status: status
    }).then(() => {
      ElMessage.success('已同意该申请')
      fetchApplyLogs()
    }).catch(err => {
      ElMessage.error('同意操作失败: ' + (err.response?.data?.message || err.message))
    })
  }
}

onMounted(() => {
  fetchApplyLogs()
})
</script>

<style scoped>
.apply-list {
  padding: 0;
}
</style>
