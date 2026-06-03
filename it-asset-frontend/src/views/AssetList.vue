<template>
  <div class="asset-list">
    <el-card>
      <template #header>
        <div class="clearfix">
          <span>资产列表</span>
          <el-button style="float: right; margin-left: 10px;" type="primary" size="small" @click="fetchAssets">刷新</el-button>
          <el-button v-if="user.role_id === 2 || user.role_id === 3" style="float: right;" type="success" size="small" @click="showAddDialog = true">添加资产</el-button>
        </div>
      </template>

      <el-table :data="groupedAssets" style="width: 100%" v-loading="loading" border stripe>
        <!-- 折叠展开列 -->
        <el-table-column type="expand">
          <template #default="props">
            <div style="padding: 10px 20px; background-color: #fafafa;">
              <div style="font-weight: bold; margin-bottom: 10px; color: #606266; font-size: 13px;">
                📋 单件明细列表（共 {{ props.row.total_quantity }} 件，基础编号：{{ props.row.base_no }}）
              </div>
              <el-table :data="props.row.children" border size="small" style="width: 100%; background: #fff;">
                <el-table-column prop="asset_no" label="资产单件编号" width="220" align="center"></el-table-column>
                <el-table-column prop="status" label="设备状态" width="120" align="center">
                  <template #default="subScope">
                    <el-tag :type="getStatusType(subScope.row.status)" size="small">
                      {{ getStatusName(subScope.row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="持有人" min-width="150" align="center">
                  <template #default="subScope">
                    <span v-if="subScope.row.user_name">{{ subScope.row.user_name }}</span>
                    <span v-else style="color: #909399;">-</span>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="220" align="center">
                  <template #default="subScope">
                    <!-- 员工：持有该单件时可报修/归还 -->
                    <el-button v-if="subScope.row.status === 1 && subScope.row.user_id === user.id" size="small" type="warning" @click="handleRepair(subScope.row)">一键报修</el-button>
                    <el-button v-if="subScope.row.status === 1 && subScope.row.user_id === user.id" size="small" type="success" @click="handleReturn(subScope.row)">申请归还</el-button>
                    <!-- 管理员：闲置时可删除单个单件 -->
                    <el-button v-if="(user.role_id === 2 || user.role_id === 3) && subScope.row.status === 0" size="small" type="danger" plain @click="handleDelete(subScope.row)">删除单件</el-button>
                    <span v-else-if="!(subScope.row.status === 1 && subScope.row.user_id === user.id)">-</span>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="name" label="资产名称" min-width="180"></el-table-column>
        <el-table-column prop="base_no" label="资产编号" width="180"></el-table-column>
        <el-table-column label="可用数量 / 总数量" width="180" align="center">
          <template #default="scope">
            <span style="font-weight: bold; color: #409eff;">{{ scope.row.remaining_quantity }}</span>
            <span style="color: #909399;"> / </span>
            <span>{{ scope.row.total_quantity }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getBatchStatusType(scope.row.remaining_quantity)">
              {{ getBatchStatusName(scope.row.remaining_quantity) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" align="center">
          <template #default="scope">
            <!-- 员工/管理员：有闲置时可申请领用，后台会自动按单件编号顺序分配闲置设备 -->
            <el-button
              v-if="scope.row.remaining_quantity > 0"
              size="small"
              type="primary"
              @click="handleApply(scope.row)"
            >申请领用</el-button>
            <el-button
              v-else
              size="small"
              type="primary"
              disabled
            >暂无闲置</el-button>

            <!-- 管理员：编辑该批次资产属性与数量 -->
            <el-button v-if="user.role_id === 2 || user.role_id === 3" size="small" type="info" @click="openEditDialog(scope.row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加资产弹窗 -->
    <el-dialog v-model="showAddDialog" title="添加资产" width="480px">
      <el-form :model="newAsset" label-width="110px">
        <el-form-item label="基础编号">
          <el-input v-model="newAsset.base_no" placeholder="如 AST20260601005"></el-input>
        </el-form-item>
        <el-form-item label="资产名称">
          <el-input v-model="newAsset.name"></el-input>
        </el-form-item>
        <el-form-item label="资产分类">
          <el-select v-model="newAsset.category_id" placeholder="请选择资产分类" style="width: 100%;">
            <el-option
              v-for="item in categories"
              :key="item.ID"
              :label="item.category_name"
              :value="item.ID"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="数量">
          <el-input-number v-model="newAsset.quantity" :min="1" :max="999" style="width: 100%;"></el-input-number>
          <div style="font-size: 12px; color: #909399; margin-top: 4px;">
            将自动生成 {{ newAsset.quantity }} 个独立资产单件，编号为
            <b>{{ newAsset.base_no || 'XXX' }}-001</b> ~ <b>{{ newAsset.base_no || 'XXX' }}-{{ String(newAsset.quantity).padStart(3, '0') }}</b>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddDialog = false">取消</el-button>
          <el-button type="primary" @click="addAsset">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑资产弹窗（批量） -->
    <el-dialog v-model="showEditDialog" title="编辑资产信息及数量" width="480px">
      <el-form :model="editingAsset" label-width="110px">
        <el-form-item label="资产编号">
          <el-input v-model="editingAsset.base_no" disabled></el-input>
        </el-form-item>
        <el-form-item label="资产名称">
          <el-input v-model="editingAsset.name"></el-input>
        </el-form-item>
        <el-form-item label="资产分类">
          <el-select v-model="editingAsset.category_id" placeholder="请选择资产分类" style="width: 100%;">
            <el-option
              v-for="item in categories"
              :key="item.ID"
              :label="item.category_name"
              :value="item.ID"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="总数量">
          <el-input-number v-model="editingAsset.quantity" :min="1" :max="999" style="width: 100%;"></el-input-number>
          <div style="font-size: 12px; color: #909399; margin-top: 4px;">
            修改数量将自动增减该批次的资产单件。<br/>
            注意：总数量不能少于当前已被借走或维修中的数量。
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="saveEditAsset">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import request from '../utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const assets = ref([])
const loading = ref(false)
const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))

const showAddDialog = ref(false)
const newAsset = ref({
  base_no: '',
  name: '',
  category_id: 1,
  quantity: 1
})

const showEditDialog = ref(false)
const editingAsset = ref({
  id: null,
  asset_no: '',
  base_no: '',
  name: '',
  category_id: 1,
  quantity: 1
})

const categories = ref([
  { ID: 1, category_name: '笔记本电脑' },
  { ID: 2, category_name: '台式电脑' },
  { ID: 3, category_name: '显示器' },
  { ID: 4, category_name: '服务器及网络设备' },
  { ID: 5, category_name: '办公外设' }
])

// 资产按 base_no 批量折叠分组
const groupedAssets = computed(() => {
  const groups = {}
  assets.value.forEach(item => {
    const key = item.base_no || item.asset_no
    if (!groups[key]) {
      groups[key] = {
        ID: item.ID,
        base_no: key,
        name: item.name,
        category_id: item.category_id,
        total_quantity: 0,
        remaining_quantity: 0,
        children: []
      }
    }
    groups[key].children.push(item)
    groups[key].total_quantity++
    if (item.status === 0) {
      groups[key].remaining_quantity++
    }
  })

  // 按单件编号升序排列每组内的单件
  Object.values(groups).forEach(g => {
    g.children.sort((a, b) => a.asset_no.localeCompare(b.asset_no))
  })

  return Object.values(groups)
})

const getBatchStatusName = (remaining) => {
  return remaining > 0 ? '有闲置' : '已领完'
}

const getBatchStatusType = (remaining) => {
  return remaining > 0 ? 'success' : 'info'
}

const fetchCategories = () => {
  if (user.value.role_id === 2 || user.value.role_id === 3) {
    request.get('/admin/categories').then(res => {
      if (res && res.length > 0) {
        categories.value = res
      }
    }).catch(err => {
      console.error('获取分类列表失败，使用预设分类:', err)
    })
  }
}

const fetchAssets = () => {
  loading.value = true
  request.get('/employee/assets/available').then(res => {
    assets.value = res || []
  }).finally(() => {
    loading.value = false
  })
}

const getStatusName = (status) => {
  const map = { 0: '闲置', 1: '使用中', 2: '维修中', 3: '已报废' }
  return map[status] || '未知'
}

const getStatusType = (status) => {
  const map = { 0: 'success', 1: 'primary', 2: 'warning', 3: 'danger' }
  return map[status] || 'info'
}

const handleApply = (row) => {
  ElMessageBox.prompt('请输入领用原因（可留空）', '申请领用', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPlaceholder: '领用原因...'
  }).then(({ value }) => {
    request.post('/employee/apply', {
      asset_id: row.ID,
      type: 1,
      reason: value || ''
    }).then(() => {
      ElMessage.success('领用申请已提交')
      fetchAssets()
    }).catch(err => {
      ElMessage.error('申请失败: ' + (err.response?.data?.message || err.message))
    })
  }).catch(() => {
    // 取消
  })
}

const handleReturn = (row) => {
  request.post('/employee/apply', {
    asset_id: row.ID,
    type: 2
  }).then(() => {
    ElMessage.success('归还申请已提交')
    fetchAssets()
  }).catch(err => {
    ElMessage.error('申请失败: ' + (err.response?.data?.message || err.message))
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
      ElMessage.success('报修申请已提交')
      fetchAssets()
    }).catch(err => {
      ElMessage.error('申请失败: ' + (err.response?.data?.message || err.message))
    })
  }).catch(() => {
    // 取消输入
  })
}

const addAsset = () => {
  if (!newAsset.value.base_no) {
    ElMessage.warning('请填写基础编号')
    return
  }
  if (!newAsset.value.name) {
    ElMessage.warning('请填写资产名称')
    return
  }
  request.post('/admin/asset', newAsset.value).then((res) => {
    const count = Array.isArray(res) ? res.length : 1
    ElMessage.success(`资产添加成功，已创建 ${count} 个单件`)
    showAddDialog.value = false
    newAsset.value = { base_no: '', name: '', category_id: 1, quantity: 1 }
    fetchAssets()
  }).catch(err => {
    ElMessage.error('添加失败: ' + (err.response?.data?.message || err.message))
  })
}

const openEditDialog = (row) => {
  const baseNo = row.base_no
  const batchAssets = assets.value.filter(a => (a.base_no || a.asset_no) === baseNo)
  editingAsset.value = {
    id: row.ID,
    asset_no: baseNo,
    base_no: baseNo,
    name: row.name,
    category_id: row.category_id,
    quantity: batchAssets.length || 1
  }
  showEditDialog.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除资产单件 ${row.asset_no} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    request.delete('/admin/asset', { params: { id: row.ID } }).then(() => {
      ElMessage.success('删除成功')
      fetchAssets()
    }).catch(err => {
      ElMessage.error('删除失败: ' + (err.response?.data?.message || err.message))
    })
  }).catch(() => {})
}

const saveEditAsset = () => {
  request.put('/admin/asset', editingAsset.value).then(() => {
    ElMessage.success('资产更新成功')
    showEditDialog.value = false
    fetchAssets()
  }).catch(err => {
    ElMessage.error('更新失败: ' + (err.response?.data?.message || err.message))
  })
}

onMounted(() => {
  fetchAssets()
  fetchCategories()
})
</script>
