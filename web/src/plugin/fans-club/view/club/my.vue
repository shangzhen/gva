<template>
  <div>
    <div class="gva-table-box">
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
      >
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
        <el-table-column align="left" label="我的角色" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.role === 'owner'" type="danger">团长</el-tag>
            <el-tag v-else-if="row.role === 'admin'" type="warning">管理员</el-tag>
            <el-tag v-else type="success">成员</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.club.status === 1" type="success">正常</el-tag>
            <el-tag v-else-if="row.club.status === 0" type="warning">待审核</el-tag>
            <el-tag v-else type="danger">禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="200">
          <template #default="{ row }">
            <el-button type="primary" link icon="view" @click="toDetail(row)">详情</el-button>
            <el-button v-if="row.role !== 'owner'" type="warning" link icon="close" @click="handleQuit(row)">退出</el-button>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { getMyClubs } from '@/plugin/fans-club/api/fansClub'
import { quitClub } from '@/plugin/fans-club/api/fansClubMember'

const router = useRouter()

// 表格数据
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 获取列表
const getTableData = async() => {
  const params = {
    page: page.value,
    pageSize: pageSize.value
  }
  const res = await getMyClubs(params)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
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

// 退出粉丝团
const handleQuit = async(row) => {
  await ElMessageBox.confirm('确定要退出该粉丝团吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })

  const res = await quitClub({ id: row.club.ID })
  if (res.code === 0) {
    ElMessage.success('退出成功')
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
