<script setup lang="ts">
import { ref, onMounted } from 'vue';
import ModalAdicionarTarefa from '@/components/ui/modal/ModalAdicionarTarefa.vue';
import ModalCadastrarCards from '@/components/ui/modal/ModalCadastrarCards.vue';
import ModalCadastrarParticipante from '@/components/ui/modal/ModalCadastrarParticpante.vue';

const showAdicionarTarefa = ref(false);
const showModalCartas = ref(false);
const showRegisterModal = ref(false);
const roomId = ref<string | null>(null);
let userId = ref<string | null>(null);
const route = useRoute();
const router = useRouter();
const name = ref("Sala de Planejamento - Planning Poker");
let roomSocket: ReturnType<typeof useRoomSocket> | null = null;

onMounted(async () => {
  document.title = name.value;
  userId.value = sessionStorage.getItem('userId');
  if (!userId.value) {
    showRegisterModal.value = true;
  }

   try {
    
    const parts = typeof window !== 'undefined'
      ? new URL(window.location.href).pathname.split('/').filter(Boolean)
      : [];
    const idFromPath = parts.length ? parts[parts.length - 1] : null;
    roomId.value =  (route.params.id as string) || idFromPath || null;

  } catch {
    roomId.value = null;
  }

  if (roomId.value && userId.value) {
    roomSocket = useRoomSocket(roomId.value, userId.value);
    console.log("Connecting to room socket with roomId:", roomId.value, "and userId:", userId.value);
    roomSocket.connect();

    watch(roomSocket.roomState, (newState) => {
      console.log("Room state updated:", newState);
        participants.value = newState.participants || [];
        tasks.value = newState.tasks || [];
    });
  }
});

onBeforeUnmount(() => {
  if (roomSocket) {
    roomSocket.disconnect();
  }
});


// --- DADOS ---
// Use 'const' para dados que não mudam
const votingCards = [1, 2, 3, 5, 8, 13, 21, '☕', '?'];
// Inicializa vazia; será populada no onMounted
const tasks = ref<string[]>([]);

// Use 'ref' para dados que podem mudar (ex: vindos de uma API ou WebSocket)
const participants = ref<any[]>([]);


const myVote = ref<string | number | null>(null);

// --- EVENTOS ---
const emit = defineEmits(['vote']);

const handleVote = (card: string | number) => {
  // Lógica para permitir desmarcar o voto
  myVote.value = myVote.value === card ? null : card;
  emit('vote', myVote.value);
};

const onRegisterCancel = () => {
  router.replace('/');
};
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 mt-2 m-7">

    <!-- COLUNA DA ESQUERDA: INFORMAÇÕES -->
    <aside class="lg:col-span-1 space-y-8">
      <section class="bg-white p-6 rounded-lg shadow-md mt-3">
        <div class="flex justify-between mb-4">
          <h2 class="text-xl font-bold mb-4">Tarefa Atual</h2>
          <button @click="showAdicionarTarefa = true" class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors">
            Adicionar Tarefa
          </button>
          <ModalAdicionarTarefa v-if="showAdicionarTarefa" @close="showAdicionarTarefa = false" />
        </div>

        <div v-if="tasks.length > 0" class="p-4 border border-gray-200 rounded-lg bg-gray-50">
          <span class="text-gray-800 text-lg">{{ tasks[0] }}</span>
        </div>
        <div v-else class="text-gray-500">Nenhuma tarefa na fila.</div>
      </section>


      <section class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-xl font-bold mb-4 border-b pb-2">Participantes</h2>
        <ul class="space-y-3">
          <li v-for="p in participants" :key="p.id" class="flex items-center justify-between p-2 rounded-md" :class="p.vote ? 'bg-green-100' : 'bg-gray-100'">
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-blue-500 text-white rounded-full flex items-center justify-center font-bold">
                {{ p.name.charAt(0).toUpperCase() }}
              </div>
              <span class="text-gray-800 font-medium">{{ p.name }}</span>
            </div>
            <span class="font-bold text-lg">{{ p.vote ? '✅' : '🤔' }}</span>
          </li>
        </ul>
      </section>
    </aside>

    <!-- COLUNA DA DIREITA: AÇÃO -->
    <main class="lg:col-span-2 space-y-8  col-2 mt-2">
      <section class="bg-white p-6 rounded-lg shadow-md">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-bold mb-4">Escolha sua carta</h2>
          <div class="flex space-x-4">
          <button @click="showModalCartas = true" class="bg-green-500 text-white px-4 mt-2 py-2 rounded-md hover:bg-blue-700 transition-colors">
          Editar Cartas
          </button>
          <button @click="handleVote(myVote!)" class="bg-blue-600 text-white px-4 mt-2 py-2 rounded-md hover:bg-blue-700 transition-colors">
          Confirmar Voto
          </button>
          </div>
          <ModalCadastrarCards v-if="showModalCartas" @close="showModalCartas = false" />
        </div>

        <div class="flex flex-wrap gap-4">
          <button
              v-for="card in votingCards"
              :key="card"
              @click="handleVote(card)"
              :class="[
              'w-16 h-24 text-2xl font-bold rounded-lg shadow-md transition-all transform hover:scale-105 border-2',
              myVote === card
                ? 'bg-blue-600 border-blue-700 text-white scale-105'
                : 'bg-white border-gray-300 text-gray-700'
            ]"
          >
            {{ card }}
          </button>
        </div>
       
      </section>
      
    </main>

    <ModalCadastrarParticipante
      v-if="showRegisterModal"
      @registered="onParticipantRegistered"
      @close="onRegisterCancel"
    />
  </div>
</template>
