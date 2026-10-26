/**
 * 资源生命周期观测。
 *
 * 为什么需要它：`performance.memory` 只能看到 JS 堆，看不到纹理、canvas 和
 * WebGL 上下文——而长会话教学场景（30 分钟、切换 20+ 个场景）真正的验证目标是
 * “回到基线”：上下文数、canvas 数、JS 堆是否回到进入全景之前的水平。
 *
 * 开发环境默认启用；生产包可以在控制台执行
 *   localStorage.setItem('pano:diag', '1') 后刷新来打开。
 *
 * 用法（浏览器控制台）：
 *   __panoDiag.report('enter-pano')    // 打一个点
 *   __panoDiag.diff()                  // 和第一个基线点对比
 *   __panoDiag.table()                 // 打印全部打点
 *   __panoDiag.contexts()              // 逐个列出存活上下文：谁创建的、多大、是否还在 DOM 里
 *
 * 验收流程（可写进论文的实验部分）：
 *   1. 进首页 → 打点 A；
 *   2. 进全景，用热点连续切 20 个场景；
 *   3. 返回首页 → 打点 B，再执行 diff()；DevTools 里手动 GC 后再 report 一次。
 *   期望：webglAttached / canvases 回到 A 的水平。webglDetached 长期不降，
 *   说明某条路径的 destroy 只摘了 DOM，没把 GL 上下文还回去。
 *
 * 已知局限：只能统计主线程 canvas；Worker / OffscreenCanvas 上的上下文不在
 * 计数范围内。库做能力探测时在临时 canvas 上创建的上下文也会被计入，
 * 这类上下文 attached=false，用 contexts() 看 origin 即可分辨。
 */

const enabled = (() => {
  try {
    if (import.meta.env.DEV) return true
    return typeof localStorage !== 'undefined' && localStorage.getItem('pano:diag') === '1'
  } catch {
    return false
  }
})()

export interface ResourceSnapshot {
  seq: number
  label: string
  at: number
  /** JS 堆，MB（Chromium 才有 performance.memory） */
  jsHeapMB: number | null
  /** 存活的 WebGL 上下文总数（含已脱离 DOM 的） */
  webglContexts: number
  /** 仍然挂在 DOM 上的 WebGL 上下文 */
  webglAttached: number
  /** 已脱离 DOM 但还没被回收的 WebGL 上下文 */
  webglDetached: number
  canvases: number
  images: number
  detail?: Record<string, unknown>
}

interface HeapMemory {
  usedJSHeapSize: number
  totalJSHeapSize: number
  jsHeapSizeLimit: number
}

interface TrackedContext {
  id: number
  type: string
  canvasRef: WeakRef<HTMLCanvasElement>
  createdAt: number
  /** 创建该上下文的调用点，用来区分“业务渲染器”和“能力探测” */
  origin: string
}

const MB = 1024 * 1024
const MAX_SAMPLES = 500

let installed = false
/**
 * 用 WeakRef 持有 canvas：观测器本身不能延长被观测对象的寿命，
 * 否则测出来的“泄漏”就是它自己造成的。
 */
const trackedContexts = new Map<number, TrackedContext>()
const contextIdByCanvas = new WeakMap<HTMLCanvasElement, number>()
let contextSeq = 0

const samples: ResourceSnapshot[] = []
let baseline: ResourceSnapshot | null = null

function readHeap(): HeapMemory | null {
  if (typeof performance === 'undefined') return null
  return (performance as Performance & { memory?: HeapMemory }).memory ?? null
}

function captureOrigin(): string {
  const stack = new Error().stack?.split('\n') ?? []
  // 跳过 "Error" 行和本模块的包装帧
  const frame = stack.find((line) => line.includes('at ') && !line.includes('resourceMonitor'))
  return frame ? frame.trim().replace(/^at\s+/, '') : 'unknown'
}

/**
 * 通过包装 getContext 统计 WebGL 上下文数量。
 * 上下文丢失（three 的 forceContextLoss / GPU 回收）时移除，这是判断
 * “viewer.destroy() 到底有没有把 GPU 资源还回去”最直接的一个指标。
 */
function install(): void {
  if (installed || typeof HTMLCanvasElement === 'undefined') return
  installed = true

  const original = HTMLCanvasElement.prototype.getContext
  HTMLCanvasElement.prototype.getContext = function (
    this: HTMLCanvasElement,
    contextId: string,
    ...args: unknown[]
  ) {
    const context = (original as any).apply(this, [contextId, ...args])
    if (context && /^(webgl|webgl2|experimental-webgl)$/.test(contextId)) {
      if (!contextIdByCanvas.has(this)) {
        const id = ++contextSeq
        contextIdByCanvas.set(this, id)
        trackedContexts.set(id, {
          id,
          type: contextId,
          canvasRef: new WeakRef(this),
          createdAt: Date.now(),
          origin: captureOrigin(),
        })
        this.addEventListener(
          'webglcontextlost',
          () => {
            // three 不会自动恢复上下文；这里也只在丢失时移除
            trackedContexts.delete(id)
            contextIdByCanvas.delete(this)
          },
          { once: true }
        )
      }
    }
    return context
  } as typeof HTMLCanvasElement.prototype.getContext

  ;(window as any).__panoDiag = resourceMonitor
}

