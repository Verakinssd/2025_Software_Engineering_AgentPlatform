<template>
  <v-card class="mx-auto pa-6" max-width="1200">
    <v-form ref="settingsFormRef" @submit.prevent="submitSettings">
      <v-card-title class="text-h5 mb-6">个人设置</v-card-title>

      <v-card-text>
        <v-row>
          <v-col cols="12" md="3">
            <v-label class="text-subtitle-1 mb-2">头像</v-label>
            <el-upload
              class="avatar-uploader"
              action="#"
              :show-file-list="false"
              :auto-upload="false"
              :on-change="handleAvatarChange"
              accept="image/*"
            >
              <img v-if="avatarPreview" :src="avatarPreview" class="avatar" alt=""/>
              <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
            </el-upload>
          </v-col>

          <v-col cols="12" md="9">
            <v-text-field
              v-model="settingsForm.nickname"
              label="昵称"
              :rules="rules.nickname"
              variant="outlined"
              density="comfortable"
              class="mb-4"
            ></v-text-field>

            <v-textarea
              v-model="settingsForm.profile"
              label="个人简介"
              :rules="rules.profile"
              variant="outlined"
              density="comfortable"
              rows="3"
              class="mb-4"
            ></v-textarea>

            <v-text-field
              v-model="settingsForm.email"
              label="邮箱"
              :rules="rules.email"
              variant="outlined"
              density="comfortable"
              class="mb-4"
            ></v-text-field>
          </v-col>
        </v-row>

        <v-divider class="my-6"></v-divider>

        <v-row>
          <v-col cols="12">
            <v-text-field
              v-model="settingsForm.old_password"
              label="当前密码"
              :rules="rules.old_password"
              type="password"
              variant="outlined"
              density="comfortable"
              class="mb-4"
            ></v-text-field>

            <v-text-field
              v-model="settingsForm.new_password"
              label="新密码"
              :rules="rules.new_password"
              type="password"
              variant="outlined"
              density="comfortable"
              class="mb-4"
            ></v-text-field>

            <v-text-field
              v-model="settingsForm.confirm_password"
              label="确认新密码"
              :rules="rules.confirm_password"
              type="password"
              variant="outlined"
              density="comfortable"
            ></v-text-field>
          </v-col>
        </v-row>
      </v-card-text>

      <v-card-actions class="px-4 pb-4">
        <v-spacer></v-spacer>
        <v-btn
          variant="text"
          color="grey-darken-1"
          @click="resetForm"
          class="mr-4"
        >
          重置
        </v-btn>
        <v-btn
          type="submit"
          color="primary"
          :loading="loading"
          variant="flat"
        >
          保存设置
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>

<script lang="ts" setup>
import { reactive, ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getProfile, updateUserProfile, uploadUserAvatar } from '@/api/user'
import { user } from '@/stores/user'

const settingsFormRef = ref()
const loading = ref(false)
const avatarPreview = ref('')

interface UserInfo {
  id: string
  username: string
  profile: string
  email: string
  avatar: string
  role: string
}

const settingsForm = reactive({
  nickname: '',
  profile: '',
  email: '',
  avatar_file: null as File | null,
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const fetchUserProfile = async () => {
  try {
    const res = await getProfile()
    if (res.data) {
      settingsForm.nickname = res.data.username || ''
      settingsForm.profile = res.data.profile || ''
      settingsForm.email = res.data.email || ''
      avatarPreview.value = res.data.avatar || '/images/default_avatar.jpg'
    }
  } catch (error) {
    ElMessage.error('获取用户信息失败')
    window.location.href = '/login'
    console.error(error)
  }
}

const handleAvatarChange = (file: { raw: File }) => {
  const isImage = file.raw.type.startsWith('image/')
  const isLt5M = file.raw.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('头像图片大小不能超过5MB!')
    return false
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    avatarPreview.value = e.target?.result as string
    settingsForm.avatar_file = file.raw
  }
  reader.readAsDataURL(file.raw)
}

const rules = reactive({
  nickname: [
    (v: string) => !v || v.length <= 10 || '昵称长度不能超过10个字符'
  ],
  profile: [
    (v: string) => !v || v.length <= 50 || '简介不能超过50个字'
  ],
  email: [
    (v: string) => !v || /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v) || '请输入有效的邮箱地址'
  ],
  old_password: [
    (v: string) => !v || v.length <= 15 || '密码长度不能超过15位'
  ],
  new_password: [
    (v: string) => !v || v.length <= 15 || '密码长度不能超过15位',
    (v: string) => !v || v !== settingsForm.old_password || '新密码不能与旧密码相同'
  ],
  confirm_password: [
    (v: string) => !v || v === settingsForm.new_password || '两次输入的密码不一致'
  ]
})

const submitSettings = async () => {
  const { valid } = await settingsFormRef.value.validate()
  if (!valid) return

  try {
    loading.value = true

    let newAvatarUrl = ''
    if (settingsForm.avatar_file) {
      newAvatarUrl = (await uploadUserAvatar(settingsForm.avatar_file)).data.avatar
    }

    const updateData: {
      new_username?: string
      profile?: string
      new_email?: string
      old_password?: string
      new_password?: string
    } = {}

    if (settingsForm.nickname !== user.username) {
      updateData.new_username = settingsForm.nickname
    }
    if (settingsForm.profile !== user.profile) {
      updateData.profile = settingsForm.profile
    }
    if (settingsForm.email !== user.email) {
      updateData.new_email = settingsForm.email
    }
    if (settingsForm.new_password) {
      if (!settingsForm.old_password) {
        throw new Error('修改密码需要提供当前密码')
      }
      updateData.old_password = settingsForm.old_password
      updateData.new_password = settingsForm.new_password
    }

    if (Object.keys(updateData).length > 0) {
      await updateUserProfile(updateData)
    }

    ElMessage.success('设置保存成功')
  } catch (error) {
    ElMessage.error(error.message)
    console.error('保存失败:', error)
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  fetchUserProfile()
  settingsForm.avatar_file = null
  settingsForm.old_password = null
  settingsForm.new_password = null
  settingsForm.confirm_password = null
}

onMounted(() => {
  fetchUserProfile()
})
</script>

<style scoped>
.avatar-uploader {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
.avatar-uploader :deep(.el-upload) {
  border: none;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}
.avatar-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.avatar {
  width: 200px;
  height: 200px;
  border-radius: 50%;
  object-fit: cover;
}

.v-divider {
  opacity: 0.6;
}
</style>
