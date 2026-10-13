<template>
  <div
    style="
      position: fixed;
      bottom: 20px;
      left: 50%;
      transform: translateX(-50%);
      width: 300px;
      height: 60px;
      background-color: #fff;
      border-radius: 12px;
      box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
      display: flex;
      justify-content: space-around;
      align-items: center;
      z-index: 1000;
    "
  >
    <v-tooltip :text="isAuthor ? '保存' : 'Fork'" location="bottom">
      <template v-slot:activator="{ props }">
        <v-btn
          class="ml-3"
          v-bind="props"
          color="primary"
          variant="text"
          :icon="isAuthor ? 'mdi-content-save' : 'mdi-source-fork'"
          :loading="!(isAuthor) && isForking"
          :disabled="isForking"
          @click="isAuthor ? saveWorkflow() : forkWorkflow()"
        />
      </template>
    </v-tooltip>

    <v-tooltip text="运行" location="bottom">
      <template v-slot:activator="{ props }">
        <v-btn
          class="ml-3"
          v-bind="props"
          color="primary"
          variant="text"
          @click="runWholeWorkflow"
          icon="mdi-play"
          :loading="running"
        ></v-btn>
      </template>
    </v-tooltip>
  </div>

  <div class="combined-editor">

    <!-- 他人工作流提示横幅 -->
    <v-alert
      v-if="!isAuthor && workflow"
      type="warning"
      variant="tonal"
      closable
      prominent
      border="start"
      class="others-workflow-alert ma-4"
      style="background-color: rgba(255, 255, 224, 0.9); border-radius: 8px; box-shadow: 0 2px 12px rgba(0, 0, 0, 0.2); backdrop-filter: blur(4px); position: fixed; top: 20px; left: 50%; transform: translateX(-50%); z-index: 2000; max-width: 800px; width: auto; min-width: 400px;"
      >
      <template v-slot:prepend>
        <v-icon color="warning">mdi-alert-circle-outline</v-icon>
      </template>
      <div class="d-flex align-center flex-wrap">
        <span class="text-subtitle-1 font-weight-medium">正在查看他人的工作流</span>
        <v-divider vertical class="mx-3"></v-divider>
        <span class="text-subtitle-1">知识库文件或引用的工作流可能无法访问，节点仅供查看无法编辑</span>
      </div>
    </v-alert>

    <v-navigation-drawer
      :mobile="false"
    >
      <v-list-item class="pa-2">
        <template #prepend>
          <v-btn
            variant="text"
            @click="() => { isAuthor ? saveWorkflow() : null; goBack() }"
            icon="mdi-arrow-left"
          ></v-btn>
        </template>
        
        <!-- 调整布局：工作流名字在上方，拥有者在下方 -->
        <div class="d-flex flex-column align-start ml-4">
          <v-list-item-title class="text-subtitle-1 font-weight-medium mb-1">
            {{ workflow?.name }}
          </v-list-item-title>
          <!-- 添加作者标识 -->
          <v-chip 
            v-if="!isAuthor && workflow" 
            size="small" 
            color="primary" 
            variant="outlined"
            class="text-caption"
          >
            <v-avatar size="14" class="mr-1">
              <v-img :src="workflow.belongsTo?.avatar || '/images/default_avatar.jpg'"></v-img>
            </v-avatar>
            {{ workflow.belongsTo?.username }}
          </v-chip>
        </div>
      </v-list-item>

      <v-expansion-panels
        v-if="isAuthor"
        multiple
        variant="accordion"
        v-model="panel"
      >
        <v-expansion-panel value="inout">
          <v-expansion-panel-title>
            <v-icon icon="mdi-playlist-edit" class="mr-2"></v-icon>
            输入/输出
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <div
              v-for="type in availableNodeTypes.filter(t => ['myinput', 'myoutput'].includes(t.name))"
              class="node-item"
              draggable="true"
              @dragstart="onDragStart($event, type.name)"
            >
              <img :src="type.icon" class="node-icon" />
              {{ type.label }}
            </div>
          </v-expansion-panel-text>
        </v-expansion-panel>

        <v-expansion-panel value="function">
          <v-expansion-panel-title>
            <v-icon icon="mdi-function" class="mr-2"></v-icon>
            功能
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <div
              v-for="type in availableNodeTypes.filter(t => ['code', 'llm', 'concat'].includes(t.name))"
              class="node-item"
              draggable="true"
              @dragstart="onDragStart($event, type.name)"
            >
              <img :src="type.icon" class="node-icon" />
              {{ type.label }}
            </div>
          </v-expansion-panel-text>
        </v-expansion-panel>

        <v-expansion-panel value="knowledge">
          <v-expansion-panel-title>
            <v-icon icon="mdi-database" class="mr-2"></v-icon>
             知识库
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <div
              v-for="type in availableNodeTypes.filter(t => ['knowledge'].includes(t.name))"
              class="node-item"
              draggable="true"
              @dragstart="onDragStart($event, type.name)"
            >
              <img :src="type.icon" class="node-icon" />
              {{ type.label }}
            </div>
          </v-expansion-panel-text>
        </v-expansion-panel>

        <v-expansion-panel value="workflow">
          <v-expansion-panel-title>
            <v-icon icon="mdi-file-tree" class="mr-2"></v-icon>
            工作流
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <div
              v-for="type in availableNodeTypes.filter(t => ['workflow'].includes(t.name))"
              class="node-item"
              draggable="true"
              @dragstart="onDragStart($event, type.name)"
            >
              <img :src="type.icon" class="node-icon" />
              {{ type.label }}
            </div>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-navigation-drawer>

    <v-container style="padding: 0; margin:0; max-width: none; width: 100%" @drop="onDrop">
      <VueFlow
        v-model:nodes="nodes"
        v-model:edges="edges"
        :class="{ dark }"
        :apply-default="false"
        @node-click="onNodeClick"
        @dragover="onDragOver"
        @dragleave="onDragLeave"
        :min-zoom="0.7"
        :max-zoom="4"
        class="basic-flow"
      >

        <Background variant="dots" gap="12" size="1" />

        <Controls position="top-left">
          <ControlButton title="Reset Transform" @click="resetTransform">
            <Icon name="reset" />
          </ControlButton>
        </Controls>

        <InteractionControls v-if="showInteractionControls"/>

        <MiniMap pannable zoomable/>
        <DropzoneBackground
          :style="{
            backgroundColor: isDragOver
              ? dark
                ? '#2c3e50'
                : '#e7f3ff'
              : 'transparent',
            transition: 'background-color 0.2s ease'
          }"
        />

        <template #node-start="{ id, data, selected }">
          <div
            style="position: relative; width: 250px; overflow: visible;"
            @mouseover = "hoverId = id"
            @mouseleave = "hoverId = null"
          >
          <Handle
              type="source"
              position="right"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#90CAF9',
                border: '1px solid #90CAF9',
                top: '50%',
                right: 0,
                transform: 'translate(50%, -50%)'
              }"
            />
          <v-card
            :style="{
              width: '100%' ,
              border: (selected)?
              '1.5px solid #90CAF9'
              : '1px solid #90CAF9',
              transition: 'box-shadow 0.2s ease',
              boxShadow: (hoverId === id || selected)
              ? '0 0 0 0.5px #90CAF9' : 'none'
              }"
            flat
          >
            <v-card-title class="d-flex align-center">
              <v-avatar size="30" style="margin-right: 10px;">
                <v-img :src="startIcon" />
              </v-avatar>
              <span class="text-subtitle-2">开始</span>
            </v-card-title>
            <v-card-text class="px-3 py-1 text-subtitle-2 align-center">
              <p>输入</p>
              <span class="ml-3 font-weight-light" v-if="!data.input.length">未配置输入</span>
              <div class="mt-2">
                <v-chip
                  v-for="(attr, i) in data.input" :key="i"
                    class="ma-1"
                    :color="attr.value? 'blue' : 'grey'"
                    size="small"
                    variant="outlined"
                    :style="{ borderColor: '#CE93D8' }"
                  >
                    <strong> {{ simplifyType(attr.value_type) }}.
                    </strong>{{ attr.value_name }}
                    <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                </v-chip>
              </div>
            </v-card-text>
          </v-card>
          </div>
        </template>

        <template #node-end="{ id, data , selected }">
          <div
            style="position: relative; width: 250px; overflow: visible;"
            @mouseover = "hoverId = id"
            @mouseleave = "hoverId = null"
          >
            <Handle
              type="target"
              position="left"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#F48FB1',
                border: '1px solid #F48FB1',
                top: '50%',
                left: 0,
                transform: 'translate(-50%, -50%)'
              }"
            />
            <v-card
              :style="{
                width: '100%' ,
                border: (selected)?
                '1.5px solid #F48FB1'
                : '1px solid #F48FB1',
                transition: 'box-shadow 0.2s ease',
                boxShadow: (hoverId === id || selected)
                ? '0 0 0 0.5px #F48FB1' : 'none'
                }"
              flat
            >
              <v-card-title class="d-flex align-center">
                <v-avatar size="30" style="margin-right: 10px;">
                  <v-img :src="endIcon" />
                </v-avatar>
                <span class="text-subtitle-2">结束</span>
              </v-card-title>
              <v-card-text class="px-3 py-1 text-subtitle-2 align-center">
                <p>输出</p>
                <span class="ml-3 font-weight-light" v-if="!data.output.length">未配置输出</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.output" :key="i"
                      class="ma-1"
                      :color="attr.value? 'blue' : 'grey'"
                      size="small"
                      variant="outlined"
                      :style="{ borderColor: '#CE93D8' }"
                    >
                      <strong> {{ simplifyType(attr.value_type) }}.
                      </strong>{{ attr.value_name }}
                      <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
            </v-card>
          </div>
        </template>

        <template #node-myinput="{ id, data , selected}">
          <div
            style="position: relative; width: 250px; overflow: visible;"
            @mouseover = "hoverId = id"
            @mouseleave = "hoverId = null"
          >
            <Handle
              type="target"
              position="left"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#CE93D8',
                border: '1px solid #CE93D8',
                top: '50%',
                left: 0,
                transform: 'translate(-50%, -50%)'
              }"
            />
            <Handle
              type="source"
              position="right"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#CE93D8',
                border: '1px solid #CE93D8',
                top: '50%',
                right: 0,
                transform: 'translate(50%, -50%)'
              }"
            />
            <v-card
              :style="{
                width: '100%' ,
                border: (selected)?
                '1.5px solid #CE93D8'
                : '1px solid #CE93D8',
                transition: 'box-shadow 0.2s ease',
                boxShadow: (hoverId === id || selected)
                ? '0 0 0 0.5px #CE93D8' : 'none'
                }"
              flat
            >
              <v-card-title class="d-flex align-center">
                <v-avatar size="30" class="mr-2">
                  <v-img :src="inputIcon" />
                </v-avatar>

                <span
                  class="text-subtitle-2"
                  v-if="!editingName || id !== editingName_id"
                  style="overflow: hidden;"
                  >
                  {{ data.name }}
                </span>
                <v-text-field
                  v-else
                  ref = "nameInputRef"
                  variant="outlined"
                  density="compact"
                  hide-details
                  class="text-subtitle-2"
                  v-model="data.name"
                  @blur = "finishEditingName"
                  @keyup.enter = "finishEditingName"
                  style="
                    flex: 1 1 auto;
                    min-width: 0;
                  "
                >
                </v-text-field>
                <v-spacer></v-spacer>
                <v-menu>
                  <template v-slot:activator="{props}">
                    <v-btn
                      icon="mdi-dots-vertical"
                      v-bind="props"
                      variant="plain"
                      size="x-small"
                    >
                    </v-btn>
                  </template>
                  <v-list density="compact">
                    <v-list-item
                      v-for="(item, i) in nodeModifyItems"
                      :key="i"
                      :value="i"
                      @click="
                        selectedNode=nodes.find((node) => node.id === id);
                        handleModify(i)"
                      >
                      <v-list-item-title>{{ item.title }}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
              </v-card-title>
              <v-card-text class="px-3 py-1 text-subtitle-2 align-center">
                <p>输入</p>
                <span class="ml-3 font-weight-light" v-if="!data.input.length">未配置输入</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.input" :key="i"
                      class="ma-1"
                      :color="attr.value? 'blue' : 'grey'"
                      size="small"
                      variant="outlined"
                      :style="{ borderColor: '#CE93D8' }"
                    >
                      <strong> {{ simplifyType(attr.value_type) }}.
                      </strong>{{ attr.value_name }}
                      <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
            </v-card>
          </div>
        </template>

        <template #node-myoutput="{ id, data , selected}">
          <div
            style="position: relative; width: 250px; overflow: visible;"
            @mouseover = "hoverId = id"
            @mouseleave = "hoverId = null"
          >
          <Handle
              type="target"
              position="left"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#B39DDB',
                border: '1px solid #B39DDB',
                top: '50%',
                left: 0,
                transform: 'translate(-50%, -50%)'
              }"
            />
            <Handle
              type="source"
              position="right"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#B39DDB',
                border: '1px solid #B39DDB',
                top: '50%',
                right: 0,
                transform: 'translate(50%, -50%)'
              }"
            />
            <v-card
              :style="{
                width: '100%' ,
                border: (selected)?
                '1.5px solid #B39DDB'
                : '1px solid #B39DDB',
                transition: 'box-shadow 0.2s ease',
                boxShadow: (hoverId === id || selected)
                ? '0 0 0 0.5px #B39DDB' : 'none'
                }"
              flat
            >
              <v-card-title class="d-flex align-center">
                <v-avatar size="30" style="margin-right: 10px;">
                  <v-img :src="outputIcon" />
                </v-avatar>
                <span
                  class="text-subtitle-2"
                  v-if="!editingName || id !== editingName_id"
                  style="overflow: hidden;"
                  >
                  {{ data.name }}
                </span>
                <v-text-field
                  v-else
                  ref = "nameOutputRef"
                  variant="outlined"
                  density="compact"
                  hide-details
                  class="text-subtitle-2"
                  v-model="data.name"
                  @blur = "finishEditingName"
                  @keyup.enter = "finishEditingName"
                  style="
                    flex: 1 1 auto;
                    min-width: 0;
                  "
                >
                </v-text-field>
                <v-spacer></v-spacer>
                <v-menu>
                  <template v-slot:activator="{props}">
                    <v-btn
                      icon="mdi-dots-vertical"
                      v-bind="props"
                      variant="plain"
                      size="x-small"
                    >
                    </v-btn>
                  </template>
                  <v-list density="compact">
                    <v-list-item
                      v-for="(item, i) in nodeModifyItems"
                      :key="i"
                      :value="i"
                      @click="
                        selectedNode=nodes.find((node) => node.id === id);
                        handleModify(i)"
                      >
                      <v-list-item-title>{{ item.title }}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
              </v-card-title>
              <v-card-text class="px-3 py-1 text-subtitle-2 align-center">
                <p>输出</p>
                <span class="ml-3 font-weight-light" v-if="!data.output.length">未配置输出</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.output" :key="i"
                      class="ma-1"
                      :color="attr.value? 'blue' : 'grey'"
                      size="small"
                      variant="outlined"
                      :style="{ borderColor: '#CE93D8' }"
                    >
                      <strong> {{ simplifyType(attr.value_type) }}.
                      </strong>{{ attr.value_name }}
                      <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
            </v-card>
          </div>
        </template>

        <template #node-llm="{ id, data , selected}">
          <div
            style="position: relative; width: 250px; overflow: visible;"
            @mouseover = "hoverId = id"
            @mouseleave = "hoverId = null"
          >
            <Handle
              type="target"
              position="left"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#9FA8DA',
                border: '1px solid #9FA8DA',
                top: '50%',
                left: 0,
                transform: 'translate(-50%, -50%)'
              }"
            />
            <Handle
              type="source"
              position="right"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#9FA8DA',
                border: '1px solid #9FA8DA',
                top: '50%',
                right: 0,
                transform: 'translate(50%, -50%)'
              }"
            />
            <v-card
              :style="{
                width: '100%' ,
                border: (selected)?
                '1.5px solid #9FA8DA'
                : '1px solid #9FA8DA',
                transition: 'box-shadow 0.2s ease',
                boxShadow: (hoverId === id || selected)
                ? '0 0 0 0.5px #9FA8DA' : 'none'
                }"
              flat
            >
              <v-card-title class="d-flex align-center">
                <v-avatar size="30" style="margin-right: 10px;">
                  <v-img :src="agentIcon" />
                </v-avatar>
                <span
                  class="text-subtitle-2"
                  v-if="!editingName || id !== editingName_id"
                  style="overflow: hidden;"
                  >
                  {{ data.name }}
                </span>
                <v-text-field
                  v-else
                  ref = "nameAgentRef"
                  variant="outlined"
                  density="compact"
                  hide-details
                  class="text-subtitle-2"
                  v-model="data.name"
                  @blur = "finishEditingName"
                  @keyup.enter = "finishEditingName"
                  style="
                    flex: 1 1 auto;
                    min-width: 0;
                  "
                >
                </v-text-field>
                <v-spacer></v-spacer>
                <v-btn 
                  icon="mdi-play" 
                  variant="plain" 
                  size="x-small"
                  @click="testAgent()"
                >
                </v-btn>
                <v-menu>
                  <template v-slot:activator="{props}">
                    <v-btn
                      icon="mdi-dots-vertical"
                      v-bind="props"
                      variant="plain"
                      size="x-small"
                    >
                    </v-btn>
                  </template>
                  <v-list density="compact">
                    <v-list-item
                      v-for="(item, i) in nodeModifyItems"
                      :key="i"
                      :value="i"
                      @click="
                        selectedNode=nodes.find((node) => node.id === id);
                        handleModify(i)"
                      >
                      <v-list-item-title>{{ item.title }}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
              </v-card-title>
              <v-card-text class="px-3 text-subtitle-2 align-center">
                <p>输入</p>
                <span class="ml-3 font-weight-light" v-if="!data.input.length">未配置输入</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.input" :key="i"
                      class="ma-1"
                      :color="attr.value? 'blue' : 'grey'"
                      size="small"
                      variant="outlined"
                      :style="{ borderColor: '#CE93D8' }"
                    >
                      <strong> {{ simplifyType(attr.value_type) }}.
                      </strong>{{ attr.value_name }}
                      <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
              <v-card-text class="px-3 py-1 text-subtitle-2 align-center">
                <p>输出</p>
                <span class="ml-3 font-weight-light" v-if="!data.output.length">未配置输出</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.output" :key="i"
                      class="ma-1"
                      :color="attr.value? 'blue' : 'grey'"
                      size="small"
                      variant="outlined"
                      :style="{ borderColor: '#CE93D8' }"
                    >
                      <strong> {{ simplifyType(attr.value_type) }}.
                      </strong>{{ attr.value_name }}
                      <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
            </v-card>
          </div>
        </template>

        <template #node-code="{ id, data , selected }">
          <div
            style="position: relative; width: 250px; overflow: visible;"
            @mouseover = "hoverId = id"
            @mouseleave = "hoverId = null"
          >
            <Handle
              type="target"
              position="left"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#80DEEA',
                border: '1px solid #80DEEA',
                top: '50%',
                left: 0,
                transform: 'translate(-50%, -50%)'
              }"
            />
            <Handle
              type="source"
              position="right"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#80DEEA',
                border: '1px solid #80DEEA',
                top: '50%',
                right: 0,
                transform: 'translate(50%, -50%)'
              }"
            />
            <v-card
              :style="{
                width: '100%' ,
                border: (selected)?
                '1.5px solid #80DEEA'
                : '1px solid #80DEEA',
                transition: 'box-shadow 0.2s ease',
                boxShadow: (hoverId === id || selected)
                ? '0 0 0 0.5px #80DEEA' : 'none'
                }"
              flat
            >
              <v-card-title class="d-flex align-center">
                <v-avatar size="30" style="margin-right: 10px;">
                  <v-img :src="codeIcon" />
                </v-avatar>
                <span
                  class="text-subtitle-2"
                  v-if="!editingName || id !== editingName_id"
                  style="overflow: hidden;"
                  >
                  {{ data.name }}
                </span>
                <v-text-field
                  v-else
                  ref = "nameCodeRef"
                  variant="outlined"
                  density="compact"
                  hide-details
                  class="text-subtitle-2"
                  v-model="data.name"
                  @blur = "finishEditingName"
                  @keyup.enter = "finishEditingName"
                  style="
                    flex: 1 1 auto;
                    min-width: 0;
                  "
                >
                </v-text-field>
                <v-spacer></v-spacer>
                <v-btn 
                  icon="mdi-play" 
                  variant="plain" 
                  size="x-small"
                  @click="testCode()"
                >
                </v-btn>
                <v-menu>
                  <template v-slot:activator="{props}">
                    <v-btn
                      icon="mdi-dots-vertical"
                      v-bind="props"
                      variant="plain"
                      size="x-small"
                    >
                    </v-btn>
                  </template>
                  <v-list density="compact">
                    <v-list-item
                      v-for="(item, i) in nodeModifyItems"
                      :key="i"
                      :value="i"
                      @click="
                        selectedNode=nodes.find((node) => node.id === id);
                        handleModify(i)"
                      >
                      <v-list-item-title>{{ item.title }}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
              </v-card-title>
              <v-card-text class="px-3 text-subtitle-2 align-center">
                <p>输入</p>
                <span class="ml-3 font-weight-light" v-if="!data.input.length">未配置输入</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.input" :key="i"
                      class="ma-1"
                      :color="attr.value? 'blue' : 'grey'"
                      size="small"
                      variant="outlined"
                      :style="{ borderColor: '#CE93D8' }"
                    >
                      <strong> {{ simplifyType(attr.value_type) }}.
                      </strong>{{ attr.value_name }}
                      <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
              <v-card-text class="px-3 py-1 text-subtitle-2 align-center">
                <p>输出</p>
                <span class="ml-3 font-weight-light" v-if="!data.output.length">未配置输出</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.output" :key="i"
                      class="ma-1"
                      :color="attr.value? 'blue' : 'grey'"
                      size="small"
                      variant="outlined"
                      :style="{ borderColor: '#CE93D8' }"
                    >
                      <strong> {{ simplifyType(attr.value_type) }}.
                      </strong>{{ attr.value_name }}
                      <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
            </v-card>
          </div>
        </template>

        <template #node-concat="{ id, data , selected }">
          <div
            style="position: relative; width: 250px; overflow: visible;"
            @mouseover = "hoverId = id"
            @mouseleave = "hoverId = null"
          >
            <Handle
              type="target"
              position="left"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#80CBC4',
                border: '1px solid #80CBC4',
                top: '50%',
                left: 0,
                transform: 'translate(-50%, -50%)'
              }"
            />
            <Handle
              type="source"
              position="right"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#80CBC4',
                border: '1px solid #80CBC4',
                top: '50%',
                right: 0,
                transform: 'translate(50%, -50%)'
              }"
            />
            <v-card
              :style="{
                width: '100%' ,
                border: (selected)?
                '1.5px solid #80CBC4'
                : '1px solid #80CBC4',
                transition: 'box-shadow 0.2s ease',
                boxShadow: (hoverId === id || selected)
                ? '0 0 0 0.5px #80CBC4' : 'none'
                }"
              flat
            >
              <v-card-title class="d-flex align-center">
                <v-avatar size="30" style="margin-right: 10px;">
                  <v-img :src="concatIcon" />
                </v-avatar>
                <span
                  class="text-subtitle-2"
                  v-if="!editingName || id !== editingName_id"
                  style="overflow: hidden;"
                  >
                  {{ data.name }}
                </span>
                <v-text-field
                  v-else
                  ref = "nameConcatRef"
                  variant="outlined"
                  density="compact"
                  hide-details
                  class="text-subtitle-2"
                  v-model="data.name"
                  @blur = "finishEditingName"
                  @keyup.enter = "finishEditingName"
                  style="
                    flex: 1 1 auto;
                    min-width: 0;
                  "
                >
                </v-text-field>
                <v-spacer></v-spacer>
                <v-btn 
                  icon="mdi-play" 
                  variant="plain" 
                  size="x-small"
                  @click="testConcat()"
                >
                </v-btn>
                <v-menu>
                  <template v-slot:activator="{props}">
                    <v-btn
                      icon="mdi-dots-vertical"
                      v-bind="props"
                      variant="plain"
                      size="x-small"
                    >
                    </v-btn>
                  </template>
                  <v-list density="compact">
                    <v-list-item
                      v-for="(item, i) in nodeModifyItems"
                      :key="i"
                      :value="i"
                      @click="
                        selectedNode=nodes.find((node) => node.id === id);
                        handleModify(i)"
                      >
                      <v-list-item-title>{{ item.title }}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
              </v-card-title>
              <v-card-text class="px-3 text-subtitle-2 align-center">
                <p>输入</p>
                <span class="ml-3 font-weight-light" v-if="!data.input.length">未配置输入</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.input" :key="i"
                      class="ma-1"
                      :color="attr.value? 'blue' : 'grey'"
                      size="small"
                      variant="outlined"
                      :style="{ borderColor: '#CE93D8' }"
                    >
                      <strong> {{ simplifyType(attr.value_type) }}.
                      </strong>{{ attr.value_name }}
                      <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
              <v-card-text class="px-3 py-1 text-subtitle-2 align-center">
                <p>输出</p>
                <span class="ml-3 font-weight-light" v-if="!data.output.length">未配置输出</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.output" :key="i"
                      class="ma-1"
                      :color="attr.value? 'blue' : 'grey'"
                      size="small"
                      variant="outlined"
                      :style="{ borderColor: '#CE93D8' }"
                    >
                      <strong> {{ simplifyType(attr.value_type) }}.
                      </strong>{{ attr.value_name }}
                      <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
            </v-card>
          </div>
        </template>

        <template #node-knowledge="{ id, data, selected }">
          <div
            style="position: relative; width: 250px; overflow: visible;"
            @mouseover="hoverId = id"
            @mouseleave="hoverId = null"
          >
            <Handle
              type="target"
              position="left"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#EF9A9A',
                border: '1px solid #EF9A9A',
                top: '50%',
                left: 0,
                transform: 'translate(-50%, -50%)'
              }"
            />
            <Handle
              type="source"
              position="right"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#EF9A9A',
                border: '1px solid #EF9A9A',
                top: '50%',
                right: 0,
                transform: 'translate(50%, -50%)'
              }"
            />
            <v-card
              :style="{
                width: '100%',
                border: (selected) ?
                '1.5px solid #EF9A9A'
                : '1px solid #EF9A9A',
                transition: 'box-shadow 0.2s ease',
                boxShadow: (hoverId === id || selected)
                ? '0 0 0 0.5px #EF9A9A' : 'none'
              }"
              flat
            >
              <v-card-title class="d-flex align-center">
                <v-avatar size="30" style="margin-right: 10px;">
                  <v-img :src="knowledgeIcon" />
                </v-avatar>
                <span
                  class="text-subtitle-2"
                  v-if="!editingName || id !== editingName_id"
                  style="overflow: hidden;"
                >
                  {{ data.name }}
                </span>
                <v-text-field
                  v-else
                  ref="nameKnowledgeRef"
                  variant="outlined"
                  density="compact"
                  hide-details
                  class="text-subtitle-2"
                  v-model="data.name"
                  @blur="finishEditingName"
                  @keyup.enter="finishEditingName"
                  style="
                    flex: 1 1 auto;
                    min-width: 0;
                  "
                >
                </v-text-field>
                <v-spacer></v-spacer>
                <v-btn 
                  icon="mdi-play" 
                  variant="plain" 
                  size="x-small"
                  @click="testKnowledge()"
                >
                </v-btn>
                <v-menu>
                  <template v-slot:activator="{props}">
                    <v-btn
                      icon="mdi-dots-vertical"
                      v-bind="props"
                      variant="plain"
                      size="x-small"
                    >
                    </v-btn>
                  </template>
                  <v-list density="compact">
                    <v-list-item
                      v-for="(item, i) in nodeModifyItems"
                      :key="i"
                      :value="i"
                      @click="
                        selectedNode=nodes.find((node) => node.id === id);
                        handleModify(i)"
                    >
                      <v-list-item-title>{{ item.title }}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
              </v-card-title>
              <v-card-text class="px-3 text-subtitle-2 align-center">
                <p>输入</p>
                <span class="ml-3 font-weight-light" v-if="!data.input.length">未配置输入</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.input" :key="i"
                    class="ma-1"
                    :color="attr.value? 'blue' : 'grey'"
                    size="small"
                    variant="outlined"
                    :style="{ borderColor: '#CE93D8' }"
                  >
                    <strong>{{ simplifyType(attr.value_type) }}.
                    </strong>{{ attr.value_name }}
                    <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
              <v-card-text class="px-3 py-1 text-subtitle-2 align-center">
                <p>输出</p>
                <span class="ml-3 font-weight-light" v-if="!data.output.length">未配置输出</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.output" :key="i"
                    class="ma-1"
                    :color="attr.value? 'blue' : 'grey'"
                    size="small"
                    variant="outlined"
                    :style="{ borderColor: '#CE93D8' }"
                  >
                    <strong>{{ simplifyType(attr.value_type) }}.
                    </strong>{{ attr.value_name }}
                    <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
              <v-card-text v-if="data.knowledgeFileId" class="px-3 py-1 text-subtitle-2 align-center">
                <v-chip size="small" color="primary" variant="outlined">
                  <v-icon size="small" class="mr-1">mdi-file-document-outline</v-icon>
                  知识库文件: {{ getKnowledgeFileDisplayName(data.knowledgeFileId) }}
                </v-chip>
              </v-card-text>
            </v-card>
          </div>
        </template>

        <template #node-workflow="{ id, data, selected }">
          <div
            style="position: relative; width: 250px; overflow: visible;"
            @mouseover="hoverId = id"
            @mouseleave="hoverId = null"
          >
            <Handle
              type="target"
              position="left"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#5C6BC0',
                border: '1px solid #5C6BC0',
                top: '50%',
                left: 0,
                transform: 'translate(-50%, -50%)'
              }"
            />
            <Handle
              type="source"
              position="right"
              class="custom-handle"
              :connectable="true"
              :style="{
                backgroundColor: '#5C6BC0',
                border: '1px solid #5C6BC0',
                top: '50%',
                right: 0,
                transform: 'translate(50%, -50%)'
              }"
            />
            <v-card
              :style="{
                width: '100%',
                border: (selected) ?
                '1.5px solid #5C6BC0'
                : '1px solid #5C6BC0',
                transition: 'box-shadow 0.2s ease',
                boxShadow: (hoverId === id || selected)
                ? '0 0 0 0.5px #5C6BC0' : 'none'
              }"
              flat
            >
              <v-card-title class="d-flex align-center">
                <v-avatar size="30" style="margin-right: 10px;">
                  <v-img :src="workflowIcon" />
                </v-avatar>
                <span
                  class="text-subtitle-2"
                  v-if="!editingName || id !== editingName_id"
                  style="overflow: hidden;"
                >
                  {{ data.name }}
                </span>
                <v-text-field
                  v-else
                  ref="nameWorkflowRef"
                  variant="outlined"
                  density="compact"
                  hide-details
                  class="text-subtitle-2"
                  v-model="data.name"
                  @blur="finishEditingName"
                  @keyup.enter="finishEditingName"
                  style="
                    flex: 1 1 auto;
                    min-width: 0;
                  "
                >
                </v-text-field>
                <v-spacer></v-spacer>
                <v-btn 
                  icon="mdi-play" 
                  variant="plain" 
                  size="x-small"
                  @click="testWorkflow()"
                >
                </v-btn>
                <v-menu>
                  <template v-slot:activator="{props}">
                    <v-btn
                      icon="mdi-dots-vertical"
                      v-bind="props"
                      variant="plain"
                      size="x-small"
                    >
                    </v-btn>
                  </template>
                  <v-list density="compact">
                    <v-list-item
                      v-for="(item, i) in nodeModifyItems"
                      :key="i"
                      :value="i"
                      @click="
                        selectedNode=nodes.find((node) => node.id === id);
                        handleModify(i)"
                    >
                      <v-list-item-title>{{ item.title }}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
              </v-card-title>
              <v-card-text class="px-3 text-subtitle-2 align-center">
                <p>输入</p>
                <span class="ml-3 font-weight-light" v-if="!data.input.length">未配置输入</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.input" :key="i"
                    class="ma-1"
                    :color="attr.value? 'blue' : 'grey'"
                    size="small"
                    variant="outlined"
                    :style="{ borderColor: '#CE93D8' }"
                  >
                    <strong>{{ simplifyType(attr.value_type) }}.
                    </strong>{{ attr.value_name }}
                    <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
              <v-card-text class="px-3 py-1 text-subtitle-2 align-center">
                <p>输出</p>
                <span class="ml-3 font-weight-light" v-if="!data.output.length">未配置输出</span>
                <div class="mt-2">
                  <v-chip
                    v-for="(attr, i) in data.output" :key="i"
                    class="ma-1"
                    :color="attr.value? 'blue' : 'grey'"
                    size="small"
                    variant="outlined"
                    :style="{ borderColor: '#CE93D8' }"
                  >
                    <strong>{{ simplifyType(attr.value_type) }}.
                    </strong>{{ attr.value_name }}
                    <v-chip v-if="attr.value" size="x-small" color="green" variant="text" label><strong>{{ attr.value }}</strong></v-chip>
                  </v-chip>
                </div>
              </v-card-text>
              <v-card-text v-if="data.workflowId" class="px-3 py-1 text-subtitle-2 align-center">
                <v-chip size="small" color="primary" variant="outlined">
                  <v-icon size="small" class="mr-1">mdi-file-tree</v-icon>
                  工作流: {{ getWorkflowDisplayName(data.workflowId) }}
                </v-chip>
              </v-card-text>
            </v-card>
          </div>
        </template>
      </VueFlow>

      <el-drawer
        v-model="runworkflow"
        append-to-body
        :z-index="2000"
        :modal = "true"
        size="500px"
        style="border-top-left-radius: 8px; border-bottom-left-radius: 8px;"
        direction=rtl
      >
        <template #header>
          <div style="display: flex ; align-items: center;">
            <h3 class="ml-3">
              试运行
            </h3>
          </div>
        </template>
        <v-divider></v-divider>
        <WorkFlowRunner />
      </el-drawer>

      <el-drawer
        v-model="drawerShow"
        append-to-body
        :z-index="2000"
        :modal = "true"
        v-if="selectedNode != null && editingName === false"
        size="500px"
        style="border-top-left-radius: 8px; border-bottom-left-radius: 8px;"
        direction=rtl
      >
        <template #header>
          <div style="display: flex ; align-items: center;">
            <el-image style="width: 50px; height: 50px;"
            :src="selectedNode.type === 'llm' ? agentIcon :
            selectedNode.type === 'myinput' ? inputIcon :
            selectedNode.type === 'myoutput' ? outputIcon :
            selectedNode.type === 'select' ? selectIcon :
            selectedNode.type === 'end' ? endIcon :
            selectedNode.type === 'start' ? startIcon :
            selectedNode.type === 'code' ? codeIcon : 
            selectedNode.type === 'concat' ? concatIcon :
            selectedNode.type === 'knowledge' ? knowledgeIcon :
            selectedNode.type === 'workflow' ? workflowIcon :
            ''" />
            <h4
              style="
              margin-left: 10px;
              max-width: 100px;
              display:inline-block;
              overflow: hidden;
              white-space: nowrap;
              "
             v-if="editingNameInDrawer == false || !isAuthor">
              {{ selectedNode.data.name }}
            </h4>
            <v-text-field
              v-else-if="isAuthor"
              ref = "nameDrawerRef"
              variant="outlined"
              density="compact"
              hide-details
              class="text-subtitle-2"
              v-model="selectedNode.data.name"
              @blur = "finishEditingNameInDrawer"
              @keyup.enter = "finishEditingNameInDrawer"
              style="
                flex: 1 1 auto;
                min-width: 0;
              "
            ></v-text-field>
            
            <!-- 只读模式标识 -->
            <v-chip 
              v-if="!isAuthor" 
              size="small" 
              color="warning" 
              variant="outlined" 
              class="ml-2"
            >
              <v-icon size="14" class="mr-1">mdi-eye</v-icon>
              只读
            </v-chip>
            
            <v-spacer></v-spacer>
            <v-tooltip :attach="true" text="测试该节点" location="bottom">
              <template v-slot:activator="{ props }">
                <v-btn 
                  v-bind="props"
                  v-if="
                    (selectedNode.type === 'code' 
                  || selectedNode.type === 'llm'
                  || selectedNode.type === 'concat'
                  || selectedNode.type === 'knowledge'
                  || selectedNode.type === 'workflow') && isAuthor" 
                  icon="mdi-play" 
                  variant="plain" 
                  @click="selectedNode.type === 'code' ? testCode() : 
                  selectedNode.type === 'agent' ?  testAgent() : 
                  selectedNode.type === 'knowledge' ? testKnowledge() :
                  selectedNode.type === 'workflow' ? testWorkflow() :
                  testConcat()"
                >
                </v-btn>
              </template>
            </v-tooltip>
            <el-drawer
              v-model="showCodeTest"
              append-to-body
              :z-index="3000"
              :modal = "true"
              size="500px"
              style="border-top-left-radius: 8px; border-bottom-left-radius: 8px; height: calc(100vh - 122px); bottom: 0; top: auto;"
              direction=rtl
            >
              <template #header>
                <h4>试运行</h4>
              </template>
              
              <v-divider></v-divider>
              
              <CodeNodeTest></CodeNodeTest>
            </el-drawer>
            <el-drawer
              v-model="showAgentTest"
              append-to-body
              :z-index="3000"
              :modal = "true"
              size="500px"
              style="border-top-left-radius: 8px; border-bottom-left-radius: 8px; height: calc(100vh - 122px); bottom: 0; top: auto;"
              direction=rtl
            >
              <template #header>
                <h4>试运行</h4>
              </template>

              <v-divider></v-divider>
              
              <AgentNodeTest></AgentNodeTest>
            </el-drawer>
            <el-drawer
              v-model="showConcatTest"
              append-to-body
              :z-index="3000"
              :modal = "true"
              size="500px"
              style="border-top-left-radius: 8px; border-bottom-left-radius: 8px; height: calc(100vh - 122px); bottom: 0; top: auto;"
              direction=rtl
            >
              <template #header>
                <h4>试运行</h4>
              </template>

              <v-divider></v-divider>
              
              <ConcatNodeTest></ConcatNodeTest>
            </el-drawer>
            <el-drawer
              v-model="showKnowledgeTest"
              append-to-body
              :z-index="3000"
              :modal = "true"
              size="500px"
              style="border-top-left-radius: 8px; border-bottom-left-radius: 8px; height: calc(100vh - 122px); bottom: 0; top: auto;"
              direction=rtl
            >
              <template #header>
                <h4>试运行</h4>
              </template>

              <v-divider></v-divider>
              
              <KnowledgeNodeTest></KnowledgeNodeTest>
            </el-drawer>
            <el-drawer
              v-model="showWorkflowTest"
              append-to-body
              :z-index="3000"
              :modal = "true"
              size="500px"
              style="border-top-left-radius: 8px; border-bottom-left-radius: 8px; height: calc(100vh - 122px); bottom: 0; top: auto;"
              direction=rtl
            >
              <template #header>
                <h4>试运行</h4>
              </template>

              <v-divider></v-divider>
              
              <WorkflowNodeTest></WorkflowNodeTest>
            </el-drawer>
            <v-menu :attach="true" v-if="selectedNode.type !== 'start' && selectedNode.type !== 'end' && isAuthor">
              <template v-slot:activator="{props}">
                <v-btn
                  icon="mdi-dots-vertical"
                  v-bind="props"
                  variant="plain"
                  class="mr-2"
                >
                </v-btn>
              </template>
              <v-list density="compact">
                <v-list-item @click="renameInDrawer">
                  <v-list-item-title>重命名</v-list-item-title>
                </v-list-item>
                <v-list-item @click="removeNodeInDrawer">
                  <v-list-item-title>删除</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </div>
        </template>

        <v-divider></v-divider>

        <StartNodeForm v-if="selectedNode.type ==='start'" :readonly="!isAuthor"></StartNodeForm>
        <EndNodeForm v-if="selectedNode.type === 'end'" :readonly="!isAuthor"></EndNodeForm>
        <AgentNodeForm v-if="selectedNode.type === 'llm'" :readonly="!isAuthor"></AgentNodeForm>
        <OutputNodeForm v-if = "selectedNode.type === 'myoutput'" :readonly="!isAuthor"></OutputNodeForm>
        <InputNodeForm v-if="selectedNode.type ==='myinput'" :readonly="!isAuthor"></InputNodeForm>
        <CodeNodeForm v-if="selectedNode.type === 'code'" :readonly="!isAuthor"></CodeNodeForm>
        <ConcatNodeForm v-if="selectedNode.type === 'concat'" :readonly="!isAuthor"></ConcatNodeForm>
        <KnowledgeNodeForm v-if="selectedNode.type === 'knowledge'" :readonly="!isAuthor"></KnowledgeNodeForm>
        <WorkflowNodeForm v-if="selectedNode.type === 'workflow'" :readonly="!isAuthor"></WorkflowNodeForm>
      </el-drawer>
    </v-container>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { initialNodes } from '@/plugins/initial-elements'
