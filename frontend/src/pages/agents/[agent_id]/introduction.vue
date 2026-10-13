<template>
  <v-btn
      class="back-button"
      icon="mdi-chevron-left"
      variant="flat"
      color="white"
      @click="goBack"
  ></v-btn>

  <div class="action-sidebar">
    <div class="action-buttons">
      <div class="action-item" @click="handlePostLike">
        <FontAwesomeIcon :icon="isLiked ? ['fas', 'heart'] : ['far', 'heart']"
                         :color="isLiked ? 'red' : ''"
                         class="action-icon"/>
        <span class="action-count">{{ post.likes }}</span>
      </div>

      <div class="action-item" @click="handleCollect">
        <FontAwesomeIcon :icon="isCollected ? ['fas', 'star'] : ['far', 'star']"
                         :color="isCollected ? 'gold' : ''"
                         class="action-icon"/>
        <span class="action-count">{{ post.collections }}</span>
      </div>

      <div class="action-item" @click="scrollToComment">
        <v-icon>mdi-comment-processing-outline</v-icon>
        <span class="action-count">{{ comments.length }}</span>
      </div>
    </div>
  </div>
  <div class="bottom-back-to-top">
    <v-btn
        variant="flat"
        color="white"
        class="back-to-top"
        @click="scrollToTop"
    >
      <div class="back-to-top-content">
        <v-icon color="grey-darken-1">mdi-chevron-up</v-icon>
        <span class="back-to-top-text">回到顶部</span>
      </div>
    </v-btn>
  </div>

  <div class="main-container">
    <div class="left-content">
      <v-container class="post-container">
        <v-card flat class="pa-6">
          <div class="d-flex justify-space-between align-center mb-6">
            <div class="d-flex align-center">
              <v-avatar size="50" class="mr-3">
                <v-img :src="post.author.avatar" :lazy-src="'/images/default_avatar.jpg'"></v-img>
              </v-avatar>
              <div ref="leftScrollContainer">
                <div class="text-subtitle-1 font-weight-medium">{{ post.author.name }}</div>
                <div class="text-caption text-medium-emphasis">{{ formatTime(post.createdAt) }}</div>
              </div>
            </div>
            <v-menu>
              <template v-slot:activator="{ props }">
                <v-btn
                    icon="mdi-dots-vertical"
                    variant="text"
                    v-bind="props"
                ></v-btn>
              </template>
              <v-list>
                <v-list-item
                    value="workflow"
                    @click="goToWorkflow"
                >
                  <v-list-item-title>查看工作流</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </div>

          <v-img
              :src="post.cover"
              :lazy-src="'/images/default_cover.jpg'"
              cover
              class="rounded-lg mb-6"
              max-height="400"
          ></v-img>

          <div class="mb-6">
            <h2 class="text-h4 mb-4">{{ post.title }}</h2>
            <div class="markdown-body" v-html="renderMarkdown(post.description)"></div>
          </div>

          <v-divider class="my-6"></v-divider>
          <h3 class="text-h6 mb-4">评论 ({{ comments.length }})</h3>

          <div class="d-flex mb-6" ref="commentInput">
            <v-avatar size="40" class="mr-3">
              <v-img :src="user.avatar !== '' ? user.avatar : '/images/default_avatar.jpg'" :lazy-src="'/images/default_avatar.jpg'"></v-img>
            </v-avatar>
            <v-textarea
                v-model="newComment"
                label="写下你的评论..."
                variant="outlined"
                rows="2"
                auto-grow
                hide-details
                class="flex-grow-1 mr-2"
            ></v-textarea>
            <v-btn
                color="white"
                @click="submitComment"
                height="40"
                class="align-self-end"
            >发送
            </v-btn>
          </div>

          <CommentItem
            v-for="comment in organizedComments"
            :key="comment.id"
            :comment="comment"
            :is-root="true"
            :render-markdown="renderMarkdown"
            :format-time="formatTime"
            :can-delete="canDelete"
            :active-reply="activeReply"
            :reply-contents="replyContents"
            @toggle-reply="toggleReply"
            @submit-reply="submitReply"
            @confirm-delete="confirmDelete"
            @handle-like="handleCommentLike"
          />

        </v-card>
        <v-dialog v-model="deleteDialog" max-width="400">
          <v-card>
            <v-card-title class="text-h6">确认删除</v-card-title>
            <v-card-text>确定要删除这条评论吗？</v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
              <v-btn color="error" variant="text" @click="confirmDeleteAction">删除</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-container>
    </div>

    <div class="right-chat">
      <v-card height="100%" flat class="chat-wrapper">
        <v-toolbar color="white" density="compact" class="chat-header">
          <v-toolbar-title>{{ post.title }}</v-toolbar-title>
          <v-spacer></v-spacer>
        </v-toolbar>

        <v-card-text class="chat-container-wrapper">
          <div ref="chatContainer" class="chat-container">
            <div class="message-item ai-message">
              <v-avatar size="36" color="primary" class="message-avatar">
                <v-img
                  :src="post.cover"
                ></v-img>
              </v-avatar>
              <div class="message-content">
                <div class="message-name">{{ post.title }}</div>
                <div class="message-bubble ai-bubble">
                  <div class="message-text markdown-body" v-html="renderMarkdown(post.description)"></div>
                  <div class="message-time">{{ new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) }}</div>
                </div>
              </div>
            </div>

            <div
              v-for="(message, index) in chatMessages"
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
              <v-img
                v-else
                :src="post.cover || '/images/default_cover.jpg'"
              ></v-img>
            </v-avatar>
              <div class="message-content">
                <div class="message-name">{{ message.role === 'user' ? user.username : post.title }}</div>
                <div class="message-bubble" :class="message.role === 'user' ? 'user-bubble' : 'ai-bubble'">
                  <div class="message-text markdown-body" v-html="renderMarkdown(message.content)"></div>
                  <div class="message-time">{{ message.time }}</div>
                </div>
              </div>
            </div>

            <div v-if="loading" class="message-item ai-message">
              <v-avatar size="36" color="primary" class="message-avatar">
                <v-img
                  :src="post.cover"
                ></v-img>
              </v-avatar>
              <div class="message-content">
                <div class="message-name">{{ post.title }}</div>
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

        <v-btn
          variant="text"
          data-guide="run-workflow"
          @click="() => { inputMessage = '/run-workflow'; sendMessage(); } "
        >运行工作流</v-btn>

        <div class="input-area-wrapper">
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
        </div>
      </v-card>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref, onMounted, computed, nextTick} from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import server from '@/api/server';
