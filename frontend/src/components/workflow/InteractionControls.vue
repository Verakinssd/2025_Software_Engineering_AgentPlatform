<script setup>
import { Panel, useVueFlow } from '@vue-flow/core'
import { ref } from 'vue'

const {
  nodesDraggable,
  nodesConnectable,
  onConnect,
  onNodeDragStop,
  onPaneClick,
  onPaneScroll,
  onPaneContextMenu,
  onNodeDragStart,
  onMoveEnd,
  addEdges,
} = useVueFlow()

const captureZoomClick = ref(false)

const captureZoomScroll = ref(false)

onConnect((params) => addEdges(params))

onNodeDragStart((e) => console.log('drag start', e))

onNodeDragStop((e) => console.log('drag stop', e))

onPaneClick((event) => captureZoomClick.value && console.log('pane click', event))

onPaneScroll((event) => captureZoomScroll.value && console.log('pane scroll', event))

onPaneContextMenu((event) => captureZoomClick.value && console.log('pane ctx menu', event))

onMoveEnd((flowTransform) => console.log('move end', flowTransform))
</script>

<template>
  <Panel position="top-right">
    <div>
      <label class="label" for="draggable">
        nodesDraggable
        <input id="draggable" v-model="nodesDraggable" type="checkbox" class="vue-flow__draggable" />
      </label>
    </div>
    <div>
      <label class="label" for="connectable">
        nodesConnectable
        <input id="connectable" v-model="nodesConnectable" type="checkbox" class="vue-flow__connectable" />
      </label>
    </div>
  </Panel>
</template>
