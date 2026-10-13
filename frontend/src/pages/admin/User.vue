<template>
  <div class="agent-management-container">
    <div class="agent-list-container">
      <div class="header">
        <h2>{{ selectedView === 'agents' ? '智能体广场' : '用户列表' }}</h2>
        <div class="right-tools">
          <v-select
            v-model="selectedView"
            :items="viewOptions"
            item-title="label"
            item-value="value"
            @update:modelValue="handleViewChange"
            density="compact"
            hide-details
            style="width: 200px; margin-right: 10px;"
          />
          <v-text-field
            v-model="searchQuery"
            label="搜索"
            clearable
            density="compact"
            hide-details
            @click:clear="handleSearchClear"
            @keyup.enter="handleSearch"
            style="width: 160px;"
          />
        </div>
      </div>

      <div class="agent-list" v-if="selectedView === 'agents'">
        <template v-if="filteredAgents.length > 0">
        <div
          v-for="agent in filteredAgents"
          :key="agent.id"
          class="agent-item"
        >
          <v-checkbox
            v-model="agent.checked"
            @change="updateSelectedAgents(agent)"
            color="blue"
          ></v-checkbox>
          <div class="agent-content">
            <div class="agent-cover">
              <v-img
                :src="agent.cover || '/images/default_cover.jpg'"
                aspect-ratio="16/9"
                class="cover-image"
                cover
                :lazy-src="agent.cover || '/images/default_cover.jpg'"
                @error="handleImageError"
              ></v-img>
            </div>
            <div class="agent-info">
              <h3>{{ truncate(agent.name, 10) }}</h3>
              <p>{{ truncate(agent.description || '暂无简介', 35) }}</p>
              <v-chip
                :color="getStatusTagType(agent.status)"
                class="agent-status"
                size="small"
              >
              {{ formatStatus(agent.status) }}
              </v-chip>
            </div>
          </div>
        </div>
        </template>
        <template v-else>
          <div style="text-align: center; color: #999; padding: 40px;">
              暂无智能体
          </div>
        </template>
      </div>

      <div class="user-list" v-if="selectedView === 'users'">
        <div
          v-for="user in filteredUsers"
          :key="user.id"
          class="user-item"
        >
          <v-checkbox
            v-model="user.checked"
            @change="updateSelectedUsers(user)"
            color="blue"
          ></v-checkbox>

          <div class="user-content">
            <v-avatar size="64" class="mr-3">
              <v-img :src="user.avatar || '/images/default_avatar.jpg'" />
            </v-avatar>

            <div class="agent-info">
                <div style="display: flex; gap: 10px; margin-bottom: 15px;">
                  <h3>{{ truncate(user.username, 10) }}</h3>
                  <v-chip
                    :color="getRoleColor(user.role)"
                    class="agent-status"
                    size="small"
                  >
                    {{ formatRole(user.role) }}
                  </v-chip>
                </div>
                <p>{{ truncate(user.profile || '暂无简介', 35) }}</p>
            </div>
          </div>
          <v-icon size="32" class="detail-icon" @click="showUserDetail(user)">mdi-chevron-right</v-icon>
        </div>
      </div>
    </div>

    <v-dialog v-model="userDetailDialog" max-width="600px">
      <v-card class="user-detail-card">
        <div class="user-basic-info">
          <v-avatar size="80" class="user-avatar">
            <v-img :src="selectedUser?.avatar || '/images/default_avatar.jpg'" />
          </v-avatar>
          <div class="user-text-info">
            <div class="username-row">
              <h3>{{ selectedUser?.username }}</h3>
              <v-chip :color="getRoleColor(selectedUser?.role)" size="small" class="ml-2">
                {{ formatRole(selectedUser?.role) }}
              </v-chip>
            </div>
            <p class="user-profile">{{ selectedUser?.profile || '暂无简介' }}</p>
          </div>
        </div>

        <v-divider class="my-4"></v-divider>

        <div class="penalty-info-section">
          <h4 class="penalty-title">处罚记录</h4>
          <template v-if="penaltyHistory.length > 0">
            <div class="penalty-record">
              <div class="penalty-type">
                <v-icon :color="getPenaltyIconColor(editingPenalty.Type)" left>
                  {{ getPenaltyIcon(editingPenalty.Type) }}
                </v-icon>
                <v-select
                  v-model="editingPenalty.Type"
                  :items="actionOptions"
                  item-title="label"
                  item-value="value"
                  density="compact"
                  hide-details
                  v-if="isEditingPenalty"
                />
                <span v-else>{{ formatPenaltyType(penaltyHistory[0].Type) }}</span>
              </div>
              <div class="penalty-details">
                <div>
                </div>
                <div v-if="penaltyHistory[0].ExpireAt">
                  <span class="detail-label">到期:</span>
                  <v-text-field
                    v-model.number="editingPenalty.duration"
                    type="number"
                    min="1"
                    density="compact"
                    hide-details
                    v-if="isEditingPenalty"
                  />
                  <span v-else>{{ formatDate(penaltyHistory[0].ExpireAt) }}</span>
                </div>
                <div>
                  <span class="detail-label">原因:</span>
                  <v-text-field
                    v-model="editingPenalty.Reason"
                    density="compact"
                    hide-details
                    v-if="isEditingPenalty"
                  />
                  <span v-else>{{ penaltyHistory[0].Reason || '无' }}</span>
                </div>
              </div>
            </div>
          </template>
          <template v-else>
            <div class="no-penalty">
              <v-icon color="grey" class="mr-2">mdi-information-outline</v-icon>
              该用户暂无处罚记录
            </div>
          </template>
        </div>

        <v-card-actions class="justify-end pa-4">
          <v-btn
            color="warning"
            variant="text"
            @click="handleEditPenalty"
            v-if="penaltyHistory.length > 0"
          >
            {{ isEditingPenalty ? '取消编辑' : '修改' }}
          </v-btn>
          <v-btn
            color="primary"
            variant="text"
            @click="handleSavePenalty"
            v-if="isEditingPenalty"
          >
            保存
          </v-btn>
          <v-btn color="primary" variant="text" @click="userDetailDialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>



    <div class="selected-owners-container">
      <div class="header">
          <h2>已选用户列表 ({{ selectedOwners.size }})</h2>
      </div>

      <div
        v-for="ownerId in Array.from(selectedOwners)"
        :key="ownerId"
        class="owner-item"
      >
        <template v-if="getUserById(ownerId)">
          <div style="display: flex; align-items: center; gap: 10px;">
            <v-avatar size="32">
              <v-img :src="getUserById(ownerId).avatar || '/images/default_avatar.jpg'" />
            </v-avatar>
            <span>{{ getUserById(ownerId).username }}</span>
          </div>
        </template>
          <template v-else>
              <span>{{ ownerId }}</span>
          </template>
      </div>

      <div class="action-panel">
        <v-btn
          color="primary"
          class="submit-btn"
          @click="showActionDialog"
          :disabled="selectedOwners.size === 0"
        >
          提交操作
        </v-btn>
      </div>
    </div>

      <v-dialog v-model="actionDialogVisible" max-width="500px" persistent>
        <v-card>
          <v-card-title class="text-h6">用户管理操作</v-card-title>
          <v-card-text>
            <v-form ref="actionFormRef" v-model="actionFormValid">
              <v-select
                v-model="actionForm.actionType"
                :items="actionOptions"
                item-title="label"
                item-value="value"
                label="操作类型"
                required
              ></v-select>
              <v-text-field
                v-if="showDurationField"
                v-model.number="actionForm.duration"
                label="封禁天数"
                type="number"
                min="1"
                required
              ></v-text-field>
              <v-textarea
                v-model="actionForm.reason"
                label="原因说明（可选）"
                rows="3"
              ></v-textarea>
            </v-form>
          </v-card-text>
          <v-card-actions class="justify-end">
            <v-btn text @click="actionDialogVisible = false">取消</v-btn>
            <v-btn color="primary" @click="handleSubmit">提交</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
  </div>