import { VueFlow, useVueFlow ,Handle } from '@vue-flow/core'
import { Controls, ControlButton } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import { Background } from '@vue-flow/background'
import Icon from '@/components/workflow/Icon.vue'
import InteractionControls from '@/components/workflow/InteractionControls.vue'
import agentIcon from '@/assets/agent.png'
import selectIcon from '@/assets/choice.png'
import startIcon from '@/assets/start.png'
import endIcon from '@/assets/end.png'
import inputIcon from '@/assets/input.png'
import outputIcon from '@/assets/output.png'
import codeIcon from '@/assets/code.png'
import concatIcon from '@/assets/concat.png'
import knowledgeIcon from '@/assets/knowledge.png'
import workflowIcon from '@/assets/workflow.png'
import DropzoneBackground from '@/components/workflow/DropzoneBackground.vue'
import { useRoute, useRouter } from 'vue-router'
import { MarkerType  } from '@vue-flow/core'
import useDragAndDrop from '@/components/workflow/useDnD'
import { onClickOutside } from '@vueuse/core'
import { 
  selectedNode, nodes, edges, fetchWorkflowContent, 
  saveWorkflowContent, runWorkflow, ForkWorkflow, runWorkflowWithInputs,
  getKnowledgeFileName, getWorkflowName, showCodeTest, showAgentTest, 
  showConcatTest, showWorkflowTest, showKnowledgeTest
} from '@/api/workflow'
import { ElMessage } from 'element-plus'
import server from '@/api/server';
import { user } from '@/stores/user'
import ConcatNodeTest from '@/components/workflow/ConcatNodeTest.vue'
import AgentNodeTest from '@/components/workflow/AgentNodeTest.vue'
import CodeNodeTest from '@/components/workflow/CodeNodeTest.vue'
import WorkflowNodeForm from '@/components/workflow/WorkflowNodeForm.vue'
import WorkflowNodeTest from '@/components/workflow/WorkflowNodeTest.vue'
import KnowledgeNodeTest from '@/components/workflow/KnowledgeNodeTest.vue'
import StartNodeForm from '@/components/workflow/StartNodeForm.vue'
import EndNodeForm from '@/components/workflow/EndNodeForm.vue'
import AgentNodeForm from '@/components/workflow/AgentNodeForm.vue'
import OutputNodeForm from '@/components/workflow/OutputNodeForm.vue'
import InputNodeForm from '@/components/workflow/InputNodeForm.vue'
import CodeNodeForm from '@/components/workflow/CodeNodeForm.vue'
import ConcatNodeForm from '@/components/workflow/ConcatNodeForm.vue'
import KnowledgeNodeForm from '@/components/workflow/KnowledgeNodeForm.vue'