import { user } from "@/stores/user"
import CommentItem from '@/components/comment/CommentItem.vue'
import {marked} from "marked";

interface Author {
  id: string
  name: string
  avatar: string
}

interface Comment {
  id: string
  content: string
  likes: number
  isLiked?: boolean
  createdAt: Date
  author: Author
  parentId?: string
}

interface CommentWithReplies extends Comment {
  replies?: Comment[]
}

const router = useRouter()
const route = useRoute()

const { agent_id } = route.params

const post = ref({
  id: '',
  title: '',
  description: '',
  cover: '',
  likes: 0,
  collections: 0,
  comments: 0,
  createdAt: new Date(),
  author: {
    id: '',
    name: '',
    avatar: ''
  }
})

// 左边社区功能部分
const isLiked = ref(false)
const isCollected = ref(false)
const comments = ref<Comment[]>([])
const deleteDialog = ref(false)
const deletingCommentId = ref('')
const showCommentButton = ref(false)
const newComment = ref('')
const activeReply = ref('')
const replyContents = ref<Record<string, string>>({})

const confirmDelete = (id: string) => {
  deletingCommentId.value = id
  deleteDialog.value = true
}

const confirmDeleteAction = async () => {
  await deleteComment(deletingCommentId.value)
  deleteDialog.value = false
}

const organizedComments = computed(() => {
  type Node = CommentWithReplies & { replyToUserName?: string }

  const map = new Map<string, Node>()
  comments.value.forEach(c => {
    map.set(c.id, { ...c, replies: [] })
  })

  const roots: Node[] = []

  map.forEach(node => {
    if (node.parentId) {
      const parent = map.get(node.parentId)
      if (parent) {
        if (node.replyTo && node.replyTo !== node.parentId) {
          const target = map.get(node.replyTo)
          if (target) {
            node.replyToUserName = target.author.name
          }
        }
        parent.replies!.push(node)
      } else {
        roots.push(node)
      }
    } else {
      roots.push(node)
    }
  })

  function sortTree(nodes: Node[]) {
    nodes.sort((a, b) =>
      new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
    )
    nodes.forEach(n => {
      if (n.replies && n.replies.length) {
        sortTree(n.replies)
      }
    })
  }
  sortTree(roots)

  return roots
})