</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import server from '@/api/server.ts'
import {ElMessage} from "element-plus";

const agents = ref([])
const users = ref([])
const selectedAgents = ref([])
const selectedUsers = ref([])
const searchQuery = ref('')
const actionDialogVisible = ref(false)
const actionForm = ref({ actionType: '', reason: '', duration: undefined })
const actionFormValid = ref(false)
const selectedView = ref('agents')
const userDetailDialog = ref(false)
const selectedUser = ref(null)
const penaltyHistory = ref([])
const isEditingPenalty = ref(false)
const editingPenalty = ref({
  Type: '',
  Reason: '',
  duration: 0
})

const actionOptions = [
    { value: 'warning', label: '警告', requiresDuration: true },
    { value: 'restrict_comment', label: '违规评论', requiresDuration: true },
    { value: 'restrict_agent_publish', label: '违规智能体', requiresDuration: true },
    { value: 'temporary_ban', label: '临时封禁', requiresDuration: true },
    { value: 'permanent_ban', label: '永久封禁' }
]

const viewOptions = [
    { value: 'agents', label: '智能体广场' },
    { value: 'users', label: '用户列表' }
]

const getActionTypeTitle = (actionType) => {
  const titleMap = {
    'warning': '警告',
    'restrict_comment': '违规评论',
    'restrict_agent_publish': '违规智能体',
    'temporary_ban': '临时封禁',
    'permanent_ban': '永久封禁'
  };
  return titleMap[actionType] || '违规通知';
};