const workflow = ref(null)

function testCode() {
  showCodeTest.value = true;
  // console.log("testCode" , selectedNode.value.data.code)
}

function testAgent() {
  showAgentTest.value = true;
}

function testConcat() {
  showConcatTest.value = true;
}

function testKnowledge() {
  showKnowledgeTest.value = true;
}

function testWorkflow() {
  showWorkflowTest.value = true;
}

const router = useRouter()
const route = useRoute()

const { agent_id } = route.params 

const hoverId = ref(null)

const showInteractionControls = ref(false);

const availableNodeTypes = [
  { name: 'llm', label: '大模型节点', icon: agentIcon },
  { name: 'select', label: '条件分支节点', icon: selectIcon },
  { name: 'myinput', label: '输入节点', icon: inputIcon },
  { name: 'myoutput', label: '输出节点', icon: outputIcon },
  { name: 'code' , label: '代码节点', icon: codeIcon },
  { name: 'concat' , label: '拼接节点', icon: concatIcon },
  { name: 'knowledge' , label: '知识库检索', icon: knowledgeIcon },
  { name: 'workflow' , label: '工作流节点', icon: workflowIcon },
]

const { onDragStart , onDrop , onDragOver , onDragLeave , isDragOver } = useDragAndDrop()


const { onInit, onNodeDragStop, onConnect, addEdges,
   setViewport, onNodesChange, applyNodeChanges,
   onEdgesChange , applyEdgeChanges, removeNodes } = useVueFlow()

