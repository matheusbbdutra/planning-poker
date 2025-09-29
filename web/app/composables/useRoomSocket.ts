import { ref, watch } from 'vue'
import { useWebSocket } from '@vueuse/core'
import type { Participant, Task, Room } from '~/types/entities'

interface Message<T = any> {
  type: string
  payload: T
}

export function useRoomSocket(roomId: string, userId: string) {
    const roomState = ref<Room | null>(null)
    const isConnected = ref(false)

    const config = useRuntimeConfig();

    const baseUrl = config.public.apiBase || 'http://localhost:8080'


    const wsUrl = `${baseUrl.replace(/^http/, 'ws')}/room/ws/${roomId}/${userId}`
    const { status, data, send, open, close } = useWebSocket(wsUrl, {
        autoReconnect: true,
    })

    watch(data, (newMessage) => {
        if (typeof newMessage === 'string') {
          try {
            const parsedMessage: Message = JSON.parse(newMessage) as Room
    
            // Usamos um switch para lidar com diferentes tipos de mensagens do servidor
            switch (parsedMessage.type) {
              case 'INITIAL_STATE':
                roomState.value = parsedMessage.payload
                break

              case 'ROOM_STATE_UPDATED':
                roomState.value = parsedMessage.payload
                break
            
            }
          } catch (error) {
            console.error('Erro ao processar mensagem do WebSocket:', error)
          }
        }
    })
    
    watch(status, (newStatus) => {
        isConnected.value = newStatus === 'OPEN';
    })


    const sendMessage = (type: string, payload: any = {}) => {
        if (status.value === 'OPEN') {
          const message: Message = { type, payload }
          send(JSON.stringify(message))
        } else {
          console.warn('WebSocket não está conectado. Mensagem não enviada.')
        }
    }

    return {
        isConnected,
        roomState,
        sendMessage,
        connect: open,
        disconnect: close,
    }
}
