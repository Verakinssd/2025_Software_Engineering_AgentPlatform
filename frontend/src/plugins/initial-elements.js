import { MarkerType } from '@vue-flow/core'


export const initialNodes = [
  {
    id: '1',
    type: 'start',
    data: { 
      name: '开始',
      input: []
    },
    position: { x: 100, y: 100 },
  },
  {
    id: '2',
    type: 'end',
    data: { 
      name: '结束' ,
      output: []
    },
    position: { x: 700, y: 100 },
  },
]