const formatTime = (date: Date) => {
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})
}

const goToWorkflow = () => {
  router.push(`/agents/${agent_id}/workflow`)
  router.push({
    path: `/agents/${agent_id}/workflow`,
  })
}

const handleCollect = async () => {
  try {
    const response = await server.Post(`/workflow/${agent_id}/favorite`, {
      workflow_id: agent_id,  // Pass workflow_id in the request body
    });

    if (response.data.is_favorite !== undefined) {
      if (response.data.is_favorite === true) {
        isCollected.value = true;
        post.value.collections = post.value.collections + 1;
        ElMessage.success('已收藏');
      } else {
        isCollected.value = false;
        post.value.collections = post.value.collections - 1;
        ElMessage.success('已取消收藏');
      }
    }
  } catch (error) {
    console.error('收藏失败', error);
    ElMessage.error('收藏失败');
  }
};

const handlePostLike = async () => {
  try {
    const response = await server.Post(`/workflow/${agent_id}/like`, {
      workflow_id: agent_id,  // Pass workflow_id in the request body
    });

    if (response.data.is_liked !== undefined) {
      if (response.data.is_liked === true) {
        isLiked.value = true;
        post.value.likes = post.value.likes + 1;
        ElMessage.success('已点赞');
      } else {
        isLiked.value = false;
        post.value.likes = post.value.likes - 1;
        ElMessage.success('已取消点赞');
      }
    }
  } catch (error) {
    console.error('点赞失败', error);
    ElMessage.error('点赞失败');
  }
};

const handleCommentLike = async (target: Comment) => {
  try {
    const response = await server.Post(`/workflow/comment/${target.id}/like`);
    console.log(response.data)
    if (response.data.like_result !== undefined) {
      const index = comments.value.findIndex(c => c.id === target.id);
      if (index !== -1) {
        const comment = comments.value[index];
        comment.isLiked = response.data.like_result;
        comment.likes += comment.isLiked ? 1 : -1;
        ElMessage.success(comment.isLiked ? '已点赞' : '已取消点赞');
      }
    }
  } catch (error) {
    console.error('评论点赞失败', error);
    ElMessage.error('评论点赞失败');
  }
};

const submitComment = async (parentId: number | null = null, replyTo: number | null = null) => {
  if (!newComment.value.trim()) {
    ElMessage.warning('评论内容不能为空');
    return;
  }

  if (!agent_id) {
    ElMessage.error('工作流 ID 无效');
    return;
  }

  // 创建 FormData 对象
  const formData = new FormData();
  formData.append('comment_type', '1');
  formData.append('content', newComment.value.trim());
  formData.append('workflow_id', agent_id as string);

  console.log('发送的表单数据:', Object.fromEntries(formData));

  try {
    const response = await server.Post('/workflow/comment', formData);

    console.log('后端返回:', response);

    if (response.data && response.data.ID) {
      const newCommentResponse = response.data;

      const comment = {
        id: newCommentResponse.ID,
        content: newCommentResponse.Content,
        likes: newCommentResponse.Likes,
        isLiked: newCommentResponse.LikeStatus,
        createdAt: new Date(newCommentResponse.CreatedAt),
        author: {
          id: newCommentResponse.CommenterID,
          name: user.username,
          avatar: user.avatar !== '' ? user.avatar : '/images/default_avatar.jpg',
        },
        parentId: newCommentResponse.ParentID || null,
        replies: newCommentResponse.Children || [],
      };

      comments.value.unshift(comment); // 将评论添加到顶部
      post.value.comments += 1; // 增加评论数

      newComment.value = ''; // 清空评论框
      showCommentButton.value = false;
      ElMessage.success('评论发布成功');

      try {
        await server.Post('/notification/send', {
          title: `@${user.username} 评论了你的工作流 ${post.value.title}`,
          content: newCommentResponse.Content,
          to_user_id: post.value.author.id
        });
      } catch (notificationError) {
        console.error('发送评论通知失败:', notificationError);
      }

    } else {
      console.log('后端返回的数据无 ID', response.data);
      ElMessage.error(`评论发布失败，后端返回数据无效: ${JSON.stringify(response.data)}`);
    }
  } catch (error) {
    console.error('发布评论失败', error);
    ElMessage.error('评论发布失败');
  }
};