const drawerShow = ref(false)
const panel = ref(['inout', 'function', 'knowledge', 'workflow', 'branch'])


function simplifyType(type) {
  if (type == 'string') return 'str';
  if (type == 'number') return 'num';
  if (type == 'boolean') return 'bool';
  if (type == 'integer') return 'int';
}

const goBack = () => {
  router.push(`/agents/${agent_id}/introduction`)
}

const nodeModifyItems = [
  { title: '重命名'} ,
  { title: '删除'  } ,
]

const editingName = ref(false)
const editingName_id = ref(null)

function handleModify(index) {
  console.log('index = ' + index);
  console.log('selectedNode = ' + selectedNode.value.data.name)
  switch (index) {
    case 0:
      editingName.value = true
      editingName_id.value = selectedNode.value.id
      break;
    case 1:
      removeNodes(selectedNode.value.id)
      selectedNode.value = null
      break;
  }
}

const editingNameInDrawer = ref(false)

function renameInDrawer() {
  editingNameInDrawer.value = true;
  console.log("ATTENTION " + editingNameInDrawer.value);
}

function removeNodeInDrawer() {
  removeNodes(selectedNode.value.id)
  selectedNode.value = null
}

const nameInputRef = ref(null)
const nameOutputRef = ref(null)
const nameAgentRef = ref(null)
const nameSelectRef = ref(null)
const nameCodeRef = ref(null)
const nameDrawerRef = ref(null)
const nameConcatRef = ref(null)
const nameKnowledgeInputRef = ref(null)
const nameKnowledgeSearchRef = ref(null)
const nameKnowledgeDeleteRef = ref(null)
const nameWorkflowRef = ref(null)

