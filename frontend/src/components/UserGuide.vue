<template>
  <div></div>
</template>

<script setup>
import { driver } from 'driver.js'
import { nextTick } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const startGuide = () => {
  const driverObj = driver({
    showProgress: true,
    popoverClass: 'driverjs-theme',
    steps: [
      {
        element: '[data-guide="logo"]',
        popover: {
          title: '欢迎来到 LingXi 平台！👋',
          description: '这是您的智能体创作和交流平台，让我们通过一个快速的新手指引，帮助您了解平台的核心功能！',
          side: 'right',
          align: 'start'
        }
      },
      {
        element: '[data-guide="home"]',
        popover: {
          title: '🏠 主页 - 发现精彩内容',
          description: '在主页您可以浏览平台上最热门的智能体作品，发现其他用户创建的精彩工作流。每个卡片都代表一个可以交互的智能体！',
          side: 'right',
          align: 'start'
        }
      },
      {
        element: '[data-guide="search"]',
        popover: {
          title: '🔍 搜索功能',
          description: '使用搜索栏快速找到您感兴趣的智能体或特定类型的工作流。支持按名称、描述等多种方式搜索！',
          side: 'bottom',
          align: 'center'
        }
      },
      {
        element: '[data-guide="agents-grid"]',
        popover: {
          title: '🤖 智能体展示区',
          description: '这里展示了平台上的公开智能体。点击任意卡片可以查看详情、与智能体对话，或查看其工作流设计。让我们先试试第一个智能体！',
          side: 'top',
          align: 'center',
          onNextClick: () => {
            driverObj.destroy()
            // 改进的智能体点击逻辑
            setTimeout(() => {
              clickFirstAgent()
            }, 300)
          }
        }
      }
    ],
    onDestroyStarted: () => {
      driverObj.destroy()
    }
  })

  // 确保页面渲染完成后再开始指引
  nextTick(() => {
    driverObj.drive()
  })
}

// 在智能体页面继续指引
const continueWithAgentGuide = () => {
  const agentDriver = driver({
    showProgress: true,
    popoverClass: 'driverjs-theme',
    steps: [
      {
        popover: {
          title: '🎉 太好了！',
          description: '您现在在智能体的详情页面。这里您可以看到智能体的介绍，并且可以与它进行对话交互。'
        }
      },
      {
        element: '.right-chat',
        popover: {
          title: '💬 对话区域',
          description: '右侧是对话区域，您可以在这里与智能体进行交互对话。',
          side: 'left',
          align: 'center'
        }
      },
      {
        element: '.right-chat .input-area',
        popover: {
          title: '💭 输入框',
          description: '在这里输入您的问题或指令，与智能体进行对话。',
          side: 'top',
          align: 'center'
        }
      },
      {
        element: '[data-guide="run-workflow"]',
        popover: {
          title: '▶️ 运行工作流',
          description: '点击这个按钮可以直接运行智能体的工作流！这是体验智能体核心功能的最快方式。试试点击它吧！',
          side: 'top',
          align: 'center',
          onNextClick: () => {
            agentDriver.destroy()
            setTimeout(() => {
              continueWithWorkspaceGuide()
            }, 500)
          }
        }
      }
    ],
    onDestroyStarted: () => {
      agentDriver.destroy()
    }
  })

  agentDriver.drive()
}