const toggleReply = (commentId: string) => {
  if (activeReply.value === commentId) {
    activeReply.value = ''
  } else {
    activeReply.value = commentId
    replyContents.value[commentId] = ''
  }
}

const submitReply = async (commentId: string) => {
  const replyContent = replyContents.value[commentId]?.trim()
  if (!replyContent) {
    return ElMessage.warning('回复内容不能为空')
  }

  let rootId = commentId
  let current = comments.value.find(c => c.id === commentId)
  while (current && current.parentId) {
    current = comments.value.find(c => c.id === current.parentId)
    if (current) {
      rootId = current.id
    }
  }

  const fd = new FormData()
  fd.append('comment_type', '2')
  fd.append('content', replyContent)
  fd.append('parent_id', rootId)
  fd.append('reply_to', commentId)

  try {
    const { data } = await server.Post('/workflow/comment', fd)
    if (!data?.ID) throw new Error()

    const newReply: Comment = {
      id: data.ID.toString(),
      content: data.Content,
      likes: data.Likes ?? 0,
      isLiked: data.LikeStatus ?? false,
      createdAt: new Date(data.CreatedAt),
      author: {
        id: data.CommenterID.toString(),
        name: user.username,
        avatar: user.avatar || '/images/default_avatar.jpg'
      },
      parentId: rootId,
      replyTo: commentId
    }

    comments.value.push(newReply)
    post.value.comments += 1
    replyContents.value[commentId] = ''
    activeReply.value = ''
    ElMessage.success('回复发布成功')

    try {
      const repliedComment = comments.value.find(c => c.id === commentId)
      console.log("hi: ", current.author.id)
      if (repliedComment) {
        await server.Post('/notification/send', {
          title: `@${user.username} 回复了你的评论`,
          content: data.Content,
          to_user_id: Number(current.author.id)
        })
      }
    } catch (notificationError) {
      console.error('发送回复通知失败:', notificationError)
    }
  } catch {
    ElMessage.error('回复发布失败')
  }
}

const fetchPostData = async () => {
  try {
    const response = await server.Get(`/workflow/${agent_id}`);
    console.log('Fetched workflow data:', response);

    const workflow = response.data;
    post.value = {
      id: workflow.id,
      title: workflow.name,
      description: workflow.description,
      cover: workflow.cover || '/images/default_cover.jpg',
      likes: workflow.likes,
      collections: workflow.collects,
      comments: 0,
      createdAt: new Date(workflow.createdAt),
      author: {
        id: workflow.belongsTo.id,
        name: workflow.belongsTo.username,
        avatar: workflow.belongsTo.avatar || '/images/default_avatar.jpg',
      }
    };
    isLiked.value = workflow.likeStatus;
    isCollected.value = workflow.collectStatus;
  } catch (error) {
    console.error('获取工作流数据失败', error);
    ElMessage.error('加载工作流数据失败');
  }
}

const fetchComments = async (order?: number, page?: number, pageSize?: number) => {
  try {
    const response = await server.Get(`/workflow/${agent_id}/comment`, {
      params: { order, page, pageSize }
    });
    const commentData = response.data;
    const flatComments: Comment[] = [];

    function processNode(node: any, parentId?: string) {
      flatComments.push({
        id: node.ID.toString(),
        content: node.Content,
        likes: node.Likes ?? 0,
        isLiked: node.LikeStatus ?? false,
        createdAt: new Date(node.CreatedAt),
        author: node.Commenter
          ? {
              id: node.Commenter.id.toString(),
              name: node.Commenter.username,
              avatar: node.Commenter.avatar || '/images/default_avatar.jpg'
            }
          : {
              id: '',
              name: '已删除用户',
              avatar: '/images/default_avatar.jpg'
            },
        parentId: parentId,
        replyTo: node.ReplyTo != null ? node.ReplyTo.toString() : parentId
      });

      if (Array.isArray(node.Children)) {
        node.Children.forEach((child: any) =>
          processNode(child, node.ID.toString())
        );
      }
    }

    if (Array.isArray(commentData)) {
      commentData.forEach((root: any) => processNode(root, undefined));
    }

    comments.value = flatComments;
  } catch (error) {
    console.error('获取评论失败', error);
    ElMessage.error('加载评论数据失败');
  }
};