onClickOutside(
  nameDrawerRef, () => {
    if (editingNameInDrawer.value) {
      finishEditingNameInDrawer()
    }
  }
)

onClickOutside(
  nameInputRef, () => {
    if (editingName.value) {
      finishEditingName()
    }
  }
)

onClickOutside(
  nameOutputRef, () => {
    if (editingName.value) {
      finishEditingName()
    }
  }
)

onClickOutside(
  nameAgentRef, () => {
    if (editingName.value) {
      finishEditingName()
    }
  }
)

onClickOutside(
  nameSelectRef, () => {
    if (editingName.value) {
      finishEditingName()
    }
  }
)

onClickOutside(
  nameCodeRef, () => {
    if (editingName.value) {
      finishEditingName()
    }
  }
)

onClickOutside(
  nameConcatRef, () => {
    if (editingName.value) {
      finishEditingName()
    }
  }
)

onClickOutside(
  nameKnowledgeInputRef, () => {
    if (editingName.value) {
      finishEditingName()
    }
  }
)

onClickOutside(
  nameKnowledgeSearchRef, () => {
    if (editingName.value) {
      finishEditingName()
    }
  }
)

onClickOutside(
  nameKnowledgeDeleteRef, () => {
    if (editingName.value) {
      finishEditingName()
    }
  }
)

