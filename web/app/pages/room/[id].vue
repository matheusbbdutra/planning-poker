<script setup lang="ts">
import ModalAdicionarTarefa from '@/components/ui/modal/ModalAdicionarTarefa.vue';
import ModalCadastrarCards from '@/components/ui/modal/ModalCadastrarCards.vue';
import ModalCadastrarParticipante from '@/components/ui/modal/ModalCadastrarParticpante.vue';
import { Room, Task, Participant, votingStatus } from '~/types/entities';

const showAdicionarTarefa = ref(false);
const showModalCartas = ref(false);
const showRegisterModal = ref(false);
const roomId = ref<string | null>(null);
let userId = ref<string | null>(null);
const route = useRoute();
const router = useRouter();
const name = ref("Sala de Planejamento - Planning Poker");
let roomSocket: ReturnType<typeof useRoomSocket> | null = null;
const votingCards = ref<Array<string | number>>([]);
const tasks = ref<Task[]>([]);
const participants = ref<Participant[]>([]);
const myVote = ref<string | number | null>(null);
const isScrumMaster = ref(false);
const votingTask = ref<Task | null>(null);
const socketStops: Array<() => void> = [];
const defaultCards = [1, 2, 3, 5, 8, 13, 21, '☕', '?'];
const visibleTaskLimit = 3;
const taskWindowStart = ref(0);
const visibleTasks = computed(() =>
  tasks.value.slice(taskWindowStart.value, taskWindowStart.value + visibleTaskLimit)
);
const showTaskNavigation = computed(() => tasks.value.length > visibleTaskLimit);
const canScrollUp = computed(() => taskWindowStart.value > 0);
const canScrollDown = computed(
  () => taskWindowStart.value + visibleTaskLimit < tasks.value.length
);

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

  if (userId.value) {
    initializeRoomSocket(userId.value);
  }

  
});


onBeforeUnmount(() => {
  socketStops.forEach(stop => stop())
  if (roomSocket) {
    roomSocket.disconnect();
  }
});

const handleVote = (card: string | number) => {
  myVote.value = myVote.value === card ? null : card;  
};

const confirmVote = () => {
  if (!votingTask.value) {
    console.warn('Nenhuma tarefa ativa para receber votos.');
    return;
  }

  if (!roomSocket || myVote.value === null) {
    return;
  }

  roomSocket.sendMessage('USER_VOTE', {
    taskId: votingTask.value.id,
    vote: String(myVote.value),
    userId: userId.value,
  });
}

const resetVote = () => {
  myVote.value = null;
};

const onRegisterCancel = () => {
  router.replace('/');
};

const addTask = (task: string) => {
  if (roomSocket) {
    roomSocket.sendMessage('ADD_TASK', { title: task });
  } else {
    console.warn('WebSocket não está conectado. Tarefa não enviada.');
    tasks.value.push(task);
  }
};

const editarCards = (cards: Array<string | number>) => {
  if (roomSocket) {
    roomSocket.sendMessage('UPDATE_CARDS', { numberOfCards: cards });
  }
};

const selectTaskForVoting = (taskId: string) => {
  if (roomSocket) {
    roomSocket.sendMessage(
      'ON_VOTING', 
      { taskId: taskId, votingStatus: votingStatus.VOTING }
    );
  }
};

const isTaskDisabled = (task: Task) => {
  if (task.isCompleted) return true;
  if (!votingTask.value) return false;
  return String(votingTask.value.id) !== String(task.id);
}

const revealCards = () => {
  if (roomSocket && votingTask.value) {
    roomSocket.sendMessage(
      'ON_VOTING', 
      { taskId: votingTask.value.id, votingStatus: votingStatus.COMPLETED }
    );
    resetVote();
  }
};

const getTaskAverage = (task: Task): string | null => {
  if (!task || !task.votes) return null;
  const numericVotes = Object.values(task.votes)
    .map(value => Number(value))
    .filter(value => Number.isFinite(value));

  if (numericVotes.length === 0) {
    return null;
  }

  const total = numericVotes.reduce((sum, value) => sum + value, 0);
  const average = total / numericVotes.length;

  return Number.isInteger(average) ? String(average) : average.toFixed(1);
};

const taskWithAverage = computed(() => {
  if (votingTask.value) {
    return null;
  }
  const completedTasks = [...tasks.value].reverse();
  return completedTasks.find(task => task.isCompleted && getTaskAverage(task)) ?? null;
});

const lastAverage = computed(() => {
  const task = taskWithAverage.value;
  return task ? getTaskAverage(task) : null;
});