const selectedOwners = computed(() => {
    const owners = new Set()
    if (selectedView.value === 'agents') {
        selectedAgents.value.forEach(agentId => {
            const agent = agents.value.find(a => a.id === agentId)
            if (agent.belongsTo?.id) {
                owners.add(agent.belongsTo.id)
            }
        })
    } else if (selectedView.value === 'users') {
        selectedUsers.value.forEach(userId => {
            owners.add(userId)
        })
    }
    return owners
})

const updateSelectedAgents = (agent) => {
    if (agent.checked) {
        selectedAgents.value.push(agent.id)
    } else {
        selectedAgents.value = selectedAgents.value.filter(id => id !== agent.id)
    }
}

const updateSelectedUsers = (user) => {
    if (user.checked) {
        selectedUsers.value.push(user.id)
    } else {
        selectedUsers.value = selectedUsers.value.filter(id => id !== user.id)
    }
}

const fetchAgents = async () => {
    console.log("agents")
    try {
        const response = await server.Get('/workflow/public', {
          params: {
            page: 1,
            pageSize: 200
          }
        })
      console.log(response.data)
        agents.value = response.data.workflows || []
    } catch (error) {
        console.error('Error fetching agents:', error)
    }
}

const fetchUsers = async () => {
    console.log("users")
    try {
        const response = await server.Get('/admin/all_user')
        console.log(response.data)
        users.value = response.data || []
    } catch (error) {
        console.error('Error fetching users:', error)
    }
}

const showUserDetail = async (user) => {
  selectedUser.value = user
  try {
    const response = await server.Get(`/admin/penalty/${user.id}`)
    if (response) {
      penaltyHistory.value = response.data ? [response.data] : []
    }
  } catch (error) {
    penaltyHistory.value = []
  }
  userDetailDialog.value = true
}

const formatPenaltyType = (type) => {
  const typeMap = {
    'warning': '警告',
    'restrict_comment': '违规评论',
    'restrict_agent_publish': '违规智能体',
    'temporary_ban': '临时封禁',
    'permanent_ban': '永久封禁'
  }
  return typeMap[type] || type
}

const getPenaltyIcon = (type) => {
  const iconMap = {
    'warning': 'mdi-alert-outline',
    'restrict_comment': 'mdi-comment-remove-outline',
    'restrict_agent_publish': 'mdi-note-remove-outline',
    'temporary_ban': 'mdi-account-clock-outline',
    'permanent_ban': 'mdi-account-lock-outline'
  }
  return iconMap[type] || 'mdi-alert-circle'
}

const getPenaltyIconColor = (type) => {
  const colorMap = {
    'warning': 'warning',
    'restrict_comment': 'orange',
    'restrict_agent_publish': 'deep-orange',
    'temporary_ban': 'red',
    'permanent_ban': 'deep-purple'
  }
  return colorMap[type] || 'grey'
}

const getPenaltyCardClass = (type) => {
  return `penalty-card-${type}`
}

