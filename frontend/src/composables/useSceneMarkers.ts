import { ref } from 'vue'
import { Marker, Popup } from '@antv/l7'
import type { Scene as L7Scene } from '@antv/l7'

export interface SceneMarkerData {
  id: number
  title: string
  longitude: number
  latitude: number
  thumbnail_url?: string
  scene_code: string
  view_count?: number
  hotspots?: any[]
}

const MARKER_COLORS = [
  '#0EA5E9', '#2563EB', '#10B981', '#F59E0B',
  '#8B5CF6', '#EC4899', '#DC2626', '#059669',
]

export function useSceneMarkers() {
  const markers = ref<Marker[]>([])
  const activeMarkerId = ref<number | null>(null)

  function getMarkerColor(index: number): string {
    return MARKER_COLORS[index % MARKER_COLORS.length]
  }

  function createMarkerElement(scene: SceneMarkerData, index: number): HTMLDivElement {
    const color = getMarkerColor(index)
    const el = document.createElement('div')
    el.className = 'scene-marker-wrapper'
    el.innerHTML = `
      <div class="scene-marker" style="--marker-color: ${color}">
        <div class="marker-pin">
          <svg width="28" height="38" viewBox="0 0 28 38">
            <path d="M14 0C6.268 0 0 6.268 0 14c0 10.5 14 24 14 24s14-13.5 14-24C28 6.268 21.732 0 14 0z"
                  fill="var(--marker-color)" stroke="rgba(255,255,255,0.8)" stroke-width="1.5"/>
            <circle cx="14" cy="14" r="6" fill="rgba(255,255,255,0.95)"/>
            <circle cx="14" cy="14" r="3" fill="var(--marker-color)"/>
          </svg>
        </div>
        <div class="marker-pulse"></div>
      </div>
      <div class="marker-label">${scene.title}</div>
    `
    return el
  }

  function createPopupContent(scene: SceneMarkerData): string {
    return `
      <div class="scene-popup">
        <div class="popup-header">
          <h3 class="popup-title">${scene.title}</h3>
          <span class="popup-coords">${scene.longitude.toFixed(6)}, ${scene.latitude.toFixed(6)}</span>
        </div>
        ${scene.thumbnail_url ? `<div class="popup-thumbnail"><img src="${scene.thumbnail_url}" alt="${scene.title}" /></div>` : ''}
        <div class="popup-info">
          <div class="info-row">
            <span class="info-label">场景 ID</span>
            <span class="info-value">${scene.id}</span>
          </div>
          ${scene.hotspots && scene.hotspots.length > 0 ? `
          <div class="info-row">
            <span class="info-label">连接热点</span>
            <span class="info-value">${scene.hotspots.length} 个</span>
          </div>` : ''}
          <div class="info-row">
            <span class="info-label info-link" data-scene-id="${scene.id}">查看全景 →</span>
          </div>
        </div>
      </div>
    `
  }

  function addMarker(
    scene: L7Scene,
    sceneData: SceneMarkerData,
    index: number,
    onClick?: (scene: SceneMarkerData) => void
  ): Marker {
    const el = createMarkerElement(sceneData, index)

    const marker = new Marker({
      element: el,
      offsets: [-14, -38],
      anchor: 'bottom-left' as any,
    }).setLnglat([sceneData.longitude, sceneData.latitude] as any)

    const popup = new Popup({
      offsets: [0, -42],
      closeButton: true,
      maxWidth: '300px',
      title: '',
      html: createPopupContent(sceneData),
    })
    marker.setPopup(popup)

    el.addEventListener('click', (e) => {
      e.stopPropagation()
      activeMarkerId.value = sceneData.id
      marker.togglePopup()
      onClick?.(sceneData)
    })

    scene.addMarker(marker)
    markers.value.push(marker)

    return marker
  }

  function clearMarkers(scene: L7Scene) {
    scene.removeAllMarkers()
    markers.value = []
    activeMarkerId.value = null
  }

  function showMarker(marker: Marker) {
    marker.show()
  }

  function hideMarker(marker: Marker) {
    marker.hide()
  }

  function showAllMarkers() {
    markers.value.forEach(m => m.show())
  }

  function hideAllMarkers() {
    markers.value.forEach(m => m.hide())
  }

  return {
    markers,
    activeMarkerId,
    createMarkerElement,
    createPopupContent,
    addMarker,
    clearMarkers,
    showMarker,
    hideMarker,
    showAllMarkers,
    hideAllMarkers,
  }
}
