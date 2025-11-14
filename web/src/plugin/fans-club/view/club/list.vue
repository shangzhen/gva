<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="粉丝团名称" prop="keyword">
          <el-input v-model="searchInfo.keyword" placeholder="搜索粉丝团名称" clearable />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="searchInfo.status" placeholder="请选择状态" clearable>
            <el-option label="全部" :value="0" />
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">创建粉丝团</el-button>
      </div>

      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="头像" width="80">
          <template #default="{ row }">
            <el-avatar v-if="row.club.avatar" :src="row.club.avatar" :size="50" />
            <el-avatar v-else :size="50">{{ row.club.name?.charAt(0) }}</el-avatar>
          </template>
        </el-table-column>
        <el-table-column align="left" label="粉丝团名称" prop="club.name" min-width="150" />
        <el-table-column align="left" label="描述" prop="club.description" min-width="200" show-overflow-tooltip />
        <el-table-column align="left" label="成员数" prop="club.memberCount" width="100" />
        <el-table-column align="left" label="等级" prop="club.level" width="80" />
        <el-table-column align="left" label="状态" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.club.status === 1" type="success">正常</el-tag>
            <el-tag v-else-if="row.club.status === 0" type="warning">待审核</el-tag>
            <el-tag v-else type="danger">禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="是否成员" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.isMember" type="success">已加入</el-tag>
            <el-tag v-else type="info">未加入</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="240">
          <template #default="{ row }">
            <el-button type="primary" link icon="view" @click="toDetail(row)">详情</el-button>
            <el-button v-if="!row.isMember" type="success" link icon="plus" @click="handleJoin(row)">加入</el-button>
            <el-button v-if="row.isOwner" type="primary" link icon="edit" @click="handleEdit(row)">编辑</el-button>
            <el-button v-if="row.isOwner" type="danger" link icon="delete" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogFormVisible"
      :title="dialogTitle"
      width="600px"
    >
      <el-form ref="elFormRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="粉丝团名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入粉丝团名称" />
        </el-form-item>
        <el-form-item label="粉丝团描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="4"
            placeholder="请输入粉丝团描述"
          />
        </el-form-item>
        <el-form-item label="头像URL" prop="avatar">
          <el-input v-model="formData.avatar" placeholder="请输入头像URL" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeDialog">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import {
  getFansClubList,
  createFansClub,
  updateFansClub,
  deleteFansClub
} from '@/plugin/fans-club/api/fansClub'
import { joinClub } from '@/plugin/fans-club/api/fansClubMember'

const router = useRouter()

// 搜索条件
const searchInfo = ref({
  keyword: '',
  status: 0
})

// 表格数据
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框
const dialogFormVisible = ref(false)
const dialogTitle = ref('创建粉丝团')
const formData = ref({
  id: 0,
  name: '',
  description: '',
  avatar: ''
})

const elFormRef = ref(null)

// 表单验证规则
const rules = reactive({
  name: [
    { required: true, message: '请输入粉丝团名称', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
  ]
})

// 获取列表
const getTableData = async() => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    keyword: searchInfo.value.keyword,
    status: searchInfo.value.status
  }
  const res = await getFansClubList(params)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

// 查询
const onSubmit = () => {
  page.value = 1
  getTableData()
}

// 重置
const onReset = () => {
  searchInfo.value = {
    keyword: '',
    status: 0
  }
  onSubmit()
}

// 分页
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 打开对话框
const openDialog = () => {
  dialogTitle.value = '创建粉丝团'
  formData.value = {
    id: 0,
    name: '',
    description: '',
    avatar: ''
  }
  dialogFormVisible.value = true
}

// 关闭对话框
const closeDialog = () => {
  dialogFormVisible.value = false
  elFormRef.value?.clearValidate()
}

// 编辑
const handleEdit = (row) => {
  dialogTitle.value = '编辑粉丝团'
  formData.value = {
    id: row.club.ID,
    name: row.club.name,
    description: row.club.description,
    avatar: row.club.avatar
  }
  dialogFormVisible.value = true
}

// 提交表单
const handleSubmit = async() => {
  await elFormRef.value?.validate(async(valid) => {
    if (!valid) return

    let res
    if (formData.value.id > 0) {
      res = await updateFansClub(formData.value)
    } else {
      res = await createFansClub(formData.value)
    }

    if (res.code === 0) {
      ElMessage.success(formData.value.id > 0 ? '更新成功' : '创建成功')
      closeDialog()
      getTableData()
    }
  })
}

// 加入粉丝团
const handleJoin = async(row) => {
  const res = await joinClub({ clubId: row.club.ID })
  if (res.code === 0) {
    ElMessage.success('加入成功')
    getTableData()
  }
}

// 删除
const handleDelete = async(row) => {
  await ElMessageBox.confirm('此操作将永久删除该粉丝团及所有相关数据, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })

  const res = await deleteFansClub({ id: row.club.ID })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    getTableData()
  }
}

// 跳转详情
const toDetail = (row) => {
  router.push({ name: 'clubDetail', params: { id: row.club.ID } })
}

onMounted(() => {
  getTableData()
})
</script>

<style scoped>
</style>