onClickOutside(
  nameWorkflowRef, () => {
    if (editingName.value) {
      finishEditingName()
    }
  }
)

function finishEditingName() {
  editingName.value = false
  editingName_id.value = null
  selectedNode.value = null
}

function finishEditingNameInDrawer() {
  editingNameInDrawer.value = false
}

function onNodeClick({ node }) {
  // 如果不是作者，只允许查看不允许编辑
  if (!isAuthor.value) {
    selectedNode.value = node
    drawerShow.value = true
    console.log("查看他人工作流节点", selectedNode.value)
    // 显示提示信息
    ElMessage.info('该工作流为只读模式，无法编辑节点')
    return
  }
  
  selectedNode.value = node
  drawerShow.value = true
  console.log("SelectedNode " , selectedNode.value)
}

onInit((vueFlowInstance) => {
  vueFlowInstance.setViewport({ x: 150 , y : 150 , zoom: 1 })
})

onNodeDragStop(({ event, nodes, node }) => {
  console.log('Node Drag Stop', { event, nodes, node })
})

onConnect((connection) => {
  addEdges({
    ...connection,
    // animated: true,
    markerEnd: MarkerType.ArrowClosed,
  })
})

function resetTransform() {
  setViewport({ x: 150, y: 150, zoom: 1.2 })
}