function canDelete(comment: Comment) {
  return comment.author.id.toString() === user.id.toString()
      || post.value.author.id.toString() === user.id.toString()
}
async function deleteComment(id: string) {
  try {
    await server.Post(`/workflow/comment/${id}/delete`)
    const toRemove = new Set<string>([id])
    comments.value.forEach(c => {
      if (c.parentId === id) toRemove.add(c.id)
    })
    comments.value = comments.value.filter(c => !toRemove.has(c.id))
    post.value.comments = Math.max(0, post.value.comments - toRemove.size)
    comments.value = [...comments.value]
    ElMessage.success('删除成功')
  } catch {
    ElMessage.error('删除失败，请稍后重试')
  }
}

marked.use({ async: false });

const renderMarkdown = (content: string) => {
  if (!content) return ''
  return marked.parse(content).trim()
}

// 滚轮跳转
const commentInput = ref<HTMLElement | null>(null);
const leftScrollContainer = ref<HTMLElement | null>(null);

const scrollToComment = () => {
  if (commentInput.value) {
    commentInput.value.scrollIntoView({ behavior: 'smooth' });
  }
};

const scrollToTop = () => {
  if (leftScrollContainer.value) {
    leftScrollContainer.value.scrollIntoView({ behavior: 'smooth' });
  }
};

// chat功能部分
const chatMessages = ref<Array<{role: string, content: string, time: string}>>([])
const inputMessage = ref('')
const loading = ref(false)
const chatContainer = ref<HTMLElement | null>(null)
const isProcessing = ref(false)

const showCommandHint = ref(false);
const messageInput = ref<HTMLInputElement | null>(null);

const insertCommand = (command: string) => {
  inputMessage.value = command;
  showCommandHint.value = false;
  nextTick(() => {
    messageInput.value?.focus();
  });
};

const goBack = () => {
  router.push('/');
};

const scrollToBottom = () => {
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight + 10000;
  }
};

// const sendMessage = async () => {
//   if (!inputMessage.value.trim()) return
//
//   const userMessage = {
//     role: 'user',
//     content: inputMessage.value,
//     time: new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})
//   }
//   chatMessages.value.push(userMessage)
//
//   const userInput = inputMessage.value
//   inputMessage.value = ''
//   loading.value = true
//   isProcessing.value = true
//
//   // 模拟AI回复（用户输入重复3遍）
//   setTimeout(() => {
//     const aiMessage = {
//       role: 'assistant',
//       content: `${userInput} ${userInput} ${userInput}`,
//       time: new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})
//     }
//     chatMessages.value.push(aiMessage)
//     loading.value = false
//     isProcessing.value = false
//
//     nextTick(() => {
//       if (chatContainer.value) {
//         chatContainer.value.scrollTop = chatContainer.value.scrollHeight
//       }
//     })
//   }, 2000)
// }

const workflowParams = ref({});
const currentParamIndex = ref(0);
const isCollectingParams = ref(false);
import {nodes, edges, runWorkflow, runWorkflowWithInputs, saveWorkflowContent, fetchWorkflowContent} from '@/api/workflow'
const startNode = ref()
const endNode = ref()
const running = ref(false)
const finishRunning = ref(false)

const sendMessage = async () => {
  if (inputMessage.value.trim() === '/run-workflow') {
    chatMessages.value.push({
      role: 'user',
      content: inputMessage.value,
      time: new Date().toLocaleTimeString()
    })
    await new Promise(resolve => setTimeout(resolve, 100));
    scrollToBottom();
    startParamCollection();
    inputMessage.value = '';
    return;
  }

  if (!inputMessage.value.trim()) return;

  const userMessage = {
    role: 'user',
    content: inputMessage.value,
    time: new Date().toLocaleTimeString()
  }
  chatMessages.value.push(userMessage);
  await new Promise(resolve => setTimeout(resolve, 100))
  scrollToBottom();

  if (isCollectingParams.value) {
    const paramName = startNode.value.data.input[currentParamIndex.value].value_name;
    workflowParams.value = {
      ...workflowParams.value,
      [paramName]: inputMessage.value
    };
    inputMessage.value = '';

    if (currentParamIndex.value < startNode.value.data.input.length - 1) {
      currentParamIndex.value++;
      askForNextParam();
    } else {
      await runWithCollectedParams();
    }
  }
  else {
    const aiMessage = {
      role: 'assistant',
      content: '暂时没有聊天功能哦~运行工作流请输入"/run-workflow"',
      time: new Date().toLocaleTimeString()
    };
    chatMessages.value.push(aiMessage);
    await new Promise(resolve => setTimeout(resolve, 100));
    scrollToBottom();
  }

  inputMessage.value = '';
}