const formatExpiration = (expireAt) => {
  if (!expireAt) return '永久'
  const now = new Date()
  const expireDate = new Date(expireAt)
  const diffDays = Math.ceil((expireDate - now) / (1000 * 60 * 60 * 24))
  return diffDays > 0 ? `剩余${diffDays}天` : '已过期'
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toISOString().split('T')[0]
}

const handleEditPenalty = () => {
  if (isEditingPenalty.value) {
    isEditingPenalty.value = false
    if (penaltyHistory.value.length > 0) {
      editingPenalty.value = {
        Type: penaltyHistory.value[0].Type,
        Reason: penaltyHistory.value[0].Reason || '',
        duration: calculateDuration(penaltyHistory.value[0].ExpireAt)
      }
    }
  } else {
    isEditingPenalty.value = true
    if (penaltyHistory.value.length > 0) {
      editingPenalty.value = {
        Type: penaltyHistory.value[0].Type,
        Reason: penaltyHistory.value[0].Reason || '',
        duration: calculateDuration(penaltyHistory.value[0].ExpireAt)
      }
    }
  }
}

const calculateDuration = (expireAt) => {
  if (!expireAt) return 0
  const now = new Date()
  const expireDate = new Date(expireAt)
  return Math.ceil((expireDate - now) / (1000 * 60 * 60 * 24))
}

const handleSavePenalty = async () => {
  try {
    const payload = {
      user_id: selectedUser.value.id,
      penalty_id: penaltyHistory.value[0].ID,
      new_penalty_type: editingPenalty.value.Type,
      new_reason: editingPenalty.value.Reason,
      new_duration: editingPenalty.value.duration.toString()
    }

    const response = await server.Post('/admin/update_penalty', payload)
    ElMessage.success('处罚记录更新成功')

    const refreshResponse = await server.Get(`/admin/penalty/${selectedUser.value.id}`)
    penaltyHistory.value = refreshResponse.data ? [refreshResponse.data] : []
    isEditingPenalty.value = false
  } catch (error) {
    console.error('更新处罚记录失败:', error)
    ElMessage.error(`更新失败: ${error.response?.data?.message || error.message}`)
  }
}


const getUserById = (id) => {
    return users.value.find(user => user.id === id)
}

const filteredAgents = computed(() => {
    if (!searchQuery.value) return agents.value
    const query = searchQuery.value.toLowerCase()
    return agents.value.filter(agent =>
        agent.name.toLowerCase().includes(query) ||
        (agent.description && agent.description.toLowerCase().includes(query))
    )
})

const filteredUsers = computed(() => {
    if (!searchQuery.value) return users.value
    const query = searchQuery.value.toLowerCase()
    return users.value.filter(user =>
        user.username.toLowerCase().includes(query) ||
        (user.profile && user.profile.toLowerCase().includes(query))
    )
})

const truncate = (text, length) => {
    if (!text) return ''
    return text.length <= length ? text : text.slice(0, length) + '...'
}

const snackbar = ref({
    show: false,
    message: ''
})

const handleSearch = () => {}
const handleSearchClear = () => { searchQuery.value = '' }
const showActionDialog = () => { actionDialogVisible.value = true }

const showDurationField = computed(() => {
  const selectedAction = actionOptions.find(opt => opt.value === actionForm.value.actionType)
  return selectedAction?.requiresDuration || false
})

const handleActionTypeChange = (newType) => {
  if (!actionOptions.find(opt => opt.value === newType)?.requiresDuration) {
    actionForm.value.duration = undefined
  }
}

const handleSubmit = async () => {
  try {
    const payload = {
      users_id: Array.from(selectedOwners.value),
      penalty_type: actionForm.value.actionType,
      reason: actionForm.value.reason || ''
    };

    if (showDurationField.value && !actionForm.value.duration) {
      ElMessage.error('请填写封禁天数')
      return
    }

    if (showDurationField.value) {
      payload.duration = String(actionForm.value.duration)
    }

    const notificationTitle = getActionTypeTitle(actionForm.value.actionType);

    await Promise.all([
      server.Post('/admin/penalty', payload),

      Promise.all(
          Array.from(selectedOwners.value).map(userId =>
              server.Post('/notification/send', {
                title: notificationTitle,
                content: actionForm.value.reason || '您的账号存在违规行为，请遵守社区规则',
                to_user_id: userId
              }).catch(e => {
                console.error(`用户 ${userId} 通知发送失败:`, e);
                return { success: false, userId };
              })
          )
      )
    ]);
    ElMessage.success("违规处理发送成功")

    actionDialogVisible.value = false;
    actionForm.value = { actionType: '', reason: '' };
    selectedUsers.value = [];
    selectedAgents.value = [];
    users.value.forEach(user => user.checked = false);
    agents.value.forEach(agent => agent.checked = false);

  } catch (error) {
    console.error('操作失败:', error);
    ElMessage.error(`操作失败: ${error.response?.data?.message || error.message}`);
  }
};