function countContexts(): { total: number; attached: number } {
  let total = 0
  let attached = 0
  for (const [id, ctx] of trackedContexts) {
    const canvas = ctx.canvasRef.deref()
    if (!canvas) {
      // canvas 已被 GC，说明上下文也随之释放了（库不会自动派发 contextlost）
      trackedContexts.delete(id)
      continue
    }
    total += 1
    if (canvas.isConnected) attached += 1
  }
  return { total, attached }
}

function takeSnapshot(label: string, detail?: Record<string, unknown>): ResourceSnapshot {
  const heap = readHeap()
  const { total, attached } = countContexts()
  const snapshot: ResourceSnapshot = {
    seq: samples.length,
    label,
    at: Date.now(),
    jsHeapMB: heap ? Number((heap.usedJSHeapSize / MB).toFixed(1)) : null,
    webglContexts: total,
    webglAttached: attached,
    webglDetached: total - attached,
    canvases: typeof document === 'undefined' ? 0 : document.querySelectorAll('canvas').length,
    images: typeof document === 'undefined' ? 0 : document.images.length,
    detail,
  }
  samples.push(snapshot)
  if (samples.length > MAX_SAMPLES) samples.shift()
  return snapshot
}

function report(label: string, detail?: Record<string, unknown>): ResourceSnapshot | null {
  if (!enabled) return null
  install()
  const snapshot = takeSnapshot(label, detail)
  if (!baseline) baseline = snapshot
  const detailText = detail ? ` ${JSON.stringify(detail)}` : ''
  console.log(
    `[panoDiag] ${label} | heap=${snapshot.jsHeapMB ?? 'n/a'}MB ` +
      `webgl=${snapshot.webglContexts}(attached ${snapshot.webglAttached}/detached ${snapshot.webglDetached}) ` +
      `canvas=${snapshot.canvases} img=${snapshot.images}${detailText}`
  )
  return snapshot
}

/** 逐个列出存活上下文，用来区分业务渲染器和库的能力探测 */
function contexts(): void {
  if (!enabled) return
  install()
  countContexts()
  console.table(
    Array.from(trackedContexts.values()).map((ctx) => {
      const canvas = ctx.canvasRef.deref()
      return {
        id: ctx.id,
        type: ctx.type,
        attached: canvas?.isConnected ?? false,
        size: canvas ? `${canvas.width}x${canvas.height}` : 'collected',
        ageMs: Date.now() - ctx.createdAt,
        origin: ctx.origin,
      }
    })
  )
}

function table(): void {
  if (!enabled) return
  console.table(samples)
}

function diff(): void {
  if (!enabled) return
  if (!baseline) {
    console.warn('[panoDiag] 还没有基线，先调用 report()')
    return
  }
  const now = takeSnapshot('now')
  console.table([
    {
      metric: 'jsHeapMB',
      baseline: baseline.jsHeapMB,
      now: now.jsHeapMB,
      delta: deltaOrNull(baseline.jsHeapMB, now.jsHeapMB),
    },
    { metric: 'webglContexts', baseline: baseline.webglContexts, now: now.webglContexts, delta: now.webglContexts - baseline.webglContexts },
    { metric: 'webglAttached', baseline: baseline.webglAttached, now: now.webglAttached, delta: now.webglAttached - baseline.webglAttached },
    { metric: 'webglDetached', baseline: baseline.webglDetached, now: now.webglDetached, delta: now.webglDetached - baseline.webglDetached },
    { metric: 'canvases', baseline: baseline.canvases, now: now.canvases, delta: now.canvases - baseline.canvases },
    { metric: 'images', baseline: baseline.images, now: now.images, delta: now.images - baseline.images },
  ])
  if (now.webglAttached > baseline.webglAttached) {
    console.warn(
      `[panoDiag] 还有 ${now.webglAttached - baseline.webglAttached} 个 WebGL 上下文挂在 DOM 上没释放，` +
        '说明某条路径的 destroy 没有被调用'
    )
  }
  if (now.webglDetached > 0) {
    console.warn(
      `[panoDiag] 有 ${now.webglDetached} 个上下文已脱离 DOM 但仍存活。` +
        '先手动 GC 再 report 一次：若仍不降，说明库把上下文或 renderer 存进了全局缓存'
    )
  }
}

function deltaOrNull(from: number | null, to: number | null): number | null {
  if (from === null || to === null) return null
  return Number((to - from).toFixed(1))
}

function reset(): void {
  samples.length = 0
  baseline = null
}

export const resourceMonitor = {
  get enabled() {
    return enabled
  },
  report,
  contexts,
  table,
  diff,
  reset,
}

if (enabled) {
  install()
}