// 改进的点击第一个智能体的函数
const clickFirstAgent = () => {
  let attempts = 0
  const maxAttempts = 3
  
  const tryClickAgent = () => {
    attempts++
    
    // 尝试多种选择器来找到第一个智能体卡片
    const selectors = [
      '[data-guide="agents-grid"] .v-card[to]',
      '[data-guide="agents-grid"] .v-card',
      '.agents-grid .v-card',
      '.v-row .v-card'
    ]
    
    let firstAgent = null
    let agentUrl = null
    
    for (const selector of selectors) {
      const agentCards = document.querySelectorAll(selector)
      if (agentCards.length > 0) {
        firstAgent = agentCards[0]
        // 尝试获取路由链接
        agentUrl = firstAgent.getAttribute('to') || firstAgent.querySelector('a')?.getAttribute('href')
        break
      }
    }
    
    if (firstAgent && agentUrl) {
      console.log('找到智能体卡片，准备跳转:', agentUrl)
      
      // 首先尝试程序化导航
      try {
        router.push(agentUrl).then(() => {
          console.log('程序化导航成功')
          setTimeout(() => {
            continueWithAgentGuide()
          }, 1000)
        }).catch(() => {
          // 如果程序化导航失败，尝试点击
          tryClickElement()
        })
        return
      } catch (error) {
        console.log('程序化导航失败，尝试点击')
        tryClickElement()
        return
      }
    } else if (firstAgent) {
      // 没有找到URL但找到了卡片，尝试点击
      tryClickElement()
      return
    }
    
    function tryClickElement() {
      const clickableElements = [
        firstAgent,
        firstAgent.querySelector('a'),
        firstAgent.querySelector('[role="button"]')
      ].filter(Boolean)
      
      for (const element of clickableElements) {
        try {
          element.click()
          console.log('成功点击智能体元素')
          setTimeout(() => {
            continueWithAgentGuide()
          }, 1500)
          return
        } catch (error) {
          console.log('点击失败:', error.message)
        }
      }
      
      // 所有点击尝试都失败
      retryOrSkip()
    }
    
    function retryOrSkip() {
      // 如果没有找到或点击失败，重试
      if (attempts < maxAttempts) {
        console.log(`第${attempts}次尝试失败，1秒后重试`)
        setTimeout(tryClickAgent, 1000)
      } else {
        console.log('多次尝试后仍无法点击智能体，直接跳转到工作空间')
        // 直接跳转到工作空间继续指引
        continueWithWorkspaceGuide()
      }
    }
    
    // 如果什么都没找到，直接重试或跳过
    retryOrSkip()
  }
  
  tryClickAgent()
}

// 跳转到工作空间继续指引
const continueWithWorkspaceGuide = () => {
  router.push('/workspace').then(() => {
    setTimeout(() => {
      const workspaceDriver = driver({
        showProgress: true,
        popoverClass: 'driverjs-theme',
        steps: [
          {
            element: '[data-guide="workspace-title"]',
            popover: {
              title: '🛠️ 工作空间 - 您的创作天地',
              description: '现在我们来到了工作空间！这里是您创建和管理智能体的地方。您可以设计工作流、配置智能体行为。',
              side: 'bottom',
              align: 'start'
            }
          },
          {
            element: '[data-guide="create-workflow"]',
            popover: {
              title: '✨ 创建您的第一个工作流',
              description: '点击这个按钮可以创建新的智能体工作流！您可以：\n• 设计对话逻辑\n• 配置AI模型参数\n• 发布您的作品供他人使用\n\n建议您尝试创建一个简单的工作流并发布到社区！',
              side: 'bottom',
              align: 'center'
            }
          }
        ],
        onDestroyStarted: () => {
          workspaceDriver.destroy()
        },
        onDestroyed: () => {
          continueWithKnowledgeGuide()
        }
      })
      
      workspaceDriver.drive()
    }, 500)
  })
}

