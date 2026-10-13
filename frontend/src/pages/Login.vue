<template>
  <div class="login-bg">
    <div class="login-container">
      <el-card class="login-card">
        <h2 class="login-title">灵犀AI智能体社区</h2>
        <div class="register-link">
          <span class="register-hint">没有帐户？</span>
          <span class="register-text" @click="router.push('/register')">注册</span>
        </div>
        <el-form
          ref="ruleFormRef"
          :rules="rules"
          :model="loginForm"
          label-width="auto"
          :label-position="labelPosition"
          :size="size"
          @submit.prevent
        >
          <el-form-item label="用户名" prop="username" class="form-label">
            <el-input v-model="loginForm.username"
                      ref="usernameInputRef"
                      @keydown.down="focusPasswordInput" />
          </el-form-item>
          <el-form-item label="密码" prop="password" class="form-label">
            <el-input
              v-model="loginForm.password"
              :type="showPassword ? 'text' : 'password'"
              ref="passwordInputRef"
              @keyup.enter="submitLogin"
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
          <el-form-item>
            <el-button type="primary" @click="submitLogin()" :loading="isloading" class="login-btn">
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
              登录
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<style scoped>
.login-bg {
  width: 100vw;
  height: 100vh;
  background-image: url('/images/login_bg.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  z-index: 1;
}

.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.login-card {
  position: relative;
  width: 400px;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  background-color: rgba(255, 255, 255, 0.9);
}

.login-title {
  color: #454546;
  margin-bottom: 30px;
  font-size: 24px;
  font-weight: 600;
}

.register-link {
  position: absolute;
  top: 30px;
  right: 30px;
  font-size: 14px;
}

.register-hint {
  color: #606266;
  font-weight: bold;
}

.register-text {
  color: #e6658f;
  cursor: pointer;
  font-weight: bold;
  margin-left: 5px;
}

.register-text:hover {
  text-decoration: underline;
}

.form-label :deep(.el-form-item__label) {
  font-weight: bold;
}

.login-btn {
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


.login-btn:hover {
  background-color: #4fa3d1;
  border-color: #4fa3d1;
}

.login-btn:active {
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
import {reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import type {
  ComponentSize,
  FormProps,
  FormRules,
  FormInstance
} from 'element-plus'
import { login } from '@/api/user'
import { View, Hide } from '@element-plus/icons-vue'

const router = useRouter()
const size = ref<ComponentSize>('large');
const labelPosition = ref<FormProps['labelPosition']>('top');
const loginForm = reactive({
  username: '',
  password: '',
})
const isloading = ref(false)
const showPassword = ref(false)

const ruleFormRef = ref<FormInstance>()

const usernameInputRef = ref<HTMLInputElement>()
const passwordInputRef = ref<HTMLInputElement>()

const focusPasswordInput = () => {
  passwordInputRef.value?.focus()
}

const validateUserName = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请输入用户名'))
  }
}

const validatePass = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  }
}

const rules = reactive<FormRules<typeof loginForm>>({
  username: [{validator: validateUserName, trigger: 'blur'}],
  password: [{validator: validatePass, trigger: 'blur'}],
})

const submitLogin = async () => {
  if (loginForm.username === '' || loginForm.password === '') {
    ElMessage.error('用户名或密码不能为空')
    return
  }
  isloading.value = true

  // await new Promise(resolve => setTimeout(resolve, 5000))

  const res = await login(loginForm.username, loginForm.password)
  //console.log("res.status = " + res.status)
  if (res.status === 200) {
    //ElMessage.success('登录成功')
    router.push('/')
  } else {
    ElMessage.error('登录失败，请检查用户名或密码')
    isloading.value = false
  }
  //console.log(res)
}
</script>


