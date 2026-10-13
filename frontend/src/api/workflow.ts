import { ref } from "vue";
import server from "./server";
import { ElMessage } from "element-plus";

export interface InputItem {
  value_name: string; 
  value_type: 'string' | 'integer' | 'number' | 'boolean'; 
  value: string | number | boolean; 
  assigned? : boolean;
  assignedNodeId? : string;
  assignedNodeType? : string;
  assignedNodeName? : string;
  assignedAttributeName? : string;
  assignedAttributeType? : string;
}

export type OutputItem = InputItem;

export interface NodeData {
  name: string;
  input?: InputItem[];
  output?: OutputItem[];
  codePython?: string;
  codeJavaScript?: string;
  codeType?: string;
  userPrompt?: string;
  systemPrompt?: string;
  concatPrompt?: string;
  knowledgeFileId?: number;
  queryTemplate?: string;
  workflowId?: number;
}

export interface FlowNode {
  id: string;
  type: string;
  position: { x: number; y: number };
  data: NodeData;
}

export const showCodeTest = ref(false);
export const showAgentTest = ref(false);
export const showConcatTest = ref(false);
export const showKnowledgeTest = ref(false);
export const showWorkflowTest = ref(false);

export const selectedNode = ref<FlowNode | null>(null);

export function getAttributes(node: FlowNode) {
  if (node.type === 'myinput') {
    return node.data.input;
  } else if (node.type === 'myoutput') {
    return node.data.output;
  } else if (node.type === 'start') {
    return node.data.input;
  } else if (node.type === 'llm') {
    return node.data.output;
  } else if (node.type === 'code') {
    return node.data.output;
  } else if (node.type === 'concat') {
    return node.data.output;
  } else if (node.type === 'knowledge') {
    return node.data.output;
  } else if (node.type === 'workflow') {
    return node.data.output;
  } else {
    return [];
  }
}

export function simplifyType(type: string) {
  switch (type) {
    case "string":
      return "str";
    case "integer":
      return "int";
    case "number":
      return "num";
    case "boolean":
      return "bool";
    default:
      return type;
  } 
}

export const saveWorkflowContent = async (workflowId: number, content: any) => {
  return await server.Post(`/workflow/${workflowId}/content`, content);
}

export const fetchWorkflowContent = async (workflowId: number) => {
  const res = await server.Get(`/workflow/${workflowId}`);
  if (res.status === 200) {
    return res.data;
  } else {
    ElMessage.error("获取工作流内容失败");
    return null;
  }
}

export const runWorkflow = async (workflowId: number, endNodeID?: string) => {
  if (endNodeID) {
    return await server.Post(`/workflow/${workflowId}/run?node_id=${endNodeID}`);
  } else {
    return await server.Post(`/workflow/${workflowId}/run`);
  }
}

export const runWorkflowWithInputs = async (
  workflowId: number, 
  inputs: Record<string, any>, 
  endNodeID?: string
) => {
  const url = endNodeID 
    ? `/workflow/${workflowId}/run_with_inputs?node_id=${endNodeID}` 
    : `/workflow/${workflowId}/run_with_inputs`;
  
  return await server.Post(url, inputs);
};

export const ForkWorkflow = async (workflowId: number) => {
  return await server.Post(`/workflow/${workflowId}/fork`);
}

export const nodes = ref([])
export const edges = ref([])

// 获取可用工作流列表（个人工作流和公开工作流）
export const getAvailableWorkflows = async () => {
  return await server.Get('/workflow/available');
};

// 获取公开工作流列表
export const getPublicWorkflows = async () => {
  return await server.Get('/workflow/public2');
};

// 获取个人工作流列表
export const getPersonalWorkflows = async () => {
  return await server.Get('/workflow/personal');
};

// 获取工作流的输入输出结构
export const getWorkflowIOStructure = async (workflowId: number) => {
  return await server.Get(`/workflow/${workflowId}/io`);
};

// 根据ID获取知识库文件名
export const getKnowledgeFileName = async (fileId: number) => {
  try {
    const res = await server.Get(`/knowledge/file/${fileId}`);
    if (res.status === 200) {
      return res.data.fileName;
    }
  } catch (error) {
    console.error('获取知识库文件名失败:', error);
  }
  return null;
};

// 根据ID获取工作流名称
export const getWorkflowName = async (workflowId: number) => {
  try {
    const res = await server.Get(`/workflow/${workflowId}`);
    if (res.status === 200) {
      return res.data.name;
    }
  } catch (error) {
    console.error('获取工作流名称失败:', error);
  }
  return null;
};