const onParticipantRegistered = async (name: string) => {
  showRegisterModal.value = false;
  const { data, error } = await useApi(`/room/${roomId.value}/join`, {
      method: 'POST',
      body: { userName: name },
  });

  if (error.value) {
    console.error("Erro ao entrar na sala:", error.value);
    return;
  }

  const createdUserId = data.value?.id;

  if (!createdUserId) {
    console.error("Resposta inesperada ao entrar na sala:", data.value);
    return;
  }
  userId.value = createdUserId;
  sessionStorage.setItem("userId", createdUserId);
  initializeRoomSocket(createdUserId);
};

const initializeRoomSocket = (connectedUserId: string) => {
  if (!roomId.value) return;

  socketStops.forEach(stop => stop());
  socketStops.length = 0;

  roomSocket = useRoomSocket(roomId.value, connectedUserId);

  socketStops.push(
    watch(roomSocket.roomState, (newState: Room | null) => {
      if (!newState) return;
      participants.value = newState.participants.filter(p => p.isConnected);
      tasks.value = newState.tasks ?? [];
      votingCards.value = newState.numberOfCards ?? defaultCards;
      votingTask.value = tasks.value.find(t => t.votingStatus === votingStatus.VOTING) ?? null;
      isScrumMaster.value = newState.participants.some(
        p => String(p.id) === String(connectedUserId) && p.isScrumMaster
      );
    }, { immediate: true })
  );

  socketStops.push(
    watch(roomSocket.isConnected, state => {
      console.log('socket status', state ? 'OPEN' : 'CLOSED');
    }, { immediate: true })
  );
};

const scrollTasksUp = () => {
  if (!canScrollUp.value) return;
  taskWindowStart.value = Math.max(0, taskWindowStart.value - 1);
};

const scrollTasksDown = () => {
  if (!canScrollDown.value) return;
  taskWindowStart.value = Math.min(
    tasks.value.length - visibleTaskLimit,
    taskWindowStart.value + 1
  );
};

watch(tasks, () => {
  const maxStart = Math.max(0, tasks.value.length - visibleTaskLimit);
  if (taskWindowStart.value > maxStart) {
    taskWindowStart.value = maxStart;
  }
});

watch(votingTask, newTask => {
  if (!newTask) {
    return;
  }
  const index = tasks.value.findIndex(
    task => String(task.id) === String(newTask.id)
  );
  if (index === -1) {
    return;
  }
  if (index < taskWindowStart.value) {
    taskWindowStart.value = index;
  } else if (index >= taskWindowStart.value + visibleTaskLimit) {
    taskWindowStart.value = index - visibleTaskLimit + 1;
  }
});
</script>