const askForNextParam = async () => {
    const nextParam = startNode.value.data.input[currentParamIndex.value];
    const aiMessage = {
        role: 'assistant',
        content: `请输入 ${nextParam.value_name} 的值 (类型: ${simplifyType(nextParam.value_type)})`,
        time: new Date().toLocaleTimeString()
    };
    chatMessages.value.push(aiMessage);
    await new Promise(resolve => setTimeout(resolve, 100))
    scrollToBottom();
};

const runWithCollectedParams = async () => {
  loading.value = true;

  await new Promise(resolve => setTimeout(resolve, 200))
  scrollToBottom();

  await runWholeWorkflow();

  loading.value = false;

  endNode.value.data.output.forEach(output => {
    chatMessages.value.push({
      role: 'assistant',
      content: `${output.value_name}: ${output.value}\n`,
      time: new Date().toLocaleTimeString()
    });
  });
  await new Promise(resolve => setTimeout(resolve, 100))
  scrollToBottom();

  isCollectingParams.value = false;
  workflowParams.value = {};
  currentParamIndex.value = 0;
};

const startParamCollection = async () => {
    if (!startNode.value?.data?.input?.length) {
        chatMessages.value.push({
            role: 'assistant',
            content: '该工作流不需要输入参数',
            time: new Date().toLocaleTimeString()
        });
        await new Promise(resolve => setTimeout(resolve, 100))
        scrollToBottom();
        return;
    }

    isCollectingParams.value = true;
    askForNextParam();
};

const simplifyType = (type) => {
  const typeMap = {
    'string': 'str',
    'number': 'num',
    'boolean': 'bool',
    'integer': 'int'
  };
  return typeMap[type] || type;
};

const runWholeWorkflow = async () => {
  running.value = true;
  finishRunning.value = false;

  try {
    console.log('runWorkflowWithInputs' , workflowParams.value)
    // 直接使用收集到的参数，无需保存工作流
    const res = await runWorkflowWithInputs(Number(agent_id), workflowParams.value);

    if (res.status === 200) {
      const { workflow } = res.data;
      const { content } = workflow;

      //这一堆确保引用更新
      nodes.value = [...content.nodes];
      edges.value = [...(content.edges || [])];
      startNode.value = {...nodes.value.find(node => node.id == 1)};
      endNode.value = {...nodes.value.find(node => node.id == 2)};
      nodes.value = content.nodes;
      edges.value = content.edges || [];

      endNode.value = nodes.value.find(node => node.id == 2);
      return true;
    } else {
      throw new Error('运行失败');
    }
  } catch (error) {
    chatMessages.value.push({
      role: 'assistant',
      content: '工作流运行失败: ' + error.message,
      time: new Date().toLocaleTimeString()
    });
    await new Promise(resolve => setTimeout(resolve, 100));
    scrollToBottom();
    return false;
  } finally {
    running.value = false;
    finishRunning.value = true;
  }
};

const handleEnter = (e) => {
  if (!isProcessing.value && !loading.value && inputMessage.value.trim()) {
    sendMessage();
  } else {
    e.preventDefault();
  }
};

onMounted(async () => {
  fetchPostData()
  fetchComments(1, 1, 100)

  const res = await fetchWorkflowContent(agent_id)
  if (res) {
    const { content } = res
    nodes.value = content.nodes
    edges.value = content.edges || []
  }
  startNode.value = nodes.value.find(node => node.id == 1)
  endNode.value = nodes.value.find(node => node.id == 2)
})
</script>

<style scoped>
/* 整体布局 */
.main-container {
  display: flex;
  width: 100%;
  max-width: 1800px;
  margin: 0 auto;
  padding: 20px;
  gap: 30px;
  height: 100vh;
}

.left-content {
  flex: 5;
  min-width: 0;
  overflow-y: auto;
  height: 100%;
  position: relative;
}

.left-content::-webkit-scrollbar {
    display: none;
}

.left-content {
    -ms-overflow-style: none;
    scrollbar-width: none;
}

