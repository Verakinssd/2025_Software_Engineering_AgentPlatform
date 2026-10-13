<template>
  <v-app>
    <v-main>
      <v-container>
        <v-btn
          color="primary"
          @click="drawer = true"
          class="open-chat-btn"
          variant="flat"
          size="large"
        >
          <v-icon start>mdi-forum</v-icon>
          测试
        </v-btn>

        <v-navigation-drawer
          v-model="drawer"
          temporary
          location="right"
          width="400"
        >
          <v-card height="100%" flat class="chat-wrapper">
            <v-toolbar color="primary" density="compact">
              <v-toolbar-title class="text-white">测试与预览</v-toolbar-title>
              <v-spacer></v-spacer>
              <v-btn icon @click="drawer = false">
                <v-icon color="white">mdi-close</v-icon>
              </v-btn>
            </v-toolbar>

            <v-card-text class="chat-container-wrapper">
              <div ref="chatContainer" class="chat-container">
                <div class="message-item ai-message">
                  <v-avatar size="36" color="primary" class="message-avatar">
                    <v-icon color="white">mdi-robot</v-icon>
                  </v-avatar>
                  <div class="message-content">
                    <div class="message-name">测试用LLM</div>
                    <div class="message-bubble ai-bubble">
                      <div class="message-text">你好！我是LLM助手，有什么可以帮您的吗？</div>
                      <div class="message-time">{{ new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) }}</div>
                    </div>
                  </div>
                </div>

                <div
                  v-for="(message, index) in messages"
                  :key="index"
                  :class="['message-item', message.role === 'user' ? 'user-message' : 'ai-message']"
                >
                  <v-avatar
                    size="36"
                    :color="message.role === 'user' ? 'blue-lighten-4' : 'primary'"
                    class="message-avatar"
                  >
                    <v-img
                      v-if="message.role === 'user' && user.avatar"
                      :src="user.avatar"
                    ></v-img>
                    <v-icon
                      v-else
                      color="white"
                    >
                      mdi-robot
                    </v-icon>
                  </v-avatar>
                  <div class="message-content">
                    <div class="message-name">{{ message.role === 'user' ? user.username : '测试用LLM' }}</div>
                    <div class="message-bubble" :class="message.role === 'user' ? 'user-bubble' : 'ai-bubble'">
                      <div class="message-text">{{ message.content }}</div>
                      <div class="message-time">{{ message.time }}</div>
                    </div>
                  </div>
                </div>

                <div v-if="loading" class="message-item ai-message">
                  <v-avatar size="36" color="primary" class="message-avatar">
                    <v-icon color="white">mdi-robot</v-icon>
                  </v-avatar>
                  <div class="message-content">
                    <div class="message-name">测试用LLM</div>
                    <div class="message-bubble ai-bubble">
                      <div class="message-text">
                        <span class="typing-dot"></span>
                        <span class="typing-dot"></span>
                        <span class="typing-dot"></span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </v-card-text>

            <v-divider></v-divider>

            <v-card-actions class="input-area">
              <v-text-field
                v-model="inputMessage"
                placeholder="输入问题..."
                variant="outlined"
                density="comfortable"
                hide-details
                single-line
                @keyup.enter="handleEnter"
                class="message-input"
              ></v-text-field>
              <v-btn
                color="primary"
                icon
                :disabled="!inputMessage || loading || isProcessing"
                @click="sendMessage"
                class="send-btn"
              >
                <v-icon size="32">mdi-arrow-right</v-icon>
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-navigation-drawer>
      </v-container>
    </v-main>
  </v-app>
</template>


<script setup>
import {nextTick, ref, watch} from 'vue'
import { user } from "@/stores/user"

const drawer = ref(false)
const inputMessage = ref('')
const messages = ref([])
const loading = ref(false)
const chatContainer = ref(null)
const isProcessing = ref(false);


//之前没用完的api，好像还剩十块钱，先随便用用
const DEEPSEEK_API_KEY = ""

watch([messages, drawer], async (newVal, oldVal) => {
  if (drawer.value) {
    await nextTick();
    console.log("hello")
    scrollToBottom();
  }
}, { deep: true });

const scrollToBottom = () => {
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
  }
};

