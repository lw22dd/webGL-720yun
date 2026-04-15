export const NODE_SIZE_PANO = 30
export const NODE_SIZE_NO_PANO = 20
export const NODE_WIDTH = 120
export const NODE_HEIGHT = 60

export const PORT_RADIUS = 6

export const COLORS = {
  node: {
    default: '#5F95FF',
    hasPano: '#5F95FF',
    noPano: '#A2B1C3',
    selected: '#4f46e5',
    hover: '#7aa3ff',
  },
  edge: {
    default: '#A2B1C3',
    selected: '#4f46e5',
    hover: '#7aa3ff',
  },
  port: {
    stroke: '#5F95FF',
    fill: '#ffffff',
    active: '#31d0c6',
  },
  toolbar: {
    bg: '#ffffff',
    border: '#e5e7eb',
  },
  canvas: {
    bg: '#f8f9fa',
    grid: '#e0e0e0',
  },
} as const

export const ZOOM_LIMITS = {
  min: 0.2,
  max: 3,
  step: 0.1,
  factor: 1.1,
}

export const GRID_CONFIG = {
  visible: true,
  type: 'dot' as const,
  args: {
    color: '#ddd',
    thickness: 1,
  },
}

export const HISTORY_STACK_SIZE = 50

export const SNAPLINE_TOLERANCE = 10