const dark = ref(false)
function toggleDarkMode() {
  dark.value = !dark.value
}

onNodesChange(async (changes) => {
  const nextChanges = changes.filter((change) => {
    if (change.type === 'remove') {
      const node = nodes.value.find((n) => n.id === change.id)
      if (node && node.type !== 'start' && node.type !== 'end') {
        drawerShow.value = false
        return true;
      } else {
        return false;
      }
    }
    return true
  })
  applyNodeChanges(nextChanges)
})

onEdgesChange(async (changes) => {
  const nextChanges = [];
  for (const change of changes) {
    nextChanges.push(change);
  }
  applyEdgeChanges(nextChanges)
})

const saveWorkflow = async () => {
  const content = {
    content: {
      nodes: nodes.value,
      edges: edges.value,
    }
  }
  const res = await saveWorkflowContent(agent_id, content)
  if (res.status === 200) {
    ElMessage.success('保存成功')
  } else {
    ElMessage.error('保存失败')
  }
}

const post = ref(null)

// 存储知识库文件名和工作流名称的缓存
const knowledgeFileNames = ref(new Map())
const workflowNames = ref(new Map())

// 获取知识库文件名的函数
const getKnowledgeFileNameById = async (fileId) => {
  if (knowledgeFileNames.value.has(fileId)) {
    return knowledgeFileNames.value.get(fileId)
  }
  
  const fileName = await getKnowledgeFileName(fileId)
  if (fileName) {
    knowledgeFileNames.value.set(fileId, fileName)
  }
  return fileName || `ID: ${fileId}`
}

