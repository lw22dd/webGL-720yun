import type { WebSocketMessage } from '@/models/upload.model'

type EventHandler = (data: any) => void

class WebSocketClient {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000
  private listeners: Map<string, EventHandler[]> = new Map()
  private isConnected = false
  private token: string = ''
  private heartbeatInterval: number | null = null

  connect(token: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this.token = token
      // 构建 WebSocket URL
      const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'
      const wsProtocol = apiBase.startsWith('https:') ? 'wss:' : 'ws:'
      const wsHost = apiBase.replace(/^https?:\/\//, '')
      const wsUrl = `${wsProtocol}//${wsHost}/api/v1/ws?token=${token}`
      console.log('[WebSocket] Connecting to:', wsUrl)

      if (this.ws) {
        this.ws.onopen = null
        this.ws.onmessage = null
        this.ws.onclose = null
        this.ws.onerror = null
        this.ws.close()
        this.ws = null
      }

      try {
        this.ws = new WebSocket(wsUrl)

        this.ws.onopen = () => {
          console.log('[WebSocket] Connected successfully')
          this.isConnected = true
          this.reconnectAttempts = 0
          this.startHeartbeat()
          resolve()
        }

        this.ws.onmessage = (event) => {
          try {
            const message: WebSocketMessage = JSON.parse(event.data)
            this.handleMessage(message)
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error)
          }
        }

        this.ws.onclose = (event) => {
          console.log('[WebSocket] Closed:', event.code, event.reason)
          this.isConnected = false
          this.stopHeartbeat()
          this.emit('disconnected', { code: event.code, reason: event.reason })

          const shouldReconnect = this.reconnectAttempts < this.maxReconnectAttempts
          this.ws = null
          if (shouldReconnect) {
            this.scheduleReconnect()
          }
        }

        this.ws.onerror = (error) => {
          console.error('[WebSocket] Error:', error)
          this.emit('error', error)
          if (!this.isConnected) {
            reject(error)
          }
        }
      } catch (error) {
        reject(error)
      }
    })
  }

  private startHeartbeat() {
    this.heartbeatInterval = window.setInterval(() => {
      if (this.ws && this.isConnected) {
        this.ws.send(JSON.stringify({ type: 'ping' }))
      }
    }, 30000)
  }

  private stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
  }

  private scheduleReconnect() {
    this.reconnectAttempts++
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1)
    console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`)

    setTimeout(() => {
      if (this.token) {
        this.connect(this.token).catch(console.error)
      }
    }, delay)
  }

  private handleMessage(message: WebSocketMessage) {
    console.log('[WebSocket] Received message:', message.type, message)
    this.emit(message.type, message.data)
    this.emit('message', message)
  }

  on(event: string, handler: EventHandler) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, [])
    }
    this.listeners.get(event)!.push(handler)
  }

  off(event: string, handler: EventHandler) {
    const handlers = this.listeners.get(event)
    if (handlers) {
      const index = handlers.indexOf(handler)
      if (index > -1) {
        handlers.splice(index, 1)
      }
    }
  }

  private emit(event: string, data: any) {
    const handlers = this.listeners.get(event)
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(data)
        } catch (error) {
          console.error(`Error in event handler for ${event}:`, error)
        }
      })
    }
  }

  send(data: any) {
    if (this.ws && this.isConnected) {
      this.ws.send(JSON.stringify(data))
    } else {
      console.error('WebSocket is not connected')
    }
  }

  disconnect() {
    this.stopHeartbeat()
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.isConnected = false
    this.listeners.clear()
  }

  getConnectionStatus(): boolean {
    return this.isConnected
  }
}

export const wsClient = new WebSocketClient()
export default wsClient