const formatStatus = (status) => {
    const statusMap = {
        private: '私有',
        public: '公开'
    }
    return statusMap[status] || status
}

const getStatusTagType = (status) => {
    const typeMap = {
        private: 'red',
        public: 'green'
    }
    return typeMap[status] || 'grey'
}

const getRoleColor = (role) => {
    return role === 'admin' ? 'red' : 'green'
}

const formatRole = (role) => {
    const roleMap = {
        admin: '管理员',
        user: '用户'
    }
    return roleMap[role] || role
}

const handleViewChange = (newView) => {
    if (newView === 'users') {
        fetchUsers()
    } else {
        fetchAgents()
    }
}

const handleImageError = (event) => {
  if (event?.target) {
    event.target.src = '/images/default_cover.jpg'
  }
}

onMounted(() => {
  fetchAgents()
  fetchUsers()
})
</script>

<style scoped>
.agent-management-container {
    display: flex;
    height: calc(100vh - 10px);
    padding: 20px;
    gap: 20px;
}

.agent-list-container,
.selected-owners-container {
    background-color: white;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
    display: flex;
    flex-direction: column;
}

.agent-list-container { flex: 3; }
.selected-owners-container { flex: 1; }

.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.right-tools {
    display: flex;
    align-items: center;
}

.agent-list, .user-list { flex: 1; overflow-y: auto; }

.agent-item,
.user-item {
    display: flex;
    align-items: flex-start;
    padding: 15px;
    border-bottom: 1px solid #ebeef5;
    gap: 15px;
}

.agent-content,
.user-content {
    display: flex;
    flex: 1;
    gap: 15px;
}

.agent-cover {
    aspect-ratio: 16 / 9;
    width: 160px;
    flex-shrink: 0;
    border-radius: 4px;
    overflow: hidden;
}

.cover-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 4px;
}

.agent-info,
.user-info {
    display: flex;
    flex-direction: column;
    flex: 1;
    min-width: 0;
}

.agent-status {
    align-self: flex-start;
    margin-top: 5px;
}

.owners-list {
    flex: 1;
    overflow-y: auto;
    margin-bottom: 20px;
}

.owner-item {
    padding: 10px 15px;
    border-bottom: 1px solid #ebeef5;
    font-size: 14px;
    color: #606266;
}

.action-panel {
    margin-top: auto;
    padding-top: 20px;
    border-top: 1px solid #ebeef5;
    background-color: white;
}

.submit-btn {
    width: 100%;
}

.user-item {
  cursor: pointer;
  position: relative;
}

.detail-icon {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: #999;
}

.user-detail-card {
  padding: 20px;
}

.user-basic-info {
  display: flex;
  align-items: center;
  padding: 16px;
}

.user-avatar {
  flex-shrink: 0;
}

.user-text-info {
  margin-left: 20px;
}

.username-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.user-profile {
  color: #666;
  margin: 0;
}

.penalty-info-section {
  padding: 0 16px;
}

.penalty-title {
  margin-bottom: 16px;
  font-weight: 500;
}

.penalty-record {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 16px;
}

.penalty-type {
  display: flex;
  align-items: center;
  font-weight: 500;
  margin-bottom: 12px;
  gap: 8px;
}

.penalty-details {
  font-size: 14px;
}

.penalty-details > div {
  margin-bottom: 8px;
  display: flex;
  align-items: center;
}

.detail-label {
  color: #666;
  display: inline-block;
  width: 50px;
  flex-shrink: 0;
}
.no-penalty {
  display: flex;
  align-items: center;
  color: #666;
  padding: 16px;
  background-color: #f9f9f9;
  border-radius: 8px;
}
</style>
