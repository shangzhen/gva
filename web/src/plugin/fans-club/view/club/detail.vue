<template>
  <div>
    <!-- 粉丝团信息 -->
    <el-card class="box-card" style="margin-bottom: 20px">
      <template #header>
        <div class="card-header">
          <span>粉丝团信息</span>
          <el-button v-if="!clubInfo.isMember" type="success" @click="handleJoin">加入粉丝团</el-button>
          <el-button v-if="clubInfo.isOwner" type="primary" @click="handleEdit">编辑</el-button>
        </div>
      </template>
      <div class="club-info">
        <div class="avatar-section">
          <el-avatar v-if="clubInfo.club.avatar" :src="clubInfo.club.avatar" :size="100" />
          <el-avatar v-else :size="100">{{ clubInfo.club.name?.charAt(0) }}</el-avatar>
        </div>
        <div class="info-section">
          <h2>{{ clubInfo.club.name }}</h2>
          <p class="description">{{ clubInfo.club.description }}</p>
          <div class="stats">
            <el-tag>成员数: {{ clubInfo.club.memberCount }}</el-tag>
            <el-tag type="warning" style="margin-left: 10px">等级: {{ clubInfo.club.level }}</el-tag>
            <el-tag v-if="clubInfo.club.status === 1" type="success" style="margin-left: 10px">正常</el-tag>
            <el-tag v-if="clubInfo.isMember" type="success" style="margin-left: 10px">我的角色: {{ roleText }}</el-tag>
          </div>
        </div>
      </div>
    </el-card>

    <!-- Tab切换 -->
    <el-tabs v-model="activeTab" @tab-click="handleTabClick">
      <el-tab-pane label="粉丝动态" name="posts">
        <div class="gva-btn-list" style="margin-bottom: 20px">
          <el-button v-if="clubInfo.isMember" type="primary" icon="plus" @click="openPostDialog">发布动态</el-button>
        </div>

        <!-- 动态列表 -->
        <div class="post-list">
          <el-card v-for="post in postList" :key="post.ID" class="post-item">
            <template #header>
              <div class="post-header">
                <span>用户ID: {{ post.userId }}</span>
                <span class="post-time">{{ formatTime(post.CreatedAt) }}</span>
              </div>
            </template>
            <div class="post-content">{{ post.content }}</div>
            <div v-if="post.images" class="post-images">
              <!-- 这里可以展示图片 -->
            </div>
            <div class="post-actions">
              <el-button link @click="handleLike(post)">
                <el-icon><StarFilled /></el-icon>
                点赞 ({{ post.likeCount }})
              </el-button>
              <el-button v-if="canDeletePost(post)" type="danger" link @click="handleDeletePost(post)">删除</el-button>
            </div>
          </el-card>

          <el-empty v-if="postList.length === 0" description="暂无动态" />
        </div>

        <div class="gva-pagination">
          <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="postPage"
            :page-size="postPageSize"
            :page-sizes="[10, 20, 50]"
            :total="postTotal"
            @current-change="handlePostPageChange"
            @size-change="handlePostSizeChange"
          />
        </div>
      </el-tab-pane>

      <el-tab-pane label="成员管理" name="members">
        <el-table :data="memberList" style="width: 100%">
          <el-table-column prop="ID" label="ID" width="80" />
          <el-table-column prop="userId" label="用户ID" width="100" />
          <el-table-column prop="role" label="角色" width="120">
            <template #default="{ row }">
              <el-tag v-if="row.role === 'owner'" type="danger">团长</el-tag>
              <el-tag v-else-if="row.role === 'admin'" type="warning">管理员</el-tag>
              <el-tag v-else type="success">成员</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="level" label="等级" width="80" />
          <el-table-column prop="points" label="积分" width="100" />
          <el-table-column prop="joinedAt" label="加入时间" min-width="180">
            <template #default="{ row }">
              {{ formatTime(row.joinedAt) }}
            </template>
          </el-table-column>
          <el-table-column v-if="clubInfo.isOwner" label="操作" min-width="200">
            <template #default="{ row }">
              <el-button
                v-if="row.role === 'member'"
                type="primary"
                link
                @click="handleSetAdmin(row)"
              >设为管理员</el-button>
              <el-button
                v-if="row.role === 'admin'"
                type="warning"
                link
                @click="handleSetMember(row)"
              >取消管理员</el-button>
              <el-button
                v-if="row.role !== 'owner'"
                type="danger"
                link
                @click="handleRemoveMember(row)"
              >移除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="gva-pagination">
          <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="memberPage"
            :page-size="memberPageSize"
            :page-sizes="[10, 20, 50]"
            :total="memberTotal"
            @current-change="handleMemberPageChange"
            @size-change="handleMemberSizeChange"
          />
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 发布动态对话框 -->
    <el-dialog v-model="postDialogVisible" title="发布动态" width="600px">
      <el-form ref="postFormRef" :model="postFormData" :rules="postRules" label-width="100px">
        <el-form-item label="动态内容" prop="content">
          <el-input
            v-model="postFormData.content"
            type="textarea"
            :rows="6"
            placeholder="分享你的想法..."
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="postDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitPost">发布</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router'
import { getFansClub } from '@/plugin/fans-club/api/fansClub'
import { joinClub, getMemberList, updateMemberRole, removeMember } from '@/plugin/fans-club/api/fansClubMember'
import { getPostList, createPost, deletePost, likePost } from '@/plugin/fans-club/api/fansClubPost'

