import { ref, watch } from 'vue'
import { useWebSocket } from '@vueuse/core'

interface Message {
  type: string
  payload: any
}

export function useRoomSocket(roomId: string, userId: string) {
    const roomState = ref<any>({})
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
            const parsedMessage: Message = JSON.parse(newMessage)
    
            // Usamos um switch para lidar com diferentes tipos de mensagens do servidor
            switch (parsedMessage.type) {
              case 'INITIAL_STATE':
                roomState.value = parsedMessage.payload
                break

              case 'ROOM_STATE_UPDATED':
                roomState.value = parsedMessage.payload
                break
              
              case 'USER_JOINED':
                // Aqui você poderia adicionar um novo participante à lista existente
                console.log('Novo participante:', parsedMessage.payload);
                break
    
              // Adicione outros 'case' para outras mensagens do servidor aqui
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
