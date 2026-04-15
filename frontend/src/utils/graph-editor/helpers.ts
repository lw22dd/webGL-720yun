export function generateId(prefix: string = 'id'): string {
  return `${prefix}-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`
}

export function formatZoom(zoom: number): string {
  return `${Math.round(zoom * 100)}%`
}
