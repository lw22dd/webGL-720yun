import { type InjectionKey } from 'vue'
import type { Graph } from '@antv/x6'
import { Dnd } from '@antv/x6'
import type { SceneNodeData } from '@/models/graphEditor/node'
import { NODE_SIZE_PANO, NODE_SIZE_NO_PANO, PORT_RADIUS, COLORS } from '@/utils/graphEditor/constants'

export interface DndActions {
  startDrag: (nodeData: SceneNodeData, e: MouseEvent) => void
}

export const DND_ACTIONS_KEY: InjectionKey<DndActions> = Symbol('dndActions')

function getPortGroups() {
  return {
    top: {
      position: 'top',
      attrs: {
        circle: {
          r: PORT_RADIUS,
          magnet: true,
          stroke: COLORS.port.stroke,
          strokeWidth: 2,
          fill: COLORS.port.fill,
          style: { visibility: 'hidden' },
        },
      },
    },
    bottom: {
      position: 'bottom',
      attrs: {
        circle: {
          r: PORT_RADIUS,
          magnet: true,
          stroke: COLORS.port.stroke,
          strokeWidth: 2,
          fill: COLORS.port.fill,
          style: { visibility: 'hidden' },
        },
      },
    },
    left: {
      position: 'left',
      attrs: {
        circle: {
          r: PORT_RADIUS,
          magnet: true,
          stroke: COLORS.port.stroke,
          strokeWidth: 2,
          fill: COLORS.port.fill,
          style: { visibility: 'hidden' },
        },
      },
    },
    right: {
      position: 'right',
      attrs: {
        circle: {
          r: PORT_RADIUS,
          magnet: true,
          stroke: COLORS.port.stroke,
          strokeWidth: 2,
          fill: COLORS.port.fill,
          style: { visibility: 'hidden' },
        },
      },
    },
  }
}

export function createDnd(graph: Graph): DndActions {
  const dndInstance = new Dnd({
    target: graph,
    getDragNode: (node: any) => node.clone({ keepId: true }),
    getDropNode: (node: any) => node.clone({ keepId: true }),
  })

  function createNodeForDnd(data: SceneNodeData) {
    const size = data.hasPano ? NODE_SIZE_PANO : NODE_SIZE_NO_PANO
    const color = data.hasPano ? COLORS.node.hasPano : COLORS.node.noPano
    return graph.createNode({
      id: data.id,
      shape: 'circle',
      x: 0,
      y: 0,
      width: size,
      height: size,
      data: { ...data },
      attrs: {
        body: {
          fill: color,
          stroke: 'transparent',
          strokeWidth: 0,
        },
        label: {
          text: data.name || '',
          fill: '#374151',
          fontSize: 11,
          fontWeight: 500,
          refX: 0,
          refY: size / 2 + 14,
          textAnchor: 'middle',
          textVerticalAnchor: 'top',
        },
      },
      ports: {
        groups: getPortGroups(),
        items: [
          { group: 'top', id: 'port-top' },
          { group: 'bottom', id: 'port-bottom' },
          { group: 'left', id: 'port-left' },
          { group: 'right', id: 'port-right' },
        ],
      },
    })
  }

  function startDrag(nodeData: SceneNodeData, e: MouseEvent) {
    const node = createNodeForDnd(nodeData)
    dndInstance.start(node, e)
  }

  return { startDrag }
}