// 获取工作流名称的函数
const getWorkflowNameById = async (workflowId) => {
  if (workflowNames.value.has(workflowId)) {
    return workflowNames.value.get(workflowId)
  }
  
  const workflowName = await getWorkflowName(workflowId)
  if (workflowName) {
    workflowNames.value.set(workflowId, workflowName)
  }
  return workflowName || `ID: ${workflowId}`
}

// 监听节点变化，自动加载新的文件名和工作流名称
watch(
  () => nodes.value,
  async (newNodes) => {
    const knowledgePromises = []
    const workflowPromises = []
    
    for (const node of newNodes) {
      if (node.data.knowledgeFileId && !knowledgeFileNames.value.has(node.data.knowledgeFileId)) {
        knowledgePromises.push(getKnowledgeFileNameById(node.data.knowledgeFileId))
      }
      if (node.data.workflowId && !workflowNames.value.has(node.data.workflowId)) {
        workflowPromises.push(getWorkflowNameById(node.data.workflowId))
      }
    }
    
    // 并行执行所有请求，不阻塞界面
    Promise.allSettled([...knowledgePromises, ...workflowPromises]).catch(error => {
      console.warn('加载节点名称失败:', error)
    })
  },
  { deep: true }
)

// 预加载所有节点的文件名和工作流名称
const preloadNames = async () => {
  // 收集所有需要加载的ID
  const knowledgeFileIds = []
  const workflowIds = []
  
  for (const node of nodes.value) {
    if (node.data.knowledgeFileId && !knowledgeFileNames.value.has(node.data.knowledgeFileId)) {
      knowledgeFileIds.push(node.data.knowledgeFileId)
    }
    if (node.data.workflowId && !workflowNames.value.has(node.data.workflowId)) {
      workflowIds.push(node.data.workflowId)
    }
  }
  
  // 并行加载所有知识库文件名
  const knowledgePromises = knowledgeFileIds.map(async (fileId) => {
    try {
      await getKnowledgeFileNameById(fileId)
    } catch (error) {
      console.warn(`加载知识库文件名失败 (ID: ${fileId}):`, error)
    }
  })
  
  // 并行加载所有工作流名称
  const workflowPromises = workflowIds.map(async (workflowId) => {
    try {
      await getWorkflowNameById(workflowId)
    } catch (error) {
      console.warn(`加载工作流名称失败 (ID: ${workflowId}):`, error)
    }
  })
  
  // 等待所有并行请求完成
  await Promise.allSettled([...knowledgePromises, ...workflowPromises])
}

const fetchPostData = async () => {
  try {
    const response = await server.Get(`/workflow/${agent_id}`)
    console.log('Fetched workflow data:', response);
    const workflow = response.data

    post.value = {
      id: workflow.id,
      title: workflow.name,
      description: workflow.description,
      cover: workflow.cover || '/images/default_cover.jpg'
    }
  } catch (error) {
    console.error('加载工作流数据失败', error)
    ElMessage.error('加载工作流数据失败')
  }
}

// 获取知识库文件名用于显示
const getKnowledgeFileDisplayName = (fileId) => {
  return knowledgeFileNames.value.get(fileId) || `ID: ${fileId}`
}

// 获取工作流名称用于显示
const getWorkflowDisplayName = (workflowId) => {
  return workflowNames.value.get(workflowId) || `ID: ${workflowId}`
}

const isForking = ref(false)

const forkWorkflow = async () => {
  const res = await ForkWorkflow(agent_id)
  if (res.status === 200) {
    const { workflow } = res.data
    ElMessage.success('Fork成功')
    await router.push(`/agents/${workflow.id}/workflow`)
    window.location.reload()  
  } else {
    ElMessage.error('Fork失败')
  }
}

const running = ref(false)
const runworkflow = ref(false)

const runWholeWorkflow = () => {
  runworkflow.value = true;
}

onMounted(async () => {
  fetchPostData()
  const res = await fetchWorkflowContent(agent_id)
  workflow.value = res
  if (res) {
    const { content } = res
    nodes.value = content.nodes || initialNodes
    edges.value = content.edges || []
    // 非阻塞地预加载所有节点的文件名和工作流名称
    preloadNames().catch(error => {
      console.warn('预加载节点名称失败:', error)
    })
  }
})

const isAuthor = computed(() => {
  if (workflow.value === null) return false
  return workflow.value.belongsTo.id === user.id
})
</script>

<style scoped>
.header-back {
  position: fixed;
  bottom: 15px;
  color: white;
  left: 20px;
  z-index: 1100;
  margin: 1px;
}

.back-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(25, 118, 210, 0.4);
}

.combined-editor {
  display: flex;
  height: 100vh;
}

.node-panel {
  width: 300px;
  padding: 10px;
  background-color: #f7f7f7;
  border-right: 1px solid #ddd;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.node-panel :deep(.v-expansion-panel) {
  margin: 0;
  border-radius: 0;
  box-shadow: none;
  border-bottom: 1px solid #eee;
}
.node-panel :deep(.v-expansion-panel:last-child) {
  border-bottom: none;
}

.node-item {
  background: #fff;
  border: 1px solid #ccc;
  padding: 8px;
  margin: 2px;
  width: 100%;
  text-align: center;
  cursor: grab;
  display: flex;
  align-items: center;
  border-radius: 8px;
  transition: 0.2s;
}

.node-item:hover {
  background: #e6f7ff;
  border-color: #91d5ff;
}

.node-icon {
  width: 30px;
  height: 30px;
}

.properties-panel {
  width: 300px;
  padding: 20px;
  background-color: #fafafa;
  border-left: 1px solid #ddd;
}

.panel-header {
  display: flex;
  margin-bottom: 20px;
  text-align: center;
}

.custom-handle {
  position: absolute;
  min-width: 10px;
  min-height: 10px;
  border-radius: 100%;
  z-index: 1000;
}

:deep(.basic-flow .vue-flow__controls) {
    justify-content:center
}

:deep(.basic-flow .vue-flow__controls .vue-flow__controls-button) {
    border:none;
    border-right:1px solid #eee
}

:deep(.basic-flow .vue-flow__controls .vue-flow__controls-button svg) {
    height:100%;
    width:100%
}

:deep(.basic-flow .vue-flow__minimap) {
  background-color: rgba(240, 242, 243, 0.7);
}

:deep(.el-drawer__header) {
  margin-bottom: 0 !important;
}

.attrs-scroll-wrapper {
  position: relative;
  display: inline-block;
  flex: 1;
  margin: 0 4px;
  height: 24px;
  overflow: hidden;
  white-space: nowrap;
  vertical-align: middle;
}

.attr-text {
  position: relative;
  padding-right: 32px;
}

.attrs-inner {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.attrs-fade {
  position: absolute;
  top: 0; right: 0;
  width: 32px; height: 100%;
  pointer-events: none;
  background: linear-gradient(to right, rgba(255,255,255,0), #fff);
}

.attr-chip {
  height: 20px;
  line-height: 20px;
  font-size: 11px;
  padding: 0 4px !important;
  background: #fafafa !important;
}
</style>
