<template>
  <section class="min-h-screen bg-gradient-to-br from-slate-50 via-white to-slate-100 py-16">
    <div class="mx-auto flex w-full max-w-6xl flex-col gap-14 px-6">
      <div class="grid gap-12 lg:grid-cols-[minmax(0,1fr)_420px] lg:items-start">
        <div class="space-y-8">
          <p class="inline-flex items-center rounded-full bg-blue-100 px-4 py-1 text-sm font-semibold text-blue-700">
            Estime em minutos
          </p>
          <div class="space-y-6">
            <h1 class="text-4xl font-bold text-slate-900 md:text-5xl">
              Planning Poker online, colaborativo e sem burocracia
            </h1>
            <p class="text-lg text-slate-600">
              Crie uma sala, convide sua equipe e conduza estimativas em tempo real. Não precisa de login e todos conectam apenas com o código da sala.
            </p>
          </div>
          <img
            src="~/assets/images/group.svg"
            alt="Ilustração de uma equipe colaborando"
            class="mx-auto w-full max-w-xl"
          >
        </div>

        <div class="space-y-6">
          <div class="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm">
            <CreateRoomCard @create-room="handleNewSession" />
          </div>
          <div class="rounded-2xl border border-slate-200 bg-white/90 p-6 shadow-sm">
            <JoinRoomCard @join-room="handleJoinRoom" />
          </div>
        </div>
      </div>

      <section class="rounded-3xl border border-slate-200 bg-white/80 p-8 shadow-sm backdrop-blur">
        <header class="mb-6 flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
          <h2 class="text-2xl font-semibold text-slate-900">Como funciona a rodada</h2>
          <span class="text-sm font-medium uppercase tracking-wide text-slate-500">Regras rápidas</span>
        </header>
        <div class="grid gap-6 md:grid-cols-3">
          <div class="space-y-3">
            <h3 class="text-lg font-semibold text-slate-800">1. Criar ou entrar na sala</h3>
            <p class="text-sm text-slate-600">
              O facilitador cria a sala e compartilha o código. Os demais acessam a tela inicial, informam o ID e escolhem um nome para entrar.
            </p>
          </div>
          <div class="space-y-3">
            <h3 class="text-lg font-semibold text-slate-800">2. Selecionar tarefa e votar</h3>
            <p class="text-sm text-slate-600">
              Com todos conectados, o Scrum Master define a tarefa em votação. Cada participante escolhe sua carta e pode alterar antes de confirmar.
            </p>
          </div>
          <div class="space-y-3">
            <h3 class="text-lg font-semibold text-slate-800">3. Revelar médias e decidir</h3>
            <p class="text-sm text-slate-600">
              Ao revelar, o sistema exibe todos os votos e calcula a média automaticamente. Use o resultado para debater e ajustar o esforço final.
            </p>
          </div>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import CreateRoomCard from "@/components/ui/cards/CreateRoomCard.vue";
import { useRouter } from '#vue-router';
import JoinRoomCard from "@/components/ui/cards/JoinRoomCard.vue"; 

onMounted(() => {
  document.title = "Planning Poker - Estime Tarefas com Sua Equipe";
});

const router = useRouter();

const handleNewSession = async (sessionName: string, userName: string) => {
  const { data, error } = await useApi('/room', {
      method: 'POST',
      body: { sessionName, userName },
  });

  if (error.value) {
    console.error("Erro ao criar a sala:", error.value);
    return;
  }

  const room = data.value.room;
  const roomId = room?.id ?? room?.ID;
  const userId = room?.participants?.find(p => p.name === userName)?.id ?? room?.Participants?.find(p => p.Name === userName)?.ID;

  if (!roomId || !userId) {
    console.error("Resposta inesperada ao criar sala:", room);
    return;
  }

  sessionStorage.setItem("userId", userId);
  sessionStorage.setItem("isScrumMaster", "true");
  router.push(`/room/${room.id}`);
};

const handleJoinRoom = async({ roomId, userName }) => {
  const { data, error } = await useApi(`/room/${roomId}/join`, {
      method: 'POST',
      body: { userName },
  });

  if (error.value) {
    console.error("Erro ao entrar na sala:", error.value);
    return;
  }

  const userId = data.value.id;

  if (!userId) {
    console.error("Resposta inesperada ao entrar na sala:", data.value);
    return;
  }
  sessionStorage.setItem("userId", userId);
  sessionStorage.setItem("isScrumMaster", "false");
  router.push(`/room/${roomId}`);
};
</script>
