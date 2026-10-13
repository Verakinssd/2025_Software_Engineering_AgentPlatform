<template>
  <div class="register-bg">
    <div class="register-container">
      <el-card class="register-card">
        <h2 class="register-title">灵犀AI智能体社区</h2>
        <div class="login-link">
          <span class="login-hint">已有帐户？</span>
          <span class="login-text" @click="router.push('/login')">登录</span>
        </div>
        <el-form
          ref="ruleFormRef"
          :rules="rules"
          :model="registerForm"
          label-width="auto"
          :label-position="labelPosition"
          :size="size"
          @submit.prevent
        >
          <el-form-item label="用户名" prop="username" class="form-label">
            <el-input v-model="registerForm.username" @keyup.enter="submitRegister"/>
          </el-form-item>
          <el-form-item label="密码" prop="password" class="form-label">
            <el-input
              v-model="registerForm.password"
              :type="showPassword ? 'text' : 'password'"
              @keyup.enter="submitRegister"
            >
              <template #suffix>
                <el-icon
                  class="password-icon"
                  @click="showPassword = !showPassword"
                >
                  <View v-if="showPassword" />
                  <Hide v-else />
                </el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="确认密码" prop="confirmPassword" class="form-label">
            <el-input
              v-model="registerForm.confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              @keyup.enter="submitRegister"
            >
              <template #suffix>
                <el-icon
                  class="password-icon"
                  @click="showConfirmPassword = !showConfirmPassword"
                >
                  <View v-if="showConfirmPassword" />
                  <Hide v-else />
                </el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="邮箱" prop="email" class="form-label">
            <el-input v-model="registerForm.email" @keyup.enter="submitRegister"/>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="submitRegister()"
              :loading="isloading"
              class="register-btn"
            >
              <template #loading>
                <div class="custom-loading">
                  <svg class="circular" viewBox="-10, -10, 50, 50">
                    <path
                      class="path"
                      d="
                      M 30 15
                      L 28 17
                      M 25.61 25.61
                      A 15 15, 0, 0, 1, 15 30
                      A 15 15, 0, 1, 1, 27.99 7.5
                      L 15 15
                      "
                      style="stroke-width: 4px; fill: rgba(0, 0, 0, 0)"
                    />
                  </svg>
                </div>
              </template>
              注册
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<style scoped>
.register-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-image: url('/images/login_bg.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  z-index: 0;
}

.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.register-card {
  position: relative;
  width: 400px;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  background-color: rgba(255, 255, 255, 0.9);
}

.register-title {
  color: #454546;
  margin-bottom: 30px;
  font-size: 24px;
  font-weight: 600;
}

.login-link {
  position: absolute;
  top: 30px;
  right: 30px;
  font-size: 14px;
}

.login-hint {
  color: #606266;
  font-weight: bold;
}

.login-text {
  color: #e6658f;
  cursor: pointer;
  font-weight: bold;
  margin-left: 5px;
}

.login-text:hover {
  text-decoration: underline;
}

.form-label :deep(.el-form-item__label) {
  font-weight: bold;
}

.register-btn {
  position: relative;
  padding-left: 12px;
  padding-right: 12px;

  border-radius: 20px;
  width: 150px;
  margin: 10px auto 0;
  display: block;
  background-color: #56afe2;
  border-color: #56afe2;
  color: white;
  font-size: 15px;
  font-weight: 100;
  letter-spacing: 2px;
}


.circular {
  width: 16px;
  height: 16px;
  vertical-align: middle;
  animation: rotate 1s linear infinite;
}

.custom-loading {
  display: inline-block;
  width: 24px;
  margin-right: 8px;
  vertical-align: middle;
}


.register-btn:hover {
  background-color: #4fa3d1;
  border-color: #4fa3d1;
}

.register-btn:active {
  background-color: #3498db;
  border-color:#3498db;
}

.password-icon {
  cursor: pointer;
  color: #606266;
  font-size: 18px;
}

.el-button .custom-loading .circular {
  margin-right: 6px;
  width: 18px;
  height: 18px;
  animation: loading-rotate 2s linear infinite;
}

.el-button .custom-loading .circular .path {
  animation: loading-dash 1.5s ease-in-out infinite;
  stroke-dasharray: 90, 150;
  stroke-dashoffset: 0;
  stroke-width: 2;
  stroke: var(--el-button-text-color);
  stroke-linecap: round;
}
</style>

<script lang="ts" setup>
import {reactive ,ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from  'vue-router'
import { View, Hide } from '@element-plus/icons-vue'
import type {
  ComponentSize ,
  FormProps ,
  FormRules ,
  FormInstance
} from 'element-plus'
import { register } from '@/api/user'
const router = useRouter()
const size = ref<ComponentSize>('large');
const labelPosition = ref<FormProps['labelPosition']>('top');
const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  email: ''
})
const isloading = ref(false)

const ruleFormRef = ref<FormInstance>()

const showPassword = ref(false)
const showConfirmPassword = ref(false)

const validateUserName = (rule : any , value: any , callback: any) => {
  if (value === '' ) {
    callback(new Error('请输入用户名'))
  }
}

const validatePass = (rule: any , value :any , callback : any) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  }
}

const validatePass2 = (rule: any , value :any , callback : any) => {
  if (value === '') {
    callback(new Error('请输入确认密码'))
  } else if (value !== registerForm.password) {
    callback(new Error('两次密码输入不一致'))
  }
}

const validateEmail = (rule: any , value :any , callback : any) => {
  if (value === '') {
    callback(new Error('请输入邮箱'))
  } else if (!/\S+@\S+\.\S+/.test(value)) {
    callback(new Error('邮箱格式不正确'))
  }
}

const rules = reactive<FormRules<typeof registerForm>> ({
  username: [ {validator : validateUserName , trigger : 'blur' }],
  password: [ {validator : validatePass , trigger: 'blur'} ],
  confirmPassword: [ {validator : validatePass2 , trigger: 'blur'} ],
  email: [ {validator : validateEmail , trigger: 'blur'} ]
})

const submitRegister = async () => {
  if (registerForm.username === '' || registerForm.password === '' || registerForm.email === ''
    || registerForm.confirmPassword === '') {
    ElMessage.error('请填写完整信息')
    return
  }
  else if (registerForm.password !== registerForm.confirmPassword) {
    ElMessage.error('两次密码输入不一致')
    return
  }
  else if (!/\S+@\S+\.\S+/.test(registerForm.email)) {
    ElMessage.error('邮箱格式不正确')
    return
  }
  isloading.value = true

  // await new Promise(resolve => setTimeout(resolve, 5000))

  const res = await register(registerForm.username , registerForm.password , registerForm.email)
  console.log("res.status = " + res.status)
  if (res.status === 200) {
    ElMessage.success('注册成功')
    router.push('/login')
  } else {
    const { response } = res
    const { status, data } = response as { status: number; data: any }
    ElMessage.error(data.message)
    isloading.value = false
  }
  console.log(res)
}
</script>