.right-chat {
  flex: 3;
  min-width: 0;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.post-container {
  max-width: 100%;
  margin: 0;
  position: relative;
}

/* 左边功能栏 */
.back-button {
  position: fixed;
  top: 40px;
  left: 100px;
  right: 40px;
  width: 50px;
  height: 50px;
  min-width: 0;
  padding: 0;
  border-radius: 4px;
  z-index: 1000;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  &:hover {
    transform: translateX(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }
}

.action-sidebar {
  position: fixed;
  left: 100px;
  top: 40%;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  z-index: 1;
}

.action-buttons {
  width: 64px;
  display: flex;
  flex-direction: column;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 12px 8px;
}

.bottom-back-to-top {
  position: fixed;
  bottom: 20px;
  left: 100px;
  z-index: 100;
}

.bottom-back-to-top.v-show {
  opacity: 1;
  transform: translateY(0);
}

.back-to-top {
  bottom: 20px;
  margin-top: 64px;
  height: auto;
  min-height: 64px;
  min-width: unset;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 8px 0;
}

.back-to-top-content {
  width: 64px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4px 0;
}

.back-to-top-text {
  font-size: 10px;
  color: rgba(0, 0, 0, 0.6);
  margin-top: 4px;
  line-height: 1.2;
}


.back-to-top {
  bottom: 20px;
  margin-top: 64px;
  height: auto;
  min-width: unset;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 8px 0;
}

.back-to-top-content {
  width: 64px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4px 0;
}

.v-card-title {
  padding: 16px 24px 12px;
  font-weight: 500;
}

.v-card-actions {
  padding: 8px 16px 16px;
}

.comment-input-wrapper {
  transition: all 0.3s ease;
}

/* 右边chat框 */
.chat-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  background-color: #fafafa;
}

.chat-header {
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.chat-container-wrapper {
  flex: 1;
  padding: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background-color: transparent;
}

.chat-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background-color: transparent;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.input-area {
  padding: 12px 16px;
  background-color: white;
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

.post-container {
  max-width: 800px;
  margin: 0 auto;
  position: relative;
}

.action-icon {
  font-size: 20px !important;
  width: 20px !important;
  height: 20px !important;
  margin-bottom: 4px;
}


.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-item:hover {
  background: rgba(0, 0, 0, 0.02);
}

.action-item .v-icon {
  font-size: 24px;
  margin-bottom: 4px;
}

.action-count {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.6);
}


.bottom-back-to-top.v-show {
  opacity: 1;
  transform: translateY(0);
}

.v-textarea textarea {
  min-height: 40px;
}

.v-textarea textarea {
  padding-bottom: 20px !important;
}

.markdown-body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
  font-size: 16px;
  line-height: 1.5;
  word-wrap: break-word;
  white-space: pre-wrap;
  overflow-wrap: break-word;
  word-break: break-word;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3),
.markdown-body :deep(h4),
.markdown-body :deep(h5),
.markdown-body :deep(h6) {
  margin-top: 24px;
  margin-bottom: 4px;
  font-weight: 600;
  line-height: 1.25;
}

.markdown-body :deep(p) {
  margin-top: 0;
  margin-bottom: 4px;
}

.markdown-body :deep(a) {
  color: #1976D2;
  text-decoration: none;
}

.markdown-body :deep(a:hover) {
  text-decoration: underline;
}

.markdown-body :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27, 31, 35, 0.05);
  border-radius: 3px;
}

.markdown-body :deep(pre) {
  word-wrap: normal;
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 3px;
}

.markdown-body :deep(blockquote) {
  padding: 0 1em;
  color: #6a737d;
  border-left: 0.25em solid #dfe2e5;
  margin: 0 0 4px 0;
}

.markdown-body :deep(img) {
  max-width: 100%;
  box-sizing: initial;
  background-color: #fff;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 2em;
  margin-top: 0;
  margin-bottom: 4px;
}

.markdown-body :deep(li) {
  margin-bottom: 4px;
  line-height: 1.5;
}

.markdown-body :deep(ul ul),
.markdown-body :deep(ol ol),
.markdown-body :deep(ul ol),
.markdown-body :deep(ol ul) {
  padding-left: 1.5em;
}

.markdown-body :deep(ul) {
  list-style-type: disc;
}

.markdown-body :deep(ol) {
  list-style-type: decimal;
}
</style>
