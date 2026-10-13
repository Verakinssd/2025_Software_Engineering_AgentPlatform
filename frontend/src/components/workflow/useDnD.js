import { useVueFlow } from '@vue-flow/core'
import { ref, watch } from 'vue'
import { v4 as uuidv4 } from 'uuid'

const typeCounts = {}

const state = {
  draggedType: ref(null),
  isDragOver: ref(false),
  isDragging: ref(false),
}

export default function useDragAndDrop() {
  const { draggedType, isDragOver, isDragging } = state

  const { addNodes, screenToFlowCoordinate, onNodesInitialized, updateNode } = useVueFlow()

  watch(isDragging, (dragging) => {
    document.body.style.userSelect = dragging ? 'none' : ''
  })

  function onDragStart(event, type) {
    if (event.dataTransfer) {
      event.dataTransfer.setData('application/vueflow', type)
      event.dataTransfer.effectAllowed = 'move'
    }

    draggedType.value = type
    isDragging.value = true

    document.addEventListener('drop', onDragEnd)
  }

  function onDragOver(event) {
    event.preventDefault()

    if (draggedType.value) {
      isDragOver.value = true

      if (event.dataTransfer) {
        event.dataTransfer.dropEffect = 'move'
      }
    }
  }

  function onDragLeave() {
    isDragOver.value = false
  }

  function onDragEnd() {
    isDragging.value = false
    isDragOver.value = false
    draggedType.value = null
    document.removeEventListener('drop', onDragEnd)
  }

  function onDrop(event) {
    const position = screenToFlowCoordinate({
      x: event.clientX,
      y: event.clientY,
    })

    const type = draggedType.value

    if (!typeCounts[type]) {
      typeCounts[type] = 0
    }
    typeCounts[type] ++

    const nodeId = uuidv4()

    const newNode = {
      id: nodeId,
      type: draggedType.value,
      position,
      data: { 
        name: `${type==='myinput'?'输入':
          type==='myoutput'?'输出':
          type==='llm'?'大模型':
          type==='select'?'条件分支':
          type==='end'?'结束':
          type==='start'?'开始':
          type==='code'?'代码':
          type==='concat'?'拼接':
          type==='knowledge'?'知识库检索':
          type==='workflow'?'工作流':
          ''}_${typeCounts[type]}`,
        input: type === 'knowledge' ? [
          {
            value_name: 'query',
            value_type: 'string',
            value: '',
            assigned: false
          }
        ] : [],
        output: type === 'knowledge' ? [
          {
            value_name: 'result',
            value_type: 'string',
            value: '',
            assigned: false
          }
        ] : [],
        codePython:
`${type==='code' ? `def main(args: Args) -> dict:
    # 获取传入参数
    a = args.params['a']
    # 构建输出对象
    return {
        "ret": a
    }` : ''}`,
        codeJavaScript:
`${type==='code' ? `async function main({ params }: Args): Promise<o> {
  // 构建输出对象
  const ret = {
    "key0": params.input + params.input, // 拼接两次入参 input 的值
    "key1": ["hello", "world"], // 输出一个数组
    "key2": { // 输出一个Object
      "key21": "hi"
    },
  };
  return ret;
}` : ''}`,
        codeType: `${type ==='code' ? 'JavaScript' : ''}`,
        userPrompt: '',
        systemPrompt: '',
        concatPrompt: '',
        knowledgeFileId: type === 'knowledge' ? 0 : undefined,
        queryTemplate: type === 'knowledge' ? '请回答关于{{query}}的问题' : '',
        workflowId: type === 'workflow' ? 0 : undefined
      },
    }

    const { off } = onNodesInitialized(() => {
      updateNode(nodeId, (node) => ({
        position: { x: node.position.x - node.dimensions.width / 2, y: node.position.y - node.dimensions.height / 2 },
      }))

      off()
    })

    addNodes(newNode)
  }

  return {
    draggedType,
    isDragOver,
    isDragging,
    onDragStart,
    onDragLeave,
    onDragOver,
    onDrop,
  }
}