// 最后指引知识库功能
const continueWithKnowledgeGuide = () => {
  router.push('/knowledge').then(() => {
    setTimeout(() => {
      const knowledgeDriver = driver({
        showProgress: true,
        popoverClass: 'driverjs-theme',
        steps: [
          {
            element: '[data-guide="knowledge-title"]',
            popover: {
              title: '📚 知识库 - 智能体的大脑',
              description: '知识库是智能体的重要组成部分！在这里您可以上传文档、资料，为您的智能体提供专业知识。',
              side: 'bottom',
              align: 'start'
            }
          },
          {
            element: '[data-guide="upload-file"]',
            popover: {
              title: '📝 上传您的知识文件',
              description: '点击这个按钮上传知识文件！支持多种格式：\n• PDF文档\n• Word文档\n• 文本文件\n• Markdown文件\n• PPT文件\n• Excel文件\n• zip压缩包文件\n\n上传后系统会自动处理，让您的智能体拥有这些知识！',
              side: 'bottom',
              align: 'center'
            }
          },
          {
            element: '[data-guide="profile"]',
            popover: {
              title: '👤 个人中心',
              description: '最后，个人中心是管理您账户信息、查看创作历史的地方。',
              side: 'right',
              align: 'start'
            }
          },
          {
            popover: {
              title: '🎊 新手指引完成！',
              description: '恭喜您完成了LingXi平台的新手指引！\n\n现在您可以：\n✅ 浏览和体验智能体\n✅ 创建自己的工作流\n✅ 上传知识库文件\n✅ 与社区互动\n\n开始您的AI创作之旅吧！'
            }
          }
        ],
        onDestroyStarted: () => {
          knowledgeDriver.destroy()
        }
      })
      
      knowledgeDriver.drive()
    }, 500)
  })
}

// 导出startGuide方法供外部使用
defineExpose({
  startGuide
})
</script>

<style>
/* Driver.js 自定义样式 */
.driverjs-theme .driver-popover {
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
  border: 1px solid #e0e0e0;
  min-width: 320px;
  max-width: 400px;
}

/* 重置所有按钮的重影效果 */
.driverjs-theme .driver-popover button {
  box-shadow: none !important;
  text-shadow: none !important;
  filter: none !important;
}

.driverjs-theme .driver-popover-title {
  font-size: 18px;
  font-weight: 600;
  color: #1976d2;
  margin-bottom: 8px;
  line-height: 1.4;
}

.driverjs-theme .driver-popover-description {
  font-size: 14px;
  color: #424242;
  line-height: 1.6;
  white-space: pre-line;
}

.driverjs-theme .driver-popover-footer {
  margin-top: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.driverjs-theme .driver-popover-next-btn {
  background: #1976d2;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  box-shadow: none !important;
  text-shadow: none !important;
  filter: none !important;
}

.driverjs-theme .driver-popover-next-btn:hover {
  background: #1565c0;
  box-shadow: none !important;
  text-shadow: none !important;
  filter: none !important;
}

.driverjs-theme .driver-popover-done-btn {
  background: #1976d2;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  box-shadow: none !important;
  text-shadow: none !important;
}

.driverjs-theme .driver-popover-done-btn:hover {
  background: #1565c0;
  box-shadow: none !important;
}

.driverjs-theme .driver-popover-prev-btn {
  background: transparent;
  color: #1976d2;
  border: 1px solid #1976d2;
  border-radius: 6px;
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  margin-right: 8px;
  transition: all 0.2s;
  box-shadow: none !important;
  text-shadow: none !important;
  filter: none !important;
}

.driverjs-theme .driver-popover-prev-btn:hover {
  background: #1976d2;
  color: white;
  box-shadow: none !important;
  text-shadow: none !important;
  filter: none !important;
}

.driverjs-theme .driver-popover-close-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  background: none;
  border: none;
  font-size: 20px;
  color: #999;
  cursor: pointer;
  line-height: 1;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.driverjs-theme .driver-popover-close-btn:hover {
  color: #666;
  background: #f5f5f5;
  border-radius: 50%;
}

.driverjs-theme .driver-popover-progress-text {
  font-size: 12px;
  color: #666;
  margin-left: auto;
  font-weight: 500;
}

/* 高亮效果增强 */
.driver-highlighted-element {
  box-shadow: 0 0 0 4px rgba(25, 118, 210, 0.3), 0 0 0 8px rgba(25, 118, 210, 0.1) !important;
  border-radius: 8px !important;
}
</style> 