const handleEnter = (e) => {
  if (!isProcessing.value && !loading.value && inputMessage.value.trim()) {
    sendMessage();
  } else {
    e.preventDefault();
  }
};

const sendMessage = async () => {
  if (!inputMessage.value.trim()) return

  const userMessage = {
    role: 'user',
    content: inputMessage.value,
    time: new Date().toLocaleTimeString()
  }
  messages.value.push(userMessage)

  const userInput = inputMessage.value
  inputMessage.value = ''
  loading.value = true
  isProcessing.value = true

  try {
    const response = await fetch('https://api.deepseek.com/v1/chat/completions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${DEEPSEEK_API_KEY}`
      },
      body: JSON.stringify({
        model: 'deepseek-chat',
        messages: [
          ...messages.value.map(msg => ({
            role: msg.role,
            content: msg.content
          }))
        ],
        temperature: 0.7,
        max_tokens: 1000
      })
    })

    if (!response.ok) {
      throw new Error(`API请求失败: ${response.status}`)
    }

    const data = await response.json()
    const aiResponse = data.choices[0].message

    const aiMessage = {
      role: 'assistant',
      content: aiResponse.content,
      displayContent: '',
      time: new Date().toLocaleTimeString(),
      tokens: data.usage?.completion_tokens
    }
    messages.value.push(aiMessage)

    await nextTick();
    scrollToBottom();
  } catch (error) {
    console.error('调用API出错:', error)
    messages.value.push({
      role: 'assistant',
      content: '抱歉，处理您的请求时出错了。请稍后再试。',
      time: new Date().toLocaleTimeString()
    })
  } finally {
    loading.value = false
    isProcessing.value = false
  }
}
</script>

<style scoped>
.rounded-drawer {
  border-top-left-radius: 12px !important;
  border-bottom-left-radius: 12px !important;
  overflow: hidden;
}

.rounded-card {
  border-radius: 0 !important;
}

.rounded-toolbar {
  border-top-left-radius: 12px !important;
}

.chat-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
}

.chat-container-wrapper {
  flex: 1;
  padding: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.chat-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.message-item {
  display: flex;
  margin-bottom: 16px;
  max-width: 100%;
}

.user-message {
  justify-content: flex-end;
}

.ai-message {
  justify-content: flex-start;
}

.message-avatar {
  flex-shrink: 0;
  margin-right: 8px;
}

.user-message .message-avatar {
  order: 1;
  margin-right: 0;
  margin-left: 8px;
}

.message-content {
  max-width: calc(100% - 60px);
  display: flex;
  flex-direction: column;
}

.user-message .message-content {
  align-items: flex-end;
}

.ai-message .message-content {
  align-items: flex-start;
}

.message-name {
  font-size: 0.75rem;
  color: #888;
  margin-bottom: 4px;
}

.message-bubble {
  padding: 10px 14px;
  border-radius: 18px;
  position: relative;
  max-width: 100%;
  word-break: break-word;
  line-height: 1.4;
}

.ai-bubble {
  background: white;
  border-top-left-radius: 4px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.user-bubble {
  background: #395ca7;
  color: white;
  border-top-right-radius: 4px;
}

.message-text {
  font-size: 0.875rem;
}

.message-time {
  font-size: 0.6875rem;
  margin-top: 4px;
  opacity: 0.8;
  text-align: right;
}

.ai-bubble .message-time {
  color: #999;
}

.user-bubble .message-time {
  color: rgba(255, 255, 255, 0.7);
}

.input-area {
  background: white;
}

.message-input {
  border-radius: 24px;
}

.send-btn {
  transform: rotate(0deg);
  transition: transform 0.2s;
}

.send-btn:hover {
  transform: rotate(-90deg);
}

.typing-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #999;
  margin: 0 2px;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-dot:nth-child(1) {
  animation-delay: 0s;
}

.typing-dot:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-dot:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.6;
  }
  30% {
    transform: translateY(-4px);
    opacity: 1;
  }
}

.chat-container::-webkit-scrollbar {
  width: 6px;
}

.chat-container::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.chat-container::-webkit-scrollbar-thumb {
  background: #d1d1d1;
  border-radius: 3px;
}
</style>