const route = useRoute()
const clubId = ref(Number(route.params.id))

// 粉丝团信息
const clubInfo = ref({
  club: {},
  isMember: false,
  isOwner: false,
  role: ''
})

// Tab
const activeTab = ref('posts')

// 动态列表
const postList = ref([])
const postPage = ref(1)
const postPageSize = ref(10)
const postTotal = ref(0)

// 成员列表
const memberList = ref([])
const memberPage = ref(1)
const memberPageSize = ref(10)
const memberTotal = ref(0)

// 发布动态对话框
const postDialogVisible = ref(false)
const postFormData = ref({
  clubId: clubId.value,
  content: '',
  images: []
})
const postFormRef = ref(null)

const postRules = reactive({
  content: [
    { required: true, message: '请输入动态内容', trigger: 'blur' },
    { min: 1, max: 1000, message: '长度在 1 到 1000 个字符', trigger: 'blur' }
  ]
})

// 角色文本
const roleText = computed(() => {
  const roleMap = {
    owner: '团长',
    admin: '管理员',
    member: '成员'
  }
  return roleMap[clubInfo.value.role] || ''
})

// 获取粉丝团信息
const getClubInfo = async() => {
  const res = await getFansClub({ id: clubId.value })
  if (res.code === 0) {
    clubInfo.value = res.data
  }
}

// 获取动态列表
const getPostData = async() => {
  const params = {
    clubId: clubId.value,
    page: postPage.value,
    pageSize: postPageSize.value
  }
  const res = await getPostList(params)
  if (res.code === 0) {
    postList.value = res.data.list || []
    postTotal.value = res.data.total
  }
}

// 获取成员列表
const getMemberData = async() => {
  const params = {
    clubId: clubId.value,
    page: memberPage.value,
    pageSize: memberPageSize.value
  }
  const res = await getMemberList(params)
  if (res.code === 0) {
    memberList.value = res.data.list || []
    memberTotal.value = res.data.total
  }
}

// Tab切换
const handleTabClick = (tab) => {
  if (tab.props.name === 'posts') {
    getPostData()
  } else if (tab.props.name === 'members') {
    getMemberData()
  }
}

// 加入粉丝团
const handleJoin = async() => {
  const res = await joinClub({ clubId: clubId.value })
  if (res.code === 0) {
    ElMessage.success('加入成功')
    getClubInfo()
  }
}

// 编辑粉丝团
const handleEdit = () => {
  // TODO: 跳转到编辑页面或打开编辑对话框
  ElMessage.info('编辑功能待实现')
}

// 打开发布动态对话框
const openPostDialog = () => {
  postFormData.value = {
    clubId: clubId.value,
    content: '',
    images: []
  }
  postDialogVisible.value = true
}

// 提交动态
const handleSubmitPost = async() => {
  await postFormRef.value?.validate(async(valid) => {
    if (!valid) return

    const res = await createPost(postFormData.value)
    if (res.code === 0) {
      ElMessage.success('发布成功')
      postDialogVisible.value = false
      getPostData()
    }
  })
}

// 点赞
const handleLike = async(post) => {
  const res = await likePost({ id: post.ID })
  if (res.code === 0) {
    ElMessage.success('点赞成功')
    getPostData()
  }
}

// 判断是否可以删除动态
const canDeletePost = (post) => {
  // 自己的动态或管理员可以删除
  return clubInfo.value.isOwner || clubInfo.value.role === 'admin'
}

// 删除动态
const handleDeletePost = async(post) => {
  await ElMessageBox.confirm('确定要删除这条动态吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })

  const res = await deletePost({ id: post.ID })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    getPostData()
  }
}

// 设为管理员
const handleSetAdmin = async(row) => {
  const res = await updateMemberRole({ id: row.ID, role: 'admin' })
  if (res.code === 0) {
    ElMessage.success('设置成功')
    getMemberData()
  }
}

// 取消管理员
const handleSetMember = async(row) => {
  const res = await updateMemberRole({ id: row.ID, role: 'member' })
  if (res.code === 0) {
    ElMessage.success('设置成功')
    getMemberData()
  }
}

// 移除成员
const handleRemoveMember = async(row) => {
  await ElMessageBox.confirm('确定要移除该成员吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })

  const res = await removeMember({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success('移除成功')
    getMemberData()
  }
}

// 分页
const handlePostPageChange = (val) => {
  postPage.value = val
  getPostData()
}

const handlePostSizeChange = (val) => {
  postPageSize.value = val
  getPostData()
}

const handleMemberPageChange = (val) => {
  memberPage.value = val
  getMemberData()
}

const handleMemberSizeChange = (val) => {
  memberPageSize.value = val
  getMemberData()
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  getClubInfo()
  getPostData()
})
</script>

<style scoped>
.club-info {
  display: flex;
  gap: 20px;
}

.avatar-section {
  flex-shrink: 0;
}

.info-section {
  flex: 1;
}

.info-section h2 {
  margin: 0 0 10px 0;
}

.description {
  color: #666;
  margin-bottom: 15px;
}

.stats {
  display: flex;
  gap: 10px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.post-list {
  margin-bottom: 20px;
}

.post-item {
  margin-bottom: 15px;
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.post-time {
  color: #999;
  font-size: 12px;
}

.post-content {
  margin-bottom: 15px;
  line-height: 1.6;
}

.post-actions {
  display: flex;
  gap: 20px;
  border-top: 1px solid #eee;
  padding-top: 10px;
}
</style>