<template>
  <div class="min-h-screen bg-slate-50 py-10 px-4 sm:px-6 lg:px-8">
    <div class="mx-auto grid max-w-6xl gap-6 lg:grid-cols-[340px_1fr]">

      <!-- COLUNA DA ESQUERDA: INFORMAÇÕES -->
      <aside class="space-y-6">
        <section class="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm">
          <div class="mb-5 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-slate-900">Tarefa Atual</h2>
          <button 
            v-if="isScrumMaster" 
            @click="showAdicionarTarefa = true" 
            class="inline-flex items-center justify-center rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold leading-none text-white shadow-sm transition-colors hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/60 focus:ring-offset-1"
          >
            <span>Adicionar</span>
          </button>
          <ModalAdicionarTarefa v-if="showAdicionarTarefa" @close="showAdicionarTarefa = false" @add-task="addTask" />
        </div>

        <div class="space-y-3">
          <button
            v-if="showTaskNavigation"
            type="button"
            class="flex w-full items-center justify-center rounded-lg border border-slate-200 bg-white py-1 text-sm font-medium text-slate-600 transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="!canScrollUp"
            @click="scrollTasksUp"
            aria-label="Mostrar tarefas anteriores"
          >
            ↑
          </button>

          <div class="space-y-4">
            <div
              v-for="task in visibleTasks"
              :key="task.id"
              class="flex items-center justify-between rounded-xl border border-slate-200 bg-white px-4 py-4 shadow-sm transition-all duration-200 hover:-translate-y-1 hover:shadow-lg"
              :class="{
                'border-blue-200 bg-blue-50': task.votingStatus === votingStatus.VOTING
              }"
            >
              <div class="min-w-0">
                <h3 class="truncate text-base font-semibold text-slate-900">{{ task.title }}</h3>
                <span
                  v-if="task.votingStatus === votingStatus.VOTING"
                  class="mt-1 inline-flex items-center rounded-full bg-blue-100 px-2 py-0.5 text-[11px] font-semibold uppercase tracking-wide text-blue-700"
                >
                  Votando...
                </span>
                <p
                  v-if="task.isCompleted && getTaskAverage(task)"
                  class="mt-2 text-sm font-medium text-slate-600"
                >
                  Média: <span class="text-slate-900">{{ getTaskAverage(task) }}</span>
                </p>
              </div>
              <button
              v-if="isScrumMaster"
              type="button"
              @click="selectTaskForVoting(task.id)"
              :disabled="isTaskDisabled(task)"
              :class="[
                'flex items-center justify-center rounded-lg px-4 py-2 text-sm font-semibold leading-none transition-colors text-center',
                isTaskDisabled(task)
                  ? 'cursor-not-allowed bg-slate-300 text-slate-500 opacity-60'
                  : 'bg-emerald-500 text-white shadow-sm hover:bg-emerald-400'
              ]">
                Votar
              </button>
            </div>
          </div>

          <button
            v-if="showTaskNavigation"
            type="button"
            class="flex w-full items-center justify-center rounded-lg border border-slate-200 bg-white py-1 text-sm font-medium text-slate-600 transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="!canScrollDown"
            @click="scrollTasksDown"
            aria-label="Mostrar próximas tarefas"
          >
            ↓
          </button>
        </div>
        <div v-if="tasks.length === 0" class="text-sm text-slate-500">Nenhuma tarefa selecionada.</div>
      </section>


        <section class="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm">
          <div class="mb-5 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-slate-900">Participantes</h2>
          <button
            v-if="isScrumMaster"
            @click="revealCards"
            class="inline-flex items-center justify-center rounded-lg bg-violet-600 px-4 py-2 text-sm font-semibold leading-none text-white shadow-sm transition-colors hover:bg-violet-500 focus:outline-none focus:ring-2 focus:ring-violet-500/60 focus:ring-offset-1"
          >
            <span>Revelar</span>
          </button>
          </div>
          <ul class="space-y-3">
            <li
              v-for="p in participants"
              :key="p.id"
              class="flex items-center justify-between rounded-xl border border-slate-200 bg-white px-4 py-3"
            >
              <div class="flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-blue-500 to-indigo-500 text-sm font-semibold text-white">
                  {{ p.name.charAt(0).toUpperCase() }}
                </div>
                <span class="text-sm font-medium text-slate-700">{{ p.name }}</span>
              </div>
              <span
                class="inline-flex h-8 w-8 items-center justify-center rounded-full"
                :class="votingTask?.votes[p.id] ? 'bg-emerald-100 text-emerald-600' : 'bg-slate-100 text-slate-500'"
              >
                {{ votingTask?.votes[p.id] ? '✓' : '…' }}
              </span>
            </li>
          </ul>
        </section>
      </aside>

      <!-- COLUNA DA DIREITA: AÇÃO -->
      <main class="rounded-2xl border border-transparent bg-white/90 p-8 shadow-lg">
        <section class="space-y-6">
          <div class="flex flex-wrap items-center justify-between gap-4">
            <h2 class="text-2xl font-semibold text-slate-900">Escolha sua carta</h2>
            <div class="flex flex-wrap items-center gap-3">
          <button 
            v-if="isScrumMaster" 
            @click="showModalCartas = true" 
            class="rounded-lg bg-emerald-500 px-4 py-2 text-sm font-semibold leading-none text-white shadow-sm transition-colors hover:bg-emerald-400 focus:outline-none focus:ring-2 focus:ring-emerald-500/60 focus:ring-offset-1"
          >
          Editar Cartas
          </button>
          <button 
            :disabled="!myVote || !votingTask" 
            @click="confirmVote" 
            class="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/60 focus:ring-offset-1 disabled:cursor-not-allowed disabled:border disabled:border-slate-200 disabled:bg-slate-200 disabled:text-slate-500"
          >
          Confirmar Voto
          </button>
          </div>
          <ModalCadastrarCards @save="editarCards" v-if="showModalCartas" @close="showModalCartas = false" />
        </div>
        
        <div class="grid grid-cols-2 gap-4 p-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
          <button
              v-for="card in votingCards"
              :key="card"
              @click="handleVote(card)"
              :class="[
              'flex aspect-[3/4] items-center justify-center rounded-2xl border border-slate-200 bg-white text-2xl font-semibold text-slate-800 shadow-sm transition-all duration-200 hover:-translate-y-1 hover:shadow-lg hover:border-indigo-200',
              myVote === card
                ? 'border-indigo-500 bg-indigo-600 text-white shadow-lg'
                : ''
            ]"
          >
            {{ card }}
          </button>
        </div>

        <div
          v-if="!votingTask && lastAverage"
          class="mt-4 rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-center text-sm font-semibold text-slate-700"
        >
          Média da última votação: <span class="text-slate-900">{{ lastAverage }}</span>
        </div>
        
       
      </section>
      </main>

      <ModalCadastrarParticipante
        v-if="showRegisterModal"
        @registered="onParticipantRegistered"
        @close="onRegisterCancel"
      />
    </div>
  </div>
</template>
